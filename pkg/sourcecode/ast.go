package sourcecode

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"github.com/elliotchance/pie/v2"
)

type callSignature struct {
	functionChain      string         // 必传
	pkg                string         // 可选
	argCount           int            // 可选, -1不去检查
	argIndex2Signature map[int]string // 可选, nil不去检查
}

type functionSignature struct {
	namePattern *regexp.Regexp
	params      []string
	results     []string
}

func astMatchCall(n *ast.CallExpr, sig *callSignature, imports map[string]string) bool {
	if astGetFunctionChain(n) != sig.functionChain {
		return false
	}

	if sig.argCount > 0 {
		if len(n.Args) != sig.argCount {
			return false
		}
	}

	// 如果传入了pkg，则调用类似: dapr.New
	if sig.pkg != "" {
		caller, found := astGetCaller(n)
		if !found {
			return false
		}
		if imports[caller] != sig.pkg {
			return false
		}
	}

	for argIndex, signature := range sig.argIndex2Signature {
		if argIndex >= len(n.Args) {
			return false
		}

		if astGetExprTypeName(n.Args[argIndex]) != signature {
			return false
		}
	}

	return true
}

// 收集所有变量声明中的类型信息
func astGetVarTypes(node ast.Node) map[string]string {
	results := make(map[string]string)
	ast.Inspect(node, func(n ast.Node) bool {
		switch nn := n.(type) {
		case *ast.AssignStmt:
			// 处理短声明（如 v := &v2_captcha{}）
			if nn.Tok == token.DEFINE {
				for i, lhs := range nn.Lhs {
					if i >= len(nn.Rhs) {
						break
					}
					varName := lhs.(*ast.Ident).Name
					typeName := astResolveVarType(nn.Rhs[i])
					if typeName != "" {
						results[varName] = typeName
					}
				}
			}
		case *ast.ValueSpec:
			// 处理普通声明（如 var v = &v2_captcha{}）
			for i, name := range nn.Names {
				if i >= len(nn.Values) {
					break
				}
				varName := name.Name
				typeName := astResolveVarType(nn.Values[i])
				if typeName != "" {
					results[varName] = typeName
				}
			}
		}
		return true
	})
	return results
}

// 解析表达式的实际类型名（如 v2_captcha）
func astResolveVarType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.UnaryExpr:
		// 处理指针类型（如 &v2_captcha{}）
		if t.Op == token.AND {
			return astResolveVarType(t.X)
		}
	case *ast.CompositeLit:
		// 处理结构体初始化（如 v2_captcha{}）
		if ident, ok := t.Type.(*ast.Ident); ok {
			return ident.Name
		}
	case *ast.CallExpr:
		// 处理构造函数（如 new(v2_captcha)）
		if fun, ok := t.Fun.(*ast.Ident); ok && fun.Name == "new" {
			if arg, ok := t.Args[0].(*ast.Ident); ok {
				return arg.Name
			}
		}
	}
	return ""
}

func astMatchFunction(fn *ast.FuncDecl, sig *functionSignature) bool {
	// 检查参数数量
	if fn.Type.Params == nil || len(fn.Type.Params.List) != len(sig.params) || fn.Type.Results == nil || len(fn.Type.Results.List) != len(sig.results) {
		return false
	}

	if sig.namePattern != nil {
		if !sig.namePattern.MatchString(fn.Name.Name) {
			return false
		}
	}

	// 检查参数
	for i, param := range fn.Type.Params.List {
		if paramTypeName := astGetExprTypeName(param.Type); paramTypeName != sig.params[i] {
			return false
		}
	}

	for i, result := range fn.Type.Results.List {
		if resultTypeName := astGetExprTypeName(result.Type); resultTypeName != sig.results[i] {
			return false
		}
	}

	return true
}

// astGetReceiverTypeName 获取函数的receiver类型, 如果ignorePointer为true, 则去除前面的*号
// e,g: (*Person) hello() {}, 传入hello的ast.FuncDecl, 返回Person
func astGetReceiverTypeName(fn *ast.FuncDecl, ignorePointer ...bool) string {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return ""
	}

	v := astGetExprTypeName(fn.Recv.List[0].Type)
	if len(ignorePointer) > 0 && ignorePointer[0] {
		return strings.ReplaceAll(v, "*", "")
	}
	return v
}

