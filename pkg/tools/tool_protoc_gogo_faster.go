package tools

type protocGogofasterTool struct {
	*toolImpl
}

func ProtocGogoFaster() Tool {
	return &protocGogofasterTool{
		toolImpl: &toolImpl{
			name: "protoc-gen-gogofaster",
		},
	}
}

func (t *protocGogofasterTool) LinuxInstall() error {
	return t.installFromSourceCode()
}

func (t *protocGogofasterTool) WindowsInstall() error {
	return t.installFromSourceCode()
}

// installProtocGenGogofaster 尝试安装 protoc-gen-gogofaster
func (t *protocGogofasterTool) installFromSourceCode() error {
	return NoArch().GoInstall("github.com/gogo/protobuf/protoc-gen-gogofaster")
	//
	//cmd := "go install github.com/gogo/protobuf/protoc-gen-gogofaster@latest"
	//output, err := script.Exec(cmd).String()
	//if err != nil {
	//	return errors.New(output)
	//}
	//return nil
}
