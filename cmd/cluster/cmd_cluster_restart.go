package cluster

import (
	"github.com/hdget/hd/pkg/cluster"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	subCmdRestartCluster = &cobra.Command{
		Use:   "restart",
		Short: "restart cluster",
		Run: func(cmd *cobra.Command, args []string) {
			restartCluster()
		},
	}
)

func init() {
	subCmdRestartCluster.PersistentFlags().StringVarP(&argClusterIp, "cluster-ip", "", "", "--cluster-ip 192.168.0.1")
	subCmdRestartCluster.PersistentFlags().IntVarP(&argClusterSize, "cluster-size", "", 0, "--cluster-size 1")
	subCmdRestartCluster.PersistentFlags().BoolVarP(&argNeedClean, "clean", "", false, "--clean")
}

func restartCluster() {
	instance, err := cluster.New(getClusterOptions()...)
	if err != nil {
		utils.Fatal("new cluster", err)
	}

	// stop cluster
	_ = instance.Stop()

	if err = instance.Start(); err != nil {
		utils.Fatal("start cluster", err)
	}
}

func getClusterOptions() []cluster.Option {
	options := make([]cluster.Option, 0)
	if argClusterSize > 0 {
		options = append(options, cluster.WithClusterSize(argClusterSize))
	}
	if argClusterIp != "" {
		options = append(options, cluster.WithClusterIp(argClusterIp))
	}
	if argNeedClean {
		options = append(options, cluster.WithClean(argNeedClean))
	}

	return options
}
