package cluster

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/tools"
	"github.com/hdget/hd/pkg/utils"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"time"
)

func (impl *clusterImpl) Start() error {
	if impl.isConsulStared() {
		return errors.New("cluster is already started")
	}

	// 检查依赖的工具是否安装
	if impl.clusterIp != "" && impl.clusterSize > 1 {
		return impl.startMultiNodeConsul()
	}

	return impl.startStandaloneConsul()
}

func (impl *clusterImpl) startStandaloneConsul() error {
	fmt.Printf("start consul in standard mode\n")

	localIp, err := utils.GetLocalIP()
	if err != nil {
		return err
	}

	if g.Debug {
		fmt.Println("local ip: ", localIp)
	}

	if err = impl.startConsul(localIp, localIp); err != nil {
		return err
	}

	return nil
}

func (impl *clusterImpl) startMultiNodeConsul() error {
	fmt.Printf("start consul in multi node mode, join to: %s\n", impl.clusterIp)

	if impl.clusterIp == "" {
		return errors.New("cluster ip is empty")
	}

	if impl.clusterSize <= 1 {
		return errors.New("cluster size must be greater than 1")
	}

	localIp, err := utils.GetLocalIP()
	if err != nil {
		return err
	}

	if g.Debug {
		fmt.Println("local ip: ", localIp)
		fmt.Println("cluster ip: ", impl.clusterIp)
		fmt.Println("cluster size: ", impl.clusterSize)
	}

	if err = impl.startConsul(localIp, impl.clusterIp); err != nil {
		return err
	}

	return nil
}

func (impl *clusterImpl) isConsulStared() bool {
	_, err := script.Exec(`consul members`).String()
	return err == nil
}

func (impl *clusterImpl) startConsul(localIp, clusterIp string) error {
	dataDir := filepath.Join(os.TempDir(), "consul")

	cmd := fmt.Sprintf("consul agent -server -data-dir %s -ui -client 127.0.0.1 -bootstrap-expect %d -bind %s -retry-join %s",
		dataDir, impl.clusterSize, localIp, clusterIp)

	healthFunc := func() bool {
		_, err := script.Exec("consul members").Stdout()
		return err == nil
	}

	err := tools.RunDaemon("consul", cmd, healthFunc, 10*time.Second)
	if err != nil {
		return err
	}
	return nil
}
