package protorefine

import (
	"fmt"
	"github.com/hdget/hd/pkg/utils"
	"path/filepath"
)

type Argument struct {
	OutputDir    string // 必须是相对路径
	ProtoDir     string
	absProtoDir  string
	absOutputDir string
}

func (a *Argument) validate(srcDir, matchFile string) error {
	if !utils.IsValidRelativePath(a.OutputDir) {
		return fmt.Errorf("outputDir must be relative sub dir of src dir, srcDir: %s, outputDir: %s", srcDir, a.OutputDir)
	}
	a.absOutputDir, _ = filepath.Abs(a.OutputDir)

	// 处理protoDir
	// 如果没有指定protoDir，尝试智能匹配protoDir
	if a.ProtoDir == "" {
		found, err := utils.FindDirContainingFiles(srcDir, a.absOutputDir, matchFile)
		if err != nil {
			return err
		}
		a.absProtoDir = found
	} else {
		a.absProtoDir, _ = filepath.Abs(a.ProtoDir)
	}

	// validate arguments
	if err := utils.IsDirReadableAndWithFiles(a.absProtoDir, ".proto"); err != nil {
		return err
	}

	return nil
}
