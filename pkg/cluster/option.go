package cluster

type Option func(impl *clusterImpl)

func WithClusterIp(clusterIp string) Option {
	return func(impl *clusterImpl) {
		impl.clusterIp = clusterIp
	}
}

func WithClusterSize(clusterSize int) Option {
	return func(impl *clusterImpl) {
		impl.clusterSize = clusterSize
	}
}

func WithClean(needClean bool) Option {
	return func(impl *clusterImpl) {
		impl.needClean = needClean
	}
}
