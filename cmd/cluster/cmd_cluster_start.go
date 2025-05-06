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
	instance, err := cluster.New(getClusterOptions()...)
	if err != nil {
		utils.Fatal("new cluster", err)
	}

	if err = instance.Start(); err != nil {
		utils.Fatal("start cluster", err)
	}
}
