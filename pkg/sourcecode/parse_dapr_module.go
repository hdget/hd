package sourcecode

import (
	"fmt"
	"go/ast"
	"path/filepath"
	"regexp"

	"github.com/hdget/common/protobuf"
)

type parsedDaprModuleInfo struct {
	kind       protobuf.DaprModuleKind
	name       string
	pkgRelPath string
}

var (
	regexModule           = regexp.MustCompile(`^&{\w+\s+(\w+)}$`)
	moduleExpr2moduleKind = map[string]protobuf.DaprModuleKind{
		"InvocationModule": protobuf.DaprModuleKind_DaprModuleKindInvocation, // 服务调用模块
		"EventModule":      protobuf.DaprModuleKind_DaprModuleKindEvent,      // 事件模块
		"DelayEventModule": protobuf.DaprModuleKind_DaprModuleKindDelayEvent, // 延迟事件模块
		"HealthModule":     protobuf.DaprModuleKind_DaprModuleKindHealth,     // 健康检测模块
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
//		module.InvocationModule
//	}
func (p *parserImpl) parseDaprModule(structName string, astStructType *ast.StructType, srcDir, path string) *parsedDaprModuleInfo {
	// 检查第一个field是否是匿名引入的模块，
	// e,g: type A struct {
	//	module.InvocationModule
	//}
	for _, field := range astStructType.Fields.List {
		// 匿名字段Names为空
		if field.Names != nil {
			continue
		}

		fieldTypeExpr := fmt.Sprintf("%s", field.Type)
		matches := regexModule.FindStringSubmatch(fieldTypeExpr)
		if len(matches) == 2 {
			moduleName := matches[1]
			if moduleKind, exists := moduleExpr2moduleKind[moduleName]; exists && moduleKind != protobuf.DaprModuleKind_DaprModuleKindUnknown {
				relPath, _ := filepath.Rel(srcDir, filepath.Dir(path))
				return &parsedDaprModuleInfo{
					kind:       moduleKind,
					name:       structName,
					pkgRelPath: filepath.ToSlash(relPath),
				}
			}
		}
	}
	return nil
}
