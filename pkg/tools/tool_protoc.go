package tools

import (
	"fmt"
	"github.com/bitfield/script"
	"os"
)

type protocTool struct {
	*toolImpl
}

const (
	urlProtoc = "https://github.com/protocolbuffers/protobuf/releases/download/v%s/%s"
)

func Protoc() Tool {
	version := "30.2"
	winfile := fmt.Sprintf("protoc-%s-win64.zip", version)
	linuxFile := fmt.Sprintf("protoc-%s-linux-x86_64.zip", version)

	return &protocTool{
		toolImpl: &toolImpl{
			name:            "protoc",
			version:         version,
			urlLinuxRelease: fmt.Sprintf(urlProtoc, version, linuxFile),
			urlWinRelease:   fmt.Sprintf(urlProtoc, version, winfile),
		},
	}
}

func (t *protocTool) IsInstalled() bool {
	_, err := script.Exec(`protoc --version`).String()
	return err == nil
}

func (t *protocTool) LinuxInstall() error {
	fmt.Printf("downloading %s...\n", t.name)

	tempDir, zipFile, err := AllPlatform().Download(t.urlLinuxRelease)
	if err != nil {
		return err
	}
	defer func() {
		if e := os.RemoveAll(tempDir); e != nil {
			fmt.Printf("delete temp dir failed: %v, dir: %s", e, tempDir)
		}
	}()

	// 获取GOPATH
	installDir, err := AllPlatform().GetGoBinDir()
	if err != nil {
		return err
	}

	// 解压zip文件
	if err = AllPlatform().UnzipSpecific(zipFile, "bin/protoc", installDir); err != nil {
		return err
	}

	return nil
}

func (t *protocTool) WindowsInstall() error {
	fmt.Printf("downloading %s...\n", t.name)

	tempDir, zipFile, err := AllPlatform().Download(t.urlWinRelease)
	if err != nil {
		return err
	}
	defer func() {
		if e := os.RemoveAll(tempDir); e != nil {
			fmt.Printf("delete temp dir failed: %v, dir: %s", e, tempDir)
		}
	}()

	// 获取GOPATH
	installDir, err := AllPlatform().GetGoBinDir()
	if err != nil {
		return err
	}

	// 解压zip文件
	if err = AllPlatform().UnzipSpecific(zipFile, "bin/protoc.exe", installDir); err != nil {
		return err
	}

	return nil
}
