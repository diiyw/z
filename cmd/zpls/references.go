package main

import (
	"github.com/diiyw/z/cmd/zpls/file"
	"path/filepath"
	"strings"

	"github.com/diiyw/z/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// ReferenceCollector 实现TraverserHandler接口，用于收集引用
type ReferenceCollector struct {
	BaseTraverserHandler
	targetIdent *parser.Ident
	locations   *[]protocol.Location
	filename    string
	content     string
	// 添加处理SelectorExpr所需字段
	targetSelector *parser.SelectorExpr
	currentDir     string
	parsedFile     *parser.File
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

// 添加一个新函数用于创建处理SelectorExpr的ReferenceCollector
func NewSelectorReferenceCollector(targetSelector *parser.SelectorExpr, targetIdent *parser.Ident, locations *[]protocol.Location, filename, content string, parsedFile *parser.File, currentDir string) *ReferenceCollector {
	return &ReferenceCollector{
		targetIdent:    targetIdent,
		targetSelector: targetSelector,
		locations:      locations,
		filename:       filename,
		content:        content,
		parsedFile:     parsedFile,
		currentDir:     currentDir,
	}
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

// 添加HandleSelectorExpr方法处理SelectorExpr的引用查找
func (rc *ReferenceCollector) HandleSelectorExpr(expr *parser.SelectorExpr, scope *scopeStack) {
	// 处理SelectorExpr，例如 module.variable 形式
	if rc.targetSelector != nil {
		if sel, ok := expr.Sel.(*parser.Ident); ok && rc.targetIdent != nil {
			if sel.Name == rc.targetIdent.Name {
				// 检查表达式的Expr部分是否为标识符
				if ident, ok := expr.Expr.(*parser.Ident); ok {
					if rc.targetSelector.Expr.(*parser.Ident).Name == ident.Name {
						// 查找模块导入
						importedFile := findImportedFile(ident.Name, rc.parsedFile.Stmts, rc.currentDir)
						if importedFile != "" {
							*rc.locations = append(*rc.locations, protocol.Location{
								URI: creatURI(importedFile),
								Range: protocol.Range{
									Start: offsetToPosition(int(sel.Pos()-1), rc.content),
									End:   offsetToPosition(int(sel.End()-1), rc.content),
								},
							})
						}
					}
				}
			}
		}
	}
}

func onReferencesFunc(context *glsp.Context, params *protocol.ReferenceParams) ([]protocol.Location, error) {
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
	var locations = make([]protocol.Location, 0)

	// 处理标识符的引用查找
	if ident, ok := expr.(*parser.Ident); ok {
		// 创建一个作用域栈来追踪变量定义
		scopeStack := newScopeStack()

		// 直接使用ReferenceCollector进行遍历
		referenceCollector := NewReferenceCollector(ident, &locations, filename, content)
		traverser := NewTraverser(referenceCollector)
		traverser.TraverseStmts(parsedFile.Stmts, scopeStack)

		locations = append(locations, protocol.Location{
			URI: creatURI(filename),
			Range: protocol.Range{
				Start: offsetToPosition(int(ident.Pos()-1), content),
				End:   offsetToPosition(int(ident.End()-1), content),
			},
		})
	}

	// 处理SelectorExpr，例如 module.variable 形式
	if selExpr, ok := expr.(*parser.SelectorExpr); ok {
		if sel, ok := selExpr.Sel.(*parser.Ident); ok {
			// 使用ReferenceCollector处理SelectorExpr引用
			scopeStack := newScopeStack()
			referenceCollector := NewSelectorReferenceCollector(selExpr, sel, &locations, filename, content, parsedFile, currentDir)
			traverser := NewTraverser(referenceCollector)
			traverser.TraverseStmts(parsedFile.Stmts, scopeStack)
		}
	}

	return locations, nil
}
