package protocompile

type Option func(impl *protobufCompilerImpl)

func WithDebug(debug bool) Option {
	return func(impl *protobufCompilerImpl) {
		impl.debug = debug
	}
}
