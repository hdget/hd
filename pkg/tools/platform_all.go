package tools

import (
	"archive/zip"
	"fmt"
	"github.com/bitfield/script"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/hd/pkg/utils"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
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
	tempDir, err := os.MkdirTemp(os.TempDir(), "hd-download-*")
	if err != nil {
		return "", "", errors.Wrap(err, "create temp dir")
	}

	// 获取文件大小
	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(5)) // 跟随最多5次重定向

	resp, err := client.R().Head(url)
	if err != nil {
		return "", "", errors.New("get file size failed")
	}
	contentLength, _ := strconv.ParseInt(resp.Header().Get("Content-Length"), 10, 64)

	// 创建进度条
	downloadDesc := fmt.Sprintf("downloading: %s\n", url)
	var bar *progressbar.ProgressBar
	if contentLength > 0 {
		bar = progressbar.NewOptions64(
			contentLength,
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionSetDescription(downloadDesc),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetWidth(50),
			progressbar.OptionThrottle(100*time.Millisecond),
			progressbar.OptionShowCount(),
			progressbar.OptionOnCompletion(func() {
				_, _ = fmt.Fprint(os.Stderr, "\n")
			}),
		)
	} else {
		// 未知大小的进度条
		bar = progressbar.NewOptions64(
			-1,
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionSetDescription(downloadDesc),
			progressbar.OptionShowBytes(true),
			progressbar.OptionShowCount(),
			progressbar.OptionOnCompletion(func() {
				_, _ = fmt.Fprint(os.Stderr, "\n")
			}),
		)
	}

	// 创建输出文件
	outputPath := filepath.Join(tempDir, filepath.Base(url))
	outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", "", err
	}
	defer func() {
		if e := outputFile.Close(); e != nil {
			fmt.Println(e)
		}
	}()

	// 执行下载并显示进度
	resp, err = client.R().
		SetDoNotParseResponse(true).
		Get(url)
	if err != nil {
		return "", "", err
	}
	defer func() {
		if e := resp.RawBody().Close(); e != nil {
			fmt.Println(e)
		}
	}()

	_, err = io.Copy(io.MultiWriter(outputFile, bar), resp.RawBody())
	if err != nil {
		return "", "", err
	}

	//// 确保所有数据写入磁盘
	//if err = outputFile.Sync(); err != nil {
	//	return "", "", err
	//}

	if err = bar.Finish(); err != nil {
		return "", "", err
	}

	//_, err = script.Get(url).WriteFile(downloadFile)
	//if err != nil {
	//	return "", "", errors.Wrapf(err, "download failed, file: %s", downloadFile)
	//}
	return tempDir, outputPath, nil
}

func (platformAll) GetGoBinDir() (string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
		if gopath == "" {
			if runtime.GOOS == "windows" {
				gopath = filepath.Join(os.Getenv("USERPROFILE"), "go")
			} else {
				gopath = filepath.Join(os.Getenv("HOME"), "go")
			}
		}
	}

	if gopath == "" {
		return "", errors.New("GOPATH not found")
	}
	return filepath.Join(gopath, "bin"), nil
}

func (platformAll) UnzipMatchedFiles(zipFile, matchPattern, destDir string) error {
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
			// 如果是目录，什么都不做
			if f.FileInfo().IsDir() {
				continue
			}

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
	destPath := filepath.Join(destDir, filepath.Base(f.Name))

	// 2. 确保destDir存在
	if _, err := utils.IsDir(destDir); err != nil {
		return fmt.Errorf("invalid dest dir, dir: %s", destDir)
	}

	// 3. 打开ZIP中的文件
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
