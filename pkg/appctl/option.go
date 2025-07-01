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

func WithPbGRPC(genGRPC bool) Option {
	return func(impl *appCtlImpl) {
		impl.pbGenGRPC = genGRPC
	}
}
