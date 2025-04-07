package tools

import (
	"fmt"
	"os"
)

type consulTool struct {
	*toolImpl
}

func Consul() Tool {
	return &consulTool{
		toolImpl: newTool(
			"consul",
			defaultConsulVersion,
			fmt.Sprintf(defaultUrlConsulWinRelease, defaultConsulVersion, defaultConsulVersion),
			fmt.Sprintf(defaultUrlConsulUnixRelease, defaultConsulVersion, defaultConsulVersion),
		),
	}
}

const (
	defaultConsulVersion        = "1.20.5"
	defaultUrlConsulWinRelease  = "https://releases.hashicorp.com/consul/%s/consul_%s_windows_amd64.zip"
	defaultUrlConsulUnixRelease = "https://releases.hashicorp.com/consul/%s/consul_%s_linux_amd64.zip"
)

func (t *consulTool) IsInstalled() bool {
	return t.success("consul --version")
}

func (t *consulTool) LinuxInstall() error {
	return t.run(`/bin/cp -f ${FILES_DIR}/repo/hashicorp.repo /etc/yum.repos.d/ && yum install -y consul`)
}

func (t *consulTool) WindowsInstall() error {
	fmt.Printf("downloading %s...\n", t.name)

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
	if err = AllPlatform().UnzipSpecific(zipFile, "consul.exe", t.GetSystemBinDir()); err != nil {
		return err
	}

	return nil
}
