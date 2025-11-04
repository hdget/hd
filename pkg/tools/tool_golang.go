package tools

import (
	"fmt"
	"path/filepath"

	"github.com/hdget/hd/g"
)

type golangTool struct {
	*toolImpl
}

func Golang() Tool {
	return &golangTool{
		toolImpl: &toolImpl{
			&g.ToolConfig{
				Name:            "go",
				Version:         "1.25.3",
				UrlWinRelease:   "https://golang.google.cn/dl/go1.25.3.windows-amd64.zip",
				UrlLinuxRelease: "https://golang.google.cn/dl/go1.25.3.linux-amd64.tar.gz",
			},
		},
	}
}

const (
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
	cmd := fmt.Sprintf(cmdLinuxInstallGolang, t.config.UrlLinuxRelease, filepath.Base(t.config.UrlLinuxRelease))
	return t.run(cmd)
}

func (t *golangTool) WindowsInstall() error {
	fmt.Println("Please install Golang manually, e,g: https://golang.google.cn/dl")
	return nil
}
