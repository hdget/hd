package protorefine

import "path/filepath"

type ProtoRefiner interface {
	Refine(pbImportPath, protoDir, outputDir string, skipDirs ...string) error // 精简proto files
}

type protoRefineImpl struct {
	srcDir string
}

func New(srcDir string) ProtoRefiner {
	return &protoRefineImpl{
		srcDir: srcDir,
	}
}

// Refine 精简proto下的文件
func (impl *protoRefineImpl) Refine(pbImportPath, protoDir, outputDir string, skipDirs ...string) error {
	// 从golang源代码中找到protobuf类型的变量类型
	golangPkgName, golangTypeNames, err := newGolangParser().parse(impl.srcDir, pbImportPath, skipDirs...)
	if err != nil {
		return err
	}

	// 去匹配proto文件中和源文件中匹配的类型的声明
	protoDeclares, err := newProtoParser().findProtoDeclares(protoDir, golangPkgName, golangTypeNames)
	if err != nil {
		return err
	}

	// 生成新的proto文件以及其依赖文件
	outputDir = filepath.Join(outputDir, filepath.Base(protoDir))
	outputFileName := filepath.Base(impl.srcDir)
	err = newProtoRender().Render(protoDir, outputDir, outputFileName, golangPkgName, protoDeclares)
	if err != nil {
		return err
	}

	return nil
}
