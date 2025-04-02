package tools

import (
	"archive/zip"
	"fmt"
	"github.com/bitfield/script"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type platformAll struct {
}

func AllPlatform() *platformAll {
	return &platformAll{}
}

func (platformAll) GoInstall(pkg string) error {
	envs := append(os.Environ(), []string{
		"GOPROXY=https://goproxy.cn,direct",
	}...)

	// 捕获输出
	cmd := fmt.Sprintf("go install %s", pkg)
	output, err := script.NewPipe().WithEnv(envs).Exec(cmd).String()
	if err != nil {
		return errors.Wrapf(err, "go install, err: %s", output)
	}
	return nil
}

func (platformAll) Download(url string) (string, string, error) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp(os.TempDir(), "download-*")
	if err != nil {
		return "", "", errors.Wrap(err, "创建临时目录失败")
	}

	downloadFile := filepath.Join(tempDir, filepath.Base(url))
	_, err = script.Get(url).WriteFile(downloadFile)
	if err != nil {
		fmt.Println("Windows系统请手动下载安装protoc: https://github.com/protocolbuffers/protobuf/releases")
		return "", "", errors.Wrap(err, "下载失败")
	}

	return tempDir, downloadFile, nil
}

func (platformAll) GetGoBinDir() (string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return "", errors.New("GOPATH未设置")
	}
	return filepath.Join(gopath, "bin"), nil
}

func (platformAll) UnzipSpecific(zipFile, sourcePath, targetDir string) error {
	// 打开ZIP文件
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("打开ZIP文件失败: %v", err)
	}
	defer r.Close()

	// 标准化路径（确保使用正斜杠）
	sourcePath = filepath.ToSlash(sourcePath)

	// 判断是文件还是目录
	isFile := !strings.HasSuffix(sourcePath, "/") &&
		!strings.HasSuffix(sourcePath, "\\") &&
		filepath.Ext(sourcePath) != ""

	// 遍历ZIP文件
	found := false
	for _, f := range r.File {
		zipPath := filepath.ToSlash(f.Name)

		// 检查是否匹配指定路径
		if isFile {
			// 精确匹配文件
			if zipPath == sourcePath {
				return extractFile(f, targetDir)
			}
		} else {
			// 匹配目录下的所有文件
			if strings.HasPrefix(zipPath, sourcePath) {
				relPath := zipPath[len(sourcePath):]
				if relPath == "" {
					continue // 跳过目录本身
				}

				// 如果是目录则创建，是文件则解压
				if f.FileInfo().IsDir() {
					destPath := filepath.Join(targetDir, relPath)
					if err = os.MkdirAll(destPath, 0755); err != nil {
						return fmt.Errorf("创建目录失败: %v", err)
					}
				} else {
					if err := extractFile(f, filepath.Join(targetDir, filepath.Dir(relPath))); err != nil {
						return err
					}
				}
				found = true
			}
		}
	}

	if !found {
		if isFile {
			return fmt.Errorf("ZIP文件中未找到文件: %s", sourcePath)
		}
		return fmt.Errorf("ZIP文件中未找到目录: %s", sourcePath)
	}

	return nil
}

// 解压单个文件
func extractFile(f *zip.File, targetDir string) error {
	// 创建目标文件路径
	destPath := filepath.Join(targetDir, filepath.Base(f.Name))

	// 创建目标文件
	outFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer outFile.Close()

	// 打开源文件
	inFile, err := f.Open()
	if err != nil {
		return fmt.Errorf("打开ZIP内文件失败: %v", err)
	}
	defer inFile.Close()

	// 复制内容
	if _, err := io.Copy(outFile, inFile); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}
