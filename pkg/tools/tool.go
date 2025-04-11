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

	if c, exist := g.ToolConfigs[name]; exist {
		fmt.Println("here:", c)

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
				fmt.Printf("%s installed\n", t.GetName())
			}
			continue
		}

		fmt.Printf("%s not installed\n", t.GetName())
		if err := installTool(t); err != nil {
			return fmt.Errorf("%s install failed: %v", t.GetName(), err)
		}
	}
	return nil
}

func (impl *toolImpl) run(cmd string) error {
	output, err := script.Exec(cmd).String()
	if err != nil {
		return errors.Wrapf(err, "%s install failed, err: %s", impl.name, output)
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
	fmt.Printf("try to install %s...\n", t.GetName())

	var err error
	switch runtime.GOOS {
	case "linux", "darwin":
		err = t.LinuxInstall()
	case "windows":
		err = t.WindowsInstall()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
	if err != nil {
		return err
	}

	if !t.IsInstalled() {
		return fmt.Errorf("%s is not avaialble after installation, please check execute path", t.GetName())
	}

	fmt.Printf("%s install succeed\n", t.GetName())
	return nil
}
