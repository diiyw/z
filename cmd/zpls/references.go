package main

import (
	"path/filepath"
	"strings"

	"github.com/diiyw/z/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// ReferenceCollector 实现TraverserHandler接口，用于收集引用
type ReferenceCollector struct {
	targetIdent *parser.Ident
	locations   *[]protocol.Location
	filename    string
	content     string
}

// NewReferenceCollector 创建一个新的ReferenceCollector
func NewReferenceCollector(targetIdent *parser.Ident, locations *[]protocol.Location, filename, content string) *ReferenceCollector {
	return &ReferenceCollector{
		targetIdent: targetIdent,
		locations:   locations,
		filename:    filename,
		content:     content,
	}
}

func (rc *ReferenceCollector) HandleAssignStmt(stmt *parser.AssignStmt, scope *scopeStack) {
	// 空实现 - 引用收集在HandleIdent中处理
}

func (rc *ReferenceCollector) HandleBlockStmt(stmt *parser.BlockStmt, scope *scopeStack) {
	// 空实现
}

func (rc *ReferenceCollector) HandleFuncLit(stmt *parser.FuncLit, scope *scopeStack) {
	// 空实现 - 引用收集在HandleIdent中处理
}

func (rc *ReferenceCollector) HandleForInStmt(stmt *parser.ForInStmt, scope *scopeStack) {
	// 空实现 - 引用收集在HandleIdent中处理
}

func (rc *ReferenceCollector) HandleIdent(ident *parser.Ident, scope *scopeStack) {
	if ident.Name == rc.targetIdent.Name {
		*rc.locations = append(*rc.locations, protocol.Location{
			URI: creatURI(rc.filename),
			Range: protocol.Range{
				Start: offsetToPosition(int(ident.Pos()-1), rc.content),
				End:   offsetToPosition(int(ident.End()-1), rc.content),
			},
		})
	}
}

func onReferencesFunc(context *glsp.Context, params *protocol.ReferenceParams) ([]protocol.Location, error) {
	filename := strings.ReplaceAll(params.TextDocument.URI, "file://", "")
	content := Document().GetText(params.TextDocument.URI)
	offset := params.Position.IndexIn(content)

	// 使用新的通用文件解析函数
	parsedFile, err := ParseFileContent(filename, content)
	if err != nil {
		return nil, err
	}

	currentDir := filepath.Dir(filename)
	var f = finder{currentDir: currentDir, offset: offset, filename: filename}
	var expr = f.findNodeByOffset(parsedFile.Stmts)
	var locations = make([]protocol.Location, 0)

	// 处理标识符的引用查找
	if ident, ok := expr.(*parser.Ident); ok {
		// 创建一个作用域栈来追踪变量定义
		scopeStack := newScopeStack()

		// 直接使用ReferenceCollector进行遍历
		referenceCollector := NewReferenceCollector(ident, &locations, f.filename, content)
		traverser := NewTraverser(referenceCollector)
		traverser.TraverseStmts(parsedFile.Stmts, scopeStack)

		locations = append(locations, protocol.Location{
			URI: creatURI(f.filename),
			Range: protocol.Range{
				Start: offsetToPosition(int(ident.Pos()-1), content),
				End:   offsetToPosition(int(ident.End()-1), content),
			},
		})
	}

	// 处理SelectorExpr，例如 module.variable 形式
	if selExpr, ok := expr.(*parser.SelectorExpr); ok {
		if sel, ok := selExpr.Sel.(*parser.Ident); ok {
			locations = f.findSelectorReferences(selExpr, sel, parsedFile)
		}
	}

	return locations, nil
}