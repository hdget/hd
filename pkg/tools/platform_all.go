package tools

import (
	"archive/zip"
	"fmt"
	"github.com/bitfield/script"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
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
		return "", "", errors.Wrap(err, "create temp dir")
	}
	defer func() {
		if e := os.RemoveAll(tempDir); e != nil {
			fmt.Printf("delete temp dir, err: %v, dir: %s", e, tempDir)
		}
	}()

	downloadFile := filepath.Join(tempDir, filepath.Base(url))
	_, err = script.Get(url).WriteFile(downloadFile)
	if err != nil {
		return "", "", errors.Wrapf(err, "download failed, file: %s", downloadFile)
	}

	return tempDir, downloadFile, nil
}

func (platformAll) GetGoBinDir() (string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return "", errors.New("GOPATH not set")
	}
	return filepath.Join(gopath, "bin"), nil
}

func (platformAll) UnzipSpecific(zipFile, matchPattern, destDir string) error {
	// 打开ZIP文件
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("open zip file failed: %v", err)
	}
	defer r.Close()

	// 遍历ZIP文件
	found := false
	for _, f := range r.File {
		matched, err := filepath.Match(matchPattern, f.Name)
		if err != nil {
			return err
		}

		// 检查是否匹配指定路径
		if matched {
			// 4. 处理匹配的文件
			if err = extractFile(f, destDir); err != nil {
				return fmt.Errorf("uncompress file failed, file: %s, err: %v", f.Name, err)
			}
			found = true
		}
	}

	if !found {
		return fmt.Errorf("no matched zip file: %s", matchPattern)
	}

	return nil
}

// extractFile 解压单个文件
func extractFile(f *zip.File, destDir string) error {
	// 1. 创建目标文件路径
	destPath := filepath.Join(destDir, f.Name)

	// 2. 检查目录是否存在，不存在则创建
	if f.FileInfo().IsDir() {
		return os.MkdirAll(destPath, f.Mode())
	}

	// 3. 确保父目录存在
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	// 4. 打开ZIP中的文件
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// 5. 创建目标文件
	outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 6. 复制文件内容
	_, err = io.Copy(outFile, rc)
	return err
}
