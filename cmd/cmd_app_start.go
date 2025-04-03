package cmd

import (
	"github.com/hdget/hd/pkg/appctl"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

var (
	cmdAppStart = &cobra.Command{
		Use:   "start",
		Short: "start app",
		Run: func(cmd *cobra.Command, args []string) {
			startApp(args)
		},
	}
)

func startApp(args []string) {
	if len(args) != 1 {
		utils.Fatal("Usage: start <app>")
	}

	baseDir, err := os.Getwd()
	if err != nil {
		utils.Fatal("get current dir", err)
	}

	err = appctl.New(baseDir, appctl.WithDebug(argDebug)).Start(args[0])
	if err != nil {
		utils.Fatal("start app", err)
	}
}
