package appctl

import (
	"fmt"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/appctl"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

var (
	subCmdStartApp = &cobra.Command{
		Use:   "start [app1,app2...]",
		Short: "start app",
		Run: func(cmd *cobra.Command, args []string) {
			// 1. 创建临时 FlagSet 解析所有参数
			tempFlags := pflag.NewFlagSet("temp", pflag.ContinueOnError)
			tempFlags.ParseErrorsWhitelist.UnknownFlags = true
			_ = tempFlags.Parse(os.Args[1:])

			// 2. 提取未知 flags
			var unknownArgs []string
			tempFlags.VisitAll(func(f *pflag.Flag) {
				if f.Changed && cmd.Flags().Lookup(f.Name) == nil {
					unknownArgs = append(unknownArgs, "--"+f.Name)
					if f.Value.Type() != "bool" {
						unknownArgs = append(unknownArgs, f.Value.String())
					}
				}
			})
			fmt.Println("未知 flags:", unknownArgs)

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

	for _, app := range g.Config.Project.Apps {
		err = appctl.New(baseDir).Start(app)
		if err != nil {
			utils.Fatal("start app", err)
		}
	}
}

func startApp(args []string) {
	fmt.Println("xxxxxxxxxxxxxxxx:", args)

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

	for _, app := range apps {
		if len(args) > 1 {
			err = appctl.New(baseDir).Start(app, args[1:]...)
		} else {
			err = appctl.New(baseDir).Start(app)
		}
		if err != nil {
			utils.Fatal("start app", err)
		}
	}
}
