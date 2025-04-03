package appctl

type Option func(impl *appCtlImpl)

func WithDebug(debug bool) Option {
	return func(impl *appCtlImpl) {
		impl.debug = debug
	}
}