// astGetExprType 返回带有指针指示的类型名称字符串
func astGetExprTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.BasicLit:
		return t.Kind.String() // "INT", "STRING"等
	case *ast.CompositeLit:
		return astGetExprTypeName(t.Type)
	case *ast.CallExpr:
		return astGetExprTypeName(t.Fun) + "()"
	case *ast.UnaryExpr:
		return fmt.Sprintf("%s%s", t.Op, astGetExprTypeName(t.X))
	case *ast.Ident:
		return t.Name // 基础类型
	case *ast.StarExpr:
		return "*" + astGetExprTypeName(t.X) // 指针类型加*前缀
	case *ast.SelectorExpr:
		return astGetExprTypeName(t.X) + "." + t.Sel.Name // 包.类型
	case *ast.ArrayType:
		return "[]" + astGetExprTypeName(t.Elt) // 切片类型
	case *ast.MapType:
		// 特别处理map的value部分，区分指针
		keyType := astGetExprTypeName(t.Key)
		valueType := astGetExprTypeName(t.Value)
		return fmt.Sprintf("map[%s]%s", keyType, valueType)
	case *ast.InterfaceType:
		return "interface{}" // 接口类型
	default:
		return fmt.Sprintf("%T", expr) // 其他未知类型
	}
}

// astGetStructType
func astGetStructInfo(n *ast.GenDecl) (string, *ast.StructType, bool) {
	// 仅处理类型声明
	if n.Tok == token.TYPE {
		for _, spec := range n.Specs {
			// 如果类型规范是类型别名或类型声明
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				// 如果类型规范是结构体类型
				if st, ok := typeSpec.Type.(*ast.StructType); ok {
					return typeSpec.Name.Name, st, true
				}
			}
		}
	}
	return "", nil, false
}

// astGetFunctionChain 获取完整的函数调用链
func astGetFunctionChain(n *ast.CallExpr) string {
	functions := astParseFunctionCallChain(n)
	return strings.Join(pie.Reverse(functions), ".")
}

// parseMetaData 递归解析链式函数调用，最近的Ident.Name作为包名，最先调用的函数在slice的最前面
func astParseFunctionCallChain(n *ast.CallExpr) []string {
	var methods []string

	// 递归提取方法名
	for {
		// 检查 Fun 是否是 SelectorExpr
		selectorExpr, ok := n.Fun.(*ast.SelectorExpr)
		if !ok {
			break
		}

		// 添加方法名
		methods = append(methods, selectorExpr.Sel.Name)

		// 检查 X 是否是另一个 CallExpr
		nextCallExpr, ok := selectorExpr.X.(*ast.CallExpr)
		if !ok {
			break
		}

		// 继续递归
		n = nextCallExpr
	}
	return methods
}

// astGetPackageImportPaths 获取包导入路径
func astGetPackageImportPaths(f *ast.File) map[string]string {
	importMap := make(map[string]string)
	for _, imp := range f.Imports {
		pkgNames := make([]string, 0)
		if imp.Name != nil {
			pkgNames = []string{imp.Name.Name} // 处理别名导入，如 `import alias "math/rand"`
		} else {
			// 提取完整路径（去掉引号）
			pkgPath := strings.Trim(imp.Path.Value, `"`)
			// 获取包名（路径的最后一部分）
			lastPart := pkgPath[strings.LastIndex(pkgPath, "/")+1:]
			// HOTFIX: 有时候定义包名会只使用横杠后的部分，例如: lib-dapr只会用dapr
			pkgNames = append(pkgNames, lastPart)
			possiblePkgName := lastPart[strings.LastIndex(lastPart, "-")+1:]
			if possiblePkgName != lastPart {
				pkgNames = append(pkgNames, lastPart)
			}
		}

		for _, pkgName := range pkgNames {
			importMap[pkgName] = strings.Trim(imp.Path.Value, `"`)
		}
	}
	return importMap
}

// 获取调用者
func astGetCaller(n *ast.CallExpr) (string, bool) {
	// 递归查找调用者
	for {
		// 检查 Fun 是否是 SelectorExpr
		selectorExpr, ok := n.Fun.(*ast.SelectorExpr)
		if !ok {
			break
		}

		// 检查 X 是否是另一个 CallExpr
		nextCallExpr, ok := selectorExpr.X.(*ast.CallExpr)
		if !ok {
			// 如果不是 CallExpr，则可能是调用者（如 sdk）
			if ident, ok := selectorExpr.X.(*ast.Ident); ok {
				return ident.Name, true
			}
			break
		}

		// 继续递归
		n = nextCallExpr
	}

	return "", false
}

// 查找变量声明
func astGetVarDeclsFromFile(file *ast.File) map[string]*ast.ValueSpec {
	results := make(map[string]*ast.ValueSpec)
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}
		for _, spec := range genDecl.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				for _, name := range valueSpec.Names {
					results[name.Name] = valueSpec
				}
			}
		}
	}
	return results
}

