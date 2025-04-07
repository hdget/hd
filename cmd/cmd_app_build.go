package cmd

import (
	"github.com/hdget/hd/pkg/appctl"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	cmdAppBuild = &cobra.Command{
		Use:   "build",
		Short: "build app",
		Run: func(cmd *cobra.Command, args []string) {
			buildApp(args)
		},
	}
)

func buildApp(args []string) {
	if len(args) != 2 {
		utils.Fatal("Usage: build <app1,app2...> <branch>")
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
		err = appctl.New(baseDir, appctl.WithDebug(argDebug)).Build(app, ref)
		if err != nil {
			utils.Fatal("build app", err)
		}
	}

}
