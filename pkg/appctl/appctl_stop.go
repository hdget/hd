package appctl

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/hdget/hd/pkg/tools"
	"strings"
)

type appStopperImpl struct {
	*appCtlImpl
}

func newAppStopper(appCtl *appCtlImpl) (*appStopperImpl, error) {
	// 检查依赖的工具是否安装
	if err := tools.Check(appCtl.debug,
		tools.Consul(),
	); err != nil {
		return nil, err
	}

	return &appStopperImpl{
		appCtlImpl: appCtl,
	}, nil
}

func (impl *appStopperImpl) stop(app string) error {
	if err := impl.consulDeregister(app); err != nil {
		return err
	}
	return nil
}
func (impl *appStopperImpl) consulDeregister(app string) error {
	matched, err := script.
		Get(fmt.Sprintf("http://127.0.0.1:8500/v1/agent/health/service/name/%s", impl.getAppId(app))).
		JQ(".[0].Checks[].ServiceID").String()
	if err != nil {
		return err
	}

	// 分割结果（如果是多个ID）
	for i, id := range strings.Split(strings.TrimSpace(matched), "\n") {
		fmt.Printf("ServiceID %d: %s\n", i+1, id)
	}
	return nil
}
