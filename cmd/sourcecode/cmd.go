package sourcecode

import (
	"github.com/spf13/cobra"
)

var (
	Command = &cobra.Command{
		Use: "sourcecode",
	}
)

func init() {
	Command.AddCommand(subCmdHandleSourceCode)
}
