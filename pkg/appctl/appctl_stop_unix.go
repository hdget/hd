//go:build !windows

package appctl

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"os"
	"syscall"
	"time"
)

func sendStopSignal(strDaprdPid, strAppPid string) error {
	daprdPid, appPid := cast.ToInt(strDaprdPid), cast.ToInt(strAppPid)

	if daprdPid == 0 || appPid == 0 {
		return fmt.Errorf("invalid pid, daprdPid: %s, appPid: %s", strDaprdPid, strAppPid)
	}

	appProcess, err := os.FindProcess(appPid)
	if err != nil {
		return errors.Wrapf(err, "找不到APP进程, pid: %d", appPid)
	}

	// 给app进程发SIGUSR1标识stop信号
	err = appProcess.Signal(syscall.SIGUSR1)
	if err != nil {
		return errors.Wrapf(err, "APP进程发送退出信号, pid: %d", appPid)
	}

	// 等待app stop
	time.Sleep(3 * time.Second)

	// 给daprd发送stop信号
	daprdProcess, err := os.FindProcess(daprdPid)
	if err != nil {
		return errors.Wrapf(err, "找不到Daprd进程, pid: %d", daprdPid)
	}

	err = daprdProcess.Signal(syscall.SIGTERM)
	if err != nil {
		return errors.Wrapf(err, "Daprd进程发送退出信号, pid: %d", daprdPid)
	}

	fmt.Println("app:", appProcess.Pid)
	fmt.Println("daprd:", parentProcess.Pid)

	return nil
}
