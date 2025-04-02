package app

import (
	"github.com/BurntSushi/toml"
	"github.com/hdget/hd/g"
	"log"
	"path/filepath"
)

type Controller interface {
	Start() error
	Stop() error
	Run() error
	Build(refName string, apps ...string) error
	Deploy() error
}

type appControlImpl struct {
	baseDir         string
	binDir          string
	debug           bool
	pbOutputDir     string
	pbOutputPackage string
}

func New(baseDir string, options ...Option) Controller {
	impl := &appControlImpl{
		baseDir:         baseDir,
		binDir:          filepath.Join(baseDir, "bin"),
		pbOutputDir:     "autogen",
		pbOutputPackage: "pb",
	}

	for _, apply := range options {
		apply(impl)
	}

	return impl
}

var (
	appName2appConfig = map[string]g.AppConfig{}
)

func init() {
	// 读取 TOML 文件
	if _, err := toml.DecodeFile("hd.toml", &g.Config); err != nil {
		log.Fatal(err)
	}

	for _, app := range g.Config.Apps {
		appName2appConfig[app.Name] = app
	}
}

func (a appControlImpl) Start() error {
	return nil
}

func (a appControlImpl) Stop() error {
	return nil
}

func (a appControlImpl) Run() error {
	return nil
}

func (a appControlImpl) Deploy() error {
	return nil
}
