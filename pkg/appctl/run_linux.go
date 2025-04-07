//go:build linux || darwin

package appctl

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func (a *appCtlImpl) run(appId, command string, healthCheck func() bool, timeout time.Duration) error {
	// 1. 设置日志文件
	logPath := fmt.Sprintf("/var/log/%s.log", appId)
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("无法创建日志文件: %v", err)
	}
	defer logFile.Close()

	// 2. 创建命令
	cmd := exec.Command(command)
	cmd.Env = os.Environ()

	// 3. 设置进程属性 - 完全脱离终端控制
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid:  true,  // 创建新会话
		Setctty: false, // 不控制终端
		Pgid:    0,     // 新进程组
	}

	// 4. 重定向标准输入/输出/错误
	cmd.Stdin = nil // 关闭输入
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	// 5. 启动进程（后台运行）
	fmt.Println(cmd.Environ())
	if err = cmd.Start(); err != nil {
		return fmt.Errorf("启动失败: %v", err)
	}

	if healthCheck != nil {
		return runHealthCheck(healthCheck, timeout, cmd)
	}

	fmt.Printf("进程已启动 (PID: %d), 日志: %s\n", cmd.Process.Pid, logPath)
	return nil
}
