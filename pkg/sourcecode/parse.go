package sourcecode

import (
	"fmt"
	"github.com/hdget/common/protobuf"
	"github.com/hdget/hd/g"
	"github.com/pkg/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type Parser interface {
	Parse() (*sourceCodeInfo, error)
}

// astPackage 表示解析后的包信息
type astPackage struct {
	Name    string
	Files   map[string]*ast.File
	FileSet *token.FileSet
}

type parserImpl struct {
	srcDir            string
	pkgRelPath2astPkg map[string]*astPackage
}

type sourceCodeInfo struct {
	serverEntryFilePath    string                  // appServer.Run的入口文件即appServer开始运行所在的go文件
	daprModules            []*parsedDaprModuleInfo // 获取DaprInvocationModule
	daprInvocationHandlers []*protobuf.DaprHandler // DaprInvocationModuleHandler
}

func newParser(srcDir string, excludeDirs []string) (Parser, error) {
	pkgRelPath2astPkg, err := recursiveParseDir(srcDir, excludeDirs)
	if err != nil {
		return nil, err
	}

	return &parserImpl{srcDir: srcDir, pkgRelPath2astPkg: pkgRelPath2astPkg}, nil
}

func (p *parserImpl) Parse() (*sourceCodeInfo, error) {
	fmt.Println(">>> parse source code")

	serverEntryPath, err := p.parseServerEntryFilePath()
	if err != nil {
		return nil, err
	}
	if serverEntryPath == "" {
		return nil, errors.New("can't find server.Start call")
	}

	if g.Debug {
		fmt.Println("server entry file path:", serverEntryPath)
	}

	daprModules, err := p.parseDaprModules()
	if err != nil {
		return nil, err
	}
	if len(daprModules) == 0 {
		return nil, errors.New("can't find dapr modules")
	}

	if g.Debug {
		for _, module := range daprModules {
			fmt.Printf(" * found modules, path: %s, name: %s\n", module.pkgRelPath, module.name)
		}
	}

	daprInvocationHandlers, err := p.parseDaprInvocationHandlers(daprModules)
	if err != nil {
		return nil, err
	}
	if g.Debug {
		fmt.Println("found handlers")
		for _, m := range daprModules {
			handlerNames := make([]string, 0)
			for _, h := range daprInvocationHandlers {
				if m.pkgRelPath == h.PkgPath && m.name == h.Module {
					handlerNames = append(handlerNames, h.Name)
				}
			}
			fmt.Printf(" * found handlers, path: %s, module: %s, total: %d, handlers: %v\n", m.pkgRelPath, m.name, len(handlerNames), handlerNames)
		}
	}

	scInfo := &sourceCodeInfo{
		serverEntryFilePath:    serverEntryPath,
		daprModules:            daprModules,
		daprInvocationHandlers: daprInvocationHandlers,
	}

	return scInfo, nil
}

// parseDir 递归解析整个项目目录
func recursiveParseDir(rootDir string, excludeDirs []string) (map[string]*astPackage, error) {
	// 获取根目录绝对路径
	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, err
	}

	// 转换排除目录为绝对路径
	var absExcludeDirs []string
	for _, dir := range excludeDirs {
		absDir := filepath.Join(absRoot, dir)
		absExcludeDirs = append(absExcludeDirs, absDir)
	}

	// 启动工作 goroutine
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 使用 WaitGroup 等待所有解析完成
	fileChan := make(chan string, 100)
	fset := token.NewFileSet()
	pkgRelPath2astPkg := make(map[string]*astPackage)
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range fileChan {
				pkgRelPath, _ := filepath.Rel(absRoot, filePath)
				pkgRelPath = filepath.ToSlash(pkgRelPath)

				// 解析单个文件
				pkgName2astPkg, err := parser.ParseDir(fset, filePath, nil, parser.ParseComments)
				if err != nil {
					continue
				}

				mu.Lock()
				for pkgName, astPkg := range pkgName2astPkg {
					pkgRelPath2astPkg[pkgRelPath] = &astPackage{
						Name:    pkgName,
						Files:   astPkg.Files,
						FileSet: fset,
					}
					break
				}
				mu.Unlock()
			}
		}()
	}

	// 遍历所有子目录
	err = filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否在排除目录中
		for _, excludeDir := range absExcludeDirs {
			if strings.HasPrefix(path, excludeDir) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// 只处理 .go 文件
		if info.IsDir() {
			fileChan <- path
		}

		return nil
	})

	close(fileChan)
	wg.Wait()

	if err != nil {
		return nil, errors.Wrapf(err, "traverse dir, srcDir: %s", rootDir)
	}

	return pkgRelPath2astPkg, nil
}
