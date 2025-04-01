package cmd

import (
	"fmt"
	"github.com/hdget/hd/pkg/protogen"
	"github.com/spf13/cobra"
	"os"
)

var (
	argProtoGen = protogen.Argument{}
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
	cmdProtoGen.PersistentFlags().StringVarP(&argProtoGen.OutputDir, "output-dir", "", "autogen", "--output-dir <sub_dir>")
	cmdProtoGen.PersistentFlags().BoolVarP(&argProtoGen.Debug, "debug", "", false, "--debug")
}

func runProtocGen() error {
	srcDir, err := os.Getwd()
	if err != nil {
		return err
	}

	err = protogen.New(srcDir).Generate(argProtoGen)
	if err != nil {
		return err
	}

	return nil
}
