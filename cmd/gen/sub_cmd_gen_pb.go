package gen

import (
	"fmt"
	"github.com/hdget/hd/pkg/protocompile"
	"github.com/hdget/hd/pkg/protorefine"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	argGenProtobuf = struct {
		outputDir     string
		outputPackage string
		generateAll   bool
		generateGRPC  bool
	}{}

	subCmdGenProtobuf = &cobra.Command{
		Use:   "pb",
		Short: "generate protobuf",
		Run: func(cmd *cobra.Command, args []string) {
			if err := protobufGenerate(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
)

const (
	protoRepoIdentifier = ".HD_PROTO_REPO" // 用来标识proto仓库所在的位置
)

func init() {
	subCmdGenProtobuf.PersistentFlags().StringVarP(&argGenProtobuf.outputPackage, "package", "", "pb", "--package <package>")
	subCmdGenProtobuf.PersistentFlags().StringVarP(&argGenProtobuf.outputDir, "output-dir", "", "autogen", "relative output dir, --output-dir <sub_dir>")
	subCmdGenProtobuf.PersistentFlags().BoolVarP(&argGenProtobuf.generateAll, "all", "", false, "--all")
	// 是否生成grpc代码
	subCmdGenProtobuf.PersistentFlags().BoolVarP(&argGenProtobuf.generateGRPC, "grpc", "", false, "--grpc")
}

func protobufGenerate() error {
	srcDir, err := os.Getwd()
	if err != nil {
		return err
	}

	rootGolangModule, err := utils.GetRootGolangModule()
	if err != nil {
		return err
	}

	// 尝试找到proto repository
	protoRepository, err := utils.FindDirContainingFiles(srcDir, []string{protoRepoIdentifier}, filepath.Join(srcDir, argGenProtobuf.outputDir))
	if err != nil {
		return err
	}

	// 第一步：先精简proto文件
	var protoDir string
	if argGenProtobuf.generateAll {
		protoDir = protoRepository
	} else {
		protoDir, err = protorefine.New().Refine(protorefine.Argument{
			GolangModule:        rootGolangModule,
			GolangSourceCodeDir: srcDir,
			ProtoRepository:     protoRepository,
			OutputPackage:       argGenProtobuf.outputPackage,
			OutputDir:           argGenProtobuf.outputDir,
		})
		if err != nil {
			return err
		}
	}

	// 第二步：编译protobuf
	absOutputDir := filepath.Join(srcDir, argGenProtobuf.outputDir)
	err = protocompile.New(protocompile.WithGRPC(argGenProtobuf.generateGRPC)).Compile(protoDir, absOutputDir, argGenProtobuf.outputPackage)
	if err != nil {
		return err
	}

	return nil
}
