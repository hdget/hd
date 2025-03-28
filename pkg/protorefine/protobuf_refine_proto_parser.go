package protorefine

import (
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

type protoParser struct {
}

type protoDeclareKind int

const (
	protoDeclareKindUnknown protoDeclareKind = iota
	protoDeclareKindEnum
	protoDeclareKindMessage
)

type protoDeclare struct {
	kind       protoDeclareKind
	descriptor desc.Descriptor
	dependents []string
	//origin     string // 原来的proto文件名
}

type protoParseResult struct {
	typeName2protoMessage map[string]*desc.MessageDescriptor
	typeName2protoEnum    map[string]*desc.EnumDescriptor
	dependents            []string
}

func newProtoParser() *protoParser {
	return &protoParser{}
}

// findProtoDeclares 找到和源代码中的golang类型对应的protobuf类型的元数据信息
func (p *protoParser) findProtoDeclares(protoDir string, golangPkgName string, golangTypeNames []string) ([]*protoDeclare, error) {
	file2parseResult, err := p.parse(protoDir)
	if err != nil {
		return nil, err
	}

	// filter duplicate pbTypeName corresponding descriptors
	descriptorName2protoDeclare := make(map[string]*protoDeclare)
	for _, golangTypeName := range golangTypeNames {
		// try to match protobuf type in proto files by golang type
		declares, found := p.match(file2parseResult, golangPkgName, golangTypeName)
		if !found {
			// protoc_gen may generate the protobuf type in source code with different naming convention
			// e,g: it capitalized the first character of the word
			// we should check this case, pb.TheWord ==> message theWord {}
			declares, found = p.match(file2parseResult, golangPkgName, p.lowerCaseFirstLetter(golangTypeName))
			if !found {
				return nil, fmt.Errorf("golang type: %s.%s not found", golangPkgName, golangTypeName)
			}
		}

		// filter the duplicate protoDeclares
		for _, d := range declares {
			descriptorName2protoDeclare[d.descriptor.GetName()] = d
		}
	}

	return pie.Values(descriptorName2protoDeclare), nil
}

// match split the name to words and recursive from end to begin to find if golang protobuf type can be found in proto files
func (p *protoParser) match(file2parseResult map[string]*protoParseResult, golangPkgName, golangTypeName string) ([]*protoDeclare, bool) {
	words := p.splitWords(golangTypeName)
	for i := len(words); i > 0; i-- {
		typeName := strings.Join(words[:i], "")
		for _, parseResult := range file2parseResult {
			if descriptor, exist := parseResult.typeName2protoMessage[typeName]; exist {
				// recursively search descendant message types and enums
				matcher := newProtobufTypeMatcher()
				matcher.traverse(golangPkgName, descriptor)

				return p.createProtoDeclares(protoDeclareKindMessage, parseResult.dependents, pie.Values(matcher.founds)...), true
			}

			if descriptor, exist := parseResult.typeName2protoEnum[typeName]; exist {
				return p.createProtoDeclares(protoDeclareKindEnum, parseResult.dependents, descriptor), true
			}
		}
	}

	return nil, false
}

// parseProtoMessageAndEnums 从proto目录中找到所有protobuf message和enum
// 返回proto文件名=>protoParseResult
func (p *protoParser) parse(protoDir string) (map[string]*protoParseResult, error) {
	files, err := os.ReadDir(protoDir)
	if err != nil {
		return nil, err
	}

	protoFiles := make([]string, 0)
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".proto") {
			protoFiles = append(protoFiles, f.Name())
		}
	}

	//compiler := protocompile.Compiler{
	//	Resolver: &protocompile.SourceResolver{
	//		ImportPaths: []string{filepath.Dir(protoDir)},
	//	},
	//}

	pbParser := &protoparse.Parser{
		ImportPaths: []string{protoDir},
	}

	fds, err := pbParser.ParseFiles(protoFiles...)
	if err != nil {
		return nil, err
	}

	//fds, err := compiler.Compile(context.Background(), protoFiles...)
	//if err != nil {
	//	return nil, err
	//}

	var totalTypes int
	results := make(map[string]*protoParseResult)
	for _, fd := range fds {
		r := &protoParseResult{
			typeName2protoMessage: make(map[string]*desc.MessageDescriptor),
			typeName2protoEnum:    make(map[string]*desc.EnumDescriptor),
			dependents:            make([]string, 0),
		}

		for _, t := range fd.GetMessageTypes() {
			r.typeName2protoMessage[t.GetName()] = t
			totalTypes += 1
		}

		for _, t := range fd.GetEnumTypes() {
			r.typeName2protoEnum[t.GetName()] = t
			totalTypes += 1
		}

		for _, dep := range fd.GetDependencies() {
			// ignore dependent files under the same dir
			if filepath.Dir(fd.GetName()) != filepath.Dir(dep.GetName()) {
				err = r.recursiveFindDependents(protoDir, dep.GetName())
				if err != nil {
					return nil, err
				}
			}
		}

		results[fd.GetName()] = r
	}

	if totalTypes == 0 {
		return nil, fmt.Errorf("no protobuf descriptors found")
	}

	return results, nil
}

func (ret *protoParseResult) recursiveFindDependents(protoDir, f string) error {
	ret.dependents = append(ret.dependents, f)

	pbParser := &protoparse.Parser{
		ImportPaths: []string{protoDir},
	}

	fds, err := pbParser.ParseFiles(f)
	if err != nil {
		return err
	}

	for _, fd := range fds {
		for _, dep := range fd.GetDependencies() {
			depPath := dep.GetName()
			if err = ret.recursiveFindDependents(protoDir, depPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *protoParser) createProtoDeclares(kind protoDeclareKind, dependents []string, descriptors ...desc.Descriptor) []*protoDeclare {
	declares := make([]*protoDeclare, len(descriptors))
	for i, d := range descriptors {
		declares[i] = &protoDeclare{
			kind:       kind,
			descriptor: d,
			dependents: dependents,
		}
	}
	return declares
}

func (p *protoParser) splitWords(input string) []string {
	var words []string
	if strings.Contains(input, "_") {
		// if name is snake case (例如：my_variable_name)
		words = strings.Split(input, "_")
	} else {
		reg := regexp.MustCompile(`([a-z])([A-Z])`)
		spaceSeparated := reg.ReplaceAllString(input, `${1} ${2}`)

		// 使用 strings.Fields 分隔成单词
		words = strings.Fields(spaceSeparated)
	}
	return words
}

func (p *protoParser) lowerCaseFirstLetter(s string) string {
	ss := []rune(s)
	ss[0] = unicode.ToLower(ss[0])
	return string(ss)
}
