package cmd

import (
	"fmt"
	"github.com/hdget/hd/pkg/protorefine"
	"github.com/spf13/cobra"
	"os"
)

var (
	argProtoRefine = &protorefine.Argument{}
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
	cmdProtoRefine.PersistentFlags().StringVarP(&argProtoRefine.OutputDir, "output-dir", "", "autogen", "")
	// 原始的proto文件所在的目录
	cmdProtoRefine.PersistentFlags().StringVarP(&argProtoRefine.ProtoDir, "proto-dir", "", "", "")
}

func refineProtoFiles() error {
	srcDir, err := os.Getwd()
	if err != nil {
		return err
	}

	_, _, err = protorefine.New(srcDir).Refine(argProtoRefine)
	if err != nil {
		return err
	}

	return nil
}
