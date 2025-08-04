package main

import (
	"path/filepath"
	"strings"

	"github.com/diiyw/z/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func onDefinitionFunc(context *glsp.Context, params *protocol.DefinitionParams) (any, error) {
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
	var locationLinks = make([]protocol.LocationLink, 0)
	if e, ok := expr.(*parser.ImportExpr); ok {
		locationLinks = append(locationLinks, protocol.LocationLink{
			OriginSelectionRange: &protocol.Range{
				Start: offsetToPosition(int(e.Pos()+7), content),
				End:   offsetToPosition(int(e.End()-3), content),
			},
			TargetURI:            creatURI(filepath.Join(currentDir, e.ModuleName+".z")),
			TargetRange:          createRange(0, 0, 0, 0),
			TargetSelectionRange: createRange(0, 0, 0, 0),
		})
	}
	if e, ok := expr.(*parser.Ident); ok {
		locationLinks = findDefinition(filename, content, e, parsedFile.Stmts)
	}
	return locationLinks, nil
}

func findDefinition(filename, content string, expr *parser.Ident, stmt []parser.Stmt) []protocol.LocationLink {
	locations := make([]protocol.LocationLink, 0)
	for _, s := range stmt {
		if assignStmt, ok := s.(*parser.AssignStmt); ok {
			for _, lh := range assignStmt.LHS {
				if lh.String() == expr.Name {
					locations = append(locations, protocol.LocationLink{
						OriginSelectionRange: &protocol.Range{
							Start: offsetToPosition(int(expr.Pos()-1), content),
							End:   offsetToPosition(int(expr.End()-1), content),
						},
						TargetRange: protocol.Range{
							Start: offsetToPosition(int(lh.Pos()-1), content),
							End:   offsetToPosition(int(lh.End()-1), content),
						},
						TargetSelectionRange: protocol.Range{
							Start: offsetToPosition(int(lh.Pos()-1), content),
							End:   offsetToPosition(int(lh.End()-1), content),
						},
						TargetURI: creatURI(filename), // Adjust the path as needed
					})
				}
			}
		}
	}
	return locations
}

func findNode(stmts []parser.Stmt, offset int) parser.Expr {
	for _, stmt := range stmts {
		if node := findStmtNode(stmt, offset); node != nil {
			return node
		}
	}
	return nil
}

func findStmtNode(stmt parser.Stmt, start int) parser.Expr {
	ss, se := stmt.Pos(), stmt.End()
	if int(ss) <= start && start < int(se) {
		switch s := stmt.(type) {
		case *parser.AssignStmt:
			for _, expr := range s.LHS {
				if e := findExprNode(expr, start); e != nil {
					return e
				}
			}
			for _, expr := range s.RHS {
				if e := findExprNode(expr, start); e != nil {
					return e
				}
				if int(expr.Pos()) <= start && start < int(expr.End()) {
					return expr
				}
			}
		case *parser.ExportStmt:
			return findExprNode(s.Result, start)
		case *parser.BlockStmt:
			for _, st := range s.Stmts {
				if int(st.Pos()) <= start && start < int(st.End()) {
					return findStmtNode(st, start)
				}
			}
		case *parser.ExprStmt:
			return findExprNode(s.Expr, start)
		case *parser.ForInStmt:
			if e := findExprNode(s.Iterable, start); e != nil {
				return e
			}
			if e := findExprNode(s.Key, start); e != nil {
				return e
			}
			if e := findExprNode(s.Value, start); e != nil {
				return e
			}
			return findStmtNode(s.Body, start)
		case *parser.ForStmt:
			if e := findExprNode(s.Cond, start); e != nil {
				return e
			}
			if node := findStmtNode(s.Init, start); node != nil {
				return node
			}
			if node := findStmtNode(s.Post, start); node != nil {
				return node
			}
			return findStmtNode(s.Body, start)
		case *parser.IfStmt:
			if e := findExprNode(s.Cond, start); e != nil {
				return e
			}
			if e := findStmtNode(s.Init, start); e != nil {
				return e
			}
			if e := findStmtNode(s.Else, start); e != nil {
				return e
			}
			return findStmtNode(s.Body, start)
		case *parser.IncDecStmt:
			if e := findExprNode(s.Expr, start); e != nil {
				return e
			}
		case *parser.ReturnStmt:
			if e := findExprNode(s.Result, start); e != nil {
				return e
			}
		}
	}
	return nil
}

