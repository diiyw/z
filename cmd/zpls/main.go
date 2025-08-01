package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/diiyw/z/cmd/format"
	"github.com/diiyw/z/parser"
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
	version string = "1.0.0"
	handler protocol.Handler
)

func main() {
	// This increases logging verbosity (optional)
	commonlog.Configure(1, nil)

	handler = protocol.Handler{
		Initialize:             initialize,
		Initialized:            initialized,
		Shutdown:               shutdown,
		SetTrace:               setTrace,
		TextDocumentCompletion: onCompletionfunc,
		CompletionItemResolve:  onCompletionResolveFunc,
		TextDocumentDefinition: onDefinitionfunc,
		TextDocumentFormatting: onFormattingFunc,
		TextDocumentDidChange:  onContentChange,
	}

	server := server.NewServer(&handler, lsName, false)

	server.RunStdio()
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

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

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

func onFormattingFunc(context *glsp.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	filename := strings.Replace(params.TextDocument.URI, "file://", "", -1)
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

func onCompletionfunc(context *glsp.Context, params *protocol.CompletionParams) (any, error) {
	keywordKind := protocol.CompletionItemKindKeyword
	// 简单判断触发补全的上下文
	items := []protocol.CompletionItem{}
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

func onContentChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	diagnostics := make([]protocol.Diagnostic, 0)
	filename := strings.Replace(params.TextDocument.URI, "file://", "", -1)
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
			// TODO 发送诊断
		}
	}
	return nil
}

func onDefinitionfunc(context *glsp.Context, params *protocol.DefinitionParams) (any, error) {
	filename := strings.Replace(params.TextDocument.URI, "file://", "", -1)
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	offset := params.Position.IndexIn(string(content))
	fileSet := parser.NewFileSet()
	sourceFile := fileSet.AddFile("definition.z", -1, len(content))
	p := parser.NewParser(sourceFile, content, nil)
	parsedFile, err := p.ParseFile()
	if err != nil {
		return nil, err
	}
	var expr = findNode(parsedFile.Stmts, offset)
	currentDir := filepath.Dir(filename)
	if e, ok := expr.(*parser.ImportExpr); ok {
		return protocol.LocationLink{
			OriginSelectionRange: &protocol.Range{
				Start: protocol.Position{
					Line:      params.Position.Line,
					Character: protocol.UInteger(e.Pos() + 7),
				},
				End: protocol.Position{
					Line:      params.Position.Line,
					Character: protocol.UInteger(e.Pos() + 3),
				},
			},
			TargetURI: protocol.URI("file://" + filepath.Join(currentDir, e.ModuleName+".z")),
			TargetRange: protocol.Range{
				Start: protocol.Position{},
				End:   protocol.Position{},
			},
			TargetSelectionRange: protocol.Range{
				Start: protocol.Position{},
				End:   protocol.Position{},
			},
		}, nil
	}
	return nil, nil
}

func findNode(stmts []parser.Stmt, offset int) parser.Expr {
	for _, stmt := range stmts {
		if node := findExprNode(stmt, offset); node != nil {
			return node
		}
	}
	return nil
}

func findExprNode(stmt parser.Stmt, start int) parser.Expr {
	ss, se := stmt.Pos(), stmt.End()
	if int(ss) <= start && start < int(se) {
		switch stmt := stmt.(type) {
		case *parser.AssignStmt:
			for _, expr := range stmt.LHS {
				if int(expr.Pos()) <= start && start < int(expr.End()) {
					return expr
				}
			}
			for _, expr := range stmt.RHS {
				if int(expr.Pos()) <= start && start < int(expr.End()) {
					return expr
				}
			}
		case *parser.ExportStmt:
			expr := stmt.Result
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
		case *parser.BlockStmt:
			for _, expr := range stmt.Stmts {
				if int(expr.Pos()) <= start && start < int(expr.End()) {
					return findExprNode(expr, start)
				}
			}
		case *parser.ExprStmt:
			expr := stmt.Expr
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
		case *parser.ForInStmt:
			expr := stmt.Iterable
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
			expr = stmt.Key
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
			expr = stmt.Value
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
			return findExprNode(stmt.Body, start)
		case *parser.ForStmt:
			expr := stmt.Cond
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
			if node := findExprNode(stmt.Init, start); node != nil {
				return node
			}
			if node := findExprNode(stmt.Post, start); node != nil {
				return node
			}
			return findExprNode(stmt.Body, start)
		case *parser.IfStmt:
			expr := stmt.Cond
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
			if node := findExprNode(stmt.Init, start); node != nil {
				return node
			}
			if node := findExprNode(stmt.Else, start); node != nil {
				return node
			}
			return findExprNode(stmt.Body, start)
		case *parser.IncDecStmt:
			expr := stmt.Expr
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
		case *parser.ReturnStmt:
			expr := stmt.Result
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
		}
	}
	return nil
}
