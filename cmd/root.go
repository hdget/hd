package cmd

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/env"
	"github.com/hdget/hd/pkg/utils"
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
	rootCmd.AddCommand(cmdAppBuild)
	rootCmd.AddCommand(cmdAppStart)
	rootCmd.AddCommand(cmdAppStop)
}

func Execute() {
	initialize()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(string(debug.Stack()))
		}
	}()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initialize() {
	// 读取配置
	if _, err := toml.DecodeFile(g.ConfigFile, &g.Config); err != nil {
		utils.Fatal(fmt.Sprintf("read config file, file: %s", g.ConfigFile), err)
	}

	for _, r := range g.Config.Repos {
		g.RepoConfigs[r.Name] = r
	}

	for _, t := range g.Config.Tools {
		g.ToolConfigs[t.Name] = t
	}

	// 初始化环境变量
	for k, v := range env.GetExportedEnvs() {
		if err := os.Setenv(k, v); err != nil {
			utils.Fatal("export HD environment variable", err)
		}
	}
}
