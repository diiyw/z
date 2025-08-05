package main

import (
	"path/filepath"
	"strings"

	"github.com/diiyw/z/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func onReferencesFunc(context *glsp.Context, params *protocol.ReferenceParams) ([]protocol.Location, error) {
	filename := strings.ReplaceAll(params.TextDocument.URI, "file://", "")
	content := Document().GetText(params.TextDocument.URI)
	offset := params.Position.IndexIn(content)
	fileSet := parser.NewFileSet()
	basename := filepath.Base(filename)
	sourceFile := fileSet.AddFile(basename, -1, len(content))
	p := parser.NewParser(sourceFile, []byte(content), nil)
	parsedFile, err := p.ParseFile()
	if err != nil {
		return nil, err
	}
	currentDir := filepath.Dir(filename)
	var f = finder{currentDir: currentDir, offset: offset, filename: filename}
	var expr = f.findNodeByOffset(parsedFile.Stmts)
	var locations = make([]protocol.Location, 0)
	
	// 处理标识符的引用查找
	if ident, ok := expr.(*parser.Ident); ok {
		locations = f.findReferences(ident, parsedFile)
	}
	
	// 处理SelectorExpr，例如 module.variable 形式
	if selExpr, ok := expr.(*parser.SelectorExpr); ok {
		if sel, ok := selExpr.Sel.(*parser.Ident); ok {
			locations = f.findSelectorReferences(selExpr, sel, parsedFile)
		}
	}
	
	return locations, nil
}

