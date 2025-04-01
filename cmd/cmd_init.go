package cmd

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"os"
	"path/filepath"
)

var (
	cmdInit = &cobra.Command{
		Use:   "init",
		Short: "initialize running environment",
		Run: func(cmd *cobra.Command, args []string) {
			runWindowsCommand(args[0])
		},
	}
)

func runWindowsCommand(dir string) {
	decoder := simplifiedchinese.GBK.NewDecoder()

	n, err := script.Exec(fmt.Sprintf(`cmd /c dir %s`, filepath.Dir(dir))).WithStdout(transform.NewWriter(os.Stdout, decoder)).Stdout()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(n)
}
