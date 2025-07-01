package tools

type protocGoTool struct {
	*toolImpl
}

func ProtocGo() Tool {
	return &protocGoTool{
		toolImpl: &toolImpl{
			name: "protoc-gen-go",
		},
	}
}

func (t *protocGoTool) LinuxInstall() error {
	return AllPlatform().GoInstall("google.golang.org/protobuf/cmd/protoc-gen-go@latest")
}

func (t *protocGoTool) WindowsInstall() error {
	return AllPlatform().GoInstall("google.golang.org/protobuf/cmd/protoc-gen-go@latest")
}
