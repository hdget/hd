package cluster

import (
	"github.com/spf13/cobra"
)

var (
	argNeedClean   bool
	argClusterIp   string
	argClusterSize int
	Command        = &cobra.Command{
		Use: "cluster",
	}
)

func init() {
	Command.AddCommand(subCmdRestartCluster)
	Command.AddCommand(subCmdStartCluster)
	Command.AddCommand(subCmdStopCluster)
}
