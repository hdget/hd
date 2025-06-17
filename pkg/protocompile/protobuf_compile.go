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
	Compile(sourceProtoDir, outputPbDir string) error // 生成protobuf.pb文件
}

type protobufCompilerImpl struct {
}

var (
	allTools = []tools.Tool{
		tools.Protoc(),
		tools.ProtocGogoFaster(),
	}
	cmdProtocGen = `protoc --proto_path=%s --gogofaster_out=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,:%s %s`
)

func New(options ...Option) ProtobufCompiler {
	impl := &protobufCompilerImpl{}
	for _, apply := range options {
		apply(impl)
	}
	return impl
}

func (impl *protobufCompilerImpl) Compile(sourceProtoDir, outputPbDir string) error {
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
	err = os.MkdirAll(outputPbDir, 0755)
	if err != nil {
		return err
	}

	for _, f := range protoFiles {
		fmt.Printf("Compiling: %s\n", f)
		// 构建 protoc 命令
		cmd := fmt.Sprintf(cmdProtocGen, filepath.ToSlash(sourceProtoDir), filepath.ToSlash(outputPbDir), f)
		// 执行编译
		output, err := script.Exec(cmd).String()
		if err != nil {
			fmt.Println(cmd)
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
