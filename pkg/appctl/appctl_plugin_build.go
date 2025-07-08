package appctl

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/env"
	"github.com/hdget/hd/pkg/protocompile"
	"github.com/hdget/hd/pkg/protorefine"
	"github.com/hdget/hd/pkg/utils"
	"github.com/pkg/errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type pluginBuilder struct {
	*appCtlImpl
	pbOutputDir     string
	pbOutputPackage string
	pbGenGRPC       bool
}

func newPluginBuilder(appCtl *appCtlImpl, pbOutputDir, pbOutputPackage string, pbGenGRPC bool) *pluginBuilder {
	return &pluginBuilder{
		appCtlImpl:      appCtl,
		pbOutputDir:     pbOutputDir,
		pbOutputPackage: pbOutputPackage,
		pbGenGRPC:       pbGenGRPC,
	}
}

func (b *pluginBuilder) build(name, repoUrl, refName string) error {
	// 创建临时目录
	tempDir, err := os.MkdirTemp(os.TempDir(), "hd-plugin-build-*")
	if err != nil {
		return errors.Wrap(err, "create temporary plugin build dir")
	}
	defer func() {
		if e := os.RemoveAll(tempDir); e != nil {
			fmt.Printf("delete temporary plugin build dir %v, dir: %s", e, tempDir)
		}
	}()

	if g.Debug {
		fmt.Println("temporary plugin build dir：", tempDir)
	}

	// 创建工作目录
	pluginSrcDir := filepath.Join(tempDir, name)

	// 拷贝源代码并切换到指定分支并获取git信息
	if g.Debug {
		fmt.Println("===> plugin build step: clone source code")
	}
	gitOperator := newGit(b.appCtlImpl)
	err = gitOperator.Clone(repoUrl, pluginSrcDir).Switch(refName)
	if err != nil {
		return err
	}
	gitBuildInfo, err := gitOperator.GetGitInfo()
	if err != nil {
		return err
	}

	// 编译Protobuf
	if g.Debug {
		fmt.Println("===> plugin build step: generate protobuf")
	}
	if err := b.generateProtobuf(pluginSrcDir, refName); err != nil {
		return err
	}

	// go build app
	if g.Debug {
		fmt.Println("===> plugin build step: go build")
	}
	if err := b.doBuild(pluginSrcDir, name, gitBuildInfo); err != nil {
		return err
	}

	return nil
}

func (b *pluginBuilder) doBuild(srcDir, name string, gitBuildInfo *gitInfo) error {
	// 切换到app源代码目录
	err := os.Chdir(srcDir)
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
	binFile := b.getExecutable(name)
	ldflagsValue := b.getLdFlags(name, gitBuildInfo)
	cmd := fmt.Sprintf("go build -ldflags='%s' -o %s", ldflagsValue, binFile)
	output, err = script.Exec(cmd).String()
	if err != nil {
		return errors.Wrapf(err, "go build, err: %s", output)
	}

	// move binary to outputdir
	if err = os.MkdirAll(b.getPluginOutputDir(), 0755); err != nil {
		return errors.Wrapf(err, "make bin dir, binDir: %s", b.getBinOutputDir())
	}
	if err = utils.CopyFile(binFile, filepath.Join(b.getPluginOutputDir(), binFile)); err != nil {
		return err
	}
	
	return nil
}

func (b *pluginBuilder) getLdFlags(app string, info *gitInfo) string {
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

func (b *pluginBuilder) generateProtobuf(srcDir, refName string) error {
	gitProtoRepo, exists := g.RepoConfigs[gitProtoRepoName]
	if !exists {
		return fmt.Errorf("repo not found, name: %s", gitProtoRepoName)
	}

	protoRepository := filepath.Join(srcDir, "proto")

	// 拷贝proto repository
	if err := newGit(b.appCtlImpl).Clone(gitProtoRepo.Url, protoRepository).Switch(refName, "main"); err != nil {
		return err
	}

	// 切换到app源代码目录
	err := os.Chdir(srcDir)
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
		ProtoRepository:     protoRepository,
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
