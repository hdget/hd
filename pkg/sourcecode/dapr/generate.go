package dapr

import (
	"encoding/json"
	"fmt"

	"os"
	"path/filepath"

	"github.com/elliotchance/pie/v2"
	"github.com/hdget/common/protobuf"
)

const (
	fileExposedHandlers = ".exposed_handlers.json"
)

// generate 将有annotation的dapr invocation handlers保存到文件中
func (h *daprSourceCodeHandler) generate(scInfo *souceCodeInfo) error {
	fmt.Println("===> generate invocation handler info")

	// 只保留有annotations的，当前只有@hd.route
	handlersWithAnnotations := pie.Filter(scInfo.invocationHandlers, func(handler *protobuf.DaprHandler) bool {
		return len(handler.Annotations) > 0
	})

	data, err := json.Marshal(handlersWithAnnotations)
	if err != nil {
		return err
	}

	ext := filepath.Ext(fileExposedHandlers)
	if len(ext) > 1 {
		ext = ext[1:]
	}

	outputPath := filepath.Join(h.AssetsPath, ext, fileExposedHandlers)
	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
