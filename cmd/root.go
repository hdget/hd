package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"runtime/debug"
)

var rootCmd = &cobra.Command{}

func init() {
	rootCmd.AddCommand(cmdProtoGen)
	rootCmd.AddCommand(cmdProtoRefine)
	rootCmd.AddCommand(cmdWindowsCommand)
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
