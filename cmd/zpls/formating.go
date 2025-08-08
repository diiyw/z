package main

import (
	"github.com/diiyw/z/cmd/format"
	"github.com/diiyw/z/cmd/zpls/file"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func onFormattingFunc(context *glsp.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	content := file.Document().GetText(params.TextDocument.URI)
	result, err := format.Format([]byte(content))
	if err != nil {
		return nil, err
	}
	// 替换整个文档内容
	return []protocol.TextEdit{
		{
			Range:   fullFileRange(content),
			NewText: result,
		},
	}, nil
}
