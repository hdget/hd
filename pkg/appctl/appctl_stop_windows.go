//go:build windows

package appctl

func sendStopSignal(daprdPid, appPid string) error {
	return nil
}
