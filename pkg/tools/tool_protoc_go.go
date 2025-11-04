package tools

import "github.com/hdget/hd/g"

type protocGoTool struct {
	*toolImpl
}

func ProtocGo() Tool {
	return &protocGoTool{
		toolImpl: &toolImpl{
			&g.ToolConfig{
				Name: "protoc-gen-go",
			},
		},
	}
}

func (t *protocGoTool) LinuxInstall() error {
	return AllPlatform().GoInstall("google.golang.org/protobuf/cmd/protoc-gen-go@latest")
}

func (t *protocGoTool) WindowsInstall() error {
	return AllPlatform().GoInstall("google.golang.org/protobuf/cmd/protoc-gen-go@latest")
}
