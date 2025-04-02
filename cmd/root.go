package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"runtime/debug"
)

var (
	argDebug bool
	rootCmd  = &cobra.Command{}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&argDebug, "debug", "", false, "--debug")

	rootCmd.AddCommand(cmdProtobufGen)
	rootCmd.AddCommand(cmdApp)
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
