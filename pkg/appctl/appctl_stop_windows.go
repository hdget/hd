//go:build windows

package appctl

func sendStopSignal(pid int) error {
	return nil
}
