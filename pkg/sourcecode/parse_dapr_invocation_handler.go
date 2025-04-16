package sourcecode

import (
	"fmt"
	"github.com/elliotchance/pie/v2"
	"go/ast"
	"go/token"
	"maps"
	"path/filepath"
	"regexp"
	"strings"
)

type daprInvocationHandler struct {
	pkgRelPath  string
	moduleName  string // receiver name
	alias       string
	name        string // method name
	comments    []string
	annotations map[string]string // annotationName => annotation value
}

var (
	hdAnnotationRegex  = regexp.MustCompile(`@hd\.(\S+)(?:\s+(.*))?`)
	commentMarkerRegex = regexp.MustCompile(`^\s*(?:\/\/|\/\*|\*\/?)\s*`)
	// invocation handler: func((context.Context, *common.InvocationEvent) (any, error)
	invocationHandlerParamSignatures  = []string{"context.Context", "*common.InvocationEvent"}
	invocationHandlerResultSignatures = []string{"any", "error"}

	// 模块注册的调用签名
	moduleRegisterCall = &callSignature{
		functionChain: "Register",
		argCount:      2,
		argIndex2Signature: map[int]string{
			1: "map[string]dapr.InvocationFunction",
		},
	}
)

// parseDaprInvocationHandlers 从第一次解析的结果中去获取DaprInvocationModule中所有handler的路由注解
func (p *parserImpl) parseDaprInvocationHandlers(moduleInfos []*daprModuleInfo) ([]*daprInvocationHandler, error) {
	results := make([]*daprInvocationHandler, 0)

	invocationModuleInfos := pie.Filter(moduleInfos, func(m *daprModuleInfo) bool {
		return m.kind == DaprModuleKindInvocation
	})

	allInvocationPkgRelPaths := pie.Unique(pie.Map(invocationModuleInfos, func(m *daprModuleInfo) string {
		return m.pkgRelPath
	}))

	// 获取到所有handler别名
	handlerName2handlerAlias := make(map[string]string)
	for _, pkgRelPath := range allInvocationPkgRelPaths {
		if astPkg := p.pkgRelPath2astPkg[pkgRelPath]; astPkg != nil {
			for _, f := range astPkg.Files {
				ast.Inspect(f, func(node ast.Node) bool {
					switch n := node.(type) {
					case *ast.FuncDecl:
						if n.Name.Name == "init" {
							founds := p.parseInvocationHandlerAlias(n, pkgRelPath)
							if len(founds) > 0 {
								maps.Copy(handlerName2handlerAlias, founds)
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
						if h := p.parseInvocationHandler(n, p.srcDir, fPath, handlerName2handlerAlias); h != nil {
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
func (p *parserImpl) parseInvocationHandler(fn *ast.FuncDecl, srcDir, filePath string, handlerName2handlerAlias map[string]string) *daprInvocationHandler {
	receiverTypeName := astGetReceiverTypeName(fn, true)
	// receiverTypeName为空表示为普通函数，忽略
	if receiverTypeName == "" {
		return nil
	}

	// 函数签名匹配
	// func(ctx context.Context, event *common.InvocationEvent) (*common.Content, any)
	if astMatchFunction(fn, invocationHandlerParamSignatures, invocationHandlerResultSignatures) {
		annotations, comments := p.extractAnnotationsAndComments(fn.Doc)

		pkgRelPath, _ := filepath.Rel(srcDir, filepath.Dir(filePath))
		pkgRelPath = filepath.ToSlash(pkgRelPath)

		fullHandlerName := fmt.Sprintf("%s.%s.%s", pkgRelPath, receiverTypeName, fn.Name.Name)

		return &daprInvocationHandler{
			pkgRelPath:  pkgRelPath,
			moduleName:  receiverTypeName,
			name:        fn.Name.Name,
			alias:       handlerName2handlerAlias[fullHandlerName],
			comments:    comments,
			annotations: annotations,
		}
	}

	return nil
}

func (p *parserImpl) extractAnnotationsAndComments(doc *ast.CommentGroup) (map[string]string, []string) {
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
func (p *parserImpl) extractHandlerAlias(n ast.Expr, pkgRelPath string, varTypes map[string]string) map[string]string {
	mapLit, ok := n.(*ast.CompositeLit)
	if !ok {
		return nil
	}

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

		// 替换为实际类型名（如 v2_captcha）
		actualValueType := varTypes[receiver.Name]
		if actualValueType == "" {
			actualValueType = receiver.Name // 回退为变量名
		}

		results[fmt.Sprintf("%s.%s.%s", pkgRelPath, actualValueType, val.Sel.Name)] = keyStr
	}

	return results
}

// parseInvocationHandlerAlias 在init函数中解析daprModule.Register函数，获取handlerAlias
//
// 返回： service/invocation.v2_xxx.handler => alias
func (p *parserImpl) parseInvocationHandlerAlias(n *ast.FuncDecl, pkgRelPath string) map[string]string {

	var handler2alias map[string]string

	// 获取所有变量和变量声明的映射关系
	varTypes := astGetVarTypes(n.Body)

	ast.Inspect(n.Body, func(n ast.Node) bool {

		switch nn := n.(type) {
		case *ast.CallExpr:
			if astMatchCall(nn, moduleRegisterCall, nil) && len(nn.Args) == 2 {
				handler2alias = p.extractHandlerAlias(nn.Args[1], pkgRelPath, varTypes)
				return false
			}
		}
		return true
	})

	return handler2alias
}
