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

			if argAll {
				startAllApp()
			} else {
				startApp(args)
			}

		},
	}
)

func init() {
	// 禁用Cobra的默认解析
	subCmdStartApp.DisableFlagParsing = true

	// 手动解析 flags
	flagSet := pflag.NewFlagSet("custom", pflag.ContinueOnError)
	flagSet.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{
		UnknownFlags: true, // 允许未知 flags
	}

	// 复制Cobra已注册的flags到手动解析的flagSet
	subCmdStartApp.Flags().VisitAll(func(f *pflag.Flag) {
		flagSet.AddFlag(f)
	})

	// 解析命令行参数
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		fmt.Println("解析错误:", err)
		return
	}

	// 获取未定义的flags
	var extraParams []string
	flagSet.VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			// 检查是否在 Cobra 中注册过
			if subCmdStartApp.Flags().Lookup(f.Name) == nil {
				extraParams = append(extraParams, f.Name)
			}
		}
	})

	// 7. 输出结果
	knownFlag, _ := subCmdStartApp.Flags().GetString("known-flag")
	fmt.Println("已知flag的值:", knownFlag)
	fmt.Println("未定义的flags:", extraParams)
}

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
