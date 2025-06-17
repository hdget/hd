package sourcecode

import (
	"github.com/hdget/hd/pkg/sourcecode"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

var (
	argPatchSkipDirs       []string // inspect时需要跳过的目录
	argAssetsPath          string
	subCmdHandleSourceCode = &cobra.Command{
		Use:   "handle",
		Short: "handle source code",
		Run: func(cmd *cobra.Command, args []string) {
			handleSourceCode()
		},
	}
)

func init() {
	subCmdHandleSourceCode.PersistentFlags().StringVarP(&argAssetsPath, "assets-path", "", "assets", "--assets-path assets")
	subCmdHandleSourceCode.PersistentFlags().StringSliceVarP(&argPatchSkipDirs, "skip", "", []string{"autogen"}, "--entry [import_path.func]")
}

func handleSourceCode() {
	srcDir, err := os.Getwd()
	if err != nil {
		utils.Fatal("get source code dir", err)
	}

	err = sourcecode.New(srcDir,
		sourcecode.WithSkipDirs(argPatchSkipDirs...),
		sourcecode.WithAssetPath(argAssetsPath),
	).Handle()
	if err != nil {
		utils.Fatal("handle source code", err)
	}
}
