package cluster

import (
	"github.com/hdget/hd/pkg/cluster"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	subCmdStopCluster = &cobra.Command{
		Use:   "stop",
		Short: "stop cluster",
		Run: func(cmd *cobra.Command, args []string) {
			stopCluster()
		},
	}
)

func stopCluster() {
	options := make([]cluster.Option, 0)
	if argNeedClean {
		options = append(options, cluster.WithClusterIp(argClusterIp))
	}

	instance, err := cluster.New(options...)
	if err != nil {
		utils.Fatal("new cluster", err)
	}

	if err = instance.Stop(); err != nil {
		utils.Fatal("stop cluster", err)
	}
}
