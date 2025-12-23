package appctl

import (
	"os"
	"strings"

	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/appctl"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	subCmdStartApp = &cobra.Command{
		Use:   "start [app1,app2...]",
		Short: "start app",
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		Run: func(cmd *cobra.Command, args []string) {
			if argAll {
				startAllApp()
			} else {
				startApp(args)
			}
		},
	}
)

func startAllApp() {
	baseDir, err := os.Getwd()
	if err != nil {
		utils.Fatal("get current dir", err)
	}

	for _, app := range g.Config.Apps {
		err = appctl.New(
			baseDir,
			appctl.WithBinOutputDir(argBinOutputDir),
		).Start(app.Name)
		if err != nil {
			utils.Fatal("start app", err)
		}
	}
}

func startApp(args []string) {
	if len(args) < 1 {
		utils.Fatal("Usage: start [app1,app2...]")
	}

	appList := args[0]
	apps := strings.Split(appList, ",")
	if len(apps) == 0 {
		utils.Fatal("you need specify at least one app")
	}

	var extraParam string
	if len(os.Args) > 4 {
		extraParam = strings.Join(os.Args[4:], " ")
	}

	baseDir, err := os.Getwd()
	if err != nil {
		utils.Fatal("get current dir", err)
	}

	for _, app := range apps {
		err = appctl.New(
			baseDir,
			appctl.WithBinOutputDir(argBinOutputDir),
		).Start(app, extraParam)
		if err != nil {
			utils.Fatal("start app", err)
		}
	}
}
