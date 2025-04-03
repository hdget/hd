package env

import (
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/hdget/hd/g"
	"github.com/joho/godotenv"
	"os"
)

const (
	envHdNamespace = "HD_NAMESPACE"
	envHdEnv       = "HD_ENV"
)

var (
	supportedEnvs = []string{"prod", "test"}
)

func GetExportedEnvs() map[string]string {
	return map[string]string{
		envHdNamespace: g.Config.Project.Name,
		envHdEnv:       g.Config.Project.Env,
	}
}

func GetHdEnv() (string, error) {
	env, exists := os.LookupEnv(envHdEnv)
	if !exists {
		return "", fmt.Errorf("env not found, env: %s", envHdEnv)
	}

	if !pie.Contains(supportedEnvs, env) {
		return "", fmt.Errorf("unsupported env, env: %s", env)
	}
	return env, nil
}

func GetHdNamespace() (string, bool) {
	return os.LookupEnv(envHdNamespace)
}

func WriteEnvFile(filename string, data map[string]string) error {
	/// 读取现有内容（如果文件存在）
	existing, err := godotenv.Read(filename)
	if err != nil && !os.IsNotExist(err) {
		return nil
	}

	// 合并新旧值
	for k, v := range data {
		existing[k] = v
	}

	// 写入文件
	return godotenv.Write(existing, filename)
}
