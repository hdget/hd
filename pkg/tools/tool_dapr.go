package tools

import (
	"fmt"
	"os"

	"github.com/bitfield/script"
	"github.com/hdget/hd/g"
)

type daprTool struct {
	*toolImpl
}

func Dapr() Tool {
	return &daprTool{
		toolImpl: newTool(
			&g.ToolConfig{
				Name:            "dapr",
				Version:         "1.15.3",
				UrlWinRelease:   "https://github.com/dapr/installer-bundle/releases/download/v1.15.3/daprbundle_windows_amd64.zip",
				UrlLinuxRelease: "https://github.com/dapr/installer-bundle/releases/download/v1.15.3/daprbundle_linux_amd64.tar.gz",
			},
		),
	}
}

func (t *daprTool) IsInstalled() bool {
	_, err := script.Exec(`dapr --version`).String()
	return err == nil
}

func (t *daprTool) LinuxInstall() error {
	tempDir, zipFile, err := AllPlatform().Download(t.config.UrlLinuxRelease)
	if err != nil {
		return err
	}
	defer func() {
		if e := os.RemoveAll(tempDir); e != nil {
			fmt.Printf("delete temp dir failed: %v, dir: %s", e, tempDir)
		}
	}()

	// 解压zip文件
	if err = AllPlatform().UnzipMatchedFiles(zipFile, "daprbundle/dapr", t.GetSystemBinDir()); err != nil {
		return err
	}

	return nil
}

func (t *daprTool) WindowsInstall() error {
	tempDir, zipFile, err := AllPlatform().Download(t.config.UrlWinRelease)
	if err != nil {
		return err
	}
	defer func() {
		if e := os.RemoveAll(tempDir); e != nil {
			fmt.Printf("delete temp dir failed: %v, dir: %s", e, tempDir)
		}
	}()

	// 解压zip文件
	if err = AllPlatform().UnzipMatchedFiles(zipFile, "daprbundle/dapr.exe", t.GetSystemBinDir()); err != nil {
		return err
	}

	return nil
}
