package main

import (
	"errors"
	"os"
	"strings"

	"github.com/diiyw/z/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func onContentChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	diagnostics := make([]protocol.Diagnostic, 0)
	filename := strings.ReplaceAll(params.TextDocument.URI, "file://", "")
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	fileSet := parser.NewFileSet()
	sourceFile := fileSet.AddFile("diagnostics.z", -1, len(content))
	p := parser.NewParser(sourceFile, content, nil)
	_, err = p.ParseFile()
	var errSeverity = protocol.DiagnosticSeverityError
	sourceZ := "z"
	if err != nil {
		var errList parser.ErrorList
		errors.As(err, &errList)
		for _, el := range errList {
			diagnostics = append(diagnostics, protocol.Diagnostic{
				Severity: &errSeverity,
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      protocol.UInteger(el.Pos.Line - 1),
						Character: protocol.UInteger(el.Pos.Column - 1),
					},
					End: protocol.Position{
						Line:      protocol.UInteger(el.Pos.Line - 1),
						Character: protocol.UInteger(el.Pos.Column),
					},
				},
				Message: el.Error(),
				Source:  &sourceZ,
			})
		}
		if len(diagnostics) > 0 {
			// 发送诊断
			context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
				URI:         params.TextDocument.URI,
				Diagnostics: diagnostics,
			})
		}
	}
	return nil
}
