package env

import (
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/hdget/hd/g"
	"github.com/hdget/sdk/common/constant"
	"github.com/joho/godotenv"
	"os"
)

const (
	filename = ".env"
)

var (
	supportedEnvs = []string{"prod", "test", "local"}
)

func Initialize() error {
	for k, v := range map[string]string{
		constant.EnvKeyNamespace:      g.Config.Project.Name,
		constant.EnvKeyRunEnvironment: g.Config.Project.Env,
	} {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	// 尝试加载.env
	_ = godotenv.Load()
	return nil
}

func GetHdEnv() (string, error) {
	env, exists := os.LookupEnv(constant.EnvKeyRunEnvironment)
	if !exists {
		return "", fmt.Errorf("env not found, env: %s", constant.EnvKeyRunEnvironment)
	}

	if !pie.Contains(supportedEnvs, env) {
		return "", fmt.Errorf("unsupported env, env: %s", env)
	}
	return env, nil
}

func GetHdNamespace() (string, bool) {
	return os.LookupEnv(constant.EnvKeyNamespace)
}

func GetHdWorkDir() (string, bool) {
	return os.LookupEnv(constant.EnvKeyWorkDir)
}

// WithHdWorkDir 当前环境变量加上HD_WORK_DIR
func WithHdWorkDir(workDir string) []string {
	return append(os.Environ(), []string{
		fmt.Sprintf("%s=%s", constant.EnvKeyWorkDir, workDir),
	}...)
}

func GetGitCredential() (string, string) {
	gitUser, _ := os.LookupEnv(constant.EnvKeyGitUser)
	gitPassword, _ := os.LookupEnv(constant.EnvKeyGitPassword)
	return gitUser, gitPassword
}

func SetGitCredential(username, password string) error {
	return save(map[string]string{
		constant.EnvKeyGitUser:     username,
		constant.EnvKeyGitPassword: password,
	})
}

func save(data map[string]string) error {
	/// 读取现有内容（如果文件存在）
	existing, err := godotenv.Read(filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// 合并新旧值
	for k, v := range data {
		existing[k] = v
	}

	// 写入文件
	return godotenv.Write(existing, filename)
}
