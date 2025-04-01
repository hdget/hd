package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// IsDir 检查路径是否为目录
func IsDir(path string) (string, error) {
	// 获取绝对路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", errors.Wrapf(err, "get absolute path, path: %s", path)
	}

	// 检查路径是否存在且是目录
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("path doesn't exist, path: %s", absPath)
		}
		return "", errors.Wrapf(err, "check path, path: %s", absPath)
	}

	if !fileInfo.IsDir() {
		return "", fmt.Errorf("path is not directory, path: %s", path)
	}
	return absPath, nil
}

// IsDirReadableAndWithFiles 检查目录是否可读
func IsDirReadableAndWithFiles(path, fileSuffix string) error {
	dirPath, err := IsDir(path)
	if err != nil {
		return err
	}

	// 尝试打开目录读取内容
	file, err := os.Open(dirPath)
	if err != nil {
		return errors.Wrapf(err, "path is not readable, path: %s", path)
	}
	defer file.Close()

	var hasFile bool
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), fileSuffix) {
			hasFile = true
			// 找到proto文件后可以提前终止遍历
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "traverse path, path: %s", path)
	}

	if !hasFile {
		return fmt.Errorf("path doesn't have proto file, path: %s", path)
	}

	return nil

}

// IsDirWritable 检查目录是否可写
func IsDirWritable(path string) error {
	dirPath, err := IsDir(path)
	if err != nil {
		return err
	}

	// 尝试创建临时文件测试写入权限
	tempFile, err := os.CreateTemp(dirPath, ".tmp*")
	if err != nil {
		return errors.Wrapf(err, "path is not writable, path: %s", path)
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	// 尝试写入内容
	if _, err = tempFile.WriteString("test"); err != nil {
		return errors.Wrapf(err, "path is not writable, path: %s", path)
	}
	return nil
}

// FindDirContainingFiles 查找包含指定文件的目录, 返回目录名
func FindDirContainingFiles(srcDir, skipDir string, files ...string) (string, error) {
	var foundDir string
	for currentDir := srcDir; !isRoot(currentDir); currentDir = filepath.Dir(currentDir) {
		err := filepath.WalkDir(currentDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil // 跳过无法访问的目录
			}

			if d.IsDir() {
				// 忽略.开头的目录和输出目录
				if strings.HasPrefix(d.Name(), ".") || path == skipDir {
					return fs.SkipDir
				}

				if getDirDepth(currentDir, path) > 2 {
					return fs.SkipDir
				}

				// 检查当前目录是否包含所有指定文件
				if containsAllFiles(path, files) {
					foundDir = path
					return fs.SkipAll
				}
			}
			return nil
		})
		if err != nil {
			return "", err
		}

		// 扎到了不需要递归向上查找
		if foundDir != "" {
			break
		}
	}

	if foundDir == "" {
		return "", fmt.Errorf("no matched dir found, files: %+v", files)
	}
	return foundDir, nil
}

func getDirDepth(baseDir, path string) int {
	// 计算当前深度
	rel, err := filepath.Rel(baseDir, path)
	if err != nil {
		return 0
	}
	return len(strings.Split(filepath.Clean(rel), string(filepath.Separator)))
}

// containsAllFiles 并发检查目录下是否包含所有指定文件
func containsAllFiles(dir string, files []string) bool {
	var wg sync.WaitGroup
	resultChan := make(chan bool, len(files))

	for _, file := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			_, err := os.Stat(filepath.Join(dir, f))
			resultChan <- (err == nil)
		}(file)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for exists := range resultChan {
		if !exists {
			return false
		}
	}
	return true
}

func isRoot(path string) bool {
	// 处理Windows系统
	if runtime.GOOS == "windows" {
		// Windows根目录形式如 C:\ 或 \\
		volumeName := filepath.VolumeName(path)
		if volumeName != "" {
			// 去掉卷名后的路径应该是空或只有分隔符
			rest := path[len(volumeName):]
			return rest == "" || rest == string(filepath.Separator)
		}
		return false
	}

	// Unix-like系统根目录是 "/"
	return path == "/"
}

// IsValidRelativePath 检查路径是否为有效的相对路径
// 允许当前目录(".")，但不允许上级目录("..")
func IsValidRelativePath(path string) bool {
	// 空路径无效
	if path == "" {
		return false
	}

	// 检查是否为绝对路径
	if filepath.IsAbs(path) {
		return false
	}

	// 清理路径并分割
	cleanPath := filepath.Clean(path)
	parts := strings.Split(cleanPath, string(filepath.Separator))

	// 检查每个部分
	for _, part := range parts {
		// 不允许上级目录引用
		if part == ".." {
			return false
		}
	}

	// 允许以下情况：
	// 1. 当前目录 (".")
	// 2. 不以分隔符开头 (不是绝对路径)
	// 3. 不包含上级目录引用
	return true
}
