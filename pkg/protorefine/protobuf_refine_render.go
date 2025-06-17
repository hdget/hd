package protorefine

import (
	"bytes"
	"github.com/elliotchance/pie/v2"
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

type protoRender struct {
}

const (
	templateProtoFile = `
syntax = "proto3";

package {{.Package}};

{{range .Imports}}
import "{{.}}";
{{end}}

{{range .Types}}
{{.}}
{{end}}
`
)

type protoFileOutput struct {
	Package string
	Imports []string
	Types   []string
}

func newProtoRender() *protoRender {
	return &protoRender{}
}

func (r *protoRender) Render(protoDir, outputDir, outputFileName, golangPkgName string, declares []*protoDeclare) error {
	err := os.MkdirAll(outputDir, 0644)
	if err != nil {
		return errors.Wrapf(err, "make output dir, dir: %s", outputDir)
	}

	importFiles, err := r.renderFile(outputDir, outputFileName, golangPkgName, declares)
	if err != nil {
		return err
	}

	for _, f := range importFiles {
		err = r.copyFile(filepath.Join(protoDir, f), filepath.Join(outputDir, f))
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *protoRender) renderFile(outputDir, outputFileName string, golangPkgName string, declares []*protoDeclare) ([]string, error) {
	pbPrinter := &protoprint.Printer{
		Compact:                              true,
		ShortOptionsExpansionThresholdLength: 200,
	}

	importMap := make(map[string]interface{})
	types := make([]string, 0)
	for _, declare := range declares {
		for _, dep := range declare.dependents {
			if _, exists := importMap[dep]; !exists {
				importMap[dep] = struct{}{}
			}
		}

		t, err := pbPrinter.PrintProtoToString(declare.descriptor)
		if err != nil {
			return nil, err
		}

		types = append(types, t)
	}

	imports := pie.Keys(importMap)

	output := &protoFileOutput{
		Package: golangPkgName,
		Imports: imports,
		Types:   types,
	}

	t := template.Must(template.New("proto").Parse(templateProtoFile))
	var buf bytes.Buffer
	if err := t.Execute(&buf, output); err != nil {
		return nil, err
	}

	outputPath := filepath.Join(outputDir, outputFileName)
	err := os.WriteFile(outputPath, buf.Bytes(), 0644)
	if err != nil {
		return nil, errors.Wrapf(err, "write proto file, path: %s", outputPath)
	}
	return imports, nil
}

func (r *protoRender) copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = srcFile.Close()
	}()

	err = os.MkdirAll(filepath.Dir(dest), 0644)
	if err != nil {
		return err
	}

	destFile, err := os.Create(dest) // creates if file doesn't exist
	if err != nil {
		return err
	}
	defer func() {
		_ = destFile.Close()
	}()

	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}
	return nil
}
