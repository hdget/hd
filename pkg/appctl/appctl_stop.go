package appctl

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/hdget/hd/g"
	"github.com/pkg/errors"
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
	case "windows": // windows下是强制终止
		output, err := script.Exec(fmt.Sprintf("dapr stop --app-id %s", impl.getAppId(app))).String()
		if err != nil {
			return errors.Wrapf(err, "%s stop failed, err: %s", app, output)
		}
	case "linux", "darwin":
		daprdPids, appPids, err := impl.getDaprRelatedPids(app)
		if err != nil {
			return err
		}

		for i := 0; i < len(daprdPids); i++ {
			if g.Debug {
				fmt.Printf("send stop signal to, daprdPid: %s, appPid: %s\n", daprdPids[i], appPids[i])
			}

			if err := sendStopSignal(daprdPids[i], appPids[i]); err != nil {
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

func (impl *appStopperImpl) getDaprRelatedPids(app string) ([]string, []string, error) {
	daprdPids, _ := script.Exec("dapr list -o json").
		JQ(fmt.Sprintf(".[] | select(.appId==\"%s\") | .daprdPid", impl.getAppId(app))).Slice()

	appPids, _ := script.Exec("dapr list -o json").
		JQ(fmt.Sprintf(".[] | select(.appId==\"%s\") | .appPid", impl.getAppId(app))).Slice()

	if len(daprdPids) != len(appPids) {
		return nil, nil, errors.New("pids not match")
	}

	return daprdPids, appPids, nil
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
