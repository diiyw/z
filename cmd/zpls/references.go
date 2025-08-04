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
	var expr = findNode(parsedFile.Stmts, offset)
	var locations = make([]protocol.Location, 0)
	if _, ok := expr.(*parser.SelectorExpr); ok {
		locations = findReferences(filename, content, expr.(*parser.SelectorExpr), parsedFile.Stmts)
	}
	return locations, nil
}

func findReferences(filename, content string, expr *parser.SelectorExpr, stmt []parser.Stmt) []protocol.Location {
	locations := make([]protocol.Location, 0)
	return locations
}
