package appctl

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/elliotchance/pie/v2"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/env"
	"github.com/hdget/hd/pkg/tools"
	"github.com/pkg/errors"
)

type AppController interface {
	Start(app string, extraParam ...string) error
	Stop(app string) error
	Build(app string, ref string) error
	Install(app string, ref string) error
	Run() error
}

type appCtlImpl struct {
	baseDir string
	binDir  string
	// pb options
	pbOutputDir     string
	pbOutputPackage string
	pbGenGRPC       bool
	// plugin options
	pluginDir string
	plugins   []string
}

const (
	namespaceAppSeparator = "-"
)

func New(baseDir string, options ...Option) AppController {
	impl := &appCtlImpl{
		baseDir: baseDir,
	}

	for _, apply := range options {
		apply(impl)
	}

	return impl
}

func (a *appCtlImpl) Start(app string, extraParam ...string) error {
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

	var startParam string
	if len(extraParam) > 0 {
		startParam = extraParam[0]
	}

	return instance.start(app, startParam)
}

func (a *appCtlImpl) Install(app string, ref string) error {
	fmt.Println()
	fmt.Printf("=== INSTALL app: %s ===\n", app)
	fmt.Println()

	return newAppInstaller(a).install(app, ref)
}

func (a *appCtlImpl) Build(app string, ref string) error {
	// 检查依赖的工具是否安装
	if err := tools.Check(
		tools.Protoc(),
		tools.ProtocGo(),
		tools.ProtocGoGRPC(),
		tools.Sqlboiler(),
	); err != nil {
		return err
	}

	// 获取app配置
	appRepoConfig, exist := g.RepoConfigs[app]
	if !exist {
		return fmt.Errorf("app repo config not found in hd.toml: %s", app)
	}

	// 如果指定了plugin,则只编译指定的plugin
	if len(a.plugins) > 0 {
		for _, pluginName := range a.plugins {
			fmt.Println()
			fmt.Printf("=== BUILD plugin: %s, ref: %s ===\n", pluginName, ref)
			fmt.Println()

			index := pie.FindFirstUsing(appRepoConfig.Plugins, func(v g.PluginConfig) bool {
				return v.Name == pluginName
			})

			if index == -1 {
				return fmt.Errorf("plugin: %s not found for app: %s in hd.toml", pluginName, app)
			}

			err := newPluginBuilder(a, a.pbOutputDir, a.pbOutputPackage, a.pbGenGRPC).build(appRepoConfig.Plugins[index].Name, appRepoConfig.Plugins[index].Url, ref)
			if err != nil {
				return errors.Wrap(err, "build plugin")
			}
		}
		return nil
	}

	// 如果未指定plugin, 如果该app下有关联的plugin配置，则编译所有plugins
	for _, pluginConfig := range appRepoConfig.Plugins {
		fmt.Println()
		fmt.Printf("=== BUILD plugin: %s, ref: %s ===\n", pluginConfig.Name, ref)
		fmt.Println()
		if err := newPluginBuilder(a, a.pbOutputDir, a.pbOutputPackage, a.pbGenGRPC).build(pluginConfig.Name, pluginConfig.Url, ref); err != nil {
			return errors.Wrap(err, "build plugin")
		}
	}

	// 编译app
	fmt.Println()
	fmt.Printf("=== BUILD app: %s, ref: %s ===\n", app, ref)
	fmt.Println()
	return newAppBuilder(a, a.pbOutputDir, a.pbOutputPackage, a.pbGenGRPC).build(app, ref)
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

func (a *appCtlImpl) getBinOutputDir() string {
	return filepath.Join(a.baseDir, a.binDir)
}

func (a *appCtlImpl) getPluginOutputDir() string {
	return filepath.Join(a.baseDir, a.pluginDir)
}

func (a *appCtlImpl) getExecutable(app string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s.exe", app)
	}
	return app
}

func (a *appCtlImpl) getAppId(app string) string {
	if namespace, exists := env.GetHdNamespace(); exists {
		var sb strings.Builder
		sb.Grow(len(namespace) + len(app) + 1)
		sb.WriteString(namespace)
		sb.WriteString(namespaceAppSeparator)
		sb.WriteString(app)
		return sb.String()
	}
	return app
}
