package g

type RootConfig struct {
	Project ProjectConfig `toml:"project"`
	Repos   []RepoConfig  `toml:"repos"`
}

type ProjectConfig struct {
	Name        string `toml:"name"`
	Env         string `toml:"env"`
	GatewayPort int    `toml:"gateway_port"`
}

type RepoConfig struct {
	Name string `toml:"name"`
	Url  string `toml:"url"`
}

var (
	Config = &RootConfig{}
)
