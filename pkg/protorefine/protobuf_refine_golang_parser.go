package protorefine

import (
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/pkg/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// github.com/rfancn/protorefine
type golangParser struct {
}

func newGolangParser() *golangParser {
	return &golangParser{}
}

// parse find all <protobuf package>.<type> in golang source codes
// return golang protobuf package name, golang type names, error
// e,g:
// pb.ProductKind
// golang protobuf package name: pb
// golang type name: ProductKind
func (golangParser) parse(srcDir, pbImportPath string, skipDirs ...string) ([]string, error) {
	st, err := os.Stat(srcDir)
	if err != nil {
		return nil, err
	}

	if !st.IsDir() {
		return nil, fmt.Errorf("invalid source code dir, dir: %s", srcDir)
	}

	results := make(map[string]struct{})
	//name2pkgImportPath := make(map[string]string)
	_ = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && pie.Contains(skipDirs, info.Name()) {
			return filepath.SkipDir
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				return errors.Wrapf(err, "golang ast parse file, path: %s", path)
			}

			pkgName2importPath := astGetImportPaths(f)

			ast.Inspect(f, func(node ast.Node) bool {
				switch n := node.(type) {
				case *ast.SelectorExpr:
					if ident, ok := n.X.(*ast.Ident); ok {
						// if type's package name equal to pbPkgImportPath
						if pkgName2importPath[ident.Name] == pbImportPath {
							if _, exists := results[n.Sel.Name]; !exists {
								results[n.Sel.Name] = struct{}{}
							}
						}
					}
					//case *ast.ImportSpec: // record all name2pkgImportPath
					//	var alias string
					//	if n.name != nil {
					//		alias = n.name.name
					//	}
					//	fullPath := n.Path.Value[1 : len(n.Path.Value)-1]
					//	pkgName := filepath.Base(fullPath)
					//	if alias == "" {
					//		name2pkgImportPath[pkgName] = fullPath
					//	} else {
					//		name2pkgImportPath[alias] = fullPath
					//	}
				}
				return true
			})
		}
		return nil
	})

	return pie.Keys(results), nil
}

// astGetImportPaths 获取导入路径
func astGetImportPaths(f *ast.File) map[string]string {
	importMap := make(map[string]string)
	for _, imp := range f.Imports {
		pkgName := ""
		if imp.Name != nil {
			pkgName = imp.Name.Name // 处理别名导入，如 `import alias "math/rand"`
		} else {
			// 提取完整路径（去掉引号）
			pkgPath := strings.Trim(imp.Path.Value, `"`)
			// 获取包名（路径的最后一部分）
			pkgName = pkgPath[strings.LastIndex(pkgPath, "/")+1:]
		}
		importMap[pkgName] = strings.Trim(imp.Path.Value, `"`)
	}
	return importMap
}
