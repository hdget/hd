package protocompile

type Option func(impl *protobufCompilerImpl)

func WithGRPC(genGRPC bool) Option {
	return func(impl *protobufCompilerImpl) {
		impl.genGRPC = genGRPC
	}
}
