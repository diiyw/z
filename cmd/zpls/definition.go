package main

import (
	"os"
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
	currentDir := filepath.Dir(filename)
	var f = finder{currentDir: currentDir, offset: offset, filename: filename}
	var expr = f.findNodeByOffset(parsedFile.Stmts)
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
		locationLinks = f.findDefinition(e, parsedFile)
	}

	// 处理SelectorExpr，例如 module.variable 形式
	if e, ok := expr.(*parser.SelectorExpr); ok {
		if sel, ok := e.Sel.(*parser.Ident); ok {
			locationLinks = f.findSelectorDefinition(e, sel, parsedFile)
		}
	}

	return locationLinks, nil
}

type finder struct {
	currentDir string
	filename   string
	offset     int
}

func (f *finder) findDefinition(expr *parser.Ident, file *parser.File) []protocol.LocationLink {
	// 首先在当前文件中查找定义
	locations := f.findDefinitionInScopes(expr, file.Stmts)

	// 如果在当前文件中找到了定义，则直接返回
	if len(locations) > 0 {
		return locations
	}

	// 如果在当前文件中没找到，检查是否是模块.标识符的形式
	// 这部分逻辑可以在后续扩展

	return locations
}

// findSelectorDefinition 处理 selector 表达式，例如 module.variable
func (f *finder) findSelectorDefinition(expr *parser.SelectorExpr, sel *parser.Ident, file *parser.File) []protocol.LocationLink {
	locations := make([]protocol.LocationLink, 0)

	// 检查表达式的Expr部分是否为标识符
	if ident, ok := expr.Expr.(*parser.Ident); ok {
		// 查找模块导入
		importedFile := f.findImportedFile(ident.Name, file.Stmts)
		if importedFile != "" {
			// 在导入的文件中查找定义
			definitions := f.findDefinitionInFile(sel, importedFile)
			locations = append(locations, definitions...)
		}
	}

	return locations
}

