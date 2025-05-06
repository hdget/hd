package cluster

import (
	"github.com/hdget/hd/pkg/cluster"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	subCmdStartCluster = &cobra.Command{
		Use:   "start",
		Short: "start cluster",
		Run: func(cmd *cobra.Command, args []string) {
			startCluster()
		},
	}
)

func init() {
	subCmdStartCluster.PersistentFlags().StringVarP(&argClusterIp, "cluster-ip", "", "", "--cluster-ip 192.168.0.1")
	subCmdStartCluster.PersistentFlags().IntVarP(&argClusterSize, "cluster-size", "", 0, "--cluster-size 1")
}

func startCluster() {
	options := make([]cluster.Option, 0)
	if argClusterSize > 0 {
		options = append(options, cluster.WithClusterSize(argClusterSize))
	}
	if argClusterIp != "" {
		options = append(options, cluster.WithClusterIp(argClusterIp))
	}
	if argNeedClean {
		options = append(options, cluster.WithClusterIp(argClusterIp))
	}

	instance, err := cluster.New(options...)
	if err != nil {
		utils.Fatal("new cluster", err)
	}

	if err = instance.Start(); err != nil {
		utils.Fatal("start cluster", err)
	}
}
