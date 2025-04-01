package protorefine

type Option func(impl *protoRefineImpl)

func WithDebug(debug bool) Option {
	return func(impl *protoRefineImpl) {
		impl.debug = debug
	}
}
