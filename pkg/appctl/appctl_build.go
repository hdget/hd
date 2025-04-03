package appctl

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/hdget/hd/g"
	"github.com/hdget/hd/pkg/protocompile"
	"github.com/hdget/hd/pkg/protorefine"
	"github.com/hdget/hd/pkg/tools"
	"github.com/hdget/hd/pkg/utils"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

type appBuilder struct {
	*appCtlImpl
	pbOutputDir     string
	pbOutputPackage string
}

const (
	gitProtoRepoName  = "proto"
	gitConfigRepoName = "config"
)

func newAppBuilder(appCtl *appCtlImpl) (*appBuilder, error) {
	// 检查依赖的工具是否安装
	if err := tools.Check(appCtl.debug,
		tools.Protoc(),
		tools.ProtocGogoFaster(),
		tools.Sqlboiler(),
	); err != nil {
		return nil, err
	}

	return &appBuilder{
		appCtlImpl:      appCtl,
		pbOutputDir:     "autogen",
		pbOutputPackage: "pb",
	}, nil
}

func (b *appBuilder) build(refName string, apps ...string) error {
	// 创建临时目录
	tempDir, err := os.MkdirTemp(os.TempDir(), "hd-build-*")
	if err != nil {
		return errors.Wrap(err, "创建Build临时目录失败")
	}
	defer os.Remove(tempDir)

	if b.debug {
		fmt.Println("临时目录：", tempDir)
	}

	for _, app := range apps {
		err = b.buildApp(tempDir, app, refName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *appBuilder) buildApp(tempDir, app, refName string) error {
	appRepoConfig, exist := g.RepoConfigs[app]
	if !exist {
		return fmt.Errorf("git repository not found, app: %s", app)
	}

	// 创建工作目录
	appSrcDir := filepath.Join(tempDir, app)

	// 拷贝源代码并切换到指定分支
	if err := newGit(b.appCtlImpl).Clone(appRepoConfig.Url, appSrcDir).Switch(refName); err != nil {
		return err
	}

	// 编译Protobuf
	if err := b.generateProtobuf(appSrcDir, refName); err != nil {
		return err
	}

	// 拷贝sqlboiler.toml
	if err := b.copySqlboilerConfigFile(appSrcDir, app, refName); err != nil {
		return err
	}

	// go build
	if err := b.golangBuild(appSrcDir, app); err != nil {
		return err
	}

	return nil
}

func (b *appBuilder) golangBuild(appSrcDir, app string) error {
	// 切换到app源代码目录
	err := os.Chdir(appSrcDir)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Chdir(b.baseDir)
	}()

	envs := append(os.Environ(), []string{
		fmt.Sprintf("HD_WORK_DIR=%s", b.baseDir),
	}...)

	// go generate
	output, err := script.NewPipe().WithEnv(envs).Exec(`go generate`).String()
	if err != nil {
		return errors.Wrapf(err, "go generate, err: %s", output)
	}

	// go build
	binFile := b.getExecutable(app)
	cmd := fmt.Sprintf("go build -o %s", binFile)
	output, err = script.Exec(cmd).String()
	if err != nil {
		return errors.Wrapf(err, "go build, err: %s", output)
	}

	// move binary to binDir
	if err = os.MkdirAll(b.absBinDir, 0755); err != nil {
		return errors.Wrapf(err, "make bin dir, binDir: %s", b.absBinDir)
	}
	if _, err = script.File(binFile).WriteFile(filepath.Join(b.absBinDir, binFile)); err != nil {
		return err
	}

	return nil
}

func (b *appBuilder) copySqlboilerConfigFile(appSrcDir, app, refName string) error {
	gitConfigRepo, exists := g.RepoConfigs[gitConfigRepoName]
	if !exists {
		return fmt.Errorf("repo config not found, name: %s", gitConfigRepoName)
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp(os.TempDir(), "hd-config-*")
	if err != nil {
		return errors.Wrap(err, "创建Build临时目录失败")
	}
	defer os.Remove(tempDir)

	if err = newGit(b.appCtlImpl).Clone(gitConfigRepo.Url, tempDir).Switch(refName, "main"); err != nil {
		return err
	}

	// 拷贝sqlboiler.toml
	srcPath := filepath.Join(tempDir, "app", app, "sqlboiler.toml")
	destPath := filepath.Join(appSrcDir, "sqlboiler.toml")
	if _, err = script.File(srcPath).WriteFile(destPath); err != nil {
		return err
	}
	return nil
}

func (b *appBuilder) generateProtobuf(srcDir, refName string) error {
	gitProtoRepo, exists := g.RepoConfigs[gitProtoRepoName]
	if !exists {
		return fmt.Errorf("repo config not found, name: %s", gitProtoRepoName)
	}

	protoRepository := filepath.Join(srcDir, "proto")

	// 拷贝protod repostory
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

	var prOptions []protorefine.Option
	if b.debug {
		prOptions = append(prOptions, protorefine.WithDebug(true))
	}

	protoDir, err := protorefine.New(prOptions...).Refine(protorefine.Argument{
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
	var pcOptions []protocompile.Option
	if b.debug {
		pcOptions = append(pcOptions, protocompile.WithDebug(true))
	}
	err = protocompile.New(pcOptions...).Compile(protoDir, filepath.Join(srcDir, b.pbOutputDir, b.pbOutputPackage))
	if err != nil {
		return err
	}

	return nil
}
