package protogen

type Argument struct {
	OutputDir string
	Debug     bool
}

func (a Argument) validate() error {
	return nil
}
