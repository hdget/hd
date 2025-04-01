package protorefine

import (
	"fmt"
	"github.com/hdget/hd/pkg/utils"
)

type Argument struct {
	GolangModule          string // 源代码的模块名
	GolangSourceCodeDir   string // 源代码的路径
	GolangProtobufPackage string // protobuf包名
	ProtoRepository       string // 原始proto文件所在的目录
	OutputDir             string // 必须是相对路径

}

func (a *Argument) validate() error {
	if a.GolangModule == "" {
		return fmt.Errorf("golang module not found")
	}

	if a.GolangSourceCodeDir == "" {
		return fmt.Errorf("golang source code dir not found")
	}

	if a.GolangProtobufPackage == "" {
		return fmt.Errorf("golang protobuf package not found")
	}

	if !utils.IsValidRelativePath(a.OutputDir) {
		return fmt.Errorf("outputDir must be relative dir, outputDir: %s", a.OutputDir)
	}

	if err := utils.IsDirReadableAndWithFiles(a.ProtoRepository, ".proto"); err != nil {
		return err
	}

	return nil
}
