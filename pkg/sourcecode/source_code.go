package sourcecode

type Handler interface {
	Handle() error
}

type HandlerImpl struct {
	SrcDir     string
	SkipDirs   []string
	AssetsPath string
}
