package cmd

import (
	"fmt"
	"github.com/hdget/hd/pkg/protogen"
	"github.com/spf13/cobra"
	"os"
)

var (
	protoDir      string   // 指定的proto所在的目录
	outputDir     string   // 输出目录
	protoDirFiles []string // 在没有指定protoDir的时候去通过matchFiles去找proto所在的文件目录

	cmdProtoGen = &cobra.Command{
		Use:   "gen",
		Short: "compile protobuf",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runProtocGen(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	cmdProtoGen.PersistentFlags().StringVarP(&outputDir, "output-dir", "", "autogen", "")
	cmdProtoGen.PersistentFlags().StringSliceVarP(&protoDirFiles, "proto-dir-files", "", []string{"gateway.proto"}, "")
	cmdProtoGen.PersistentFlags().StringVarP(&protoDir, "proto-dir", "", "", "")
}

func runProtocGen() error {
	srcDir, err := os.Getwd()
	if err != nil {
		return err
	}

	err = protogen.New(srcDir).Generate(protoDir, outputDir)
	if err != nil {
		return err
	}

	return nil
}
