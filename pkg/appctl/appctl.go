package appctl

import (
	"fmt"
	"github.com/hdget/hd/pkg/env"
	"github.com/hdget/hd/pkg/tools"
	"path/filepath"
	"runtime"
)

type AppController interface {
	Start(app string, extraParam string) error
	Stop(app string) error
	Build(app string, ref string) error
	Install(app string, ref string) error
	Run() error
}

type appCtlImpl struct {
	baseDir   string
	binDir    string
	absBinDir string
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

func (a *appCtlImpl) Start(app string, extraParam string) error {
	fmt.Println()
	fmt.Printf("=== START app: %s ===\n", app)
	fmt.Println()

	// 检查依赖的工具是否安装
	if err := tools.Check(
		tools.Consul(),
		tools.Dapr(),
	); err != nil {
		return err
	}

	instance, err := newAppStarter(a)
	if err != nil {
		return err
	}

	return instance.start(app, extraParam)
}

func (a *appCtlImpl) Install(app string, ref string) error {
	fmt.Println()
	fmt.Printf("=== INSTALL app: %s ===\n", app)
	fmt.Println()

	return newAppInstaller(a).install(app, ref)
}

func (a *appCtlImpl) Build(app string, ref string) error {
	fmt.Println()
	fmt.Printf("=== BUILD app: %s, ref: %s ===\n", app, ref)
	fmt.Println()

	// 检查依赖的工具是否安装
	if err := tools.Check(
		tools.Protoc(),
		tools.ProtocGogoFaster(),
		tools.Sqlboiler(),
	); err != nil {
		return err
	}

	return newAppBuilder(a).build(app, ref)
}

func (a *appCtlImpl) Stop(app string) error {
	fmt.Println()
	fmt.Printf("=== STOP app: %s ===\n", app)
	fmt.Println()

	// 检查依赖的工具是否安装
	if err := tools.Check(
		tools.Consul(),
	); err != nil {
		return err
	}

	return newAppStopper(a).stop(app)
}

func (a *appCtlImpl) Run() error {
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
