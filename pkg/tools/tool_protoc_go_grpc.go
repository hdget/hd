package tools

type protocGoGRPCTool struct {
	*toolImpl
}

func ProtocGoGRPC() Tool {
	return &protocGoGRPCTool{
		toolImpl: &toolImpl{
			name: "protoc-gen-go",
		},
	}
}

func (t *protocGoGRPCTool) LinuxInstall() error {
	return AllPlatform().GoInstall("google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest")
}

func (t *protocGoGRPCTool) WindowsInstall() error {
	return AllPlatform().GoInstall("google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest")
}
