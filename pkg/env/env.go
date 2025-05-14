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
	envGitUser     = "HD_GIT_USER"
	envGitPassword = "HD_GIT_PASSWORD"
	envHdWorkDir   = "HD_WORK_DIR"
	filename       = ".env"
)

var (
	supportedEnvs = []string{"prod", "test", "local"}
	exportedEnvs  = map[string]string{
		envHdNamespace: g.Config.Project.Name,
		envHdEnv:       g.Config.Project.Env,
	}
)

func Initialize() error {
	for k, v := range exportedEnvs {
		fmt.Println("set env, k:", k, "v:", v)
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	return godotenv.Load()
	//err := save(map[string]string{
	//	envHdNamespace: g.Config.Project.Name,
	//	envHdEnv:       g.Config.Project.Env,
	//})
	//if err != nil {
	//	return err
	//}
	//return godotenv.Load()
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

func GetHdWorkDir() (string, bool) {
	return os.LookupEnv(envHdWorkDir)
}

// WithHdWorkDir 当前环境变量加上HD_WORK_DIR
func WithHdWorkDir(workDir string) []string {
	return append(os.Environ(), []string{
		fmt.Sprintf("HD_WORK_DIR=%s", workDir),
	}...)
}

func GetGitCredential() (string, string) {
	gitUser, _ := os.LookupEnv(envGitUser)
	gitPassword, _ := os.LookupEnv(envGitPassword)
	return gitUser, gitPassword
}

func SetGitCredential(username, password string) error {
	return save(map[string]string{
		envGitUser:     username,
		envGitPassword: password,
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
