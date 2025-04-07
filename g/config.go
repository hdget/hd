package g

type RootConfig struct {
	Project ProjectConfig `toml:"project"`
	Repos   []RepoConfig  `toml:"repos"`
	Tools   []ToolConfig  `toml:"tools"`
}

type ProjectConfig struct {
	Name        string   `toml:"name"`
	Env         string   `toml:"env"`
	GatewayPort int      `toml:"gateway_port"`
	Apps        []string `toml:"apps"` // 应用启动顺序
}

type RepoConfig struct {
	Name string `toml:"name"`
	Url  string `toml:"url"`
}

type ToolConfig struct {
	Name            string `toml:"name"`
	Version         string `toml:"version"`
	UrlWinRelease   string `toml:"url_win_release"`
	UrlLinuxRelease string `toml:"url_linux_release"`
}

var (
	Config      = &RootConfig{}
	ToolConfigs = map[string]ToolConfig{}
	RepoConfigs = map[string]RepoConfig{}
)
