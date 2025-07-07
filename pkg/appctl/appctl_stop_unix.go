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

const (
	appQuitCountdown = 3
)

func sendStopSignal(strDaprdPid, strAppPid string) error {
	daprdPid, appPid := cast.ToInt(strDaprdPid), cast.ToInt(strAppPid)

	if daprdPid == 0 || appPid == 0 {
		return fmt.Errorf("invalid pid, daprdPid: %s, appPid: %s", strDaprdPid, strAppPid)
	}

	appProcess, err := os.FindProcess(appPid)
	if err != nil {
		return errors.Wrapf(err, "find app process, pid: %d", appPid)
	}

	// 给app进程发SIGUSR1标识stop信号
	err = appProcess.Signal(syscall.SIGUSR1)
	if err != nil {
		return errors.Wrapf(err, "send app process stop signal, pid: %d", appPid)
	}

	// 等待app stop
	for i := appQuitCountdown; i > 0; i-- {
		fmt.Printf("wait app stop: %d second\n", i)
		time.Sleep(1 * time.Second) // 阻塞 1 秒
	}

	// 给daprd发送stop信号
	daprdProcess, err := os.FindProcess(daprdPid)
	if err != nil {
		return errors.Wrapf(err, "find daprd process, pid: %d", daprdPid)
	}

	err = daprdProcess.Signal(syscall.SIGTERM)
	if err != nil {
		return errors.Wrapf(err, "send daprd process terminal signal, pid: %d", daprdPid)
	}

	fmt.Println("app:", appProcess.Pid)
	fmt.Println("daprd:", daprdProcess.Pid)

	return nil
}
