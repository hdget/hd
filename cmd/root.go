package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/hdget/hd/cmd/appctl"
	"github.com/hdget/hd/cmd/cluster"
	"github.com/hdget/hd/cmd/gen"
	"github.com/hdget/hd/cmd/sourcecode"
	"github.com/hdget/hd/g"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	rootCmd = &cobra.Command{
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
		},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&g.Debug, "debug", "d", false, "--debug")

	rootCmd.AddCommand(gen.Command)
	rootCmd.AddCommand(appctl.Command)
	rootCmd.AddCommand(sourcecode.Command)
	rootCmd.AddCommand(cluster.Command)
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
