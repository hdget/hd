package protorefine

import (
	"fmt"
	"github.com/hdget/hd/pkg/utils"
)

type Argument struct {
	GolangModule        string // 源代码的模块名
	GolangSourceCodeDir string // 源代码的路径
	ProtoRepository     string // 原始proto文件所在的目录
	OutputDir           string
	OutputPackage       string
}

func (arg Argument) validate() error {
	if arg.GolangModule == "" {
		return fmt.Errorf("empty golang module")
	}

	if arg.GolangSourceCodeDir == "" {
		return fmt.Errorf("empty golang source code dir")
	}

	if arg.OutputPackage == "" {
		return fmt.Errorf("empty output package")
	}

	if !utils.IsValidRelativePath(arg.OutputDir) {
		return fmt.Errorf("outputDir must be relative dir, outputDir: %s", arg.OutputDir)
	}

	if err := utils.IsDirReadableAndWithFiles(arg.ProtoRepository, ".proto"); err != nil {
		return err
	}

	return nil
}
