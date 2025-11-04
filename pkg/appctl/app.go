package appctl

import (
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitfield/script"
	"github.com/elliotchance/pie/v2"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/env"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type Apper interface {
	GetRepoUrl() (string, error)
	GetPath(binDir string) string
	GetStartCommand(binDir, extraParam string) (string, error)
	PreStart() error
	GetId() string
	GetHealthChecker() func() bool
}

type appImpl struct {
	name   string
	config g.AppConfig
}

type port struct {
	appPort      int
	externalPort int
	randomPorts  []int
}

const (
	defaultDaprPortStart = 55000
	defaultDaprPortEnd   = 59999
	//cmdExternalAppStart  = "%s run --app-address 127.0.0.1:%d --web-address :%d"
	//cmdInternalAppStart  = "%s run --app-address 127.0.0.1:%d"
	cmdRunDapr = "dapr run"
	cmdRunApp  = "%s run"
)

var (
	daprPorts = []string{
		"--app-port",
		"--dapr-grpc-port",
		"--dapr-http-port",
		"--dapr-internal-grpc-port",
		"--metrics-port",
	}
	daprArguments = map[string]any{
		"--app-id":                    "",
		"--app-protocol":              "grpc",
		"--config":                    "config/dapr/config.yaml",
		"--resources-path":            "config/dapr/components",
		"--scheduler-host-address":    "''",
		"--placement-host-address":    "''",
		"--enable-app-health-check":   "",
		"--app-health-probe-interval": 5,
	}
)

func newApp(name string) (Apper, error) {
	index := pie.FindFirstUsing(g.Config.Apps, func(v g.AppConfig) bool {
		return strings.EqualFold(v.Name, name)
	})
	if index == -1 {
		return nil, fmt.Errorf("app '%s' configt not found in hd.toml", name)
	}

	return &appImpl{
		name:   name,
		config: g.Config.Apps[index],
	}, nil
}

func (impl *appImpl) GetStartCommand(binDir, extraParam string) (string, error) {
	allocatedPort, err := impl.allocatePort()
	if err != nil {
		return "", errors.Wrap(err, "get allocatedPort")
	}

	commands := []string{
		impl.getDaprArgument(allocatedPort),
		"--",
		impl.getAppRunCommand(binDir, extraParam, allocatedPort),
	}

	return strings.Join(commands, " "), nil
}

func (impl *appImpl) PreStart() error {
	// 如果配置了外部公开的端口，在启动前需要判断端口是否可用
	if impl.config.ExternalPort > 0 {
		if !impl.isPortAvailable(impl.config.ExternalPort) {
			return fmt.Errorf("external port '%d' not available", impl.config.ExternalPort)
		}
	}
	return nil
}

func (impl *appImpl) GetPath(binDir string) string {
	return filepath.ToSlash(filepath.Join(binDir, impl.name))
}

func (impl *appImpl) GetRepoUrl() (string, error) {
	if impl.config.Repo == "" {
		return "", errors.New("empty repository")
	}
	return impl.config.Repo, nil
}

func (impl *appImpl) GetId() string {
	if namespace, exists := env.GetHdNamespace(); exists {
		var sb strings.Builder
		sb.Grow(len(namespace) + len(impl.name) + 1)
		sb.WriteString(namespace)
		sb.WriteString(namespaceAppSeparator)
		sb.WriteString(impl.name)
		return sb.String()
	}
	return impl.name
}

func (impl *appImpl) GetHealthChecker() func() bool {
	return func() bool {
		pid, _ := script.Exec("dapr list").Match(impl.GetId()).Column(11).String()
		return pid != ""
	}
}

func (impl *appImpl) allocatePort() (*port, error) {
	daprPortStart := defaultDaprPortStart
	if g.Config.Dapr.PortStart != 0 {
		daprPortStart = g.Config.Dapr.PortStart
	}

	daprPortEnd := defaultDaprPortEnd
	if g.Config.Dapr.PortEnd != 0 {
		daprPortEnd = g.Config.Dapr.PortEnd
	}

	var externalPort int
	if impl.config.ExternalPort > 0 {
		externalPort = impl.config.ExternalPort
	}

	ports, err := impl.findAvailablePorts(len(daprPorts), daprPortStart, daprPortEnd)
	if err != nil {
		return nil, errors.Wrap(err, "find system available ports")
	}

	return &port{
		appPort:      ports[0],
		externalPort: externalPort, // 对外提供访问的端口，一般是HTTP
		randomPorts:  ports,
	}, nil
}

