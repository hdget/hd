package appctl

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/tools"
	"strings"
	"time"
)

type appStartImpl struct {
	*appCtlImpl
	env string
}

const (
	// 52000-59999
	defaultGatewayPort = 52000
	cmdNormalAppStart  = "%s run --app-address 127.0.0.1:%d --env %s"
	cmdGatewayAppStart = "%s run --app-address 127.0.0.1:%d --web-address :%d --env %s"
	cmdDaprStart       = "dapr run --app-id %s %s -- %s"
)

var (
	daprPortRange   = []int{53000, 59999}
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
		"--enable-app-health-check",
		"--scheduler-host-address",
		"--app-health-probe-interval 60",
	}
)

func newAppStarter(appCtl *appCtlImpl) (*appStartImpl, error) {
	// 检查依赖的工具是否安装
	if err := tools.Check(appCtl.debug,
		tools.Dapr(),
	); err != nil {
		return nil, err
	}

	env, err := getEnv()
	if err != nil {
		return nil, err
	}

	return &appStartImpl{
		appCtlImpl: appCtl,
		env:        env,
	}, nil
}

func (a *appStartImpl) Start(app string) error {
	cmd, err := a.getStartCommand(app)
	if err != nil {
		return err
	}

	err = a.runDetached(a.getAppId(app), cmd, a.getHealthChecker(app), 1*time.Second)
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

	var subCmd string
	switch app {
	case "gateway":
		subCmd = fmt.Sprintf(cmdGatewayAppStart, a.getExecutable(app), ports[0], a.getGatewayPort(), a.env)
	default:
		subCmd = fmt.Sprintf(cmdNormalAppStart, a.getExecutable(app), ports[0], a.env)
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
		return defaultGatewayPort
	}
	return g.Config.Project.GatewayPort
}

func (a *appStartImpl) getHealthChecker(app string) func() bool {
	return func() bool {
		pid, _ := script.Echo("dapr list").Match(a.getAppId(app)).Column(11).String()
		return pid != ""
	}
}
