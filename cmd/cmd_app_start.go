package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cmdAppStart = &cobra.Command{
		Use:   "start",
		Short: "start app",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {

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
