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
		toolImpl: &toolImpl{
			name:    "consul",
			version: "1.17.0",
		},
	}
}

const (
	urlConsulWindowsBinary = "https://releases.hashicorp.com/consul/%s/consul_%s_windows_amd64.zip"
)

func (t *consulTool) IsInstalled() bool {
	return t.success("consul --version")
}

func (t *consulTool) LinuxInstall() error {
	return t.run(`/bin/cp -f ${FILES_DIR}/repo/hashicorp.repo /etc/yum.repos.d/ && yum install -y consul`)
}

func (t *consulTool) WindowsInstall() error {
	fmt.Printf("正在下载%s v%s...\n", t.name, t.version)
	url := fmt.Sprintf(urlConsulWindowsBinary, t.version, t.version)
	tempDir, zipFile, err := AllPlatform().Download(url)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// 获取GOPATH
	installDir, err := AllPlatform().GetGoBinDir()
	if err != nil {
		return err
	}

	// 解压zip文件
	if err = AllPlatform().UnzipSpecific(zipFile, "consul.exe", installDir); err != nil {
		return err
	}

	return nil
}
