package protocompile

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/tools"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

type ProtobufCompiler interface {
	Compile(sourceProtoDir, outputDir, pkgName string) error // 生成protobuf.pb文件
}

type protobufCompilerImpl struct {
	genGRPC bool
}

var (
	allTools = []tools.Tool{
		tools.Protoc(),
		tools.ProtocGogoFaster(),
	}
	// cmdProtocGen = `protoc --proto_path=%s --gogofaster_out=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,:%s %s`
	cmdProtocGen     = `protoc --proto_path=%s %s --go_out=%s --go_opt=M%s=./%s`
	cmdProtocGenGRPC = `--go-grpc_out=%s --go-grpc_opt=M%s=./%s`
)

func New(options ...Option) ProtobufCompiler {
	impl := &protobufCompilerImpl{}
	for _, apply := range options {
		apply(impl)
	}
	return impl
}

func (impl *protobufCompilerImpl) Compile(sourceProtoDir, outputDir, pkgName string) error {
	if g.Debug {
		fmt.Println("===> protobuf compiling...")
	}

	// 检查依赖的工具是否安装
	if err := tools.Check(allTools...); err != nil {
		return err
	}

	// 获取当前目录下所有 .proto 文件（不包括子目录）
	protoFiles, err := findProtoFiles(sourceProtoDir)
	if err != nil {
		return errors.Wrap(err, "查找.proto文件失败")
	}
	if len(protoFiles) == 0 {
		return nil
	}

	// 创建输出目录
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	for _, f := range protoFiles {
		fmt.Printf("Compiling: %s\n", f)
		// 构建 protoc 命令
		cmds := []string{fmt.Sprintf(cmdProtocGen, filepath.ToSlash(sourceProtoDir), f, filepath.ToSlash(outputDir), f, pkgName)}
		if impl.genGRPC {
			cmds = append(cmds, fmt.Sprintf(cmdProtocGenGRPC, filepath.ToSlash(sourceProtoDir), f, pkgName))
		}

		command := strings.Join(cmds, " ")

		fmt.Println(command)

		// 执行编译
		output, err := script.Exec(command).String()
		if err != nil {
			fmt.Println(command)
			return errors.Wrapf(err, "protoc编译失败, output: %s", output)
		}
	}

	return nil
}

func findProtoFiles(dir string) ([]string, error) {
	var protoFiles []string

	// 读取目录内容
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %v", err)
	}

	for _, entry := range entries {
		// 跳过子目录
		if entry.IsDir() {
			continue
		}

		// 检查是否是 .proto 文件
		if strings.HasSuffix(entry.Name(), ".proto") {
			protoFiles = append(protoFiles, entry.Name())
		}
	}

	return protoFiles, nil
}
