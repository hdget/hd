package tools

import (
	"fmt"
	"github.com/bitfield/script"
	"runtime"
)

type Tool interface {
	GetName() string
	IsInstalled() bool
	LinuxInstall() error
	WindowsInstall() error
}

type toolImpl struct {
	name    string
	version string
}

func (t *toolImpl) GetName() string {
	return t.name
}

func (t *toolImpl) IsInstalled() bool {
	var cmd string
	if runtime.GOOS == "windows" {
		cmd = fmt.Sprintf("where %s", t.GetName())
	} else {
		cmd = fmt.Sprintf("which %s", t.GetName())
	}

	_, err := script.Exec(cmd).String()
	return err == nil
}

func Check(tools []Tool, debug bool) error {
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
