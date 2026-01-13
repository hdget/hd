package appctl

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hdget/hd/pkg/env"
	"github.com/hdget/hd/pkg/utils"
	"github.com/pkg/errors"
)

type appInstallerImpl struct {
	*appCtlImpl
}

func newAppInstaller(appCtl *appCtlImpl) *appInstallerImpl {
	return &appInstallerImpl{
		appCtlImpl: appCtl,
	}
}

func (impl *appInstallerImpl) install(app, ref string) error {
	env, err := env.GetHdEnv()
	if err != nil {
		return err
	}

	configRepo, err := impl.getRepositoryConfig(repoConfig)
	if err != nil {
		return errors.Wrapf(err, "repository not found, name: %s", repoConfig)
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp(os.TempDir(), "hd-install-*")
	if err != nil {
		return errors.Wrap(err, "create temporary install dir")
	}
	//defer func() {
	//	if e := os.RemoveAll(tempDir); e != nil {
	//		fmt.Printf("delete temporary install dir: %v, dir: %s", e, tempDir)
	//	}
	//}()

	if err = newGit(impl.appCtlImpl).Clone(configRepo.Url, tempDir).Switch(ref, "main"); err != nil {
		return err
	}

	for srcPath, destPath := range map[string]string{
		filepath.Join(tempDir, "app", app, fmt.Sprintf("%s.%s.toml", app, env)): filepath.Join(impl.baseDir, "config", "app", app),
		filepath.Join(tempDir, "dapr", env, "*"):                                filepath.Join(impl.baseDir, "config", "dapr"),
	} {
		if err = utils.CopyWithWildcard(srcPath, destPath); err != nil {
			return errors.Wrapf(err, "copy dir, src: %s, dest: %s", srcPath, destPath)
		}
	}
	return nil
}
