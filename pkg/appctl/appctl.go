package appctl

import (
	"fmt"
	"github.com/hdget/hd/pkg/env"
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
	baseDir   string
	binDir    string
	absBinDir string
	debug     bool
}

func New(baseDir string, options ...Option) AppController {
	impl := &appCtlImpl{
		baseDir:   baseDir,
		binDir:    "bin",
		absBinDir: filepath.Join(baseDir, "bin"),
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
	if nm, exists := env.GetHdNamespace(); exists {
		return fmt.Sprintf("%s_%s", nm, app)
	}
	return app
}
