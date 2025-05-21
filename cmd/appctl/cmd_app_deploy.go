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

	ctl := appctl.New(baseDir)

	for _, app := range pie.Reverse(g.Config.Project.Apps) {
		err = ctl.Stop(app)
		if err != nil {
			utils.Fatal("stop app", err)
		}
	}

	for _, app := range g.Config.Project.Apps {
		err = ctl.Build(app, ref)
		if err != nil {
			utils.Fatal("build app", err)
		}

		err = ctl.Install(app, ref)
		if err != nil {
			utils.Fatal("install app", err)
		}

		err = ctl.Start(app)
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

	ctl := appctl.New(baseDir)

	for _, app := range apps {
		err = ctl.Stop(app)
		if err != nil {
			utils.Fatal("stop app", err)
		}

		err = ctl.Build(app, ref)
		if err != nil {
			utils.Fatal("build app", err)
		}

		err = ctl.Install(app, ref)
		if err != nil {
			utils.Fatal("install app", err)
		}

		err = ctl.Start(app, extraParam)
		if err != nil {
			utils.Fatal("start app", err)
		}
	}
}

//
//func runWindowsCommand(dir string) {
//	decoder := simplifiedchinese.GBK.NewDecoder()
//
//	n, err := script.Exec(fmt.Sprintf(`cmd /c dir %s`, filepath.Dir(dir))).WithStdout(transform.NewWriter(os.Stdout, decoder)).Stdout()
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//	fmt.Println(n)
//}
