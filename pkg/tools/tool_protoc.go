package tools

import (
	"fmt"
	"github.com/bitfield/script"
)

type protocTool struct {
	*toolImpl
}

const (
	urlProtoc = "https://github.com/protocolbuffers/protobuf/releases/download/v%s/%s"
)

func Protoc() Tool {
	version := "30.2"
	winFile := fmt.Sprintf("protoc-%s-win64.zip", version)
	linuxFile := fmt.Sprintf("protoc-%s-linux-x86_64.zip", version)

	return &protocTool{
		toolImpl: newTool("protoc", version, fmt.Sprintf(urlProtoc, version, linuxFile), fmt.Sprintf(urlProtoc, version, winFile)),
	}
}

func (t *protocTool) IsInstalled() bool {
	_, err := script.Exec(`protoc --version`).String()
	return err == nil
}

func (t *protocTool) LinuxInstall() error {
	_, zipFile, err := AllPlatform().Download(t.urlLinuxRelease)
	if err != nil {
		return err
	}
	//defer func() {
	//	if e := os.RemoveAll(tempDir); e != nil {
	//		fmt.Printf("delete temp dir failed: %v, dir: %s", e, tempDir)
	//	}
	//}()

	// 获取GOPATH
	installDir, err := AllPlatform().GetGoBinDir()
	if err != nil {
		return err
	}

	// 解压zip文件
	if err = AllPlatform().UnzipMatchedFiles(zipFile, "bin/protoc", installDir); err != nil {
		return err
	}

	return nil
}

func (t *protocTool) WindowsInstall() error {
	_, zipFile, err := AllPlatform().Download(t.urlWinRelease)
	if err != nil {
		return err
	}
	//defer func() {
	//	if e := os.RemoveAll(tempDir); e != nil {
	//		fmt.Printf("delete temp dir failed: %v, dir: %s", e, tempDir)
	//	}
	//}()

	// 获取GOPATH
	installDir, err := AllPlatform().GetGoBinDir()
	if err != nil {
		return err
	}

	// 解压zip文件
	if err = AllPlatform().UnzipMatchedFiles(zipFile, "bin/protoc.exe", installDir); err != nil {
		return err
	}

	return nil
}
