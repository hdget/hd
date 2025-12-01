package tools

import (
	"fmt"
	"os"

	"github.com/hdget/hd/g"
)

type consulTool struct {
	*toolImpl
}

func Consul() Tool {
	return &consulTool{
		toolImpl: newTool(
			&g.ToolConfig{
				Name:            "consul",
				Version:         "1.20.5",
				UrlWinRelease:   "https://releases.hashicorp.com/consul/1.20.5/consul_1.20.5_windows_amd64.zip",
				UrlLinuxRelease: "https://releases.hashicorp.com/consul/1.20.5/consul_1.20.5_linux_amd64.zip",
			},
		),
	}
}

func (t *consulTool) IsInstalled() bool {
	return t.success("consul --version")
}

func (t *consulTool) LinuxInstall() error {
	err := t.copyFile("etc/yum.repos.d/hashicorp.repo")
	if err != nil {
		return err
	}

	return t.run(`yum install -y consul`)
}

func (t *consulTool) WindowsInstall() error {
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
	if err = AllPlatform().UnzipMatchedFiles(zipFile, "consul.exe", t.GetSystemBinDir()); err != nil {
		return err
	}

	return nil
}
