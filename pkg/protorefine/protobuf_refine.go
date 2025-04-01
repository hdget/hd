package protorefine

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type ProtoRefine interface {
	Refine(*Argument) (string, string, error) // 精简proto files, 返回proto输出目录和packageName
}

type protoRefineImpl struct {
	srcDir string
	debug  bool
}

func New(srcDir string, options ...Option) ProtoRefine {
	impl := &protoRefineImpl{
		srcDir: srcDir,
		debug:  false,
	}

	for _, apply := range options {
		apply(impl)
	}

	return impl
}

// Refine 精简proto下的文件
func (impl *protoRefineImpl) Refine(arg *Argument) (string, string, error) {
	if arg == nil {
		return "", "", fmt.Errorf("argument is nil")
	}

	// obtain root module name
	rootModule, err := getRootModuleName()
	if err != nil {
		return "", "", errors.Wrap(err, "get source module name")
	}

	protoFile := fmt.Sprintf("%s.proto", filepath.Base(rootModule))

	if err := arg.validate(impl.srcDir, protoFile); err != nil {
		return "", "", err
	}

	// get proto file's golang package name
	golangPkgName, err := newProtoParser().getGolangPackageName(arg.absProtoDir)
	if err != nil {
		return "", "", errors.Wrap(err, "parse golang package name")
	}

	pbImportPath := path.Join(rootModule, arg.OutputDir, golangPkgName)
	outputDir := filepath.Join(arg.absOutputDir, filepath.Base(arg.absProtoDir))

	if impl.debug {
		fmt.Println("root module:", rootModule)
		fmt.Println("pb import path:", pbImportPath)
		fmt.Println("source proto dir:", arg.absProtoDir)
		fmt.Println("output proto dir:", outputDir)
	}

	// 实际逻辑
	// 从golang源代码中找到protobuf类型的变量类型
	golangTypeNames, err := newGolangParser().parse(impl.srcDir, pbImportPath, arg.absOutputDir)
	if err != nil {
		return "", "", err
	}

	if len(golangTypeNames) == 0 {
		return "", "", fmt.Errorf("golang protobuf type reference not found")
	}

	// 去匹配proto文件中和源文件中匹配的类型的声明
	protoDeclares, err := newProtoParser().findProtoDeclares(arg.absProtoDir, golangPkgName, golangTypeNames)
	if err != nil {
		return "", "", err
	}

	if len(protoDeclares) == 0 {
		return "", "", fmt.Errorf("empty matched protobuf types, protoDir: %s", arg.absProtoDir)
	}

	// 生成新的proto文件以及其依赖文件
	err = newProtoRender().Render(arg.absProtoDir, outputDir, protoFile, golangPkgName, protoDeclares)
	if err != nil {
		return "", "", err
	}

	return outputDir, golangPkgName, nil
}

func getRootModuleName() (string, error) {
	// 获取根模块名
	cmdOutput, err := exec.Command("go", "list", "-m").CombinedOutput()
	if err != nil {
		return "", err
	}

	// 按换行符拆分结果
	lines := bytes.Split(cmdOutput, []byte("\n"))
	if len(lines) == 0 {
		return "", errors.New("source code may not using go module")
	}

	return strings.TrimSpace(string(lines[0])), nil
}
