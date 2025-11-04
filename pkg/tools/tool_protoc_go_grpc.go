package tools

import "github.com/hdget/hd/g"

type protocGoGRPCTool struct {
	*toolImpl
}

func ProtocGoGRPC() Tool {
	return &protocGoGRPCTool{
		toolImpl: &toolImpl{
			&g.ToolConfig{
				Name: "protoc-gen-go-grpc",
			},
		},
	}
}

func (t *protocGoGRPCTool) LinuxInstall() error {
	return AllPlatform().GoInstall("google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest")
}

func (t *protocGoGRPCTool) WindowsInstall() error {
	return AllPlatform().GoInstall("google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest")
}
