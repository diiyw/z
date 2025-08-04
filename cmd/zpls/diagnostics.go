package main

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/diiyw/z/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func onTextDocumentChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	return onDiagnostic(context, params.TextDocument.URI)
}

func onDiagnostic(context *glsp.Context, uri string) error {
	filename := strings.ReplaceAll(uri, "file://", "")
	content := Document().GetText(uri)
	var name = filepath.Base(filename)
	fileSet := parser.NewFileSet()
	sourceFile := fileSet.AddFile(name, -1, len(content))
	p := parser.NewParser(sourceFile, []byte(content), nil)
	_, err := p.ParseFile()
	var errSeverity = protocol.DiagnosticSeverityError
	sourceZ := "z"
	diagnostics := make([]protocol.Diagnostic, 0)
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
	}
	// 发送诊断
	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})
	return nil
}

func onWorkspaceDidChangeWatchedFiles(context *glsp.Context, params *protocol.DidChangeWatchedFilesParams) error {
	for _, change := range params.Changes {
		if change.Type == protocol.FileChangeTypeDeleted {
			continue
		}
		if err := onDiagnostic(context, change.URI); err != nil {
			return err
		}
	}
	return nil
}
