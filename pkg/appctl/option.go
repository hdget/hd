package appctl

type Option func(impl *appCtlImpl)

func WithPbOutputDir(outputDir string) Option {
	return func(impl *appCtlImpl) {
		impl.pbOutputDir = outputDir
	}
}

func WithPbOutputPackage(outputPackage string) Option {
	return func(impl *appCtlImpl) {
		impl.pbOutputPackage = outputPackage
	}
}

func WithBinOutputDir(binDir string) Option {
	return func(impl *appCtlImpl) {
		impl.binDir = binDir
	}
}

func WithPluginOutputDir(pluginDir string) Option {
	return func(impl *appCtlImpl) {
		impl.pluginDir = pluginDir
	}
}

func WithPlugins(plugins []string) Option {
	return func(impl *appCtlImpl) {
		impl.plugins = plugins
	}
}

func WithPbGRPC(genGRPC bool) Option {
	return func(impl *appCtlImpl) {
		impl.pbGenGRPC = genGRPC
	}
}
