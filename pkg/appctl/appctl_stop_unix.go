//go:build !windows

package appctl

import (
	"github.com/pkg/errors"
	"os"
	"syscall"
)

func sendStopSignal(pid int) error {
	if pid == 0 {
		return errors.New("invalid pid")
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return errors.Wrapf(err, "找不到APP进程, pid: %d", pid)
	}

	// 给APP发SIGUSR1标识stop信号
	err := process.Signal(syscall.SIGUSR1)
	if err != nil {
		return errors.Wrapf(err, "给APP发送退出信号, pid: %d", pid)
	}

	// 给父进程daprd发送term信号
	parentProcess, err := os.FindProcess(syscall.Getppid())
	if err != nil {
		return errors.Wrapf(err, "找不到Daprd进程, pid: %d", syscall.Getppid())
	}

	err = parentProcess.Signal(syscall.SIGTERM)
	if err != nil {
		return errors.Wrapf(err, "无法终止进程, pid: %d", parentProcess.Pid)
	}

	return nil
}
