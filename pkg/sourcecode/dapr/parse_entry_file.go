package dapr

import (
	"go/ast"

	astUtils "github.com/hdget/utils/ast"
)

var (
	callSignatureNewGrpcServer = &astUtils.CallSignature{
		FunctionChain: "NewGrpcServer",
		Package:       "github.com/hdget/lib-dapr",
		ArgCount:      -1, // 不检查argCount
	}
)

// parseServerEntryFilePath 查找server.Start调用所在的文件路径
func (p *scParser) parseServerEntryFilePath() (string, error) {
	var serverEntryFilePath string

	for _, astPkg := range p.pkgRelPath2astPkg {
		for fPath, f := range astPkg.Files {
			// 新建一个记录导入别名与包名映射关系的字典
			caller2pkgImportPath := astUtils.GetPackageImportPaths(f)

			ast.Inspect(f, func(node ast.Node) bool {
				switch n := node.(type) {
				case *ast.CallExpr: // 函数调用
					if serverEntryFilePath != "" {
						return false
					}

					// 查找server.Start调用所在的文件路径
					if astUtils.MatchCall(n, callSignatureNewGrpcServer, caller2pkgImportPath) {
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
