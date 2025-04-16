package sourcecode

import (
	"fmt"
	"go/ast"
	"path/filepath"
)

type DaprModuleKind int

const (
	DaprModuleKindUnknown DaprModuleKind = iota
	DaprModuleKindInvocation
	DaprModuleKindEvent
	DaprModuleKindDelayEvent
	DaprModuleKindHealth
)

type daprModuleInfo struct {
	kind       DaprModuleKind
	name       string
	pkgRelPath string
}

var (
	moduleExpr2moduleKind = map[string]DaprModuleKind{
		"&{dapr InvocationModule}": DaprModuleKindInvocation, // 服务调用模块
		"&{dapr EventModule}":      DaprModuleKindEvent,      // 事件模块
		"&{dapr HealthModule}":     DaprModuleKindDelayEvent, // 健康检测模块
		"&{dapr DelayEventModule}": DaprModuleKindHealth,     // 延迟事件模块
	}
)

// parseDaprModules 所有DaprModule信息
func (p *parserImpl) parseDaprModules() ([]*daprModuleInfo, error) {
	results := make([]*daprModuleInfo, 0)

	for _, astPkg := range p.pkgRelPath2astPkg {
		for fPath, f := range astPkg.Files {
			ast.Inspect(f, func(node ast.Node) bool {
				switch n := node.(type) {
				case *ast.GenDecl: // import, constant, type or variable declaration
					if stName, st, found := astGetStructInfo(n); found {
						if m := p.parseDaprModule(stName, st, p.srcDir, fPath); m != nil {
							results = append(results, m)
						}
					}
				}
				return true
			})
		}
	}
	return results, nil
}

// parseDaprModule 检查是否有Dapr模块
// 判断结构声明中的field类型是否出现在moduleExpr2moduleKind中，如果出现了，说明类似下面的type声明找到了
//
//	type v1_example struct {
//		dapr.InvocationModule
//	}
func (p *parserImpl) parseDaprModule(structName string, astStructType *ast.StructType, srcDir, path string) *daprModuleInfo {
	// 检查第一个field是否是匿名引入的模块， e,g: type A struct { dapr.InvocationModule }
	for _, field := range astStructType.Fields.List {
		fieldTypeExpr := fmt.Sprintf("%s", field.Type)
		if moduleKind, exists := moduleExpr2moduleKind[fieldTypeExpr]; exists && moduleKind != DaprModuleKindUnknown {
			relPath, _ := filepath.Rel(srcDir, filepath.Dir(path))
			return &daprModuleInfo{
				kind:       moduleKind,
				name:       structName,
				pkgRelPath: filepath.ToSlash(relPath),
			}
		}
	}
	return nil
}