// findImportedFile 查找导入的文件路径
func (f *finder) findImportedFile(moduleName string, stmts []parser.Stmt) string {
	for _, stmt := range stmts {
		if imp, ok := stmt.(*parser.ExprStmt); ok {
			if importExpr, ok := imp.Expr.(*parser.ImportExpr); ok {
				if importExpr.ModuleName == moduleName {
					// 构造导入文件的完整路径
					importPath := filepath.Join(f.currentDir, moduleName+".z")
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

// findDefinitionInScopes 在作用域中查找定义，考虑块级作用域
func (f *finder) findDefinitionInScopes(expr *parser.Ident, stmts []parser.Stmt) []protocol.LocationLink {
	locations := make([]protocol.LocationLink, 0)

	// 创建一个作用域栈来追踪变量定义
	scopeStack := newScopeStack()

	// 在当前作用域中查找定义
	f.collectDefinitionsInScopes(expr, stmts, scopeStack, &locations)

	return locations
}

// scopeStack 作用域栈，用于管理变量定义的作用域
type scopeStack struct {
	scopes [][]string // 每个作用域中的变量名
}

func newScopeStack() *scopeStack {
	return &scopeStack{
		scopes: make([][]string, 0),
	}
}

func (s *scopeStack) pushScope() {
	s.scopes = append(s.scopes, make([]string, 0))
}

func (s *scopeStack) popScope() {
	if len(s.scopes) > 0 {
		s.scopes = s.scopes[:len(s.scopes)-1]
	}
}

func (s *scopeStack) addVariable(name string) {
	if len(s.scopes) > 0 {
		s.scopes[len(s.scopes)-1] = append(s.scopes[len(s.scopes)-1], name)
	}
}

func (s *scopeStack) isVariableInScope(name string) bool {
	// 从内到外检查作用域
	for i := len(s.scopes) - 1; i >= 0; i-- {
		scope := s.scopes[i]
		for _, variable := range scope {
			if variable == name {
				return true
			}
		}
	}
	return false
}

// collectDefinitionsInScopes 收集作用域中的定义
func (f *finder) collectDefinitionsInScopes(expr *parser.Ident, stmts []parser.Stmt, scopeStack *scopeStack, locations *[]protocol.LocationLink) {
	content := Document().GetText(creatURI(f.filename))

	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *parser.AssignStmt:
			// 添加左侧的变量到当前作用域
			for _, lh := range s.LHS {
				if ident, ok := lh.(*parser.Ident); ok {
					scopeStack.addVariable(ident.Name)
					// 检查是否是我们要查找的标识符
					if ident.Name == expr.Name {
						*locations = append(*locations, protocol.LocationLink{
							OriginSelectionRange: &protocol.Range{
								Start: offsetToPosition(int(expr.Pos()-1), content),
								End:   offsetToPosition(int(expr.End()-1), content),
							},
							TargetRange: protocol.Range{
								Start: offsetToPosition(int(ident.Pos()-1), content),
								End:   offsetToPosition(int(ident.End()-1), content),
							},
							TargetSelectionRange: protocol.Range{
								Start: offsetToPosition(int(ident.Pos()-1), content),
								End:   offsetToPosition(int(ident.End()-1), content),
							},
							TargetURI: creatURI(f.filename),
						})
					}
				}
			}
		case *parser.BlockStmt:
			// 推入新的作用域
			scopeStack.pushScope()
			// 在块中递归查找
			f.collectDefinitionsInScopes(expr, s.Stmts, scopeStack, locations)
			// 弹出作用域
			scopeStack.popScope()
		case *parser.ForInStmt:
			// for-in语句创建新的作用域
			scopeStack.pushScope()
			// 添加key和value到作用域
			scopeStack.addVariable(s.Key.Name)
			if s.Key.Name == expr.Name {
				content := Document().GetText(creatURI(f.filename))
				*locations = append(*locations, protocol.LocationLink{
					OriginSelectionRange: &protocol.Range{
						Start: offsetToPosition(int(expr.Pos()-1), content),
						End:   offsetToPosition(int(expr.End()-1), content),
					},
					TargetRange: protocol.Range{
						Start: offsetToPosition(int(s.Key.Pos()-1), content),
						End:   offsetToPosition(int(s.Key.End()-1), content),
					},
					TargetSelectionRange: protocol.Range{
						Start: offsetToPosition(int(s.Key.Pos()-1), content),
						End:   offsetToPosition(int(s.Key.End()-1), content),
					},
					TargetURI: creatURI(f.filename),
				})
			}

			scopeStack.addVariable(s.Value.Name)
			if s.Value.Name == expr.Name {
				content := Document().GetText(creatURI(f.filename))
				*locations = append(*locations, protocol.LocationLink{
					OriginSelectionRange: &protocol.Range{
						Start: offsetToPosition(int(expr.Pos()-1), content),
						End:   offsetToPosition(int(expr.End()-1), content),
					},
					TargetRange: protocol.Range{
						Start: offsetToPosition(int(s.Value.Pos()-1), content),
						End:   offsetToPosition(int(s.Value.End()-1), content),
					},
					TargetSelectionRange: protocol.Range{
						Start: offsetToPosition(int(s.Value.Pos()-1), content),
						End:   offsetToPosition(int(s.Value.End()-1), content),
					},
					TargetURI: creatURI(f.filename),
				})
			}

			// 在循环体中递归查找
			f.collectDefinitionsInScopes(expr, s.Body.Stmts, scopeStack, locations)
			// 弹出作用域
			scopeStack.popScope()
		}
	}
}

func (f *finder) findNodeByOffset(stmts []parser.Stmt) parser.Expr {
	for _, stmt := range stmts {
		if node := f.findStmtNode(stmt); node != nil {
			return node
		}
	}
	return nil
}

func (f *finder) findStmtNode(stmt parser.Stmt) parser.Expr {
	ss, se := stmt.Pos(), stmt.End()
	if int(ss) <= f.offset && f.offset < int(se) {
		switch s := stmt.(type) {
		case *parser.AssignStmt:
			for _, expr := range s.LHS {
				if e := f.findExprNode(expr); e != nil {
					return e
				}
			}
			for _, expr := range s.RHS {
				if e := f.findExprNode(expr); e != nil {
					return e
				}
				if int(expr.Pos()) <= f.offset && f.offset < int(expr.End()) {
					return expr
				}
			}
		case *parser.ExportStmt:
			return f.findExprNode(s.Result)
		case *parser.BlockStmt:
			for _, st := range s.Stmts {
				if int(st.Pos()) <= f.offset && f.offset < int(st.End()) {
					return f.findStmtNode(st)
				}
			}
		case *parser.ExprStmt:
			return f.findExprNode(s.Expr)
		case *parser.ForInStmt:
			if e := f.findExprNode(s.Iterable); e != nil {
				return e
			}
			if e := f.findExprNode(s.Key); e != nil {
				return e
			}
			if e := f.findExprNode(s.Value); e != nil {
				return e
			}
			return f.findStmtNode(s.Body)
		case *parser.ForStmt:
			if e := f.findExprNode(s.Cond); e != nil {
				return e
			}
			if node := f.findStmtNode(s.Init); node != nil {
				return node
			}
			if node := f.findStmtNode(s.Post); node != nil {
				return node
			}
			return f.findStmtNode(s.Body)
		case *parser.IfStmt:
			if e := f.findExprNode(s.Cond); e != nil {
				return e
			}
			if e := f.findStmtNode(s.Init); e != nil {
				return e
			}
			if e := f.findStmtNode(s.Else); e != nil {
				return e
			}
			return f.findStmtNode(s.Body)
		case *parser.IncDecStmt:
			if e := f.findExprNode(s.Expr); e != nil {
				return e
			}
		case *parser.ReturnStmt:
			if e := f.findExprNode(s.Result); e != nil {
				return e
			}
		}
	}
	return nil
}

func (f *finder) findExprNode(expr parser.Expr) parser.Expr {
	if f.offset > int(expr.End()) {
		return nil
	}
	switch e := expr.(type) {
	case *parser.ArrayLit:
		for _, element := range e.Elements {
			if e := f.findExprNode(element); e != nil {
				return e
			}
		}
	case *parser.BadExpr:
		return nil
	case *parser.BinaryExpr:
		if int(e.LHS.Pos()) <= f.offset && f.offset < int(e.LHS.End()) {
			return f.findExprNode(e.LHS)
		}
		if int(e.RHS.Pos()) <= f.offset && f.offset < int(e.RHS.End()) {
			return f.findExprNode(e.RHS)
		}
	case *parser.BoolLit:
		return e
	case *parser.CallExpr:
		if int(e.Func.Pos()) <= f.offset && f.offset < int(e.Func.End()) {
			return f.findExprNode(e.Func)
		}
		for _, a := range e.Args {
			if e := f.findExprNode(a); e != nil {
				return e
			}
		}
	case *parser.CharLit:
		return e
	case *parser.CondExpr:
		if int(e.Cond.Pos()) <= f.offset && f.offset < int(e.Cond.End()) {
			return f.findExprNode(e.Cond)
		}
		if int(e.False.Pos()) <= f.offset && f.offset < int(e.False.End()) {
			return f.findExprNode(e.False)
		}
		if int(e.True.Pos()) <= f.offset && f.offset < int(e.True.End()) {
			return f.findExprNode(e.True)
		}
	case *parser.ErrorExpr:
		return e
	case *parser.FloatLit:
		return e
	case *parser.FuncLit:
		if int(e.Type.Pos()) <= f.offset && f.offset < int(e.Type.End()) {
			return f.findExprNode(e.Type)
		}
		return f.findStmtNode(e.Body)
	case *parser.FuncType:
		for _, a := range e.Params.List {
			if e := f.findExprNode(a); e != nil {
				return e
			}
		}
		return e
	case *parser.ImmutableExpr:
		return e
	case *parser.ImportExpr:
		return e
	case *parser.IndexExpr:
		if int(e.Expr.Pos()) <= f.offset && f.offset < int(e.Expr.End()) {
			return f.findExprNode(e.Expr)
		}
		if int(e.Index.Pos()) <= f.offset && f.offset < int(e.Index.End()) {
			return f.findExprNode(e.Index)
		}
	case *parser.Ident:
		if int(e.Pos()) <= f.offset && f.offset < int(e.End()) {
			return e
		}
	case *parser.IntLit:
		return e
	case *parser.MapElementLit:
		if int(e.Key.Pos()) <= f.offset && f.offset < int(e.Key.End()) {
			return f.findExprNode(e.Key)
		}
		if int(e.Value.Pos()) <= f.offset && f.offset < int(e.Value.End()) {
			return f.findExprNode(e.Value)
		}
	case *parser.MapLit:
		for _, element := range e.Elements {
			if e := f.findExprNode(element); e != nil {
				return e
			}
		}
	case *parser.ParenExpr:
		if int(e.Expr.Pos()) <= f.offset && f.offset < int(e.Expr.End()) {
			return f.findExprNode(e.Expr)
		}
	case *parser.SelectorExpr:
		if int(e.Expr.Pos()) <= f.offset && f.offset < int(e.Expr.End()) {
			return f.findExprNode(e.Expr)
		}
		if int(e.Sel.Pos()) <= f.offset && f.offset < int(e.Sel.End()) {
			return f.findExprNode(e.Sel)
		}
	case *parser.SliceExpr:
		if int(e.Expr.Pos()) <= f.offset && f.offset < int(e.Expr.End()) {
			return f.findExprNode(e.Expr)
		}
		if int(e.Low.Pos()) <= f.offset && f.offset < int(e.Low.End()) {
			return f.findExprNode(e.Low)
		}
		if int(e.High.Pos()) <= f.offset && f.offset < int(e.High.End()) {
			return f.findExprNode(e.High)
		}
	case *parser.StringLit:
		return e
	case *parser.UnaryExpr:
		if int(e.Expr.Pos()) <= f.offset && f.offset < int(e.Expr.End()) {
			return f.findExprNode(e.Expr)
		}
	case *parser.UndefinedLit:
		return e
	}
	return nil
}

// findDefinitionInFile 在指定文件中查找标识符定义
func (f *finder) findDefinitionInFile(expr *parser.Ident, filename string) []protocol.LocationLink {
	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return []protocol.LocationLink{}
	}

	// 读取文件内容
	content, err := os.ReadFile(filename)
	if err != nil {
		return []protocol.LocationLink{}
	}

	// 解析文件
	fileSet := parser.NewFileSet()
	basename := filepath.Base(filename)
	sourceFile := fileSet.AddFile(basename, -1, len(content))
	p := parser.NewParser(sourceFile, content, nil)
	parsedFile, err := p.ParseFile()
	if err != nil {
		return []protocol.LocationLink{}
	}

	// 在解析后的文件中查找定义
	return f.findDefinitionInScopes(expr, parsedFile.Stmts)
}

// 以下是从references.go移动过来的方法
func (f *finder) findReferences(expr *parser.Ident, file *parser.File) []protocol.Location {
	locations := make([]protocol.Location, 0)
	
	// 查找当前文件中的引用
	locations = append(locations, f.findReferencesInScopes(expr, file.Stmts)...)
	
	return locations
}

// findSelectorReferences 处理 selector 表达式，例如 module.variable 的引用查找
func (f *finder) findSelectorReferences(expr *parser.SelectorExpr, sel *parser.Ident, file *parser.File) []protocol.Location {
	locations := make([]protocol.Location, 0)
	
	// 检查表达式的Expr部分是否为标识符
	if ident, ok := expr.Expr.(*parser.Ident); ok {
		// 查找模块导入
		importedFile := f.findImportedFile(ident.Name, file.Stmts)
		if importedFile != "" {
			// 在导入的文件中查找引用
			references := f.findReferencesInFile(sel, importedFile)
			locations = append(locations, references...)
		}
	}
	
	return locations
}

// findReferencesInScopes 在作用域中查找引用
func (f *finder) findReferencesInScopes(expr *parser.Ident, stmts []parser.Stmt) []protocol.Location {
	locations := make([]protocol.Location, 0)
	
	// 收集所有引用
	f.collectReferences(expr, stmts, &locations)
	
	return locations
}

// collectReferences 收集引用
func (f *finder) collectReferences(expr *parser.Ident, stmts []parser.Stmt, locations *[]protocol.Location) {
	content := Document().GetText(creatURI(f.filename))
	
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *parser.AssignStmt:
			// 检查右侧表达式中的引用
			for _, rh := range s.RHS {
				f.findExprReferences(expr, rh, locations, content)
			}
			// 检查左侧表达式中的引用（如果变量在声明前被引用）
			for _, lh := range s.LHS {
				if ident, ok := lh.(*parser.Ident); ok && ident.Name == expr.Name {
					*locations = append(*locations, protocol.Location{
						URI: creatURI(f.filename),
						Range: protocol.Range{
							Start: offsetToPosition(int(ident.Pos()-1), content),
							End:   offsetToPosition(int(ident.End()-1), content),
						},
					})
				}
			}
		case *parser.BlockStmt:
			// 在块中递归查找
			f.collectReferences(expr, s.Stmts, locations)
		case *parser.ExprStmt:
			f.findExprReferences(expr, s.Expr, locations, content)
		case *parser.ForInStmt:
			f.findExprReferences(expr, s.Iterable, locations, content)
			f.findExprReferences(expr, s.Key, locations, content)
			f.findExprReferences(expr, s.Value, locations, content)
			// 在循环体中递归查找
			f.collectReferences(expr, s.Body.Stmts, locations)
		case *parser.ForStmt:
			f.findExprReferences(expr, s.Cond, locations, content)
			f.findStmtReferences(expr, s.Init, locations, content)
			f.findStmtReferences(expr, s.Post, locations, content)
			// 在循环体中递归查找
			f.collectReferences(expr, s.Body.Stmts, locations)
		case *parser.IfStmt:
			f.findExprReferences(expr, s.Cond, locations, content)
			f.findStmtReferences(expr, s.Init, locations, content)
			f.findStmtReferences(expr, s.Else, locations, content)
			// 在条件体中递归查找
			f.collectReferences(expr, s.Body.Stmts, locations)
		case *parser.IncDecStmt:
			f.findExprReferences(expr, s.Expr, locations, content)
		case *parser.ReturnStmt:
			f.findExprReferences(expr, s.Result, locations, content)
		case *parser.ExportStmt:
			f.findExprReferences(expr, s.Result, locations, content)
		}
	}
}

