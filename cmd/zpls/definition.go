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

	// 使用新的通用文件解析函数
	parsedFile, err := ParseFileContent(filename, content)
	if err != nil {
		return nil, err
	}

	currentDir := filepath.Dir(filename)

	// 使用Traverser方式查找节点
	nodeFinder := NewNodeFinder(offset)
	scopeStack := newScopeStack()
	traverser := NewTraverser(nodeFinder)
	traverser.TraverseStmts(parsedFile.Stmts, scopeStack)

	var expr = nodeFinder.GetNode()
	var locationLinks = make([]protocol.LocationLink, 0)

	// 处理import表达式
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
		return locationLinks, nil
	}

	// 处理MapElementLit中的键
	if e, ok := expr.(*parser.MapElementLit); ok {
		if ident, ok := e.Key.(*parser.Ident); ok {
			locationLinks = append(locationLinks, protocol.LocationLink{
				OriginSelectionRange: &protocol.Range{
					Start: offsetToPosition(int(ident.Pos()), content),
					End:   offsetToPosition(int(ident.End()), content),
				},
				TargetRange: protocol.Range{
					Start: offsetToPosition(int(ident.Pos()), content),
					End:   offsetToPosition(int(ident.End()), content),
				},
				TargetSelectionRange: protocol.Range{
					Start: offsetToPosition(int(ident.Pos()), content),
					End:   offsetToPosition(int(ident.End()), content),
				},
				TargetURI: creatURI(filename),
			})
			return locationLinks, nil
		}
	}

	// 处理标识符（普通变量）
	if e, ok := expr.(*parser.Ident); ok {
		// 创建一个作用域栈来追踪变量定义
		scopeStack := newScopeStack()

		// 使用新的遍历器查找定义
		content := Document().GetText(creatURI(filename))
		locations := make([]protocol.LocationLink, 0)
		definitionCollector := NewDefinitionCollector(e, &locations, filename, content)
		traverser := NewTraverser(definitionCollector)
		traverser.TraverseStmts(parsedFile.Stmts, scopeStack)
		locationLinks = locations
		return locationLinks, nil
	}

	return locationLinks, nil
}

// NodeFinder 实现TraverserHandler接口，用于查找指定偏移位置的节点
type NodeFinder struct {
	BaseTraverserHandler
	offset int
	parent parser.Expr
	node   parser.Expr
}

// DefinitionCollector 实现TraverserHandler接口，用于收集定义
type DefinitionCollector struct {
	BaseTraverserHandler
	targetIdent *parser.Ident
	locations   *[]protocol.LocationLink
	filename    string
	content     string
}

// NewDefinitionCollector 创建一个新的DefinitionCollector
func NewDefinitionCollector(targetIdent *parser.Ident, locations *[]protocol.LocationLink, filename, content string) *DefinitionCollector {
	return &DefinitionCollector{
		targetIdent: targetIdent,
		locations:   locations,
		filename:    filename,
		content:     content,
	}
}

func (dc *DefinitionCollector) HandleAssignStmt(stmt *parser.AssignStmt, scope *scopeStack) {
	// 添加左侧的变量到当前作用域
	for _, lh := range stmt.LHS {
		if ident, ok := lh.(*parser.Ident); ok {
			scope.addVariable(ident.Name)
			// 检查是否是我们要查找的标识符
			if ident.Name == dc.targetIdent.Name {
				*dc.locations = append(*dc.locations, protocol.LocationLink{
					OriginSelectionRange: &protocol.Range{
						Start: offsetToPosition(int(dc.targetIdent.Pos()-1), dc.content),
						End:   offsetToPosition(int(dc.targetIdent.End()-1), dc.content),
					},
					TargetRange: protocol.Range{
						Start: offsetToPosition(int(ident.Pos()-1), dc.content),
						End:   offsetToPosition(int(ident.End()-1), dc.content),
					},
					TargetSelectionRange: protocol.Range{
						Start: offsetToPosition(int(ident.Pos()-1), dc.content),
						End:   offsetToPosition(int(ident.End()-1), dc.content),
					},
					TargetURI: creatURI(dc.filename),
				})
			}
		}
	}
}
