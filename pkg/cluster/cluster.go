package cluster

import (
	"github.com/hdget/hd/pkg/tools"
	"os"
	"path/filepath"
)

type Cluster interface {
	Start() error
	Stop() error
}

type clusterImpl struct {
	clusterIp   string
	clusterSize int
	needClean   bool
}

func New(options ...Option) (Cluster, error) {
	// 检查依赖的工具是否安装
	if err := tools.Check(tools.Consul()); err != nil {
		return nil, err
	}

	impl := &clusterImpl{
		clusterSize: 1,
	}

	for _, apply := range options {
		apply(impl)
	}

	return impl, nil
}

func (*clusterImpl) getConsulDataDir() string {
	return filepath.Join(os.TempDir(), "consul")
}