// astGetVarDeclsFromFunc 从函数体中提取所有 *ast.ValueSpec
// 提取函数体内所有ValueSpec（包括转换短声明）
func astGetVarDeclsFromFunc(body *ast.BlockStmt) map[string]*ast.ValueSpec {
	results := make(map[string]*ast.ValueSpec)

	for _, stmt := range body.List {
		switch node := stmt.(type) {
		case *ast.DeclStmt:
			// 处理var块声明（包含多个ValueSpec）
			if genDecl, ok := node.Decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
				for _, spec := range genDecl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						for _, name := range valueSpec.Names {
							results[name.Name] = valueSpec
						}
					}
				}
			}

		case *ast.AssignStmt:
			// 将短声明（:=）转换为ValueSpec
			if node.Tok == token.DEFINE {
				for i, lhs := range node.Lhs {
					if i >= len(node.Rhs) {
						break
					}
					if ident, ok := lhs.(*ast.Ident); ok {
						results[ident.Name] = &ast.ValueSpec{
							Names:  []*ast.Ident{ident},
							Values: []ast.Expr{node.Rhs[i]},
						}
					}
				}
			}
		}
	}
	return results
}

////
////// getFunctionChain 获取完整的函数调用链
////func gastGetFunctionChain(n *ast.CallExpr) string {
////	functions := make([]string, 0)
////	astRecursiveParseFunction(n, functions)
////	return strings.Join(pie.Reverse(functions), ".")
////}
//
//// astGetEmbedInfo 获取嵌入资源的信息，返回变量名，embed路径
//func astGetEmbedVarAndRelPath(n *ast.GenDecl) (string, string, bool) {
//	// 如果是 GenDecl 类型，则可能是 import 或者变量声明等
//	if n.Tok == token.VAR {
//		for _, spec := range n.Specs {
//			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
//				if astIsEmbedFSType(valueSpec.Type) {
//					return valueSpec.Names[0].name, astGetEmbedRelPath(n), true
//				}
//			}
//		}
//	}
//	return "", "", false
//}
//
//// 检查类型是否为 embed.FS
//func astIsEmbedFSType(expr ast.Expr) bool {
//	if selectorExpr, ok := expr.(*ast.SelectorExpr); ok {
//		if ident, ok := selectorExpr.X.(*ast.Ident); ok && ident.name == "embed" {
//			if selectorExpr.Sel.name == "FS" {
//				return true
//			}
//		}
//	}
//	return false
//}
//
//// 获取 embed 路径
//func astGetEmbedRelPath(n *ast.GenDecl) string {
//	// 如果直接定义变量
//	// //go:embed assets/*
//	// var assets embed.FS
//	if n.Doc != nil {
//		for _, comment := range n.Doc.List {
//			if strings.HasPrefix(comment.Text, "//go:embed") {
//				// 提取路径部分
//				return filepath.Dir(strings.TrimSpace(strings.TrimPrefix(comment.Text, "//go:embed")))
//			}
//		}
//	}
//	// 如果定义在var block中
//	// var (
//	//   //go:embed assets/*
//	//   assets embed.FS
//	// )
//	for _, spec := range n.Specs {
//		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
//			if valueSpec.Doc != nil {
//				for _, comment := range valueSpec.Doc.List {
//					if strings.HasPrefix(comment.Text, "//go:embed") {
//						// 提取路径部分
//						return filepath.Dir(strings.TrimSpace(strings.TrimPrefix(comment.Text, "//go:embed")))
//					}
//				}
//			}
//		}
//	}
//	return ""
//}
//
//// Parse 尝试从源代码中查找嵌入路径, 返回嵌入资源的绝对路径和相对路径
//func astParseEmbed(callerFilePath string) (string, string, error) {
//	// 创建一个新的文件集
//	fset := token.NewFileSet()
//
//	// 解析源文件，同时保留注释
//	f, err := parser.ParseFile(fset, callerFilePath, nil, parser.ParseComments)
//	if err != nil {
//		return "", "", err
//	}
//
//	// 遍历AST节点
//	count := 0
//	var foundVar, foundRelPath, embedAbsPath string
//	ast.Inspect(f, func(node ast.Node) bool {
//		switch n := node.(type) {
//		case *ast.GenDecl:
//			if varName, relPath, ok := astGetEmbedVarAndRelPath(n); ok {
//				foundVar = varName
//				foundRelPath = relPath
//				return false
//			}
//		}
//		count += 1
//		return foundVar == ""
//	})
//
//	fmt.Println("xxxxxxxxxxxxxx:", count)
//
//	if foundVar == "" {
//		return "", "", fmt.Errorf("embed.FS variable declare not found, var: %s", foundVar)
//	}
//
//	// 有可能定义了embed.FS,但是没有指定编译指令//go:embed
//	if foundRelPath == "" {
//		return "", "", fmt.Errorf("//go:embed compiler directive not found, var: %s", foundVar)
//	}
//
//	if foundRelPath == "." {
//		return "", "", fmt.Errorf("//go:embed must specify a directory, var: %s", foundVar)
//	}
//
//	embedAbsPath = filepath.Join(filepath.Dir(callerFilePath), foundRelPath)
//	return embedAbsPath, foundRelPath, nil
//}
