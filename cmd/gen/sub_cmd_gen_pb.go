package gen

import (
	"fmt"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/protocompile"
	"github.com/hdget/hd/pkg/protorefine"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	argProtobufGen = struct {
		outputDir     string
		outputPackage string
		generateAll   bool
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

func init() {
	subCmdGenProtobuf.PersistentFlags().StringVarP(&argProtobufGen.outputPackage, "package", "", "pb", "--package <package>")
	subCmdGenProtobuf.PersistentFlags().StringVarP(&argProtobufGen.outputDir, "output-dir", "", "autogen", "--output-dir <sub_dir>")
	subCmdGenProtobuf.PersistentFlags().BoolVarP(&argProtobufGen.generateAll, "all", "", false, "--all")
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
	matchFiles := []string{fmt.Sprintf("%s.proto", filepath.Base(rootGolangModule))}
	protoRepository, err := utils.FindDirContainingFiles(srcDir, matchFiles, filepath.Join(srcDir, argProtobufGen.outputDir))
	if err != nil {
		return err
	}

	// 第一步：先精简proto文件
	var protoDir string
	if argProtobufGen.generateAll {
		protoDir = protoRepository
	} else {
		var prOptions []protorefine.Option
		if g.Debug {
			prOptions = append(prOptions, protorefine.WithDebug(true))
		}
		protoDir, err = protorefine.New(prOptions...).Refine(protorefine.Argument{
			GolangModule:        rootGolangModule,
			GolangSourceCodeDir: srcDir,
			ProtoRepository:     protoRepository,
			OutputPackage:       argProtobufGen.outputPackage,
			OutputDir:           argProtobufGen.outputDir,
		})
		if err != nil {
			return err
		}
	}

	// 第二步：编译protobuf
	var pcOptions []protocompile.Option
	if g.Debug {
		pcOptions = append(pcOptions, protocompile.WithDebug(true))
	}
	outputPbDir := filepath.Join(srcDir, argProtobufGen.outputDir, argProtobufGen.outputPackage)
	err = protocompile.New(pcOptions...).Compile(protoDir, outputPbDir)
	if err != nil {
		return err
	}

	return nil
}
