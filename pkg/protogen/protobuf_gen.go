package protogen

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/hdget/hd/pkg/protorefine"
	"github.com/hdget/hd/pkg/tools"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

type ProtobufGenerator interface {
	Generate(argument Argument) error // 生成protobuf.pb文件
}

type protoGenImpl struct {
	srcDir string
}

var (
	allTools = []tools.Tool{
		tools.Protoc(),
		tools.ProtocGogoFaster(),
	}
	cmdProtocGen = `protoc --proto_path=%s --gogofaster_out=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,:%s %s`
)

func New(srcDir string) ProtobufGenerator {
	return &protoGenImpl{
		srcDir: srcDir,
	}
}

func (impl *protoGenImpl) Generate(arg Argument) error {
	err := arg.validate()
	if err != nil {
		return err
	}

	// 检查依赖的工具是否安装
	if err = tools.Check(allTools, arg.Debug); err != nil {
		return err
	}

	// 精简proto文件
	var prOptions []protorefine.Option
	if arg.Debug {
		prOptions = append(prOptions, protorefine.WithDebug(true))
	}
	protoDir, pkgName, err := protorefine.New(impl.srcDir, prOptions...).Refine(&protorefine.Argument{
		OutputDir: arg.OutputDir,
	})
	if err != nil {
		return err
	}

	err = protocGen(protoDir, filepath.Join(arg.OutputDir, pkgName))
	if err != nil {
		return err
	}
	return nil
}

func protocGen(protoDir, outputDir string) error {
	// 获取当前目录下所有 .proto 文件（不包括子目录）
	protoFiles, err := findProtoFiles(protoDir)
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
		cmd := fmt.Sprintf(cmdProtocGen, filepath.ToSlash(protoDir), filepath.ToSlash(outputDir), f)
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
