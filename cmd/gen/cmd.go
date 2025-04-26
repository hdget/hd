package gen

import (
	"github.com/spf13/cobra"
)

var (
	Command = &cobra.Command{
		Use: "gen",
	}
)

func init() {
	Command.AddCommand(subCmdGenProtobuf)
	Command.AddCommand(subCmdGenConfig)
	Command.AddCommand(subCmdGenDatabase)
}
