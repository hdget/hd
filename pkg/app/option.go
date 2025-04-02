package app

type Option func(impl *appControlImpl)

func WithDebug(debug bool) Option {
	return func(impl *appControlImpl) {
		impl.debug = debug
	}
}

func WithPbOutputPackage(pkgName string) Option {
	return func(impl *appControlImpl) {
		impl.pbOutputPackage = pkgName
	}
}

func WithPbOutputDir(outputDir string) Option {
	return func(impl *appControlImpl) {
		impl.pbOutputDir = outputDir
	}
}
