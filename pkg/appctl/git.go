package appctl

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/hdget/hd/pkg/env"
	"github.com/hdget/hd/pkg/utils"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"sync"
)

type gitImpl struct {
	*appCtlImpl
	repo *git.Repository
}

var (
	errRefNotFound = errors.New("ref not found")
	once           sync.Once
	cachedAuth     *http.BasicAuth
)

func newGit(appCtl *appCtlImpl) *gitImpl {
	// _ = script.Exec(`git config --global credential.helper store`).Wait()
	_ = script.Exec(`git config --global advice.detachedHead false`).Wait()
	return &gitImpl{
		appCtlImpl: appCtl,
	}
}

func (impl *gitImpl) Clone(url, destDir string) *gitImpl {
	if impl.debug {
		fmt.Printf("git clone, url: %s, destDir: %s\n", url, destDir)
	}

	if err := os.RemoveAll(destDir); err != nil {
		utils.Fatal("remove dest directory", err)
	}

	var err error
	impl.repo, err = git.PlainClone(destDir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
		Auth:     impl.getAuth(),
	})
	if err != nil {
		utils.Fatal("clone repository", err)
	}
	return impl
}

func (impl *gitImpl) Switch(refName string, fallbackRefName ...string) error {
	if impl.debug {
		fmt.Printf("git switch, ref: %s, fallback: %s\n", refName, fallbackRefName)
	}

	err := impl.checkout(refName)
	if err != nil {
		if errors.Is(err, errRefNotFound) && len(fallbackRefName) > 0 {
			return impl.checkout(fallbackRefName[0])
		}
		return err
	}
	return nil
}

func (impl *gitImpl) checkout(refName string) error {
	if impl.repo == nil {
		return fmt.Errorf("git仓库未初始化")
	}

	w, err := impl.repo.Worktree()
	if err != nil {
		return err
	}

	// 尝试作为分支切换
	branchRef := plumbing.ReferenceName("refs/heads/" + refName)
	if _, err := impl.repo.Reference(branchRef, true); err == nil {
		return w.Checkout(&git.CheckoutOptions{
			Branch: branchRef,
			Force:  true,
		})
	}

	// 尝试作为标签切换
	tagRef := plumbing.ReferenceName("refs/tags/" + refName)
	if ref, err := impl.repo.Reference(tagRef, true); err == nil {
		return w.Checkout(&git.CheckoutOptions{
			Hash:  ref.Hash(),
			Force: true,
		})
	}

	// 尝试作为提交哈希切换
	if hash := plumbing.NewHash(refName); !hash.IsZero() {
		if _, err := impl.repo.CommitObject(hash); err == nil {
			return w.Checkout(&git.CheckoutOptions{
				Hash:  hash,
				Force: true,
			})
		}
	}

	return errRefNotFound
}

func (impl *gitImpl) getAuth() *http.BasicAuth {
	once.Do(func() {
		gitUser, _ := os.LookupEnv("GIT_USER")
		if gitUser == "" {
			gitUser = utils.GetInput(">>> GIT用户: ")
		}

		gitPassword, _ := os.LookupEnv("GIT_PASSWORD")
		if gitPassword == "" {
			gitPassword = utils.GetInput(">>> GIT密码: ")
		}

		_ = env.WriteEnvFile(filepath.Join(impl.baseDir, ".env"), map[string]string{
			"GIT_USER":     gitUser,
			"GIT_PASSWORD": gitPassword,
		})

		cachedAuth = &http.BasicAuth{
			Username: gitUser,     // 对于GitHub，可以是任意非空字符串
			Password: gitPassword, // 实际使用你的访问令牌
		}
	})
	return cachedAuth
}
