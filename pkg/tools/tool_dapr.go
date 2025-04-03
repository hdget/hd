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
	urlDaprLinuxBinary   = "https://newaigou-public.oss-cn-shanghai.aliyuncs.com/download/dapr/%s/daprbundle_linux_amd64.tar.gz"
	urlDaprWindowsBinary = "https://newaigou-public.oss-cn-shanghai.aliyuncs.com/download/dapr/%s/daprbundle_windows_amd64.zip"
)

func Dapr() Tool {
	return &daprTool{
		toolImpl: &toolImpl{
			name: "dapr",
		},
	}
}

func (t *daprTool) IsInstalled() bool {
	_, err := script.Exec(`dapr --version`).String()
	return err == nil
}

func (t *daprTool) LinuxInstall() error {
	fmt.Printf("正在下载%s v%s...\n", t.name, t.version)
	url := fmt.Sprintf(urlDaprLinuxBinary, t.version)
	tempDir, zipFile, err := AllPlatform().Download(url)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// 解压zip文件
	if err = AllPlatform().UnzipSpecific(zipFile, "daprbundle/dapr", t.GetSystemBinDir()); err != nil {
		return err
	}

	return nil
}

func (t *daprTool) WindowsInstall() error {
	fmt.Printf("正在下载%s v%s...\n", t.name, t.version)
	url := fmt.Sprintf(urlDaprWindowsBinary, t.version)
	tempDir, zipFile, err := AllPlatform().Download(url)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// 解压zip文件
	if err = AllPlatform().UnzipSpecific(zipFile, "daprbundle/dapr.exe", t.GetSystemBinDir()); err != nil {
		return err
	}

	return nil
}
