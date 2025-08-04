package main

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	// Z语言关键字
	keywords = []string{
		"if",
		"else",
		"for",
		"return",
		"func",
		"export",
		"in",
	}

	// Z语言内置函数
	builtinFunctions = []string{
		"to_json",
		"from_json",
		"print",
		"printf",
		"sprintf",
		"format",
		"len",
		"copy",
		"append",
		"delete",
		"splice",
		"type_name",
		"int",
		"bool",
		"float",
		"char",
		"bytes",
		"error",
		"string",
		"time",
		"is_string",
		"is_bool",
		"is_float",
		"is_char",
		"is_bytes",
		"is_error",
		"is_undefined",
		"is_function",
		"is_callable",
		"is_array",
		"is_immutable_array",
		"is_map",
		"is_iterable",
		"is_time",
	}

	// 常量
	constants = []string{"true", "false", "undefined"}
)

func onCompletionFunc(context *glsp.Context, params *protocol.CompletionParams) (any, error) {
	keywordKind := protocol.CompletionItemKindKeyword
	// 简单判断触发补全的上下文
	var items []protocol.CompletionItem
	for _, keyword := range keywords {
		items = append(items, protocol.CompletionItem{
			Label: keyword,
			Kind:  &keywordKind,
			Data:  "keyword-" + keyword,
		})
	}
	funcKind := protocol.CompletionItemKindFunction
	snippetFormat := protocol.InsertTextFormatSnippet
	for _, f := range builtinFunctions {
		text := f + `($1)`
		items = append(items, protocol.CompletionItem{
			Label:            f,
			Kind:             &funcKind,
			InsertText:       &text,
			InsertTextFormat: &snippetFormat,
			Data:             "function-" + f,
		})
	}
	constKind := protocol.CompletionItemKindConstant
	for _, constant := range constants {
		items = append(items, protocol.CompletionItem{
			Label: constant,
			Kind:  &constKind,
			Data:  "constant-" + constant,
		})
	}
	return items, nil
}

func onCompletionResolveFunc(context *glsp.Context, params *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	return params, nil
}
