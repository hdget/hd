package cmd

import (
	"fmt"
	"github.com/hdget/hd/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

const (
	fileDatabase = "gateway.db"
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
	baseDir, err := os.Getwd()
	if err != nil {
		utils.Fatal("get current dir", err)
	}

	absDbFile := filepath.Join(baseDir, fileDatabase)
	if _, err := os.Stat(absDbFile); !os.IsNotExist(err) {
		fmt.Printf(" * database file: %s exists, skip...\n", absDbFile)
	} else {
		//if err = createEmptyDb(absDbFile); err != nil {
		//	return err
		//}
		//
		//if err = addGatewayRoutes(); err != nil {
		//	return err
		//}

		fmt.Printf(" * database file: %s created\n", absDbFile)
	}

}

//
//func addGatewayRoutes() error {
//	// update myself routes
//	handlers, err := getParsedExposedHandlers()
//	if err != nil {
//		return err
//	}
//
//	if len(handlers) == 0 {
//		return nil
//	}
//
//	// 更新gateway的路由
//	dbClient, err := sqlite3_sqlboiler.NewClient(absDbFile)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		_ = dbClient.Close()
//	}()
//
//	return nil
//}
//
//func createEmptyDb(absDbFile string) error {
//	// go generate运行时实在main那一级的目录
//	data, err := assets.Manager.ReadFile(path.Join("db", fileEmptyDatabase))
//	if err != nil {
//		return err
//	}
//
//	if len(data) == 0 {
//		return errors.New("empty schema.db")
//	}
//
//	// create db file in workDir
//	err = os.WriteFile(absDbFile, data, 0666)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func getParsedExposedHandlers() ([]*protobuf.DaprHandler, error) {
//	// IMPORTANT: embedfs使用的是斜杠来获取文件路径,在windows平台下如果使用filepath来处理路径会导致问题
//	data, err := assets.Manager.ReadFile(path.Join("json", fileExposedHandlers))
//	if err != nil {
//		return nil, err
//	}
//
//	var handlers []*protobuf.DaprHandler
//	err = json.Unmarshal(data, &handlers)
//	if err != nil {
//		return nil, err
//	}
//	return handlers, nil
//}
