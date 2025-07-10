//go:build !windows

package appctl

import (
	"fmt"
	"github.com/spf13/cast"
	"os"
	"syscall"
)

const (
	appQuitCountdown = 3
)

func sendStopSignal(strDaprdPid, strAppPid string) error {
	daprdPid, appPid := cast.ToInt(strDaprdPid), cast.ToInt(strAppPid)

	if daprdPid == 0 || appPid == 0 {
		return fmt.Errorf("invalid pid, daprdPid: %s, appPid: %s", strDaprdPid, strAppPid)
	}

	// 给app进程发SIGUSR1标识stop信号
	if process, _ := os.FindProcess(appPid); process != nil {
		err := process.Signal(syscall.SIGTERM)
		if err != nil {
			fmt.Printf("send app process terminal signal, pid: %d, err: %v\n", appPid, err)
		} else { // 等待app stop
			for i := appQuitCountdown; i > 0; i-- {
				fmt.Printf("wait app stop: %d seconds\n", i)
				time.Sleep(1 * time.Second) // 阻塞 1 秒
			}
		}
	}

	// 给daprd发送stop信号
	if process, _ := os.FindProcess(daprdPid); process != nil {
		err := process.Signal(syscall.SIGTERM)
		if err != nil {
			fmt.Printf("send daprd process terminal signal, pid: %d, err: %v\n", daprdPid, err)
		}
	}

	return nil
}