// FindAvailablePorts 查找指定数量的可用端口
func (impl *appImpl) findAvailablePorts(count, startPort, endPort int) ([]int, error) {
	var availablePorts []int

	for p := startPort; p <= endPort; p++ {
		if impl.isPortAvailable(p) {
			availablePorts = append(availablePorts, p)
			if len(availablePorts) >= count {
				return availablePorts, nil
			}
		}
	}

	if len(availablePorts) < count {
		return nil, fmt.Errorf("端口范围%d-%d内只找到%d个可用端口,需要 %d 个",
			startPort, endPort, len(availablePorts), count)
	}

	return availablePorts, nil
}

// isPortAvailable 双重验证端口是否可用
func (impl *appImpl) isPortAvailable(port int) bool {
	// 双重检查机制
	for i := 0; i < 2; i++ { // 检查两次减少竞态条件风险
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return false
		}
		_ = ln.Close()

		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 500*time.Millisecond)
		if err != nil {
			return true // 连接失败说明端口可能可用
		}
		if conn != nil {
			_ = conn.Close()
			return false
		}

		time.Sleep(10 * time.Millisecond) // 短暂等待
	}
	return false
}

func (impl *appImpl) getAppRunCommand(binDir, extraParam string, port *port) string {
	commands := []string{
		fmt.Sprintf(cmdRunApp, impl.GetPath(binDir)),
	}

	argMap := make(map[string]string)
	argMap["--app-address"] = fmt.Sprintf("127.0.0.1:%d", port.appPort)
	if impl.config.ExternalPort > 0 {
		argMap["--web-address"] = fmt.Sprintf(":%d", port.externalPort)
	}

	for key, value := range argMap {
		commands = append(commands, key, value)
	}

	if extraParam != "" {
		commands = append(commands, extraParam)
	}

	return strings.Join(commands, " ")
}

func (impl *appImpl) getDaprArgument(port *port) string {
	commands := []string{
		cmdRunDapr,
	}

	// IMPORTANT: app-port必须是第一个
	for i, p := range port.randomPorts {
		commands = append(commands, daprPorts[i], cast.ToString(p))
	}

	daprArguments["--app-id"] = impl.GetId()

	if g.Config.Dapr.AppProtocol != "" {
		daprArguments["--app-protocol"] = g.Config.Dapr.AppProtocol
	}

	if g.Config.Dapr.ConfigPath != "" {
		daprArguments["--config_path"] = g.Config.Dapr.ConfigPath
	}

	if g.Config.Dapr.ResourcePath != "" {
		daprArguments["--resource_path"] = g.Config.Dapr.ResourcePath
	}

	if g.Config.Dapr.SchedulerHostAddress != "" {
		daprArguments["--scheduler_host_address"] = g.Config.Dapr.SchedulerHostAddress
	}

	if g.Config.Dapr.PlacementHostAddress != "" {
		daprArguments["--placement_host_address"] = g.Config.Dapr.PlacementHostAddress
	}

	if g.Config.Dapr.AppHealthProbeInterval > 0 {
		daprArguments["--app_health_probe_interval"] = g.Config.Dapr.AppHealthProbeInterval
	}

	if g.Config.Dapr.DisableAppHealthCheck {
		delete(daprArguments, "--enable-app-health-check")
		delete(daprArguments, "--app_health_probe_interval")
	} else if g.Config.Dapr.AppHealthProbeInterval > 0 {
		daprArguments["--app_health_probe_interval"] = g.Config.Dapr.AppHealthProbeInterval
	}

	for k, v := range daprArguments {
		commands = append(commands, k, cast.ToString(v))
	}

	return strings.Join(commands, " ")
}
