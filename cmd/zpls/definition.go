package main

import (
	"path/filepath"
	"strings"

	"github.com/diiyw/z/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func onDefinitionfunc(context *glsp.Context, params *protocol.DefinitionParams) (any, error) {
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
