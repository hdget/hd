//go:build linux || darwin

package tools

import (
	"fmt"
	"github.com/pkg/errors"
	"mvdan.cc/sh/v3/shell"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func RunDaemon(name, command string, healthCheck func() bool, timeout time.Duration) error {
	// 1. 设置日志文件
	logPath := fmt.Sprintf("/var/log/%s.log", name)
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return errors.Wrapf(err, "create log file, logFile: %s", logPath)
	}
	defer logFile.Close()

	// 2. 创建命令
	args, err := shell.Fields(command, nil)
	if err != nil {
		return err
	}

	// Windows 设置
	cmd := exec.Command(args[0], args[1:]...)

	// 3. 设置进程属性 - 完全脱离终端控制
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid:     true,  // 创建新会话
		Setctty:    false, // 不控制终端
		Foreground: false, // 不在前台运行
		Pgid:       0,     // 新进程组
	}

	// 4. 重定向标准输入/输出/错误
	cmd.Stdin = nil // 关闭输入
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	// 5. 启动进程（后台运行）
	if err = cmd.Start(); err != nil {
		return errors.Wrapf(err, "failed start, command: %s", command)
	}

	if healthCheck != nil {
		return runHealthCheck(healthCheck, timeout, cmd)
	}

	fmt.Printf("%s(PID: %d) started, log path: %s\n", name, cmd.Process.Pid, logPath)
	return nil
}
