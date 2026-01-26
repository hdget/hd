package sourcecode

import (
	"os"

	"github.com/hdget/hd/pkg/sourcecode"
	"github.com/hdget/hd/pkg/sourcecode/dapr"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	argPatchSkipDirs       []string // inspect时需要跳过的目录
	argAssetsPath          string
	argHandler             string
	subCmdHandleSourceCode = &cobra.Command{
		Use:   "handle",
		Short: "handle source code",
		Run: func(cmd *cobra.Command, args []string) {
			handleSourceCode()
		},
	}
)

func init() {
	subCmdHandleSourceCode.PersistentFlags().StringVarP(&argAssetsPath, "assets", "a", "assets", "--assets assets")
	subCmdHandleSourceCode.PersistentFlags().StringSliceVarP(&argPatchSkipDirs, "skip", "s", []string{"autogen"}, "--entry [import_path.func]")
	subCmdHandleSourceCode.PersistentFlags().StringVarP(&argHandler, "handler", "h", "dapr", "--handler dapr")
}

func handleSourceCode() {
	srcDir, err := os.Getwd()
	if err != nil {
		utils.Fatal("get source code dir", err)
	}

	switch argHandler {
	case "dapr":
		err = dapr.New(srcDir,
			sourcecode.WithSkipDirs(argPatchSkipDirs...),
			sourcecode.WithAssetPath(argAssetsPath),
		).Handle()
	case "webservice":
	}

	if err != nil {
		utils.Fatal("handle source code", err)
	}
}
