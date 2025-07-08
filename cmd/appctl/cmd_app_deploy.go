package appctl

import (
	"github.com/elliotchance/pie/v2"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/appctl"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	subCmdDeployApp = &cobra.Command{
		Use:   "deploy [app1,app2...] [branch]",
		Short: "deploy app",
		Run: func(cmd *cobra.Command, args []string) {
			if argAll {
				deployAllApp(args)
			} else {
				deployApp(args)
			}
		},
	}
)

func init() {
	// protobuf编译后的包名
	subCmdDeployApp.PersistentFlags().StringVarP(&argAppBuild.pbOutputPackage, "pb-package", "", "pb", "--pb-package [package_name]]")
	// protobuf编译后的目录
	subCmdDeployApp.PersistentFlags().StringVarP(&argAppBuild.pbOutputDir, "pb-dir", "", "autogen", "relative pb output dir, --pb-dir [dir]")
	// 是否要输出grpc
	subCmdDeployApp.PersistentFlags().BoolVarP(&argAppBuild.pbGenGRPC, "grpc", "", false, "--grpc")
	// 指定编译哪个plugin
	subCmdDeployApp.PersistentFlags().StringSliceVarP(&argAppBuild.plugins, "plugins", "", nil, "--plugins [plugin1,plugin2...]")
	// plugin编译输出到哪个目录
	subCmdDeployApp.PersistentFlags().StringVarP(&argAppBuild.pluginOutputDir, "plugin-dir", "", "plugins", "relative plugin output dir, --plugin-dir [dir]")
}

func deployAllApp(args []string) {
	if len(args) < 1 {
		utils.Fatal("Usage: deploy [branch] --all")
	}

	ref := args[0]
	if ref == "" {
		utils.Fatal("Usage: deploy <branch> --all")
	}

	baseDir, err := os.Getwd()
	if err != nil {
		utils.Fatal("get current dir", err)
	}

	for _, app := range pie.Reverse(g.Config.Project.Apps) {
		err = appctl.New(baseDir).Stop(app)
		if err != nil {
			utils.Fatal("stop app", err)
		}
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
			utils.Fatal("build app", err)
		}

		err = appctl.New(baseDir).Install(app, ref)
		if err != nil {
			utils.Fatal("install app", err)
		}

		err = appctl.New(
			baseDir,
			appctl.WithBinOutputDir(argBinOutputDir),
		).Start(app)
		if err != nil {
			utils.Fatal("start app", err)
		}
	}
}

func deployApp(args []string) {
	if len(args) < 2 {
		utils.Fatal("Usage: deploy [app1,app2...] [branch]")
	}

	appList, ref := args[0], args[1]
	if ref == "" {
		utils.Fatal("you need specify branch")
	}

	apps := strings.Split(appList, ",")
	if len(apps) == 0 {
		utils.Fatal("you need specify at least one app")
	}

	baseDir, err := os.Getwd()
	if err != nil {
		utils.Fatal("get current dir", err)
	}

	var extraParam string
	if len(os.Args) > 5 {
		extraParam = strings.Join(os.Args[5:], " ")
	}

	for _, app := range apps {
		err = appctl.New(baseDir).Stop(app)
		if err != nil {
			utils.Fatal("stop app", err)
		}

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
			utils.Fatal("build app", err)
		}

		err = appctl.New(baseDir).Install(app, ref)
		if err != nil {
			utils.Fatal("install app", err)
		}

		err = appctl.New(
			baseDir,
			appctl.WithBinOutputDir(argBinOutputDir),
		).Start(app, extraParam)
		if err != nil {
			utils.Fatal("start app", err)
		}
	}
}
