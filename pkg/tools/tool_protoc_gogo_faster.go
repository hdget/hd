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
	return AllPlatform().GoInstall("github.com/gogo/protobuf/protoc-gen-gogofaster@latest")
}

func (t *protocGogofasterTool) WindowsInstall() error {
	return AllPlatform().GoInstall("github.com/gogo/protobuf/protoc-gen-gogofaster@latest")
}
