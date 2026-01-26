package dapr

import (
	"fmt"
	"go/ast"

	astUtils "github.com/hdget/utils/ast"

	"go/token"
	"maps"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/elliotchance/pie/v2"
	"github.com/hdget/common/protobuf"
)

var (
	// 调用函数的函数签名：invocation handler: func(biz.Context,[]byte) (any, error)
	functionSignatureInvocationHandler = &astUtils.FunctionSignature{
		NamePattern: regexp.MustCompile(`.*Handler`),
		Params:      []string{"biz.Context", "[]byte"},
		Results:     []string{"any", "error"},
	}

	// 模块注册的调用签名, e,g:
	// moduleRegisterCall = &callSignature{
	//	functionChain: "Register",
	//	argCount:      2,
	//	argIndex2Signature: map[int]string{
	//		1: "map[string]dapr.InvocationFunction",
	//	},
	//}
	// 模块初始化的调用签名
	callSignatureNewInvocationModule = &astUtils.CallSignature{
		FunctionChain: "NewInvocationModule",
		Package:       "github.com/hdget/lib-dapr/module",
		ArgCount:      3,
	}

	hdAnnotationRegex  = regexp.MustCompile(`@hd\.(\S+)(?:\s+(.*))?`)
	commentMarkerRegex = regexp.MustCompile(`^\s*(?:\/\/|\/\*|\*\/?)\s*`)
)

// parseDaprInvocationHandlers 从第一次解析的结果中去获取DaprInvocationModule中所有handler的路由注解
func (p *scParser) parseDaprInvocationHandlers(moduleInfos []*daprModule) ([]*protobuf.DaprHandler, error) {
	results := make([]*protobuf.DaprHandler, 0)

	invocationModuleInfos := pie.Filter(moduleInfos, func(m *daprModule) bool {
		return m.kind == protobuf.DaprModuleKind_DaprModuleKindInvocation
	})

	allInvocationPkgRelPaths := pie.Unique(pie.Map(invocationModuleInfos, func(m *daprModule) string {
		return m.pkgRelPath
	}))

	// 获取到注册的handler
	registeredHandlerPath2handlerAlias := make(map[string]string)
	for _, pkgRelPath := range allInvocationPkgRelPaths {
		if astPkg := p.pkgRelPath2astPkg[pkgRelPath]; astPkg != nil {

			// 收集整个包的变量声明和定义
			pkgVarTypes := make(map[string]string)
			pkgVarDecls := make(map[string]*ast.ValueSpec)
			for _, f := range astPkg.Files {
				maps.Copy(pkgVarTypes, astUtils.GetVarTypes(f))
				maps.Copy(pkgVarDecls, astUtils.GetVarDeclsFromFile(f))
			}

			for _, f := range astPkg.Files {
				// 新建一个记录导入别名与包名映射关系的字典
				caller2pkgImportPath := astUtils.GetPackageImportPaths(f)

				ast.Inspect(f, func(node ast.Node) bool {
					switch n := node.(type) {
					case *ast.FuncDecl:
						if n.Name.Name == "init" {
							founds := p.parseInvocationHandlerAlias(n, pkgRelPath, caller2pkgImportPath, pkgVarTypes, pkgVarDecls)
							if len(founds) > 0 {
								maps.Copy(registeredHandlerPath2handlerAlias, founds)
							}
							return false
						}
					}
					return true
				})

			}
		}
	}

	for _, pkgRelPath := range allInvocationPkgRelPaths {
		if astPkg := p.pkgRelPath2astPkg[pkgRelPath]; astPkg != nil {
			for fPath, f := range astPkg.Files {
				ast.Inspect(f, func(node ast.Node) bool {
					switch n := node.(type) {
					case *ast.FuncDecl:
						if h := p.parseInvocationHandler(n, p.srcDir, fPath, registeredHandlerPath2handlerAlias); h != nil {
							results = append(results, h)
						}
					}
					return true
				})

			}
		}
	}

	return results, nil
}

// parseInvocationHandler 解析Dapr所有invocation handlers
func (p *scParser) parseInvocationHandler(fn *ast.FuncDecl, srcDir, filePath string, registeredHandlerPath2handlerAlias map[string]string) *protobuf.DaprHandler {
	receiverTypeName := astUtils.GetReceiverTypeName(fn, true)
	// receiverTypeName为空表示为普通函数，忽略
	if receiverTypeName == "" {
		return nil
	}

	// 函数签名匹配
	// func(ctx context.Context, event *common.InvocationEvent) (*common.Content, any)
	if astUtils.MatchFunction(fn, functionSignatureInvocationHandler) {
		annotations, comments := p.extractAnnotationsAndComments(fn.Doc)

		pkgRelPath, _ := filepath.Rel(srcDir, filepath.Dir(filePath))
		pkgRelPath = filepath.ToSlash(pkgRelPath)

		handlerPath := fmt.Sprintf("%s.%s.%s", pkgRelPath, receiverTypeName, fn.Name.Name)
		if handlerAlias, exist := registeredHandlerPath2handlerAlias[handlerPath]; exist && handlerAlias != "" {
			return &protobuf.DaprHandler{
				PkgPath:     pkgRelPath,
				ModuleKind:  protobuf.DaprModuleKind_DaprModuleKindInvocation,
				Module:      receiverTypeName,
				Name:        fn.Name.Name,
				Alias:       handlerAlias,
				Comments:    comments,
				Annotations: annotations,
			}
		}
	}

	return nil
}

