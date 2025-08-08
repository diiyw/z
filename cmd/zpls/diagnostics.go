package main

import (
	"errors"
	"github.com/diiyw/z/cmd/zpls/file"
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
	content := file.Document().GetText(uri)

	// 使用新的通用文件解析函数
	_, err := ParseFileContent(filename, content)
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
