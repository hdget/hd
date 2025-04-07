package appctl

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func runHealthCheck(healthCheck func() bool, timeout time.Duration, cmd *exec.Cmd) error {
	// 创建健康检查超时上下文
	healthCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 健康检查循环
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-healthCtx.Done():
			// 超时后杀死进程
			_ = cmd.Process.Kill()
			return fmt.Errorf("health check timeout after %v, cmd: %s", timeout, strings.Join(cmd.Args, " "))

		case <-ticker.C:
			if healthCheck() {
				// 启动成功
				return nil
			}

			fmt.Println("xxxxxxx")

			// 检查进程是否仍在运行
			if cmd.Process == nil || (cmd.ProcessState != nil && cmd.ProcessState.Exited()) {
				return fmt.Errorf("process exited")
			}
		}
	}
}
