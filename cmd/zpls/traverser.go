package main

import (
	"github.com/diiyw/z/parser"
)

// TraverserHandler 定义了遍历AST时需要处理的接口
type TraverserHandler interface {
	HandleAssignStmt(*parser.AssignStmt, *scopeStack)
	HandleBlockStmt(*parser.BlockStmt, *scopeStack)
	HandleFuncLit(*parser.FuncLit, *scopeStack)
	HandleForInStmt(*parser.ForInStmt, *scopeStack)
	HandleIdent(*parser.Ident, *scopeStack)
	HandleSelectorExpr(*parser.SelectorExpr, *scopeStack)
	HandleImportExpr(*parser.ImportExpr, *scopeStack)
	HandleExportStmt(*parser.ExportStmt, *scopeStack)
}

// Traverser 用于遍历AST
type Traverser struct {
	handler TraverserHandler
}

// BaseTraverserHandler 提供TraverserHandler接口的空实现
type BaseTraverserHandler struct{}

// HandleAssignStmt 提供空实现
func (b *BaseTraverserHandler) HandleAssignStmt(stmt *parser.AssignStmt, scope *scopeStack) {}

// HandleBlockStmt 提供空实现
func (b *BaseTraverserHandler) HandleBlockStmt(stmt *parser.BlockStmt, scope *scopeStack) {}

// HandleFuncLit 提供空实现
func (b *BaseTraverserHandler) HandleFuncLit(lit *parser.FuncLit, scope *scopeStack) {}

// HandleForInStmt 提供空实现
func (b *BaseTraverserHandler) HandleForInStmt(stmt *parser.ForInStmt, scope *scopeStack) {}

// HandleIdent 提供空实现
func (b *BaseTraverserHandler) HandleIdent(ident *parser.Ident, scope *scopeStack) {}

// HandleSelectorExpr 提供空实现
func (b *BaseTraverserHandler) HandleSelectorExpr(expr *parser.SelectorExpr, scope *scopeStack) {}

// HandleImportExpr 提供空实现
func (b *BaseTraverserHandler) HandleImportExpr(expr *parser.ImportExpr, scope *scopeStack) {}

// HandleExportStmt 提供空实现
func (b *BaseTraverserHandler) HandleExportStmt(stmt *parser.ExportStmt, scope *scopeStack) {}

// NewTraverser 创建一个新的Traverser
func NewTraverser(handler TraverserHandler) *Traverser {
	return &Traverser{
		handler: handler,
	}
}

// TraverseStmts 遍历语句列表
func (t *Traverser) TraverseStmts(stmts []parser.Stmt, scope *scopeStack) {
	for _, stmt := range stmts {
		t.TraverseStmt(stmt, scope)
	}
}

// TraverseStmt 遍历单个语句
func (t *Traverser) TraverseStmt(stmt parser.Stmt, scope *scopeStack) {
	switch s := stmt.(type) {
	case *parser.AssignStmt:
		t.handler.HandleAssignStmt(s, scope)
	case *parser.BlockStmt:
		t.handler.HandleBlockStmt(s, scope)
		scope.pushScope()
		t.TraverseStmts(s.Stmts, scope)
		scope.popScope()
	case *parser.ExprStmt:
		// 检查是否是函数字面量
		if funcLit, ok := s.Expr.(*parser.FuncLit); ok {
			t.handler.HandleFuncLit(funcLit, scope)
			scope.pushScope()
			// 添加参数到作用域
			if funcLit.Type.Params != nil {
				for _, param := range funcLit.Type.Params.List {
					t.handler.HandleIdent(param, scope)
				}
			}
			// 在函数体中递归查找
			if funcLit.Body != nil {
				t.TraverseStmts(funcLit.Body.Stmts, scope)
			}
			// 弹出作用域
			scope.popScope()
		} else if importExpr, ok := s.Expr.(*parser.ImportExpr); ok {
			t.handler.HandleImportExpr(importExpr, scope)
		}
	case *parser.ForInStmt:
		t.handler.HandleForInStmt(s, scope)
		scope.pushScope()
		t.handler.HandleIdent(s.Key, scope)
		t.handler.HandleIdent(s.Value, scope)
		t.TraverseExpr(s.Iterable, scope)
		t.TraverseStmt(s.Body, scope)
		scope.popScope()
	case *parser.ForStmt:
		t.TraverseStmt(s.Init, scope)
		t.TraverseExpr(s.Cond, scope)
		t.TraverseStmt(s.Post, scope)
		t.TraverseStmt(s.Body, scope)
	case *parser.IfStmt:
		t.TraverseStmt(s.Init, scope)
		t.TraverseExpr(s.Cond, scope)
		t.TraverseStmt(s.Body, scope)
		t.TraverseStmt(s.Else, scope)
	case *parser.IncDecStmt:
		t.TraverseExpr(s.Expr, scope)
	case *parser.ReturnStmt:
		t.TraverseExpr(s.Result, scope)
	case *parser.ExportStmt:
		t.handler.HandleExportStmt(s, scope)
		t.TraverseExpr(s.Result, scope)
	}
}

// TraverseExpr 遍历表达式
func (t *Traverser) TraverseExpr(expr parser.Expr, scope *scopeStack) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *parser.ArrayLit:
		for _, element := range e.Elements {
			t.TraverseExpr(element, scope)
		}
	case *parser.BinaryExpr:
		t.TraverseExpr(e.LHS, scope)
		t.TraverseExpr(e.RHS, scope)
	case *parser.CallExpr:
		t.TraverseExpr(e.Func, scope)
		for _, arg := range e.Args {
			t.TraverseExpr(arg, scope)
		}
	case *parser.CondExpr:
		t.TraverseExpr(e.Cond, scope)
		t.TraverseExpr(e.True, scope)
		t.TraverseExpr(e.False, scope)
	case *parser.FuncLit:
		t.handler.HandleFuncLit(e, scope)
		scope.pushScope()
		if e.Type.Params != nil {
			for _, param := range e.Type.Params.List {
				t.handler.HandleIdent(param, scope)
			}
		}
		if e.Body != nil {
			t.TraverseStmts(e.Body.Stmts, scope)
		}
		scope.popScope()
	case *parser.ImportExpr:
		t.handler.HandleImportExpr(e, scope)
	case *parser.IndexExpr:
		t.TraverseExpr(e.Expr, scope)
		t.TraverseExpr(e.Index, scope)
	case *parser.MapLit:
		for _, element := range e.Elements {
			t.TraverseExpr(element, scope)
		}
	case *parser.ParenExpr:
		t.TraverseExpr(e.Expr, scope)
	case *parser.SelectorExpr:
		t.TraverseExpr(e.Expr, scope)
		t.TraverseExpr(e.Sel, scope)
		t.handler.HandleSelectorExpr(e, scope)
	case *parser.SliceExpr:
		t.TraverseExpr(e.Expr, scope)
		if e.Low != nil {
			t.TraverseExpr(e.Low, scope)
		}
		if e.High != nil {
			t.TraverseExpr(e.High, scope)
		}
	case *parser.UnaryExpr:
		t.TraverseExpr(e.Expr, scope)
	case *parser.Ident:
		t.handler.HandleIdent(e, scope)
	}
}
