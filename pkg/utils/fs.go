package utils

import (
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/pkg/errors"
	"io"
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
func FindDirContainingFiles(srcDir string, matchFiles []string, skipDirs ...string) (string, error) {
	// 记录所有已检查过的目录，防止二次检查
	visitedDirs := map[string]struct{}{}

	var foundDir string
	for currentDir := srcDir; !isRoot(currentDir); currentDir = filepath.Dir(currentDir) {
		err := filepath.WalkDir(currentDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil // 跳过无法访问的目录
			}

			if d.IsDir() {
				// 忽略.开头的目录和输出目录
				if strings.HasPrefix(d.Name(), ".") || pie.Contains(skipDirs, path) {
					return fs.SkipDir
				}

				// 过滤掉大于2级深度的目录
				if getDirDepth(currentDir, path) > 2 {
					return fs.SkipDir
				}

				// 已经访问过，跳过
				if _, visited := visitedDirs[path]; visited {
					return fs.SkipDir
				}

				// 检查当前目录是否包含所有指定文件
				if containsAllFiles(path, matchFiles) {
					foundDir = path
					return fs.SkipAll
				}

				visitedDirs[path] = struct{}{}
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
		return "", fmt.Errorf("no matched dir found, files: %+v", matchFiles)
	}
	return foundDir, nil
}

// CopyWithWildcard 支持通配符的目录复制
func CopyWithWildcard(srcPattern, dst string) error {
	// 获取匹配的所有文件和目录
	matches, err := filepath.Glob(srcPattern)
	if err != nil {
		return fmt.Errorf("通配符匹配失败: %v", err)
	}

	if len(matches) == 0 {
		return fmt.Errorf("没有找到匹配的文件或目录")
	}

	// 确保目标目录存在
	if err := os.MkdirAll(dst, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 遍历所有匹配项
	for _, src := range matches {
		// 获取文件信息
		info, err := os.Stat(src)
		if err != nil {
			return fmt.Errorf("获取文件信息失败: %v", err)
		}

		// 处理目录
		if info.IsDir() {
			// 获取目录名
			dirName := filepath.Base(src)
			// 创建目标子目录
			targetDir := filepath.Join(dst, dirName)
			if err = os.MkdirAll(targetDir, info.Mode()); err != nil {
				return fmt.Errorf("创建子目录失败: %v", err)
			}
			// 递归复制目录内容
			if err := CopyDir(src, targetDir); err != nil {
				return err
			}
		} else {
			// 处理文件
			targetFile := filepath.Join(dst, filepath.Base(src))
			if err := CopyFile(src, targetFile); err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyDir 复制整个目录
func CopyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("计算相对路径失败: %v", err)
		}

		dstPath := filepath.Join(dst, relPath)

		// 处理目录
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// 处理文件
		return CopyFile(path, dstPath)
	})
}

// CopyFile 复制单个文件
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// 复制文件权限
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}

func ExistsFile(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	// 其他类型的错误（如权限问题）也视为不存在
	return false
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
