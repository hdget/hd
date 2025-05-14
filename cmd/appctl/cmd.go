package appctl

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/env"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	argAll  bool // 是否操作所有app
	Command = &cobra.Command{
		Use: "app",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initialize()
		},
	}
)

func init() {
	Command.PersistentFlags().BoolVarP(&argAll, "all", "a", false, "--all")

	Command.AddCommand(subCmdBuildApp)
	Command.AddCommand(subCmdDeployApp)
	Command.AddCommand(subCmdRestartApp)
	Command.AddCommand(subCmdStartApp)
	Command.AddCommand(subCmdStopApp)
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
	if err := env.Initialize(); err != nil {
		utils.Fatal("environment variables not initialized", err)
	}
}
