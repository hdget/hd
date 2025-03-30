package protorefine

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path"
	"path/filepath"
)

type ProtoRefiner interface {
	Refine(argument *Argument) error // 精简proto files
}

type protoRefineImpl struct {
	srcDir string
}

const (
	envHdNamespace = "HD_NAMESPACE"
)

func New(srcDir string) ProtoRefiner {
	return &protoRefineImpl{
		srcDir: srcDir,
	}
}

// Refine 精简proto下的文件
func (impl *protoRefineImpl) Refine(arg *Argument) error {
	project, exists := os.LookupEnv(envHdNamespace)
	if !exists {
		return errors.New("project name not found")
	}

	if arg == nil {
		return errors.New("invalid argument")
	}

	if err := arg.validate(); err != nil {
		return err
	}

	pbImportPath := path.Join(project, arg.OutputDir, arg.OutputPackage)

	fmt.Println(1)

	// 从golang源代码中找到protobuf类型的变量类型
	golangPkgName, golangTypeNames, err := newGolangParser().parse(impl.srcDir, pbImportPath, arg.OutputDir)
	if err != nil {
		return err
	}

	if len(golangTypeNames) == 0 {
		return errors.New("golang protobuf type not found")
	}

	fmt.Println(2, len(golangTypeNames))

	// 去匹配proto文件中和源文件中匹配的类型的声明
	protoDeclares, err := newProtoParser().findProtoDeclares(arg.ProtoDir, golangPkgName, golangTypeNames)
	if err != nil {
		return err
	}

	fmt.Println(3)

	// 生成新的proto文件以及其依赖文件
	absOutputDir := filepath.Join(arg.OutputDir, arg.OutputPackage)
	outputFileName := filepath.Base(impl.srcDir)
	err = newProtoRender().Render(arg.ProtoDir, absOutputDir, outputFileName, golangPkgName, protoDeclares)
	if err != nil {
		return err
	}

	return nil
}
