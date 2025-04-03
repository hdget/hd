//go:build windows

package appctl

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

func (a *appCtlImpl) runDetached(appId, command string, healthCheck func() bool, timeout time.Duration) error {
	logDir := filepath.Join(os.Getenv("ProgramData"), "MyApp", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("无法创建日志目录: %v", err)
	}

	logPath := filepath.Join(logDir, fmt.Sprintf("%s.log", appId))

	// 打开日志文件
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("无法打开日志文件: %v", err)
	}
	defer logFile.Close()

	// 创建命令
	cmd := exec.Command(command)

	// Windows 设置
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP, // 显示窗口
		// HideWindow:  true // 隐藏窗口
		// CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | 0x00000008 // CREATE_NO_WINDOW
	}

	// 重定向输出
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Stdin = nil // 关闭输入

	// 启动进程
	if err = cmd.Start(); err != nil {
		return fmt.Errorf("启动命令失败: %v", err)
	}

	if healthCheck != nil {
		return runHealthCheck(healthCheck, timeout, cmd)
	}

	log.Printf("进程已启动 (PID: %d), 日志: %s", cmd.Process.Pid, logPath)
	return nil
}