// findExprReferences 在表达式中查找引用
func (f *finder) findExprReferences(target *parser.Ident, expr parser.Expr, locations *[]protocol.Location, content string) {
	if expr == nil {
		return
	}
	
	switch e := expr.(type) {
	case *parser.Ident:
		if e.Name == target.Name {
			*locations = append(*locations, protocol.Location{
				URI: creatURI(f.filename),
				Range: protocol.Range{
					Start: offsetToPosition(int(e.Pos()-1), content),
					End:   offsetToPosition(int(e.End()-1), content),
				},
			})
		}
	case *parser.ArrayLit:
		for _, element := range e.Elements {
			f.findExprReferences(target, element, locations, content)
		}
	case *parser.BinaryExpr:
		f.findExprReferences(target, e.LHS, locations, content)
		f.findExprReferences(target, e.RHS, locations, content)
	case *parser.CallExpr:
		f.findExprReferences(target, e.Func, locations, content)
		for _, arg := range e.Args {
			f.findExprReferences(target, arg, locations, content)
		}
	case *parser.CondExpr:
		f.findExprReferences(target, e.Cond, locations, content)
		f.findExprReferences(target, e.True, locations, content)
		f.findExprReferences(target, e.False, locations, content)
	case *parser.FuncLit:
		if e.Body != nil {
			f.collectReferences(target, e.Body.Stmts, locations)
		}
	case *parser.IndexExpr:
		f.findExprReferences(target, e.Expr, locations, content)
		f.findExprReferences(target, e.Index, locations, content)
	case *parser.MapElementLit:
		f.findExprReferences(target, e.Key, locations, content)
		f.findExprReferences(target, e.Value, locations, content)
	case *parser.MapLit:
		for _, element := range e.Elements {
			f.findExprReferences(target, element, locations, content)
		}
	case *parser.ParenExpr:
		f.findExprReferences(target, e.Expr, locations, content)
	case *parser.SelectorExpr:
		f.findExprReferences(target, e.Expr, locations, content)
		f.findExprReferences(target, e.Sel, locations, content)
	case *parser.SliceExpr:
		f.findExprReferences(target, e.Expr, locations, content)
		if e.Low != nil {
			f.findExprReferences(target, e.Low, locations, content)
		}
		if e.High != nil {
			f.findExprReferences(target, e.High, locations, content)
		}
	case *parser.UnaryExpr:
		f.findExprReferences(target, e.Expr, locations, content)
	}
}

