package appctl

import (
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/appctl"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	argAppBuild = struct {
		pbOutputDir     string
		pluginOutputDir string
		pbOutputPackage string
		pbGenGRPC       bool
		plugins         []string
	}{}

	subCmdBuildApp = &cobra.Command{
		Use:   "build [app1,app2...] [branch]",
		Short: "build app",
		Run: func(cmd *cobra.Command, args []string) {
			if argAll {
				buildAllApp(args)
			} else {
				buildApp(args)
			}

		},
	}
)

func init() {
	// protobuf编译后的包名
	subCmdBuildApp.PersistentFlags().StringVarP(&argAppBuild.pbOutputPackage, "pb-package", "", "pb", "--pb-package [package_name]]")
	// protobuf编译后的目录
	subCmdBuildApp.PersistentFlags().StringVarP(&argAppBuild.pbOutputDir, "pb-dir", "", "autogen", "relative pb output dir, --pb-dir [dir]")
	// 是否要输出grpc
	subCmdBuildApp.PersistentFlags().BoolVarP(&argAppBuild.pbGenGRPC, "grpc", "", false, "--grpc")
	// 指定编译哪个plugin
	subCmdBuildApp.PersistentFlags().StringSliceVarP(&argAppBuild.plugins, "plugins", "", nil, "--plugins [plugin1,plugin2...]")
	// plugin编译输出到哪个目录
	subCmdBuildApp.PersistentFlags().StringVarP(&argAppBuild.pluginOutputDir, "plugin-dir", "", "plugins", "relative plugin output dir, --plugin-dir [dir]")
}

func buildAllApp(args []string) {
	if len(args) < 1 {
		utils.Fatal("Usage: build [branch] --all")
	}

	ref := args[0]

	baseDir, err := os.Getwd()
	if err != nil {
		utils.Fatal("get current dir", err)
	}

	for _, app := range g.Config.Project.Apps {
		err = appctl.New(
			baseDir,
			appctl.WithBinOutputDir(argBinOutputDir),
			appctl.WithPluginOutputDir(argAppBuild.pluginOutputDir),
			appctl.WithPlugins(argAppBuild.plugins),
			appctl.WithPbOutputDir(argAppBuild.pbOutputDir),
			appctl.WithPbOutputPackage(argAppBuild.pbOutputPackage),
			appctl.WithPbGRPC(argAppBuild.pbGenGRPC),
		).Build(app, ref)
		if err != nil {
			utils.Fatal("build", err)
		}
	}
}

func buildApp(args []string) {
	if len(args) < 2 {
		utils.Fatal("Usage: build [app1,app2...] <branch>")
	}

	appList, ref := args[0], args[1]

	apps := strings.Split(appList, ",")
	if len(apps) == 0 {
		utils.Fatal("you need specify at least one app")
	}

	if ref == "" {
		utils.Fatal("you need specify branch")
	}

	baseDir, err := os.Getwd()
	if err != nil {
		utils.Fatal("get current dir", err)
	}

	for _, app := range apps {
		err = appctl.New(baseDir,
			appctl.WithBinOutputDir(argBinOutputDir),
			appctl.WithPluginOutputDir(argAppBuild.pluginOutputDir),
			appctl.WithPlugins(argAppBuild.plugins),
			appctl.WithPbOutputDir(argAppBuild.pbOutputDir),
			appctl.WithPbOutputPackage(argAppBuild.pbOutputPackage),
			appctl.WithPbGRPC(argAppBuild.pbGenGRPC),
		).Build(app, ref)
		if err != nil {
			utils.Fatal("build", err)
		}
	}

}
