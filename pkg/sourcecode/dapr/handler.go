package dapr

import (
	"github.com/hdget/hd/pkg/sourcecode"
	"github.com/hdget/hd/pkg/utils"
	"github.com/pkg/errors"
)

type daprSourceCodeHandler struct {
	*sourcecode.HandlerImpl
}

// New 初始化源代码管理器
func New(srcDir string, options ...sourcecode.Option) sourcecode.Handler {
	h := &daprSourceCodeHandler{
		HandlerImpl: &sourcecode.HandlerImpl{
			SrcDir: srcDir,
		},
	}

	for _, apply := range options {
		apply(h.HandlerImpl)
	}

	return h
}

func (h *daprSourceCodeHandler) Handle() error {
	if err := utils.IsDirWritable(h.AssetsPath); err != nil {
		return errors.Wrapf(err, "assets path is not writable, assetsPath: %s", h.AssetsPath)
	}

	// 第一步：先解析源代码数据
	parser, err := newParser(h.SrcDir, h.SkipDirs)
	if err != nil {
		return err
	}

	scInfo, err := parser.parse()
	if err != nil {
		return err
	}

	// 第二步：根据解析后的元数据，给源代码打补丁，保证服务启动时Dapr模块能自动注册
	err = h.patch(scInfo)
	if err != nil {
		return err
	}

	// 第三步：生成路由json文件
	err = h.generate(scInfo)
	if err != nil {
		return err
	}

	return nil
}
