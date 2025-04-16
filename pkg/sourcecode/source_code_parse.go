package sourcecode

//type sourceCodeParser struct {
//	pkgName2astPkg map[string]*ast.Package
//}
//
//func newParser() *sourceCodeParser {
//	return &sourceCodeParser{
//		//sc: &sourceCodeInfo{
//		//	DaprModules:            make([]*daprModuleInfo, 0),
//		//	DaprInvocationHandlers: make([]*daprInvocationHandler, 0),
//		//},
//	}
//}
//
//// Parse 从源代码中解析源代码信息
//// 1. dapr modules
//// 2. server.Start的入口文件
//// 3. route annotations
//func (p *sourceCodeParser) Parse(srcDir string, skipDirs []string) (*sourceCodeInfo, error) {
//	if err := p.preprocess(srcDir, skipDirs); err != nil {
//		return nil, err
//	}
//
//	daprModules :=
//
//	if err := p.parseDaprInvocationHandlers(srcDir); err != nil {
//		return nil, err
//	}
//
//	return p.sc, nil
//}
