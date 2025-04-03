package tools

import (
	"fmt"
	"github.com/hdget/hd/g"
	"os"
)

type consulTool struct {
	*toolImpl
}

func Consul() Tool {
	impl := &consulTool{
		toolImpl: &toolImpl{
			name:            "consul",
			version:         defaultConsulVersion,
			urlLinuxRelease: fmt.Sprintf(defaultUrlConsulUnixRelease, defaultConsulVersion, defaultConsulVersion),
			urlWinRelease:   fmt.Sprintf(defaultUrlConsulWinRelease, defaultConsulVersion, defaultConsulVersion),
		},
	}

	if c, exist := g.ToolConfigs["consul"]; exist {
		if c.UrlLinuxRelease != "" {
			impl.urlLinuxRelease = c.UrlLinuxRelease
		}

		if c.UrlWinRelease != "" {
			impl.urlWinRelease = c.UrlWinRelease
		}
		if impl.version != "" {
			impl.version = c.Version
		}
	}

	return impl
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
	fmt.Printf("正在下载%s...\n", t.name)

	tempDir, zipFile, err := AllPlatform().Download(t.urlWinRelease)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// 解压zip文件
	if err = AllPlatform().UnzipSpecific(zipFile, "consul.exe", t.GetSystemBinDir()); err != nil {
		return err
	}

	return nil
}
