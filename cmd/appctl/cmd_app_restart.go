package appctl

import (
	"fmt"
	"github.com/hdget/hd/pkg/appctl"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	subCmdRestartApp = &cobra.Command{
		Use:   "restart [app1,app2...]",
		Short: "restart app",
		Run: func(cmd *cobra.Command, args []string) {
			if argAll {
				fmt.Println("restart all app not supported")
				os.Exit(0)
			}

			restartApp(args)
		},
	}
)

func restartApp(args []string) {
	if len(args) < 1 {
		utils.Fatal("Usage: start [app1,app2...]")
	}

	appList := args[0]

	apps := strings.Split(appList, ",")
	if len(apps) == 0 {
		utils.Fatal("you need specify at least one app")
	}

	baseDir, err := os.Getwd()
	if err != nil {
		utils.Fatal("get current dir", err)
	}

	var extraParam string
	if len(os.Args) > 4 {
		extraParam = strings.Join(os.Args[4:], " ")
	}

	for _, app := range apps {
		err = appctl.New(baseDir).Stop(app)
		if err != nil {
			utils.Fatal("stop app", err)
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
