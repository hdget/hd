package dapr

import (
	"fmt"

	"github.com/hdget/common/protobuf"
	"github.com/hdget/hd/g"
	astUtils "github.com/hdget/utils/ast"
	"github.com/pkg/errors"
)

type souceCodeInfo struct {
	patchFile          string                  // 需要patch的文件，appServer.Run的入口文件即appServer开始运行所在的go文件
	modules            []*daprModule           // 获取DaprInvocationModule
	invocationHandlers []*protobuf.DaprHandler // DaprInvocationModuleHandler
}

type scParser struct {
	srcDir            string
	pkgRelPath2astPkg map[string]*astUtils.Package
}

func newParser(srcDir string, excludeDirs []string) (*scParser, error) {
	pkgRelPath2astPkg, err := astUtils.InspectPackage(srcDir, excludeDirs)
	if err != nil {
		return nil, err
	}

	return &scParser{srcDir: srcDir, pkgRelPath2astPkg: pkgRelPath2astPkg}, nil
}

func (p *scParser) parse() (*souceCodeInfo, error) {
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

	scInfo := &souceCodeInfo{
		patchFile:          serverEntryPath,
		modules:            daprModules,
		invocationHandlers: daprInvocationHandlers,
	}

	return scInfo, nil
}
