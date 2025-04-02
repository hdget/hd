package cmd

import (
	"github.com/hdget/hd/pkg/app"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var (
	cmdApp = &cobra.Command{
		Use:   "app",
		Short: "app related commands",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	subCmdAppDeploy = &cobra.Command{
		Use:   "deploy",
		Short: "deploy app",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	subCmdAppBuild = &cobra.Command{
		Use:   "build",
		Short: "build app",
		Run: func(cmd *cobra.Command, args []string) {
			buildApp(args)
		},
	}
	subCmdAppStart = &cobra.Command{
		Use:   "start",
		Short: "start app",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	subCmdAppStop = &cobra.Command{
		Use:   "stop",
		Short: "stop app",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	subCmdAppRun = &cobra.Command{
		Use:   "run",
		Short: "run app",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)

func init() {
	cmdApp.AddCommand(subCmdAppStart)
	cmdApp.AddCommand(subCmdAppStop)
	cmdApp.AddCommand(subCmdAppRun)
	cmdApp.AddCommand(subCmdAppBuild)
	cmdApp.AddCommand(subCmdAppDeploy)

}

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

	err = app.New(baseDir).Build(refName, apps...)
	if err != nil {
		log.Fatal(err)
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
