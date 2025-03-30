package protogen

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type ProtobufGenerator interface {
	FindProtoDir(matchFiles []string) (string, error) // 智能查找proto files的目录
	Generate(protoDir, outputDir string) error        // 生成protobuf.pb文件
}

type protobufGenImpl struct {
	srcDir string
}

func (impl *protobufGenImpl) Generate(protoDir, outputDir string) error {
	//TODO implement me
	panic("implement me")
}

func New(srcDir string) ProtobufGenerator {
	return &protobufGenImpl{
		srcDir: srcDir,
	}
}

func (impl *protobufGenImpl) FindProtoDir(protoFiles []string) (string, error) {
	protoDir, err := impl.findDirContainingFiles(protoFiles)
	if err != nil {
		return "", err
	}
	return protoDir, nil
}

// 查找包含指定文件的目录, 返回目录名
func (impl *protobufGenImpl) findDirContainingFiles(files []string) (string, error) {
	var foundDir string
	for currentDir := impl.srcDir; !impl.isRoot(currentDir); currentDir = filepath.Dir(currentDir) {
		err := filepath.WalkDir(currentDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil // 跳过无法访问的目录
			}

			if d.IsDir() {
				// 忽略起始目录
				if currentDir == impl.srcDir {
					return fs.SkipDir
				}

				// 检查当前目录是否包含所有指定文件
				if impl.containsAllFiles(path, files) {
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
		return "", fmt.Errorf("speicifc dir not found")
	}
	return foundDir, nil
}

// containsAllFiles 并发检查目录下是否包含所有指定文件
func (*protobufGenImpl) containsAllFiles(dir string, files []string) bool {
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

func (*protobufGenImpl) isRoot(path string) bool {
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
