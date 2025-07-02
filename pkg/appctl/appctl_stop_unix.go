//go:build !windows

package appctl

import (
	"fmt"
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
		return errors.Wrapf(err, "找不到进程, pid: %d", pid)
	}

	// 获取 PGID
	pgid, err := syscall.Getpgid(process.Pid)
	if err != nil {
		return errors.Wrapf(err, "获取pgid失败, pid: %d", pid)
	}

	fmt.Println("==pid:", pid)
	fmt.Println("==pgid:", pgid)

	err = syscall.Kill(-pgid, syscall.SIGTERM) // 注意负号 `-pgid`
	if err != nil {
		return errors.Wrapf(err, "无法终止进程, pgid: %d", pgid)
	}

	return nil
}
