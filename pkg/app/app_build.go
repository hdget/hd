package app

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

func (a appControlImpl) Build(refName string, apps ...string) error {
	// 检查依赖的工具是否安装
	if err := tools.Check(a.debug, tools.Protoc(), tools.ProtocGogoFaster()); err != nil {
		return err
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp(os.TempDir(), "hd-build-*")
	if err != nil {
		return errors.Wrap(err, "创建Build临时目录失败")
	}
	fmt.Println("临时目录：", tempDir)
	// defer os.Remove(tempDir)

	for _, app := range apps {
		err = a.buildApp(tempDir, app, refName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a appControlImpl) buildApp(tempDir, app, refName string) error {
	appConfig, exist := appName2appConfig[app]
	if !exist {
		return fmt.Errorf("app config not found, app: %s", app)
	}

	// 创建工作目录
	appSrcDir := filepath.Join(tempDir, app)

	// 拷贝源代码并切换到指定分支
	if err := newGit().Clone(appConfig.Repo, appSrcDir).Switch(refName); err != nil {
		return err
	}

	// 编译Protobuf
	if err := a.generateProtobuf(appSrcDir, refName); err != nil {
		return err
	}

	// 拷贝sqlboiler.toml
	if err := a.copySqlboilerFile(app, appSrcDir, refName); err != nil {
		return err
	}

	return nil
}

func (a appControlImpl) copySqlboilerFile(app, appSrcDir, refName string) error {
	destConfigDir := filepath.Join(a.baseDir, "config")
	if err := newGit().Clone(g.Config.ConfigRepo, destConfigDir).Switch(refName, "main"); err != nil {
		return err
	}

	// 拷贝sqlboiler.toml
	sourceSqlboilerFile := filepath.Join(destConfigDir, "app", app, "sqlboiler.toml")
	destSqlboilerFile := filepath.Join(appSrcDir, "sqlboiler.toml")
	if _, err := script.File(sourceSqlboilerFile).WriteFile(destSqlboilerFile); err != nil {
		return err
	}
	return nil
}

func (a appControlImpl) generateProtobuf(srcDir, refName string) error {
	protoRepository := filepath.Join(srcDir, "proto")

	// 拷贝protod repostory
	if err := newGit().Clone(g.Config.ProtoRepo, protoRepository).Switch(refName, "main"); err != nil {
		return err
	}

	// 切换到app源代码目录
	err := os.Chdir(srcDir)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Chdir(a.baseDir)
	}()

	rootGolangModule, err := utils.GetRootGolangModule()
	if err != nil {
		return err
	}

	var prOptions []protorefine.Option
	if a.debug {
		prOptions = append(prOptions, protorefine.WithDebug(true))
	}

	protoDir, err := protorefine.New(prOptions...).Refine(protorefine.Argument{
		GolangModule:        rootGolangModule,
		GolangSourceCodeDir: srcDir,
		ProtoRepository:     protoRepository,
		OutputPackage:       a.pbOutputPackage,
		OutputDir:           a.pbOutputDir,
	})
	if err != nil {
		return err
	}

	// 第二步：编译protobuf
	var pcOptions []protocompile.Option
	if a.debug {
		pcOptions = append(pcOptions, protocompile.WithDebug(true))
	}
	err = protocompile.New(pcOptions...).Compile(protoDir, filepath.Join(srcDir, a.pbOutputDir, a.pbOutputPackage))
	if err != nil {
		return err
	}

	return nil
}
