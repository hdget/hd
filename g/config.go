package g

type HDConfig struct {
	Apps       []AppConfig `toml:"apps"`
	ConfigRepo string      `toml:"config_repo"`
	ProtoRepo  string      `toml:"proto_repo"`
}

type AppConfig struct {
	Name string `toml:"name"`
	Repo string `toml:"repo"`
}

var (
	Config = &HDConfig{}
)
