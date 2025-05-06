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

func newTool(name, version, urlLinuxRelease, urlWinRelease string) *toolImpl {
	t := &toolImpl{
		name:            name,
		version:         version,
		urlLinuxRelease: urlLinuxRelease,
		urlWinRelease:   urlWinRelease,
	}

	if c, exist := g.ToolConfigs[name]; exist {
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

func Check(tools ...Tool) error {
	for _, t := range tools {
		if t.IsInstalled() {
			if g.Debug {
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

// GetSystemBinDir 获取系统标准bin目录列表
func (impl *toolImpl) GetSystemBinDir() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Microsoft", "WindowsApps")
	}
	return "/usr/local/bin"
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

func installTool(t Tool) error {
	fmt.Printf("install %s...\n", t.GetName())

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

//func RunDaemon() error {
//	cmd := exec.Command("your-command", "arg1", "arg2")
//
//	// 重定向输出
//	outFile, err := os.Create("/tmp/cmd.log")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer outFile.Close()
//
//	cmd.Stdout = outFile
//	cmd.Stderr = outFile
//
//	// 设置进程属性
//	cmd.SysProcAttr = &syscall.SysProcAttr{
//		Setsid:     true, // 创建新会话
//		Setctty:    false,
//		Foreground: false,
//	}
//
//	// 启动进程
//	err = cmd.Start()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 可选：记录PID到文件
//	pidFile, err := os.Create("/tmp/cmd.pid")
//	if err != nil {
//		log.Printf("Warning: could not create PID file: %v\n", err)
//	} else {
//		_, _ = pidFile.WriteString(fmt.Sprintf("%d", cmd.Process.Pid))
//		pidFile.Close()
//	}
//
//}
