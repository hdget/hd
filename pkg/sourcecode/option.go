package sourcecode

type Option func(*sourceCodeHandlerImpl)

// // WithEntrySignature 定义服务调用函数签名， 方便定位哪个文件需要patch来添加模块所在包的导入路径
//
//	func WithParseSignature(signature string) Option {
//		return func(m *sourceCodeHandlerImpl) error {
//			if signature == "" {
//				return errors.New("patch file signature not specified")
//			}
//
//			tokens := strings.Split(signature, ".")
//			if len(tokens) != 2 {
//				return errors.New("invalid signature format, format: caller.functionChain")
//			}
//
//			importPath, functionChain := strings.TrimSpace(tokens[0]), strings.TrimSpace(tokens[1])
//			if importPath == "" || functionChain == "" {
//				return errors.New("empty import path or function chain")
//			}
//
//			m.signature = &callSignature{
//				caller: importPath,
//				functionChain: functionChain,
//			}
//			return nil
//		}
//	}
func WithSkipDirs(dirs ...string) Option {
	return func(m *sourceCodeHandlerImpl) {
		m.skipDirs = append(m.skipDirs, dirs...)
	}
}

func WithAssetPath(assetsPath string) Option {
	return func(m *sourceCodeHandlerImpl) {
		m.assetsPath = assetsPath
	}
}
