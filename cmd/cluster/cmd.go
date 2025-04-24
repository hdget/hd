package cluster

import (
	"github.com/spf13/cobra"
)

var (
	argNeedClean bool
	Command      = &cobra.Command{
		Use: "cluster",
	}
)

func init() {
	// 是否需要清除cluster数据
	Command.PersistentFlags().BoolVarP(&argNeedClean, "clean", "", false, "--clean")

	Command.AddCommand(subCmdInitCluster)
	Command.AddCommand(subCmdStopCluster)
}
