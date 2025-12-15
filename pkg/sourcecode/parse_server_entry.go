package sourcecode

import "go/ast"

var (
	serverEntryCall = &callSignature{
		functionChain: "NewGrpcServer",
		pkg:           "github.com/hdget/lib-dapr",
		argCount:      -1, // 不检查argCount
	}
)

// parseServerEntryFilePath 查找server.Start调用所在的文件路径
func (p *parserImpl) parseServerEntryFilePath() (string, error) {
	var serverEntryFilePath string

	for _, astPkg := range p.pkgRelPath2astPkg {
		for fPath, f := range astPkg.Files {
			// 新建一个记录导入别名与包名映射关系的字典
			caller2pkgImportPath := astGetPackageImportPaths(f)

			ast.Inspect(f, func(node ast.Node) bool {
				switch n := node.(type) {
				case *ast.CallExpr: // 函数调用
					if serverEntryFilePath != "" {
						return false
					}

					// 查找server.Start调用所在的文件路径
					if astMatchCall(n, serverEntryCall, caller2pkgImportPath) {
						serverEntryFilePath = fPath
						return false
					}
				}
				return true
			})
		}
	}
	return serverEntryFilePath, nil
}
