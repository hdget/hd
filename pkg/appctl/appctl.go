package appctl

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/utils"
	"os"
	"path/filepath"
	"runtime"
)

type AppController interface {
	Start(app string) error
	Stop() error
	Run() error
	Build(refName string, apps ...string) error
	Deploy() error
}

type appCtlImpl struct {
	baseDir string
	binDir  string
	debug   bool
}

var (
	repoName2repoConfig = map[string]g.RepoConfig{}
	configFile          = "hd.toml"
)

func init() {
	// 读取 TOML 文件
	if _, err := toml.DecodeFile(configFile, &g.Config); err != nil {
		utils.Fatal(fmt.Sprintf("read config file, file: %s", configFile), err)
	}

	for _, repo := range g.Config.Repos {
		repoName2repoConfig[repo.Name] = repo
	}

	// 初始化导出环境变量
	for k, v := range getExportedEnvs() {
		if err := os.Setenv(k, v); err != nil {
			utils.Fatal("export HD environment variable", err)
		}
	}
}

func (a *appCtlImpl) Start(app string) error {
	instance, err := newAppStarter(a)
	if err != nil {
		return err
	}
	return instance.Start(app)
}

func (a *appCtlImpl) Build(refName string, apps ...string) error {
	instance, err := newAppBuilder(a)
	if err != nil {
		return err
	}

	return instance.Build(refName, apps...)
}

func New(baseDir string, options ...Option) AppController {
	impl := &appCtlImpl{
		baseDir: baseDir,
		binDir:  filepath.Join(baseDir, "bin"),
	}

	for _, apply := range options {
		apply(impl)
	}

	return impl
}

func (a *appCtlImpl) Stop() error {
	return nil
}

func (a *appCtlImpl) Run() error {
	return nil
}

func (a *appCtlImpl) Deploy() error {
	return nil
}

func (a *appCtlImpl) getExecutable(app string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s.exe", app)
	}
	return app
}

func (a *appCtlImpl) getAppId(app string) string {
	if nm, exists := os.LookupEnv(envHdNamespace); exists {
		return fmt.Sprintf("%s_%s", nm, app)
	}
	return app
}
