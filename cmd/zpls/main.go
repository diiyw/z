package main

import (
	"github.com/diiyw/z/cmd/zpls/file"
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	// Must include a backend implementation
	// See CommonLog for other options: https://github.com/tliron/commonlog
	_ "github.com/tliron/commonlog/simple"
)

const lsName = "Z language"

var (
	version = "1.0.0"
	handler protocol.Handler
)

func main() {
	// This increases logging verbosity (optional)
	commonlog.Configure(1, nil)
	handler = protocol.Handler{
		Initialize: initialize,
		Initialized: func(context *glsp.Context, params *protocol.InitializedParams) error {
			return nil
		},
		Shutdown: func(context *glsp.Context) error {
			protocol.SetTraceValue(protocol.TraceValueOff)
			return nil
		},
		SetTrace: func(context *glsp.Context, params *protocol.SetTraceParams) error {
			protocol.SetTraceValue(params.Value)
			return nil
		},
		TextDocumentDidOpen: func(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
			file.Document().Set(params.TextDocument.URI, params.TextDocument)
			return nil
		},
		TextDocumentDidClose: func(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
			file.Document().Delete(params.TextDocument.URI)
			return nil
		},
		TextDocumentDidChange: func(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
			uri := params.TextDocument.URI
			// 合并所有内容变更
			text := getTextFromVersion(uri, params.ContentChanges)
			doc := protocol.TextDocumentItem{
				URI:        params.TextDocument.URI,
				LanguageID: params.TextDocument.TextDocumentIdentifier.URI,
				Version:    params.TextDocument.Version,
				Text:       text,
			}
			file.Document().Set(uri, doc)
			return onTextDocumentChange(context, params)
		},
		TextDocumentCompletion:         onCompletionFunc,
		TextDocumentDefinition:         onDefinitionFunc,
		TextDocumentFormatting:         onFormattingFunc,
		WorkspaceDidChangeWatchedFiles: onWorkspaceDidChangeWatchedFiles,
		TextDocumentReferences:         onReferencesFunc,
	}

	lspServer := server.NewServer(&handler, lsName, false)

	if err := lspServer.RunTCP(":60066"); err != nil {
		panic(err)
	}
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}
