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
	subCmdStopApp = &cobra.Command{
		Use:   "stop [app1,app2...]",
		Short: "stop app",
		Run: func(cmd *cobra.Command, args []string) {
			if argAll {
				stopAllApp()
			} else {
				stopApp(args)
			}
		},
	}
)

func stopAllApp() {
	baseDir, err := os.Getwd()
	if err != nil {
		utils.Fatal("get current dir", err)
	}

	for _, app := range g.Config.Project.Apps {
		err = appctl.New(baseDir).Stop(app)
		if err != nil {
			utils.Fatal("stop app", err)
		}
	}
}

func stopApp(args []string) {
	if len(args) < 1 {
		utils.Fatal("Usage: stop [app1,app2...]")
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

	for _, app := range apps {
		err = appctl.New(baseDir).Stop(app)
		if err != nil {
			utils.Fatal("stop app", err)
		}
	}
}
