//go:build !windows

package appctl

import (
	"github.com/pkg/errors"
	"os"
)

func sendStopSignal(pid int) error {
	if pid == 0 {
		return errors.New("invalid pid")
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return errors.Wrapf(err, "找不到进程, pid: %d", pid)
	}

	err = process.Signal(syscall.SIGUSR1)
	if err != nil {
		return errors.Wrapf(err, "无法终止进程, pid: %d", pid)
	}

	return nil
}
