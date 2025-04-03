package cmd

import (
	"github.com/hdget/hd/pkg/appctl"
	"github.com/spf13/cobra"
	"log"
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
		log.Fatal("Usage: build <app1,app2...> <branch>")
	}

	appList, refName := args[0], args[1]

	apps := strings.Split(appList, ",")
	if len(apps) == 0 {
		log.Fatal("you need specify at least one app")
	}

	if refName == "" {
		log.Fatal("you need specify branch")
	}

	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = appctl.New(baseDir, appctl.WithDebug(argDebug)).Build(refName, apps...)
	if err != nil {
		log.Fatal(err)
	}
}
