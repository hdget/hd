package cmd

import (
	"fmt"
	"github.com/hdget/hd/cmd/appctl"
	"github.com/hdget/hd/cmd/cluster"
	"github.com/hdget/hd/cmd/gen"
	"github.com/hdget/hd/cmd/sourcecode"
	"github.com/hdget/hd/g"
	"github.com/spf13/cobra"
	"os"
	"runtime/debug"
)

var (
	rootCmd = &cobra.Command{
		Args: nil,
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&g.Debug, "debug", "d", false, "--debug")

	rootCmd.AddCommand(gen.Command)
	rootCmd.AddCommand(appctl.Command)
	rootCmd.AddCommand(sourcecode.Command)
	rootCmd.AddCommand(cluster.Command)
	rootCmd.AddCommand(cmdInitGatewayDb)
}

func Execute() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(string(debug.Stack()))
		}
	}()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
