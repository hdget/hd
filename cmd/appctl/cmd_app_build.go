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
		pbOutputPackage string
		pbGenGRPC       bool
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
	subCmdBuildApp.PersistentFlags().StringVarP(&argAppBuild.pbOutputPackage, "package", "", "pb", "--package <package>")
	subCmdBuildApp.PersistentFlags().StringVarP(&argAppBuild.pbOutputDir, "output-dir", "", "autogen", "relative output dir, --output-dir <sub_dir>")
	subCmdBuildApp.PersistentFlags().BoolVarP(&argAppBuild.pbGenGRPC, "grpc", "", false, "--grpc")
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
		err = appctl.New(baseDir,
			appctl.WithPbOutputDir(argAppBuild.pbOutputDir),
			appctl.WithPbOutputPackage(argAppBuild.pbOutputPackage),
			appctl.WithPbGRPC(argAppBuild.pbGenGRPC),
		).Build(app, ref)
		if err != nil {
			utils.Fatal("stop app", err)
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
			appctl.WithPbOutputDir(argAppBuild.pbOutputDir),
			appctl.WithPbOutputPackage(argAppBuild.pbOutputPackage),
			appctl.WithPbGRPC(argAppBuild.pbGenGRPC),
		).Build(app, ref)
		if err != nil {
			utils.Fatal("build app", err)
		}
	}

}
