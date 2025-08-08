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
	BaseTraverserHandler
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
			locations = findSelectorReferences(selExpr, sel, parsedFile, currentDir)
		}
	}

	return locations, nil
}

// findSelectorReferences 处理 selector 表达式，例如 module.variable 的引用查找
func findSelectorReferences(expr *parser.SelectorExpr, sel *parser.Ident, file *parser.File, currentDir string) []protocol.Location {
	locations := make([]protocol.Location, 0)

	// 检查表达式的Expr部分是否为标识符
	if ident, ok := expr.Expr.(*parser.Ident); ok {
		// 查找模块导入
		importedFile := findImportedFile(ident.Name, file.Stmts, currentDir)
		if importedFile != "" {
			// 在导入的文件中查找引用
			references := findReferencesInFile(sel, importedFile)
			locations = append(locations, references...)
		}
	}

	return locations
}

// findReferencesInFile 在指定文件中查找标识符引用
func findReferencesInFile(expr *parser.Ident, filename string) []protocol.Location {
	// 使用通用文件解析函数获取解析后的文件
	parsedFile, err := parseFile(filename)
	if err != nil {
		return []protocol.Location{}
	}

	// 在解析后的文件中查找引用
	return findReferencesInScopes(expr, parsedFile.Stmts, filename)
}

// findReferencesInScopes 在作用域中查找引用
func findReferencesInScopes(expr *parser.Ident, stmts []parser.Stmt, filename string) []protocol.Location {
	locations := make([]protocol.Location, 0)

	// 创建一个作用域栈来追踪变量定义
	scopeStack := newScopeStack()

	// 使用新的遍历器查找引用
	content := Document().GetText(creatURI(filename))
	referenceCollector := NewReferenceCollector(expr, &locations, filename, content)
	traverser := NewTraverser(referenceCollector)
	traverser.TraverseStmts(stmts, scopeStack)

	return locations
}
