package appctl

type Option func(impl *appCtlImpl)

func WithBinDir(binDir string) Option {
	return func(impl *appCtlImpl) {
		impl.binDir = binDir
	}
}

func WithPluginDir(pluginDir string) Option {
	return func(impl *appCtlImpl) {
		impl.pluginDir = pluginDir
	}
}

func WithPlugins(plugins []string) Option {
	return func(impl *appCtlImpl) {
		impl.pluginNames = plugins
	}
}
