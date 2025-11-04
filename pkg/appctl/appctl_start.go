package appctl

import (
	"fmt"

	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/env"
	"github.com/hdget/hd/pkg/tools"
	"github.com/pkg/errors"

	"time"
)

type appStartImpl struct {
	*appCtlImpl
	env string
}

const (
	defaultTimeout = 5 * time.Second
)

func newAppStarter(appCtl *appCtlImpl) (*appStartImpl, error) {
	e, err := env.GetHdEnv()
	if err != nil {
		return nil, err
	}

	return &appStartImpl{
		appCtlImpl: appCtl,
		env:        e,
	}, nil
}

func (a *appStartImpl) start(name string, extraParam string) error {
	app, err := newApp(name)
	if err != nil {
		return errors.Wrapf(err, "failed to new app instance '%s'", name)
	}

	// 检查app是否可启动
	if err = app.PreStart(); err != nil {
		return errors.Wrapf(err, "app '%s' can not start", name)
	}

	// 获取启动命令
	cmd, err := app.GetStartCommand(a.binDir, extraParam)
	if err != nil {
		return err
	}

	if g.Debug {
		fmt.Println(cmd)
	}

	err = tools.RunDaemon(app.GetId(), cmd, app.GetHealthChecker(), defaultTimeout)
	if err != nil {
		return err
	}

	return nil
}
