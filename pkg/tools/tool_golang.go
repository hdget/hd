package tools

import (
	"fmt"
	"path/filepath"
)

type golangTool struct {
	*toolImpl
}

func Golang() Tool {
	return &golangTool{
		toolImpl: &toolImpl{
			name: "go",
		},
	}
}

const (
	urlLinuxGolangBinary  = "https://golang.google.cn/dl/go%s.linux-amd64.tar.gz"
	cmdLinuxInstallGolang = `
wget %s && \
rm -rf /usr/local/go && \
tar -C /usr/local -xzf %s && \
mkdir -p $HOME/go/bin`
)

func (t *golangTool) IsInstalled() bool {
	return t.success("go version")
}

func (t *golangTool) LinuxInstall() error {
	url := fmt.Sprintf(urlLinuxGolangBinary, t.version)
	cmd := fmt.Sprintf(cmdLinuxInstallGolang, url, filepath.Base(url))
	return t.run(cmd)
}

func (t *golangTool) WindowsInstall() error {
	fmt.Println("请手动安装Golang, e,g: https://golang.google.cn/dl")
	return nil
}
