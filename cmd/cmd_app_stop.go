package cmd

import (
	"github.com/hdget/hd/pkg/appctl"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

var (
	cmdAppStop = &cobra.Command{
		Use:   "stop",
		Short: "stop app",
		Run: func(cmd *cobra.Command, args []string) {
			stopApp(args)
		},
	}
)

func stopApp(args []string) {
	if len(args) != 1 {
		utils.Fatal("Usage: stop <app>")
	}

	baseDir, err := os.Getwd()
	if err != nil {
		utils.Fatal("get current dir", err)
	}

	err = appctl.New(baseDir, appctl.WithDebug(argDebug)).Stop(args[0])
	if err != nil {
		utils.Fatal("stop app", err)
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
