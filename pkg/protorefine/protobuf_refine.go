package protorefine

import (
	"fmt"
	"path"
	"path/filepath"
)

type ProtoRefine interface {
	GetOutputGolangPackageName(string) (string, error) // 获取输出的golang包名
	Refine(Argument) (string, error)                   // 精简proto files, 返回proto输出目录
}

type protoRefineImpl struct {
	debug bool
}

func New(options ...Option) ProtoRefine {
	impl := &protoRefineImpl{
		debug: false,
	}

	for _, apply := range options {
		apply(impl)
	}

	return impl
}

func (impl *protoRefineImpl) GetOutputGolangPackageName(protoDir string) (string, error) {
	return newProtoParser().getGolangPackageName(protoDir)
}

// Refine 精简proto下的文件
func (impl *protoRefineImpl) Refine(arg Argument) (string, error) {
	if err := arg.validate(); err != nil {
		return "", err
	}

	absProtoRepository, _ := filepath.Abs(arg.ProtoRepository)
	absOutputDir, _ := filepath.Abs(arg.OutputDir)
	pbImportPath := path.Join(arg.GolangModule, arg.OutputDir, arg.OutputPackage)
	outputDir := filepath.Join(absOutputDir, filepath.Base(arg.ProtoRepository))

	if impl.debug {
		fmt.Println("### PROTOBUF REFINE ###")
		fmt.Println("golang module:", arg.GolangModule)
		fmt.Println("pb import path:", pbImportPath)
		fmt.Println("proto repository:", absProtoRepository)
		fmt.Println("output proto dir:", outputDir)
	}

	// 实际逻辑
	// 从golang源代码中找到所有protobuf类型
	golangTypeNames, err := newGolangParser().parse(arg.GolangSourceCodeDir, pbImportPath, absOutputDir)
	if err != nil {
		return "", err
	}

	if len(golangTypeNames) == 0 {
		return "", fmt.Errorf("golang protobuf type reference not found")
	}

	// 去匹配proto文件中和源文件中匹配的类型的声明
	protoDeclares, err := newProtoParser().findProtoDeclares(absProtoRepository, arg.OutputPackage, golangTypeNames)
	if err != nil {
		return "", err
	}

	if len(protoDeclares) == 0 {
		return "", fmt.Errorf("empty matched protobuf types, protoDir: %s", absProtoRepository)
	}

	// 生成新的proto文件以及其依赖文件
	outputFilename := fmt.Sprintf("%s.proto", filepath.Base(arg.GolangModule))
	err = newProtoRender().Render(absProtoRepository, outputDir, outputFilename, arg.OutputPackage, protoDeclares)
	if err != nil {
		return "", err
	}

	return outputDir, nil
}
