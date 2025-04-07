package appctl

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/elliotchance/pie/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"net/http"
	"os"
	"time"
)

type appStopperImpl struct {
	*appCtlImpl
}

func newAppStopper(appCtl *appCtlImpl) *appStopperImpl {
	return &appStopperImpl{
		appCtlImpl: appCtl,
	}
}

func (impl *appStopperImpl) stop(app string) error {
	if err := impl.consulDeregister(app); err != nil {
		return err
	}

	if err := impl.stopDaprdApp(app); err != nil {
		return err
	}

	return nil
}

func (impl *appStopperImpl) stopDaprdApp(app string) error {
	pids := impl.getDaprdPids(app)
	for _, pid := range pids {
		if impl.debug {
			fmt.Printf("kill dapr app: %d\n", pid)
		}
		if err := impl.kill(pid); err != nil {
			return err
		}
	}

	return nil
}

func (impl *appStopperImpl) consulDeregister(app string) error {
	svcIds := impl.getConsulRegisteredSvcIds(app)

	client := &http.Client{}
	for _, svcId := range svcIds {
		if impl.debug {
			fmt.Printf("deregister service: %s\n", svcId)
		}
		err := impl.deregister(client, svcId)
		if err != nil {
			return err
		}
	}

	remains := impl.getConsulRegisteredSvcIds(app)
	for {
		if len(remains) == 0 {
			return nil
		}
		remains = impl.getConsulRegisteredSvcIds(app)
		time.Sleep(1 * time.Second)
	}
}

func (impl *appStopperImpl) getConsulRegisteredSvcIds(app string) []string {
	svcIds, _ := script.Get(fmt.Sprintf("http://127.0.0.1:8500/v1/agent/health/service/name/%s", impl.getAppId(app))).
		JQ(".[].Checks[].ServiceID").
		Replace("\"", "").Slice()
	return svcIds
}

func (impl *appStopperImpl) getDaprdPids(app string) []int {
	daprdPids, _ := script.Exec("dapr list -o json").
		JQ(fmt.Sprintf(".[] | select(.appId==\"%s\") | .daprdPid", impl.getAppId(app))).Slice()

	return pie.Map(daprdPids, func(v string) int {
		return cast.ToInt(v)
	})
}

func (impl *appStopperImpl) deregister(client *http.Client, svcId string) error {
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("http://127.0.0.1:8500/v1/agent/service/deregister/%s", svcId), nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (impl *appStopperImpl) kill(pid int) error {
	if pid == 0 {
		return errors.New("invalid pid")
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return errors.Wrapf(err, "找不到进程, pid: %d", pid)
	}

	err = process.Kill()
	if err != nil {
		return errors.Wrapf(err, "无法终止进程, pid: %d", pid)
	}

	return nil
}
