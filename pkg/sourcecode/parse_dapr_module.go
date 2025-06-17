package sourcecode

import (
	"fmt"
	"github.com/hdget/common/protobuf"
	"go/ast"
	"path/filepath"
)

type parsedDaprModuleInfo struct {
	kind       protobuf.DaprModuleKind
	name       string
	pkgRelPath string
}

var (
	moduleExpr2moduleKind = map[string]protobuf.DaprModuleKind{
		"&{dapr InvocationModule}": protobuf.DaprModuleKind_DaprModuleKindInvocation, // 服务调用模块
		"&{dapr EventModule}":      protobuf.DaprModuleKind_DaprModuleKindEvent,      // 事件模块
		"&{dapr DelayEventModule}": protobuf.DaprModuleKind_DaprModuleKindDelayEvent, // 延迟事件模块
		"&{dapr HealthModule}":     protobuf.DaprModuleKind_DaprModuleKindHealth,     // 健康检测模块
	}
)

// parseDaprModules 所有DaprModule信息
func (p *parserImpl) parseDaprModules() ([]*parsedDaprModuleInfo, error) {
	results := make([]*parsedDaprModuleInfo, 0)

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
func (p *parserImpl) parseDaprModule(structName string, astStructType *ast.StructType, srcDir, path string) *parsedDaprModuleInfo {
	// 检查第一个field是否是匿名引入的模块， e,g: type A struct { dapr.InvocationModule }
	for _, field := range astStructType.Fields.List {
		fieldTypeExpr := fmt.Sprintf("%s", field.Type)
		if moduleKind, exists := moduleExpr2moduleKind[fieldTypeExpr]; exists && moduleKind != protobuf.DaprModuleKind_DaprModuleKindUnknown {
			relPath, _ := filepath.Rel(srcDir, filepath.Dir(path))
			return &parsedDaprModuleInfo{
				kind:       moduleKind,
				name:       structName,
				pkgRelPath: filepath.ToSlash(relPath),
			}
		}
	}
	return nil
}
