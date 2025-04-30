//go:build windows

package tools

import (
	"fmt"
	"github.com/pkg/errors"
	"mvdan.cc/sh/v3/shell"
	"os"
	"os/exec"
	"time"
)

func RunDaemon(name, command string, healthCheck func() bool, timeout time.Duration) error {
	// 创建命令
	args, err := shell.Fields(command, nil)
	if err != nil {
		return err
	}

	// Windows 设置
	cmd := exec.Command(args[0], args[1:]...)
	//cmd.SysProcAttr = &syscall.SysProcAttr{
	//	HideWindow:    true,                                          // 隐藏窗口
	//	CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | 0x00000008, // CREATE_NO_WINDOW
	//}

	// 重定向输出
	cmd.Stdin = nil // 关闭输入
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout

	// 启动进程
	if err = cmd.Start(); err != nil {
		return errors.Wrapf(err, "启动命令失败, cmd:%s", cmd)
	}

	if healthCheck != nil {
		return runHealthCheck(healthCheck, timeout, cmd)
	}

	fmt.Printf("%s(PID: %d) started.\n", name, cmd.Process.Pid)
	return nil
}
