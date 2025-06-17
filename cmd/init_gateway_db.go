package cmd

import (
	"fmt"
	"github.com/hdget/hd/assets"
	"github.com/hdget/hd/pkg/env"
	"github.com/hdget/hd/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
)

const (
	dbGatewaySchema = "gateway_schema.db"
	dbGateway       = "gateway.db"
)

var cmdInitGatewayDb = &cobra.Command{
	Use: "init_gateway_db",
	Run: func(cmd *cobra.Command, args []string) {
		initGatewayDb()
	},
	Long:  "initialize gateway db",
	Short: "initialize gateway db",
}

func initGatewayDb() {
	fmt.Println("")
	fmt.Println("=== initialize gateway database ===")
	fmt.Println("")

	workDir, found := env.GetHdWorkDir()
	if !found {
		utils.Fatal("hd work dir not set")
	}

	destDbFile := filepath.Join(workDir, dbGateway)
	if _, err := os.Stat(destDbFile); !os.IsNotExist(err) {
		fmt.Printf(" * database file: %s exists, skip...\n", destDbFile)
		return
	}

	if err := createEmptyGatewayDb(destDbFile); err != nil {
		utils.Fatal("create empty db", err)
	}
	fmt.Printf(" * database file: %s created\n", destDbFile)
}

func createEmptyGatewayDb(absDbFile string) error {
	// go generate运行时实在main那一级的目录
	data, err := assets.Manager.ReadFile(path.Join("db", dbGatewaySchema))
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return errors.New("empty gateway.db")
	}

	// create db file in workDir
	err = os.WriteFile(absDbFile, data, 0666)
	if err != nil {
		return err
	}

	return nil
}
