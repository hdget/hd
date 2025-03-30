package cmd

import (
	"fmt"
	"github.com/hdget/hd/pkg/protorefine"
	"github.com/spf13/cobra"
	"os"
)

var (
	arg            = &protorefine.Argument{}
	cmdProtoRefine = &cobra.Command{
		Use:   "refine",
		Short: "refine proto files",
		Run: func(cmd *cobra.Command, args []string) {
			if err := refineProtoFiles(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	// 输出的目录
	cmdProtoRefine.PersistentFlags().StringVarP(&arg.OutputDir, "output-dir", "", "autogen", "")
	// 输出的包名
	cmdProtoRefine.PersistentFlags().StringVarP(&arg.OutputPackage, "output-package", "", "pb", "")
	// 原始的proto文件所在的目录
	cmdProtoRefine.PersistentFlags().StringVarP(&arg.ProtoDir, "proto-dir", "", "..", "")
	// 该参数用来智能查找proto-dir
	cmdProtoRefine.PersistentFlags().StringSliceVarP(&arg.ProtoDirMatchFiles, "proto-dir-match-files", "", []string{"gateway.proto"}, "")
}

func refineProtoFiles() error {
	srcDir, err := os.Getwd()
	if err != nil {
		return err
	}

	err = protorefine.New(srcDir).Refine(arg)
	if err != nil {
		return err
	}

	return nil
}
