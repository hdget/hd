package appctl

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitfield/script"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/env"
	"github.com/hdget/hd/pkg/protocompile"
	"github.com/hdget/hd/pkg/protorefine"
	"github.com/hdget/hd/pkg/utils"
	"github.com/pkg/errors"
)

type appBuilder struct {
	*appCtlImpl
	pbOutputDir     string
	pbOutputPackage string
	pbGenGRPC       bool
}

const (
	repoProto  = "proto"
	repoConfig = "config"
)

func newAppBuilder(appCtl *appCtlImpl, pbOutputDir, pbOutputPackage string, pbGenGRPC bool) *appBuilder {
	return &appBuilder{
		appCtlImpl:      appCtl,
		pbOutputDir:     pbOutputDir,
		pbOutputPackage: pbOutputPackage,
		pbGenGRPC:       pbGenGRPC,
	}
}

func (b *appBuilder) build(name, refName string) error {
	app, err := newApp(name)
	if err != nil {
		return errors.Wrapf(err, "new app, app: %s", name)
	}

	appRepoUrl, err := app.GetRepoUrl()
	if err != nil {
		return errors.Wrapf(err, "get app git repository url, app: %s", name)
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp(os.TempDir(), "hd-build-*")
	if err != nil {
		return errors.Wrap(err, "create temporary build dir")
	}
	//defer func() {
	//	if e := os.RemoveAll(tempDir); e != nil {
	//		fmt.Printf("delete temporary build dir %v, dir: %s", e, tempDir)
	//	}
	//}()

	if g.Debug {
		fmt.Println("temporary build dir：", tempDir)
	}

	// 创建工作目录
	appSrcDir := filepath.Join(tempDir, name)

	// 拷贝源代码并切换到指定分支并获取git信息
	if g.Debug {
		fmt.Println("===> build step: clone source code")
	}
	gitOperator := newGit(b.appCtlImpl)
	err = gitOperator.Clone(appRepoUrl, appSrcDir).Switch(refName)
	if err != nil {
		return err
	}
	gitBuildInfo, err := gitOperator.GetGitInfo()
	if err != nil {
		return err
	}

	// 编译Protobuf
	if g.Debug {
		fmt.Println("===> build step: generate protobuf")
	}
	if err := b.generateProtobuf(appSrcDir, refName); err != nil {
		return err
	}

	// 拷贝sqlboiler.toml
	if g.Debug {
		fmt.Println("===> build step: copy sqlboiler config file")
	}
	if err = b.copySqlboilerConfigFile(appSrcDir, name, refName); err != nil {
		return err
	}

	// go build name
	if g.Debug {
		fmt.Println("===> build step: go build")
	}
	if err := b.golangAppBuild(appSrcDir, name, gitBuildInfo); err != nil {
		return err
	}

	return nil
}

func (b *appBuilder) getBuildLdflags(app string, info *gitInfo) string {
	impDir := path.Join(app, "cmd")
	ldflags := []string{
		"-w -s",
		fmt.Sprintf("-X %s.gitTag=%s", impDir, info.tag),
		fmt.Sprintf("-X %s.gitCommit=%s", impDir, info.commit),
		fmt.Sprintf("-X %s.gitBranch=%s", impDir, info.branch),
		fmt.Sprintf("-X %s.buildDate=%s", impDir, time.Now().Format("2006-01-02T15:04:05-0700")),
	}
	return strings.Join(ldflags, " ")
}

func (b *appBuilder) golangAppBuild(appSrcDir, app string, gitBuildInfo *gitInfo) error {
	// 切换到app源代码目录
	err := os.Chdir(appSrcDir)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Chdir(b.baseDir)
	}()

	// go generate
	output, err := script.NewPipe().WithEnv(env.WithHdWorkDir(b.baseDir)).Exec(`go generate`).String()
	if err != nil {
		return errors.Wrapf(err, "go generate, err: %s", output)
	}

	// go build
	binFile := b.getExecutable(app)
	ldflagsValue := b.getBuildLdflags(app, gitBuildInfo)
	cmd := fmt.Sprintf("go build -ldflags='%s' -o %s", ldflagsValue, binFile)
	output, err = script.Exec(cmd).String()
	if err != nil {
		return errors.Wrapf(err, "go build, err: %s", output)
	}

	// move binary to binDir
	if err = os.MkdirAll(b.getBinOutputDir(), 0755); err != nil {
		return errors.Wrapf(err, "make bin dir, binDir: %s", b.getBinOutputDir())
	}
	if err = utils.CopyFile(binFile, filepath.Join(b.getBinOutputDir(), binFile)); err != nil {
		return err
	}

	return nil
}

func (b *appBuilder) copySqlboilerConfigFile(appSrcDir, app, refName string) error {
	// 创建临时目录
	tempDir, err := os.MkdirTemp(os.TempDir(), "hd-config-*")
	if err != nil {
		return errors.Wrap(err, "create temporary config dir")
	}
	//defer func() {
	//	if e := os.RemoveAll(tempDir); e != nil {
	//		fmt.Printf("delete temporary config dir: %v", e)
	//	}
	//}()

	// clone config repo
	configRepo, err := b.getRepositoryConfig(repoConfig)
	if err != nil {
		return errors.Wrapf(err, "repository not found, name: %s", repoConfig)
	}

	if err = newGit(b.appCtlImpl).Clone(configRepo.Url, tempDir).Switch(refName, "main"); err != nil {
		return err
	}

	// 拷贝sqlboiler.toml
	srcPath := filepath.Join(tempDir, "app", app, "sqlboiler.toml")
	destPath := filepath.Join(appSrcDir, "sqlboiler.toml")
	if err = utils.CopyFile(srcPath, destPath); err != nil {
		return err
	}
	return nil
}

func (b *appBuilder) generateProtobuf(srcDir, refName string) error {
	protoRepo, err := b.getRepositoryConfig(repoProto)
	if err != nil {
		return errors.Wrapf(err, "repository not found, name: %s", repoProto)
	}

	protoOutputDir := filepath.Join(srcDir, "proto")

	// 拷贝protod repostory
	if err = newGit(b.appCtlImpl).Clone(protoRepo.Url, protoOutputDir).Switch(refName, "main"); err != nil {
		return err
	}

	// 切换到app源代码目录
	err = os.Chdir(srcDir)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Chdir(b.baseDir)
	}()

	rootGolangModule, err := utils.GetRootGolangModule()
	if err != nil {
		return err
	}

	protoDir, err := protorefine.New().Refine(protorefine.Argument{
		GolangModule:        rootGolangModule,
		GolangSourceCodeDir: srcDir,
		ProtoRepository:     protoOutputDir,
		OutputPackage:       b.pbOutputPackage,
		OutputDir:           b.pbOutputDir,
	})
	if err != nil {
		return err
	}

	// 第二步：编译protobuf
	err = protocompile.New(protocompile.WithGRPC(b.pbGenGRPC)).Compile(protoDir, filepath.Join(srcDir, b.pbOutputDir), b.pbOutputPackage)
	if err != nil {
		return err
	}

	return nil
}
