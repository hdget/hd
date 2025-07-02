package appctl

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/elliotchance/pie/v2"
	"github.com/hdget/hd/g"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"net/http"
	"runtime"
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

	switch platform := runtime.GOOS; platform {
	case "windows":
		output, err := script.Exec(fmt.Sprintf("dapr stop --app-id %s", impl.getAppId(app))).String()
		if err != nil {
			return errors.Wrapf(err, "%s stop failed, err: %s", app, output)
		}
	case "linux", "darwin":
		pids := impl.getAppPids(app)
		for _, pid := range pids {
			if g.Debug {
				fmt.Printf("send terminal signal to: %d\n", pid)
			}

			if err := sendStopSignal(pid); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("stop on: %s not supported", platform)
	}

	return nil
}

func (impl *appStopperImpl) consulDeregister(app string) error {
	svcIds := impl.getConsulRegisteredSvcIds(app)

	client := &http.Client{}
	for _, svcId := range svcIds {
		if g.Debug {
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

func (impl *appStopperImpl) getAppPids(app string) []int {
	appPids, _ := script.Exec("dapr list -o json").
		JQ(fmt.Sprintf(".[] | select(.appId==\"%s\") | .appPid", impl.getAppId(app))).Slice()

	return pie.Map(appPids, func(v string) int {
		return cast.ToInt(v)
	})
}

func (impl *appStopperImpl) deregister(client *http.Client, svcId string) error {
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("http://127.0.0.1:8500/v1/agent/service/deregister/%s", svcId), nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return nil
}
