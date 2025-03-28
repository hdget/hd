package cmd

import (
	"fmt"
	"github.com/hdget/hd/pkg/protorefine"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var (
	//protoDir      string   // 指定的proto所在的目录
	//outputDir     string   // 输出目录
	//protoDirFiles []string // 在没有指定protoDir的时候去通过matchFiles去找proto所在的文件目录

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

const (
	envHdNamespace = "HD_NAMESPACE"
)

func init() {
	cmdProtoRefine.PersistentFlags().StringVarP(&outputDir, "output-dir", "", "autogen", "")
	cmdProtoRefine.PersistentFlags().StringSliceVarP(&protoDirFiles, "proto-dir-files", "", []string{"gateway.proto"}, "")
	cmdProtoRefine.PersistentFlags().StringVarP(&protoDir, "proto-dir", "", "", "")
}

func refineProtoFiles() error {
	srcDir, err := os.Getwd()
	if err != nil {
		return err
	}

	project, exists := os.LookupEnv(envHdNamespace)
	if !exists {
		return errors.New("project name not found")
	}

	pbImportPath := path.Join(project, "autogen", "pb")

	err = protorefine.New(srcDir).Refine(pbImportPath, protoDir, outputDir, "autogen")
	if err != nil {
		return err
	}

	return nil
}
