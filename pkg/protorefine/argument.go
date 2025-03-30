package protorefine

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

type Argument struct {
	OutputDir          string   // 输出目录
	OutputPackage      string   // 包名
	ProtoDir           string   // 指定的proto所在的目录
	ProtoDirMatchFiles []string // 在没有指定protoDir的时候去通过matchFiles去找proto所在的文件目录
}

func (arg *Argument) validate() error {
	// check proto dir is readable or not
	if err := dirExistsAndReadable(arg.ProtoDir); err != nil {
		return errors.Wrap(err, "invalid proto dir")
	}

	if err := isValidDirName(arg.OutputDir); err != nil {
		return errors.Wrap(err, "invalid output dir")
	}

	if arg.OutputPackage == "" {
		return errors.New("output package name must be specified")
	}

	if arg.ProtoDir == "" && len(arg.ProtoDirMatchFiles) == 0 {
		return errors.New("proto dir match files must be specified")
	}

	return nil
}

func dirExistsAndReadable(dirPath string) error {
	if err := isValidDirName(dirPath); err != nil {
		return err
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return errors.Wrapf(err, "get abs dir, dir: %s", dirPath)
	}

	// 检查是否为目录
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.Wrapf(err, "path doesn't exist, path: %s", dirPath) // 目录不存在
		}
		return errors.Wrapf(err, "path is not accessable, path: %s", dirPath)
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("path is not a directory, path: %s", dirPath)
	}

	// 检查目录是否可读
	d, err := os.Open(absPath)
	if err != nil {
		return errors.Wrapf(err, "path is not readable, path: %s", dirPath)
	}
	defer d.Close()

	return nil
}

func isValidDirName(name string) error {
	// 检查空值或空白
	if strings.TrimSpace(name) == "" {
		return errors.New("dir is empty")
	}

	if name == "." || name == ".." {
		return nil
	}

	// 检查是否包含路径分隔符（如 `/` 或 `\`）
	if strings.ContainsAny(name, `/\`) {
		return errors.New("dir not contains path separator")
	}

	// 尝试拼接临时路径并验证
	tempPath := filepath.Join(os.TempDir(), name)
	_, err := filepath.Abs(tempPath)
	if err != nil {
		return errors.Wrap(err, "invalid dir")
	}
	return nil
}
