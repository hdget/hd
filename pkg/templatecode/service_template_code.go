package templatecode

type ServiceGenerator interface {
	Gen(destDir, name string) error
}

type serviceGeneratorImpl struct{}
