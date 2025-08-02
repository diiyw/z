package main

import (
	"os"
	"strings"

	"github.com/diiyw/z/cmd/format"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func onFormattingFunc(context *glsp.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	filename := strings.ReplaceAll(params.TextDocument.URI, "file://", "")
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	result, err := format.Format(content)
	if err != nil {
		return nil, err
	}
	// 替换整个文档内容
	return []protocol.TextEdit{
		{
			Range: protocol.Range{
				Start: protocol.Position{},
				End:   protocol.Position{}.EndOfLineIn(string(content)),
			},
			NewText: result,
		},
	}, nil
}