func (p *scParser) extractAnnotationsAndComments(doc *ast.CommentGroup) (map[string]string, []string) {
	if doc == nil || len(doc.List) == 0 {
		return nil, nil
	}

	annotations := make(map[string]string)
	comments := make([]string, 0)
	for _, comment := range doc.List {
		text := comment.Text
		// 移除注释标记
		content := strings.TrimSpace(commentMarkerRegex.ReplaceAllString(text, ""))

		// 检查是否是注解行
		if matches := hdAnnotationRegex.FindAllStringSubmatch(content, -1); len(matches) > 0 {
			for _, match := range matches {
				if len(match) >= 2 {
					name := match[1]
					value := ""
					if len(match) >= 3 {
						value = strings.TrimSpace(match[2])
					}
					annotations[name] = value
				}
			}
		} else {
			// 不是注解行，保留注释内容
			// 处理多行注释块
			if strings.HasPrefix(text, "/*") {
				lines := strings.Split(content, "\n")
				for _, line := range lines {
					trimmed := strings.TrimSpace(line)
					if trimmed != "" {
						comments = append(comments, trimmed)
					}
				}
			} else {
				// 单行注释
				comments = append(comments, content)
			}
		}
	}

	return annotations, comments
}

// extractHandlerAlias 提取map中的键值对，并替换变量名为类型名
func (p *scParser) extractHandlerAlias(mapLit *ast.CompositeLit, pkgRelPath string, fnVarTypes, pkgVarTypes map[string]string) map[string]string {
	results := make(map[string]string)
	for _, elt := range mapLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		// 提取键（如 "refresh_image"）
		key, ok := kv.Key.(*ast.BasicLit)
		if !ok || key.Kind != token.STRING {
			continue
		}
		keyStr := strings.Trim(key.Value, `"`)

		// 提取值（如 v.refreshImageCaptchaHandler）
		val, ok := kv.Value.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		// 获取接收者变量名（如 v）
		receiver, ok := val.X.(*ast.Ident)
		if !ok {
			continue
		}

		// 查找变量的实际类型名（如 v2_captcha）先尝试从函数内部查找变量名，然后包级，最后回退到变量名
		actualReceiverType := fnVarTypes[receiver.Name]
		if actualReceiverType == "" {
			actualReceiverType = pkgVarTypes[receiver.Name]
			if actualReceiverType == "" {
				actualReceiverType = receiver.Name // 回退为变量名
			}
		}

		results[fmt.Sprintf("%s.%s.%s", pkgRelPath, actualReceiverType, val.Sel.Name)] = keyStr
	}

	return results
}

// parseInvocationHandlerAlias 在init函数中解析daprModule.Register函数，获取handlerAlias
//
// 返回： service/invocation.v2_xxx.handler => alias
func (p *scParser) parseInvocationHandlerAlias(n *ast.FuncDecl, pkgRelPath string, caller2pkgImportPath, pkgVarTypes map[string]string, pkgVarDecls map[string]*ast.ValueSpec) map[string]string {
	var handler2alias map[string]string

	fnVarDecls := astUtils.GetVarDeclsFromFunc(n.Body)
	fnVarTypes := astUtils.GetVarTypes(n.Body)

	ast.Inspect(n.Body, func(n ast.Node) bool {
		switch nn := n.(type) {
		case *ast.CallExpr:
			if astUtils.MatchCall(nn, callSignatureNewInvocationModule, caller2pkgImportPath) && len(nn.Args) == 3 {
				// 处理map参数（直接内联或通过变量传递）
				switch param := nn.Args[2].(type) {
				case *ast.CompositeLit: // 直接内联map
					handler2alias = p.extractHandlerAlias(param, pkgRelPath, fnVarTypes, pkgVarTypes)
				case *ast.Ident: // 通过变量传递
					// 首先尝试从init函数中获取变量定义
					if varDecl := fnVarDecls[param.Name]; varDecl != nil {
						if mapLit, ok := varDecl.Values[0].(*ast.CompositeLit); ok {
							handler2alias = p.extractHandlerAlias(mapLit, pkgRelPath, fnVarTypes, pkgVarTypes)
						}
						// 然后尝试从整个包里面获取变量定义
					} else if varDecl = pkgVarDecls[param.Name]; varDecl != nil {
						if mapLit, ok := varDecl.Values[0].(*ast.CompositeLit); ok {
							handler2alias = p.extractHandlerAlias(mapLit, pkgRelPath, fnVarTypes, pkgVarTypes)
						}
					}
				}
				return false
			}
		}
		return true
	})

	return handler2alias
}
