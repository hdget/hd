package tools

import (
	"fmt"
	"github.com/bitfield/script"
	"os"
)

type daprTool struct {
	*toolImpl
}

const (
	defaultDaprVersion         = "1.15.3"
	defaultUrlDaprLinuxRelease = "https://github.com/dapr/installer-bundle/releases/download/v%s/daprbundle_linux_amd64.tar.gz"
	defaultUrlDaprWinRelease   = "https://github.com/dapr/installer-bundle/releases/download/v%s/daprbundle_windows_amd64.zip"
)

func Dapr() Tool {
	return &daprTool{
		toolImpl: newTool(
			"dapr",
			defaultDaprVersion,
			fmt.Sprintf(defaultUrlDaprLinuxRelease, defaultDaprVersion),
			fmt.Sprintf(defaultUrlDaprWinRelease, defaultDaprVersion),
		),
	}
}

func (t *daprTool) IsInstalled() bool {
	_, err := script.Exec(`dapr --version`).String()
	return err == nil
}

func (t *daprTool) LinuxInstall() error {
	tempDir, zipFile, err := AllPlatform().Download(t.urlLinuxRelease)
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
	tempDir, zipFile, err := AllPlatform().Download(t.urlWinRelease)
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
