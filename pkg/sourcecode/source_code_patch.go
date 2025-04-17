package sourcecode

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"os/exec"
	"path"
	"strings"
)

// patch 匿名导入invocationModule和eventModule的路径到服务运行入口文件
// MonkeyPatch 修改源代码的方式匿名导入pkg, sourceFile是相对于basePath的相对路径
func (impl *sourceCodeHandlerImpl) patch(scInfo *sourceCodeInfo) error {
	fmt.Println(">>> patch source code")

	// 获取项目模块名
	projectModuleName, err := impl.getProjectModuleName()
	if err != nil {
		return err
	}

	// 将源代码解析为抽象语法树（AST）
	fset := token.NewFileSet()
	// IMPORTANT: 这里要保证注释不被丢失
	astFile, err := parser.ParseFile(fset, scInfo.serverEntryFilePath, nil, parser.ParseComments)
	if err != nil {
		return errors.Wrapf(err, "golang ast parseMetaData file, path: %s", scInfo.serverEntryFilePath)
	}

	// 记录所有已经导入的包
	importedPaths := make(map[string]struct{})
	for _, spec := range astFile.Imports {
		importedPaths[spec.Path.Value] = struct{}{}
	}

	// 创建新的import节点匿名插入到import声明列表
	newImported := make([]string, 0)
	for _, m := range scInfo.daprModules {
		// IMPORTANT: spec.Path.Value是带了双引号的
		impPath := "\"" + path.Join(projectModuleName, m.pkgRelPath) + "\""

		// 当patch进去的路径不存在时才加入
		if _, exists := importedPaths[impPath]; !exists {
			// 创建一个新的匿名ImportSpec节点
			spec := &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: impPath,
				},
				Name: ast.NewIdent("_"), // 下划线表示匿名导入
			}

			// 创建一个新的声明并插入到文件的声明列表中
			decl := &ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					spec,
				},
			}

			astFile.Decls = append([]ast.Decl{decl}, astFile.Decls...)

			newImported = append(newImported, impPath)

			importedPaths[impPath] = struct{}{}
		}
	}

	if len(newImported) == 0 {
		fmt.Println("All modules imported. No action required!")
		fmt.Println("")
		return nil
	}

	// 使用printer包将抽象语法树（AST）打印成代码
	buf := bytes.NewBufferString("")
	err = printer.Fprint(buf, fset, astFile)
	if err != nil {
		return err
	}

	// 打开文件
	file, err := os.OpenFile(scInfo.serverEntryFilePath, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将新代码内容写入文件
	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}

	// 确保所有操作都已写入磁盘
	err = file.Sync()
	if err != nil {
		return err
	}

	for _, impPath := range newImported {
		fmt.Printf("%s added\n", impPath)
	}
	//h.print(newImported)

	return nil
}

//func (h *addModuleImportPathsHandler) print(newAdded []string) {
//	table := tablewriter.NewWriter(os.Stdout)
//	table.SetHeader([]string{"NO", "IMPORT PACKAGE"})
//	table.SetRowLine(true)
//	for i, impPath := range newAdded {
//		table.Append([]string{cast.ToString(i + 1), impPath})
//	}
//	table.Render() // Send output
//}

func (impl *sourceCodeHandlerImpl) getProjectModuleName() (string, error) {
	// 获取根模块名
	cmdOutput, err := exec.Command("go", "list", "-m").CombinedOutput()
	if err != nil {
		return "", err
	}

	// 按换行符拆分结果
	lines := bytes.Split(cmdOutput, []byte("\n"))
	if len(lines) == 0 {
		return "", errors.New("project is not using go module or not run go list -m in project root dir")
	}

	return strings.TrimSpace(string(lines[0])), nil
}
