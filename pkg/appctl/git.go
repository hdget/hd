package appctl

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bitfield/script"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/env"
	"github.com/hdget/hd/pkg/utils"
	"github.com/pkg/errors"
)

type gitImpl struct {
	*appCtlImpl
	repo *git.Repository
}

type gitInfo struct {
	tag    string
	branch string
	commit string
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
	if g.Debug {
		fmt.Printf("git clone, url: %s, destDir: %s\n", url, destDir)
	}

	if err := os.RemoveAll(destDir); err != nil {
		utils.Fatal("remove dest directory", err)
	}

	var err error
	impl.repo, err = git.PlainClone(destDir, true, &git.CloneOptions{
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
	if g.Debug {
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

func (impl *gitImpl) GetGitInfo() (*gitInfo, error) {
	latestTag, err := impl.getLatestTag()
	if err != nil {
		return nil, err
	}

	commit, err := impl.getHeadShortHash()
	if err != nil {
		return nil, err
	}

	branch, err := impl.getCurrentBranchName()
	if err != nil {
		return nil, err
	}

	return &gitInfo{
		branch: branch,
		commit: commit,
		tag:    latestTag,
	}, nil
}

func (impl *gitImpl) checkout(refName string) error {
	if impl.repo == nil {
		return fmt.Errorf("git repository not initialized")
	}

	w, err := impl.repo.Worktree()
	if err != nil {
		return err
	}

	// 尝试作为远程分支切换
	remoteRefName := plumbing.NewRemoteReferenceName("origin", refName)
	if remoteRef, err := impl.repo.Reference(remoteRefName, true); err == nil {
		localBranchRef := plumbing.NewBranchReferenceName(refName)
		if e := impl.repo.Storer.SetReference(plumbing.NewHashReference(localBranchRef, remoteRef.Hash())); e != nil {
			return errors.Wrap(err, "set local branch reference to remote branch")
		}
		return w.Checkout(&git.CheckoutOptions{
			Branch: localBranchRef,
			Create: false,
		})
	}

	// 尝试作为标签切换
	var foundHash plumbing.Hash
	tagRefs, err := impl.repo.Tags()
	if err != nil {
		return errors.Wrap(err, "get tag list")
	}
	err = tagRefs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().Short() == refName {
			foundHash = ref.Hash()
			return nil
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "iterate tag")
	}
	if !foundHash.IsZero() {
		return w.Checkout(&git.CheckoutOptions{
			Hash:  foundHash,
			Force: true,
		})
	}

	// 尝试作为提交哈希切换
	if hash := plumbing.NewHash(refName); !hash.IsZero() {
		if _, err := impl.repo.CommitObject(hash); err == nil {
			return w.Checkout(&git.CheckoutOptions{
				Hash:  hash,
				Keep:  false, // 不保留工作区更改
				Force: true,
			})
		}
	}

	return errRefNotFound
}

func (impl *gitImpl) getAuth() *http.BasicAuth {
	once.Do(func() {
		gitUser, gitPassword := env.GetGitCredential()
		if gitUser == "" {
			gitUser = utils.GetInput(">>> GIT用户")
		}
		if gitPassword == "" {
			gitPassword = utils.GetInput(">>> GIT密码")
		}
		_ = env.SetGitCredential(gitUser, gitPassword)

		cachedAuth = &http.BasicAuth{
			Username: gitUser,     // 对于GitHub，可以是任意非空字符串
			Password: gitPassword, // 实际使用你的访问令牌
		}
	})
	return cachedAuth
}

func (impl *gitImpl) getLatestTag() (string, error) {
	// 获取当前HEAD
	headRef, err := impl.repo.Head()
	if err != nil {
		return "", errors.Wrap(err, "get HEAD")
	}

	// 获取当前提交
	currentCommit, err := impl.repo.CommitObject(headRef.Hash())
	if err != nil {
		return "", errors.Wrap(err, "get current commit")
	}

	tags, err := impl.repo.Tags()
	if err != nil {
		return "", errors.Wrap(err, "get tags")
	}

	var latestTagName string
	var latestTagTime time.Time
	err = tags.ForEach(func(ref *plumbing.Reference) error {
		// 获取标签对应的提交
		var tagCommit *object.Commit

		// 尝试解析带注释标签
		if tagObj, err := impl.repo.TagObject(ref.Hash()); err == nil {
			tagCommit, _ = tagObj.Commit()
		} else { // 轻量级标签
			tagCommit, _ = impl.repo.CommitObject(ref.Hash())
		}

		// 检查是否在当前提交的历史中
		if impl.isAncestorCommit(currentCommit, tagCommit) {
			// 保留最新时间戳的标签
			if tagCommit.Committer.When.After(latestTagTime) {
				latestTagName = ref.Name().Short()
				latestTagTime = tagCommit.Committer.When
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return latestTagName, nil
}

// 检查是否在提交历史中
func (impl *gitImpl) isAncestorCommit(current, target *object.Commit) bool {
	iter := object.NewCommitPreorderIter(current, nil, nil)
	for {
		commit, err := iter.Next()
		if err != nil || commit == nil {
			break
		}
		if commit.Hash == target.Hash {
			return true
		}
	}
	return false
}

func (impl *gitImpl) getCurrentBranchName() (string, error) {
	ref, err := impl.repo.Head()
	if err != nil {
		return "", err
	}

	// 如果是分支引用
	if ref.Name().IsBranch() {
		return ref.Name().Short(), nil
	}

	// 如果不是分支（如分离头状态）
	return "", fmt.Errorf("not on a branch (detached HEAD)")
}

func (impl *gitImpl) getHeadShortHash() (string, error) {
	ref, err := impl.repo.Head()
	if err != nil {
		return "", err
	}

	return ref.Hash().String()[:7], nil // 取前7个字符作为短哈希
}
