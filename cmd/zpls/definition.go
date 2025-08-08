package main

import (
	"github.com/diiyw/z/cmd/zpls/file"
	"path/filepath"
	"strings"

	"github.com/diiyw/z/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func onDefinitionFunc(context *glsp.Context, params *protocol.DefinitionParams) (any, error) {
	filename := strings.ReplaceAll(params.TextDocument.URI, "file://", "")
	content := file.Document().GetText(params.TextDocument.URI)
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

	// 创建一个作用域栈来追踪变量定义
	scopeStack = newScopeStack()

	// 使用新的遍历器查找定义
	locations := make([]protocol.LocationLink, 0)
	definitionCollector := NewDefinitionCollector(expr, &locations, filename, content, currentDir)
	traverser = NewTraverser(definitionCollector)
	traverser.TraverseStmts(parsedFile.Stmts, scopeStack)
	locationLinks = locations

	return locationLinks, nil
}

// DefinitionCollector 实现TraverserHandler接口，用于收集定义
type DefinitionCollector struct {
	BaseTraverserHandler
	targetNode parser.Expr
	locations  *[]protocol.LocationLink
	filename   string
	content    string
	currentDir string
}

// NewDefinitionCollector 创建一个新的DefinitionCollector
func NewDefinitionCollector(targetNode parser.Expr, locations *[]protocol.LocationLink, filename, content, currentDir string) *DefinitionCollector {
	return &DefinitionCollector{
		targetNode: targetNode,
		locations:  locations,
		filename:   filename,
		content:    content,
		currentDir: currentDir,
	}
}

func (dc *DefinitionCollector) HandleAssignStmt(stmt *parser.AssignStmt, scope *scopeStack) {
	// 添加左侧的变量到当前作用域
	for _, lh := range stmt.LHS {
		if ident, ok := lh.(*parser.Ident); ok {
			scope.addVariable(ident.Name)
			// 检查是否是我们要查找的标识符
			if targetIdent, ok := dc.targetNode.(*parser.Ident); ok && ident.Name == targetIdent.Name {
				*dc.locations = append(*dc.locations, protocol.LocationLink{
					OriginSelectionRange: &protocol.Range{
						Start: offsetToPosition(int(targetIdent.Pos()-1), dc.content),
						End:   offsetToPosition(int(targetIdent.End()-1), dc.content),
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

func (dc *DefinitionCollector) HandleImportExpr(expr *parser.ImportExpr, scope *scopeStack) {
	// 检查是否是我们要查找的import表达式
	if targetImport, ok := dc.targetNode.(*parser.ImportExpr); ok && expr == targetImport {
		*dc.locations = append(*dc.locations, protocol.LocationLink{
			OriginSelectionRange: &protocol.Range{
				Start: offsetToPosition(int(expr.Pos()+7), dc.content),
				End:   offsetToPosition(int(expr.End()-3), dc.content),
			},
			TargetURI:            creatURI(filepath.Join(dc.currentDir, expr.ModuleName+".z")),
			TargetRange:          createRange(0, 0, 0, 0),
			TargetSelectionRange: createRange(0, 0, 0, 0),
		})
	}
}
