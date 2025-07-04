package sourcecode

import (
	"github.com/hdget/hd/pkg/utils"
	"github.com/pkg/errors"
)

type SourceCodeHandler interface {
	Handle() error
}

type sourceCodeHandlerImpl struct {
	srcDir     string
	skipDirs   []string
	assetsPath string
}

// New 初始化源代码管理器
func New(srcDir string, options ...Option) SourceCodeHandler {
	m := &sourceCodeHandlerImpl{
		srcDir: srcDir,
	}

	for _, apply := range options {
		apply(m)
	}

	return m
}

func (impl *sourceCodeHandlerImpl) Handle() error {
	if err := utils.IsDirWritable(impl.assetsPath); err != nil {
		return errors.Wrapf(err, "assets path is not writable, assetsPath: %s", impl.assetsPath)
	}

	// 第一步：先解析源代码数据
	parser, err := newParser(impl.srcDir, impl.skipDirs)
	if err != nil {
		return err
	}

	scInfo, err := parser.Parse()
	if err != nil {
		return err
	}

	// 第二步：根据解析后的元数据，给源代码打补丁，保证服务启动时Dapr模块能自动注册
	err = impl.patch(scInfo)
	if err != nil {
		return err
	}

	// 第三步：生成路由json文件
	err = impl.generate(scInfo)
	if err != nil {
		return err
	}

	return nil
}
