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
	arg = struct {
		all       bool // 是否操作所有app
		binDir    string
		pluginDir string
		plugins   []string
	}{}

	Command = &cobra.Command{
		Use: "app",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initialize()
		},
	}
)

func init() {
	Command.PersistentFlags().BoolVarP(&arg.all, "all", "a", false, "--all")
	// app二进制文件的保存目录
	Command.PersistentFlags().StringVarP(&arg.binDir, "bin-dir", "", "bin", "relative app binary output dir, --bin-dir [dir]")
	// 指定编译哪个plugin
	Command.PersistentFlags().StringSliceVarP(&arg.plugins, "plugins", "", nil, "--plugins [plugin1,plugin2...]")
	// plugin编译输出到哪个目录
	Command.PersistentFlags().StringVarP(&arg.pluginDir, "plugin-dir", "", "plugins", "relative plugin output dir, --plugin-dir [dir]")

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

	// 初始化环境变量
	if err := env.Initialize(); err != nil {
		utils.Fatal("environment variables not initialized", err)
	}
}
