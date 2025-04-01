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
	urlProtocWindowsBinary = "https://github.com/protocolbuffers/protobuf/releases/download/v%s/protoc-%s-win64.zip"
)

func Protoc() Tool {
	return &protocTool{
		toolImpl: &toolImpl{
			name:    "protoc",
			version: "30.2",
		},
	}
}

func (t *protocTool) IsInstalled() bool {
	_, err := script.Exec(`protoc --version`).String()
	return err == nil
}

func (t *protocTool) LinuxInstall() error {
	cmd := `(curl -L https://github.com/protocolbuffers/protobuf/releases/latest/download/protoc-*.zip -o protoc.zip && \
                     unzip -o protoc.zip -d /usr/local && \
                     rm protoc.zip) || (echo "安装失败，请手动下载: https://github.com/protocolbuffers/protobuf/releases" && exit 1)`

	output, err := script.Exec(cmd).String()
	if err != nil {
		return fmt.Errorf("protoc安装失败: %v\n输出: %s", err, output)
	}

	return nil
}

func (t *protocTool) WindowsInstall() error {
	fmt.Printf("正在下载protoc v%s...\n", t.version)

	url := fmt.Sprintf(urlProtocWindowsBinary, t.version, t.version)
	tempDir, zipFile, err := NoArch().DownloadFile(url, "protoc.zip")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// 获取GOPATH
	installDir, err := NoArch().GetGoBinDir()
	if err != nil {
		return err
	}

	// 解压zip文件
	if err = NoArch().UnzipSpecific(zipFile, "bin/protoc.exe", installDir); err != nil {
		return err
	}

	return nil
}