func findExprNode(expr parser.Expr, start int) parser.Expr {
	if start > int(expr.End()) {
		return nil
	}
	switch e := expr.(type) {
	case *parser.ArrayLit:
		for _, element := range e.Elements {
			if e := findExprNode(element, start); e != nil {
				return e
			}
		}
	case *parser.BadExpr:
		return nil
	case *parser.BinaryExpr:
		if int(e.LHS.Pos()) <= start && start < int(e.LHS.End()) {
			return findExprNode(e.LHS, start)
		}
		if int(e.RHS.Pos()) <= start && start < int(e.RHS.End()) {
			return findExprNode(e.RHS, start)
		}
	case *parser.BoolLit:
		return e
	case *parser.CallExpr:
		if int(e.Func.Pos()) <= start && start < int(e.Func.End()) {
			return findExprNode(e.Func, start)
		}
		for _, a := range e.Args {
			if e := findExprNode(a, start); e != nil {
				return e
			}
		}
	case *parser.CharLit:
		return e
	case *parser.CondExpr:
		if int(e.Cond.Pos()) <= start && start < int(e.Cond.End()) {
			return findExprNode(e.Cond, start)
		}
		if int(e.False.Pos()) <= start && start < int(e.False.End()) {
			return findExprNode(e.False, start)
		}
		if int(e.True.Pos()) <= start && start < int(e.True.End()) {
			return findExprNode(e.True, start)
		}
	case *parser.ErrorExpr:
		return e
	case *parser.FloatLit:
		return e
	case *parser.FuncLit:
		if int(e.Type.Pos()) <= start && start < int(e.Type.End()) {
			return findExprNode(e.Type, start)
		}
		return findStmtNode(e.Body, start)
	case *parser.FuncType:
		for _, a := range e.Params.List {
			if e := findExprNode(a, start); e != nil {
				return e
			}
		}
		return e
	case *parser.ImmutableExpr:
		return e
	case *parser.ImportExpr:
		return e
	case *parser.IndexExpr:
		if int(e.Expr.Pos()) <= start && start < int(e.Expr.End()) {
			return findExprNode(e.Expr, start)
		}
		if int(e.Index.Pos()) <= start && start < int(e.Index.End()) {
			return findExprNode(e.Index, start)
		}
	case *parser.Ident:
		if int(e.Pos()) <= start && start < int(e.End()) {
			return e
		}
	case *parser.IntLit:
		return e
	case *parser.MapElementLit:
		if int(e.Key.Pos()) <= start && start < int(e.Key.End()) {
			return findExprNode(e.Key, start)
		}
		if int(e.Value.Pos()) <= start && start < int(e.Value.End()) {
			return findExprNode(e.Value, start)
		}
	case *parser.MapLit:
		for _, element := range e.Elements {
			if e := findExprNode(element, start); e != nil {
				return e
			}
		}
	case *parser.ParenExpr:
		if int(e.Expr.Pos()) <= start && start < int(e.Expr.End()) {
			return findExprNode(e.Expr, start)
		}
	case *parser.SelectorExpr:
		if int(e.Expr.Pos()) <= start && start < int(e.Expr.End()) {
			return findExprNode(e.Expr, start)
		}
		if int(e.Sel.Pos()) <= start && start < int(e.Sel.End()) {
			return findExprNode(e.Sel, start)
		}
	case *parser.SliceExpr:
		if int(e.Expr.Pos()) <= start && start < int(e.Expr.End()) {
			return findExprNode(e.Expr, start)
		}
		if int(e.Low.Pos()) <= start && start < int(e.Low.End()) {
			return findExprNode(e.Low, start)
		}
		if int(e.High.Pos()) <= start && start < int(e.High.End()) {
			return findExprNode(e.High, start)
		}
	case *parser.StringLit:
		return e
	case *parser.UnaryExpr:
		if int(e.Expr.Pos()) <= start && start < int(e.Expr.End()) {
			return findExprNode(e.Expr, start)
		}
	case *parser.UndefinedLit:
		return e
	}
	return nil
}
