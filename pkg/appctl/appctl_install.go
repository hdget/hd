package appctl

import (
	"fmt"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/env"
	"github.com/hdget/hd/pkg/utils"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
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

	// 创建临时目录
	tempDir, err := os.MkdirTemp(os.TempDir(), "hd-install-*")
	if err != nil {
		return errors.Wrap(err, "创建Build临时目录失败")
	}
	defer func() {
		if e := os.RemoveAll(tempDir); e != nil {
			fmt.Printf("删除临时目录失败: %v, dir: %s", e, tempDir)
		}
	}()

	gitConfigRepo, exists := g.RepoConfigs[gitConfigRepoName]
	if !exists {
		return fmt.Errorf("repo not found, name: %s", gitConfigRepoName)
	}

	if err = newGit(impl.appCtlImpl).Clone(gitConfigRepo.Url, tempDir).Switch(ref, "main"); err != nil {
		return err
	}

	for srcPath, destPath := range map[string]string{
		filepath.Join(tempDir, "app", app, fmt.Sprintf("%s.%s.toml", app, env)): filepath.Join(impl.baseDir, "config", "app", app),
		filepath.Join(tempDir, "dapr", env, "*"):                                filepath.Join(impl.baseDir, "config", "dapr"),
	} {
		if err = utils.CopyWithWildcard(srcPath, destPath); err != nil {
			return errors.Wrapf(err, "复制目录失败, src: %s, dest: %s", srcPath, destPath)
		}
	}
	return nil
}
