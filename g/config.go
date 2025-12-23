package g

type HdConfig struct {
	Project ProjectConfig      `toml:"project"`
	Apps    []AppConfig        `toml:"apps"` // 应用启动顺序
	Repos   []RepositoryConfig `toml:"repos"`
	Tools   []ToolConfig       `toml:"tools"`
	Dapr    DaprConfig         `toml:"dapr"`
}

type ProjectConfig struct {
	Name string `toml:"name"`
	Env  string `toml:"env"`
}

type AppConfig struct {
	Name         string         `toml:"name"`
	ExternalPort int            `toml:"external_port"`
	Repo         string         `toml:"repo"`
	Protocol     string         `toml:"protocol"`
	Plugins      []PluginConfig `toml:"plugins"`
}

type RepositoryConfig struct {
	Name string `toml:"name"`
	Url  string `toml:"url"`
}

type PluginConfig struct {
	Name string `toml:"name"`
	Url  string `toml:"url"`
}

type ToolConfig struct {
	Name            string `toml:"name"`
	Version         string `toml:"version"`
	UrlWinRelease   string `toml:"url_win_release"`
	UrlLinuxRelease string `toml:"url_linux_release"`
}

type DaprConfig struct {
	PortStart              int    `toml:"port_start"`
	PortEnd                int    `toml:"port_end"`
	AppProtocol            string `toml:"app_protocol"`
	ConfigPath             string `toml:"config_path"`
	ResourcePath           string `toml:"resource_path"`
	SchedulerHostAddress   string `toml:"scheduler_host_address"`
	PlacementHostAddress   string `toml:"placement_host_address"`
	DisableAppHealthCheck  bool   `toml:"disable_app_health_check"`
	AppHealthProbeInterval int    `toml:"app_health_probe_interval"`
}
