package appctl

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/env"
	"github.com/hdget/hd/pkg/utils"
	"github.com/pkg/errors"
)

type appInstallerImpl struct {
	*appCtlImpl
	appConfig *g.AppConfig
}

func newAppInstaller(appCtl *appCtlImpl, appConfig *g.AppConfig) *appInstallerImpl {
	return &appInstallerImpl{
		appCtlImpl: appCtl,
		appConfig:  appConfig,
	}
}

func (impl *appInstallerImpl) install(app, ref string) error {
	env, err := env.GetHdEnv()
	if err != nil {
		return err
	}

	configRepoConf, err := impl.getRepositoryConfig(impl.appConfig.ConfigRepo)
	if err != nil {
		return errors.Wrapf(err, "repository not found, name: %s", defaultConfigRepo)
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp(os.TempDir(), "hd-install-*")
	if err != nil {
		return errors.Wrap(err, "create temporary install dir")
	}
	defer func() {
		if e := os.RemoveAll(tempDir); e != nil {
			fmt.Printf("delete temporary install dir: %v, dir: %s", e, tempDir)
		}
	}()

	if err = newGit(impl.appCtlImpl).Clone(configRepoConf.Url, tempDir).Switch(ref, "main"); err != nil {
		return err
	}

	// 拷贝应用配置
	srcFiles := filepath.Join(tempDir, "app", app, fmt.Sprintf("%s.%s.toml", app, env))
	destDir := filepath.Join(impl.baseDir, "config", "app", app)
	if err = utils.CopyWithWildcard(srcFiles, destDir); err != nil {
		return errors.Wrapf(err, "copy dir, src: %s, dest: %s", srcFiles, destDir)
	}

	// 如果存在dapr配置，则拷贝dapr配置
	if _, err = os.Stat(filepath.Join(tempDir, "dapr")); err != nil && !os.IsNotExist(err) {
		return err
	}

	srcFiles = filepath.Join(tempDir, "dapr", env, "*")
	destDir = filepath.Join(impl.baseDir, "config", "dapr")
	if err = utils.CopyWithWildcard(srcFiles, destDir); err != nil {
		return errors.Wrapf(err, "copy dir, src: %s, dest: %s", srcFiles, destDir)
	}

	return nil
}
