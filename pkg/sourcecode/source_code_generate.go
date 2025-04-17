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

	//invocationModules := pie.Filter(scInfo.daprModules, func(info *daprModuleInfo) bool {
	//	return info.kind == DaprModuleKindInvocation
	//})

	//routeAnnotations := make([]*types.RouteAnnotation, 0)
	//for _, m := range invocationModules {
	//	for _, h := range scInfo.daprInvocationHandlers {
	//		if m.pkgRelPath == h.pkgRelPath && m.name == h.module {
	//			r, err := impl.createRouteItem(h)
	//			if err != nil {
	//				return err
	//			}
	//
	//			if r != nil {
	//				routeAnnotations = append(routeAnnotations, r)
	//			}
	//		}
	//	}
	//}

	//if len(routeAnnotations) > 0 {
	data, err := json.Marshal(scInfo.daprInvocationHandlers)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(impl.assetsPath, fileInvocationHandlers)
	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		return err
	}
	//}

	return nil
}

//
//func (impl *sourceCodeHandlerImpl) formatInvo(h *daprInvocationHandler) (*Route, error) {
//	// 只要有@hd.route就需要生成routeItem
//	annValue, exist := h.annotations[routeAnnotationName]
//	if !exist {
//		return nil, nil
//	}
//
//	// 设置初始值
//	routeItem := &Route{
//		PackagePath: h.pkgRelPath,
//		Module:      h.module,
//		Handler:     h.alias,
//		Comments:    h.comments,
//		HttpMethods: []string{"GET"},
//		Permissions: []string{},
//	}
//
//	annValue = strings.TrimSpace(annValue)
//	if annValue == "" {
//		return routeItem, nil
//	}
//
//	// 如果已经定义routeAnnotation, 则updateRouteItem
//	var ann routeAnnotation
//	err := json.Unmarshal([]byte(annValue), &ann)
//	if err != nil {
//		return nil, err
//	}
//
//	if ann.Endpoint != "" {
//		routeItem.Endpoint = ann.Endpoint
//	}
//
//	if ann.IsPublic {
//		routeItem.IsPublic = 1
//	}
//
//	if ann.IsRawResponse {
//		routeItem.IsRawResponse = 1
//	}
//
//	if len(ann.Methods) > 0 {
//		routeItem.HttpMethods = ann.Methods
//	}
//
//	if ann.Origin != "" {
//		routeItem.AllowOrigin = ann.Origin
//	}
//
//	if len(ann.Permissions) > 0 {
//		routeItem.Permissions = ann.Permissions
//	}
//
//	return routeItem, nil
//}