// findStmtReferences 在语句中查找引用
func (f *finder) findStmtReferences(target *parser.Ident, stmt parser.Stmt, locations *[]protocol.Location, content string) {
	if stmt == nil {
		return
	}
	
	// 处理不同类型的语句
	switch s := stmt.(type) {
	case *parser.AssignStmt:
		for _, rh := range s.RHS {
			f.findExprReferences(target, rh, locations, content)
		}
	case *parser.ExprStmt:
		f.findExprReferences(target, s.Expr, locations, content)
	case *parser.IfStmt:
		f.findExprReferences(target, s.Cond, locations, content)
	}
}

// findReferencesInFile 在指定文件中查找标识符引用
func (f *finder) findReferencesInFile(expr *parser.Ident, filename string) []protocol.Location {
	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return []protocol.Location{}
	}
	
	// 读取文件内容
	content, err := os.ReadFile(filename)
	if err != nil {
		return []protocol.Location{}
	}
	
	// 解析文件
	fileSet := parser.NewFileSet()
	basename := filepath.Base(filename)
	sourceFile := fileSet.AddFile(basename, -1, len(content))
	p := parser.NewParser(sourceFile, content, nil)
	parsedFile, err := p.ParseFile()
	if err != nil {
		return []protocol.Location{}
	}
	
	// 在解析后的文件中查找引用
	return f.findReferencesInScopes(expr, parsedFile.Stmts)
}
