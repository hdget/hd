package sourcecode

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	fileInvocationHandlers = ".invocation_handlers.json"
)

// generate 生成相关文件括route.json
func (impl *sourceCodeHandlerImpl) generate(scInfo *sourceCodeInfo) error {
	fmt.Println("===> generate invocation handler info")

	data, err := json.Marshal(scInfo.daprInvocationHandlers)
	if err != nil {
		return err
	}

	ext := filepath.Ext(fileInvocationHandlers)
	if len(ext) > 1 {
		ext = ext[1:]
	}

	outputPath := filepath.Join(impl.assetsPath, ext, fileInvocationHandlers)
	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
