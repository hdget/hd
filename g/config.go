package g

type HdConfig struct {
	Project ProjectConfig      `toml:"project"`
	Apps    []AppConfig        `toml:"apps"` // 应用启动顺序
	Repos   []RepositoryConfig `toml:"repos"`
	Tools   []ToolConfig       `toml:"tools"`
}

type ProjectConfig struct {
	Name string   `toml:"name"`
	Env  string   `toml:"env"`
	Apps []string `toml:"apps"`
}

type AppConfig struct {
	Name         string          `toml:"name"`
	AppPort      int             `toml:"app_port"`
	AppExposed   bool            `toml:"app_exposed"`
	ExternalPort int             `toml:"external_port"` // 外部端口
	Repo         string          `toml:"repo"`
	Protocol     string          `toml:"protocol"`
	ConfigRepo   string          `toml:"config_repo"`
	ProtoRepo    string          `toml:"proto_repo"`
	Build        *BuildConfig    `toml:"build"`
	Plugins      []*PluginConfig `toml:"plugins"`
	Dapr         DaprConfig      `toml:"dapr"`
}

type BuildConfig struct {
	PbDir        string `toml:"pb_dir"`        // protobuf编译后保存的的目录
	PbPackage    string `toml:"pb_package"`    // protobuf编译后生成的包名
	UseGRPC      bool   `toml:"use_grpc"`      // 是否使用了GRPC, 需要编译GRPC代码
	UseProtobuf  bool   `toml:"use_protobuf"`  // 是否使用了protobuf， 需要编译protobuf文件
	UseSQLBoiler bool   `toml:"use_sqlboiler"` // 是否使用了sqlboiler， 需要自动生成sqlboiler代码
}

type RepositoryConfig struct {
	Name string `toml:"name"`
	Url  string `toml:"url"`
}

type PluginConfig struct {
	Name string `toml:"name"`
	Repo string `toml:"repo"`
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
