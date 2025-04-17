package appctl

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/env"
	"path/filepath"
	"strings"
	"time"
)

type appStartImpl struct {
	*appCtlImpl
	env string
}

const (
	cmdNormalAppStart  = "%s run --app-address 127.0.0.1:%d --env %s"
	cmdGatewayAppStart = "%s run --app-address 127.0.0.1:%d --web-address :%d --env %s"
	cmdDaprStart       = "dapr run --app-id %s %s -- %s"
	defaultTimeout     = 30 * time.Second
)

var (
	daprPortRange   = []int{55000, 59999}
	daprPortOptions = []string{
		"--app-port %d",
		"--dapr-grpc-port %d",
		"--dapr-http-port %d",
		"--dapr-internal-grpc-port %d",
		"--metrics-port %d",
	}
	daprFixedOptions = []string{
		"--app-protocol grpc",
		"--config config/dapr/config.yaml",
		"--resources-path config/dapr/components",
		"--scheduler-host-address ''",
		"--placement-host-address ''",
		"--enable-app-health-check",
		"--app-health-probe-interval 5",
	}
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

func (a *appStartImpl) start(app string) error {
	cmd, err := a.getStartCommand(app)
	if err != nil {
		return err
	}

	if g.Debug {
		fmt.Println(cmd)
	}

	err = a.run(a.getAppId(app), cmd, a.getHealthChecker(app), defaultTimeout)
	if err != nil {
		return err
	}

	return nil
}

func (a *appStartImpl) getStartCommand(app string) (string, error) {
	ports, err := findAvailablePorts(len(daprPortOptions), daprPortRange[0], daprPortRange[1])
	if err != nil {
		return "", err
	}

	appBinPath := filepath.ToSlash(filepath.Join(a.binDir, app))

	var subCmd string
	switch app {
	case "gateway":
		if !isPortAvailable(a.getGatewayPort()) {
			return "", fmt.Errorf("gateway port %d is not available", a.getGatewayPort())
		}
		subCmd = fmt.Sprintf(cmdGatewayAppStart, appBinPath, ports[0], a.getGatewayPort(), a.env)
	default:
		subCmd = fmt.Sprintf(cmdNormalAppStart, appBinPath, ports[0], a.env)
	}

	var daprOptions []string
	for i, option := range daprPortOptions {
		daprOptions = append(daprOptions, fmt.Sprintf(option, ports[i]))
	}
	daprOptions = append(daprOptions, daprFixedOptions...)

	return fmt.Sprintf(cmdDaprStart, a.getAppId(app), strings.Join(daprOptions, " "), subCmd), nil
}

func (a *appStartImpl) getGatewayPort() int {
	if g.Config.Project.GatewayPort == 0 {
		return g.DefaultGatewayPort
	}
	return g.Config.Project.GatewayPort
}

func (a *appStartImpl) getHealthChecker(app string) func() bool {
	return func() bool {
		pid, _ := script.Exec("dapr list").Match(a.getAppId(app)).Column(11).String()
		return pid != ""
	}
}
