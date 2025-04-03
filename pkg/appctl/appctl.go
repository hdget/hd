package appctl

import (
	"fmt"
	"github.com/hdget/hd/pkg/utils"
	"os"
	"path/filepath"
	"runtime"
)

type AppController interface {
	Start(app string) error
	Stop(app string) error
	Run() error
	Build(refName string, apps ...string) error
	Deploy() error
}

type appCtlImpl struct {
	baseDir string
	binDir  string
	debug   bool
}

func init() {
	// 初始化导出环境变量
	for k, v := range getExportedEnvs() {
		if err := os.Setenv(k, v); err != nil {
			utils.Fatal("export HD environment variable", err)
		}
	}
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

func (a *appCtlImpl) Start(app string) error {
	instance, err := newAppStarter(a)
	if err != nil {
		return err
	}
	return instance.start(app)
}

func (a *appCtlImpl) Build(refName string, apps ...string) error {
	instance, err := newAppBuilder(a)
	if err != nil {
		return err
	}

	return instance.build(refName, apps...)
}

func (a *appCtlImpl) Stop(app string) error {
	instance, err := newAppStopper(a)
	if err != nil {
		return err
	}

	return instance.stop(app)
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
