package gen

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	argGenDatabase = struct {
		dbname      string
		generateAll bool
	}{}

	subCmdGenDatabase = &cobra.Command{
		Use:   "database",
		Short: "generate database",
		Run: func(cmd *cobra.Command, args []string) {
			generateDatabase()
		},
	}
)

func init() {
	subCmdGenDatabase.PersistentFlags().StringVarP(&argGenDatabase.dbname, "dbname", "", "", "--dbname [dbname]")
	subCmdGenDatabase.PersistentFlags().BoolVarP(&argGenDatabase.generateAll, "all", "", false, "--all")
}

func generateDatabase() {
	fmt.Println("generate database")
}
