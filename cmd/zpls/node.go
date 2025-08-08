package main

import (
	"os"
	"path/filepath"

	"github.com/diiyw/z/parser"
)

// NodeFinder 实现TraverserHandler接口，用于查找指定偏移位置的节点
type NodeFinder struct {
	BaseTraverserHandler
	offset int
	parent parser.Expr
	node   parser.Expr
}

// NewNodeFinder 创建一个新的NodeFinder
func NewNodeFinder(offset int) *NodeFinder {
	return &NodeFinder{
		offset: offset,
	}
}

func (nf *NodeFinder) HandleAssignStmt(stmt *parser.AssignStmt, scope *scopeStack) {
	// 检查LHS表达式
	for _, expr := range stmt.LHS {
		if nf.isOffsetInNode(expr) {
			nf.node = expr
			return
		}
	}

	// 检查RHS表达式
	for _, expr := range stmt.RHS {
		if nf.isOffsetInNode(expr) {
			nf.node = expr
			return
		}
	}
}

func (nf *NodeFinder) HandleFuncLit(stmt *parser.FuncLit, scope *scopeStack) {
	// FuncLit本身可能就是我们要找的节点
	if nf.isOffsetInNode(stmt) {
		nf.parent = stmt
		nf.node = stmt
	}
}

func (nf *NodeFinder) HandleForInStmt(stmt *parser.ForInStmt, scope *scopeStack) {
	// 检查Key和Value
	if nf.isOffsetInNode(stmt.Key) {
		nf.node = stmt.Key
		return
	}

	if nf.isOffsetInNode(stmt.Value) {
		nf.node = stmt.Value
		return
	}

	// 检查Iterable
	if nf.isOffsetInNode(stmt.Iterable) {
		nf.node = stmt.Iterable
		return
	}
}

func (nf *NodeFinder) HandleIdent(ident *parser.Ident, scope *scopeStack) {
	// 检查标识符是否在偏移位置
	if nf.isOffsetInNode(ident) {
		nf.node = ident
	}
}

// HandleSelectorExpr 处理SelectorExpr节点
func (nf *NodeFinder) HandleSelectorExpr(expr *parser.SelectorExpr, scope *scopeStack) {
	// 检查Expr部分是否包含偏移位置
	if nf.isOffsetInNode(expr.Expr) {
		nf.parent = expr
		nf.node = expr.Expr
		return
	}

	// 检查Sel部分是否包含偏移位置
	if nf.isOffsetInNode(expr.Sel) {
		nf.parent = expr
		nf.node = expr.Sel
		return
	}
}

// isOffsetInNode 检查偏移位置是否在节点范围内
func (nf *NodeFinder) isOffsetInNode(node parser.Node) bool {
	if node == nil {
		return false
	}
	return int(node.Pos()) <= nf.offset && nf.offset < int(node.End())
}

// GetNode 返回找到的节点
func (nf *NodeFinder) GetNode() parser.Expr {
	return nf.node
}

// GetParent 返回父节点
func (nf *NodeFinder) GetParent() parser.Node {
	return nf.parent
}

// parseFile 是一个辅助方法，用于封装文件的读取和解析
func parseFile(filename string) (*parser.File, error) {
	// 读取文件内容
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// 使用新的通用文件解析函数
	return ParseFileContent(filename, string(content))
}

// findImportedFile 查找导入的文件路径
func findImportedFile(moduleName string, stmts []parser.Stmt, currentDir string) string {
	for _, stmt := range stmts {
		if imp, ok := stmt.(*parser.ExprStmt); ok {
			if importExpr, ok := imp.Expr.(*parser.ImportExpr); ok {
				if importExpr.ModuleName == moduleName {
					// 构造导入文件的完整路径
					importPath := filepath.Join(currentDir, moduleName+".z")
					if _, err := os.Stat(importPath); err == nil {
						return importPath
					}
					return ""
				}
			}
		}
	}
	return ""
}
