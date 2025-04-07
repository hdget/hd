package tools

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/hdget/hd/g"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"runtime"
)

type Tool interface {
	GetName() string
	IsInstalled() bool
	LinuxInstall() error
	WindowsInstall() error
}

type toolImpl struct {
	name            string
	version         string
	urlWinRelease   string
	urlLinuxRelease string
}

func newTool(name, version, urlWinRelease, urlLinuxRelease string) *toolImpl {
	t := &toolImpl{
		name:            name,
		version:         version,
		urlLinuxRelease: urlWinRelease,
		urlWinRelease:   urlLinuxRelease,
	}
	if c, exist := g.ToolConfigs["consul"]; exist {
		if c.UrlLinuxRelease != "" {
			t.urlLinuxRelease = c.UrlLinuxRelease
		}

		if c.UrlWinRelease != "" {
			t.urlWinRelease = c.UrlWinRelease
		}
		if t.version != "" {
			t.version = c.Version
		}
	}
	return t
}

func (impl *toolImpl) GetName() string {
	return impl.name
}

func (impl *toolImpl) IsInstalled() bool {
	var cmd string
	if runtime.GOOS == "windows" {
		cmd = fmt.Sprintf("where %s", impl.GetName())
	} else {
		cmd = fmt.Sprintf("which %s", impl.GetName())
	}

	_, err := script.Exec(cmd).String()
	return err == nil
}

func Check(debug bool, tools ...Tool) error {
	for _, t := range tools {
		if t.IsInstalled() {
			if debug {
				fmt.Printf("%s已安装\n", t.GetName())
			}
			continue
		}

		fmt.Printf("%s未安装\n", t.GetName())
		if err := installTool(t); err != nil {
			return fmt.Errorf("%s安装失败: %v", t.GetName(), err)
		}
	}
	return nil
}

func (impl *toolImpl) run(cmd string) error {
	output, err := script.Exec(cmd).String()
	if err != nil {
		return errors.Wrapf(err, "%s安装失败, 错误:%s", impl.name, output)
	}
	return nil
}

func (impl *toolImpl) success(cmd string) bool {
	return script.Exec(cmd).Wait() == nil
}

// GetSystemBinDir 获取系统标准bin目录列表
func (impl *toolImpl) GetSystemBinDir() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Microsoft", "WindowsApps")
	}
	return "/usr/local/bin"
}

func installTool(t Tool) error {
	fmt.Printf("尝试安装%s...\n", t.GetName())

	var err error
	switch runtime.GOOS {
	case "linux", "darwin":
		err = t.LinuxInstall()
	case "windows":
		err = t.WindowsInstall()
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
	if err != nil {
		return err
	}

	if !t.IsInstalled() {
		return fmt.Errorf("%s安装后仍然不可用，请检查PATH", t.GetName())
	}

	fmt.Printf("%s安装成功\n", t.GetName())
	return nil
}
