package cluster

import (
	"github.com/hdget/hd/pkg/cluster"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	argClusterIp      string
	argClusterSize    int
	subCmdInitCluster = &cobra.Command{
		Use:   "init",
		Short: "init cluster",
		Run: func(cmd *cobra.Command, args []string) {
			initCluster()
		},
	}
)

func init() {
	subCmdInitCluster.PersistentFlags().StringVarP(&argClusterIp, "cluster-ip", "", "", "--cluster-ip 192.168.0.1")
	subCmdInitCluster.PersistentFlags().IntVarP(&argClusterSize, "cluster-size", "", 0, "--cluster-size 1")
}

func initCluster() {
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

	if err = instance.Restart(); err != nil {
		utils.Fatal("restart cluster", err)
	}
}
