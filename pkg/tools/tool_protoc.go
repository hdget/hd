package tools

import (
	"fmt"

	"github.com/bitfield/script"
	"github.com/hdget/hd/g"

	"os"
)

type protocTool struct {
	*toolImpl
}

func Protoc() Tool {
	return &protocTool{
		toolImpl: newTool(
			&g.ToolConfig{
				Name:            "protoc",
				Version:         "30.2",
				UrlWinRelease:   "https://github.com/protocolbuffers/protobuf/releases/download/v30.2/protoc-30.2-win64.zip",
				UrlLinuxRelease: "https://github.com/protocolbuffers/protobuf/releases/download/v30.2/protoc-30.2-linux-x86_64.zip",
			},
		),
	}
}

func (t *protocTool) IsInstalled() bool {
	_, err := script.Exec(`protoc --version`).String()
	return err == nil
}

func (t *protocTool) LinuxInstall() error {
	tempDir, zipFile, err := AllPlatform().Download(t.config.UrlLinuxRelease)
	if err != nil {
		return err
	}
	defer func() {
		if e := os.RemoveAll(tempDir); e != nil {
			fmt.Printf("delete temp dir failed: %v, dir: %s", e, tempDir)
		}
	}()

	// 获取GOPATH
	installDir, err := AllPlatform().GetGoBinDir()
	if err != nil {
		return err
	}

	// 解压zip文件
	if err = AllPlatform().UnzipMatchedFiles(zipFile, "bin/protoc", installDir); err != nil {
		return err
	}

	return nil
}

func (t *protocTool) WindowsInstall() error {
	tempDir, zipFile, err := AllPlatform().Download(t.config.UrlWinRelease)
	if err != nil {
		return err
	}
	defer func() {
		if e := os.RemoveAll(tempDir); e != nil {
			fmt.Printf("delete temp dir failed: %v, dir: %s", e, tempDir)
		}
	}()

	// 获取GOPATH
	installDir, err := AllPlatform().GetGoBinDir()
	if err != nil {
		return err
	}

	// 解压zip文件
	if err = AllPlatform().UnzipMatchedFiles(zipFile, "bin/protoc.exe", installDir); err != nil {
		return err
	}

	return nil
}
