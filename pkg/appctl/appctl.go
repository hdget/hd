package appctl

import (
	"fmt"
	"github.com/hdget/hd/pkg/env"
	"path/filepath"
	"runtime"
)

type AppController interface {
	Start(apps []string) error
	Stop(apps []string) error
	Build(apps []string, refName string) error
	Run() error
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

func (a *appCtlImpl) Start(apps []string) error {
	instance, err := newAppStarter(a)
	if err != nil {
		return err
	}

	for _, app := range apps {
		err = instance.start(app)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func (a *appCtlImpl) Build(apps []string, ref string) error {
	instance, err := newAppBuilder(a)
	if err != nil {
		return err
	}

	for _, app := range apps {
		err = instance.build(app, ref)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *appCtlImpl) Stop(apps []string) error {
	instance, err := newAppStopper(a)
	if err != nil {
		return err
	}

	for _, app := range apps {
		err = instance.stop(app)
		if err != nil {
			return err
		}
	}

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
	if nm, exists := env.GetHdNamespace(); exists {
		return fmt.Sprintf("%s_%s", nm, app)
	}
	return app
}
