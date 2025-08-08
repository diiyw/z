package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/diiyw/z/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	// Z语言关键字
	keywords = []string{
		"if",
		"else",
		"for",
		"return",
		"func",
		"export",
		"in",
	}

	// Z语言内置函数
	builtinFunctions = []string{
		"to_json",
		"from_json",
		"print",
		"printf",
		"sprintf",
		"format",
		"len",
		"copy",
		"append",
		"delete",
		"splice",
		"type_name",
		"int",
		"bool",
		"float",
		"char",
		"bytes",
		"error",
		"string",
		"time",
		"is_string",
		"is_bool",
		"is_float",
		"is_char",
		"is_bytes",
		"is_error",
		"is_undefined",
		"is_function",
		"is_callable",
		"is_array",
		"is_immutable_array",
		"is_map",
		"is_iterable",
		"is_time",
	}

	// 常量
	constants = []string{"true", "false", "undefined"}
)

// CompletionCollector 实现TraverserHandler接口，用于收集变量
type CompletionCollector struct {
	BaseTraverserHandler
	variables *[]string
	varMap    map[string]bool
	imports   map[string]string
	filename  string
}

// NewCompletionCollector 创建一个新的CompletionCollector
func NewCompletionCollector(variables *[]string) *CompletionCollector {
	// 使用map去重
	varMap := make(map[string]bool)
	for _, v := range *variables {
		varMap[v] = true
	}

	return &CompletionCollector{
		variables: variables,
		varMap:    varMap,
		imports:   make(map[string]string),
	}
}

// SetFilename 设置当前处理的文件名
func (c *CompletionCollector) SetFilename(filename string) {
	c.filename = filename
}

// GetImports 获取收集到的导入模块
func (c *CompletionCollector) GetImports() map[string]string {
	return c.imports
}

func (c *CompletionCollector) HandleAssignStmt(stmt *parser.AssignStmt, scope *scopeStack) {
	// 添加左侧的变量到当前作用域
	for _, lh := range stmt.LHS {
		if ident, ok := lh.(*parser.Ident); ok {
			scope.addVariable(ident.Name)
			// 如果变量还没有加入到列表中，则添加
			if !c.varMap[ident.Name] {
				*c.variables = append(*c.variables, ident.Name)
				c.varMap[ident.Name] = true
			}
		}
	}
}

func (c *CompletionCollector) HandleFuncLit(stmt *parser.FuncLit, scope *scopeStack) {
	// 添加参数到作用域
	if stmt.Type.Params != nil {
		for _, param := range stmt.Type.Params.List {
			scope.addVariable(param.Name)
			// 如果变量还没有加入到列表中，则添加
			if !c.varMap[param.Name] {
				*c.variables = append(*c.variables, param.Name)
				c.varMap[param.Name] = true
			}
		}
	}
}

func (c *CompletionCollector) HandleForInStmt(stmt *parser.ForInStmt, scope *scopeStack) {
	// 添加key和value到作用域
	scope.addVariable(stmt.Key.Name)
	// 如果变量还没有加入到列表中，则添加
	if !c.varMap[stmt.Key.Name] {
		*c.variables = append(*c.variables, stmt.Key.Name)
		c.varMap[stmt.Key.Name] = true
	}

	scope.addVariable(stmt.Value.Name)
	// 如果变量还没有加入到列表中，则添加
	if !c.varMap[stmt.Value.Name] {
		*c.variables = append(*c.variables, stmt.Value.Name)
		c.varMap[stmt.Value.Name] = true
	}
}

// HandleImportExpr 处理import语句，收集导入的模块
func (c *CompletionCollector) HandleImportExpr(expr *parser.ImportExpr, scope *scopeStack) {
	// 处理import语句，收集导入的模块
	if c.filename != "" {
		currentDir := filepath.Dir(c.filename)
		moduleName := expr.ModuleName
		modulePath := filepath.Join(currentDir, moduleName+".z")
		if _, err := os.Stat(modulePath); err == nil {
			c.imports[moduleName] = modulePath
		}
	}
}

// HandleExportStmt 处理export语句，收集导出的变量
func (c *CompletionCollector) HandleExportStmt(stmt *parser.ExportStmt, scope *scopeStack) {
	if m, ok := stmt.Result.(*parser.MapLit); ok {
		// 解析MapLit中的所有元素
		for _, element := range m.Elements {
			if ident, ok := element.Key.(*parser.Ident); ok {
				// 如果变量还没有加入到列表中，则添加
				if !c.varMap[ident.Name] {
					*c.variables = append(*c.variables, ident.Name)
					c.varMap[ident.Name] = true
				}
			}
		}
	}
}

func onCompletionFunc(context *glsp.Context, params *protocol.CompletionParams) (any, error) {
	keywordKind := protocol.CompletionItemKindKeyword
	// 简单判断触发补全的上下文
	var items []protocol.CompletionItem
	for _, keyword := range keywords {
		items = append(items, protocol.CompletionItem{
			Label: keyword,
			Kind:  &keywordKind,
			Data:  "keyword-" + keyword,
		})
	}
	funcKind := protocol.CompletionItemKindFunction
	snippetFormat := protocol.InsertTextFormatSnippet
	for _, f := range builtinFunctions {
		text := f + `($1)`
		items = append(items, protocol.CompletionItem{
			Label:            f,
			Kind:             &funcKind,
			InsertText:       &text,
			InsertTextFormat: &snippetFormat,
			Data:             "function-" + f,
		})
	}
	constKind := protocol.CompletionItemKindConstant
	for _, constant := range constants {
		items = append(items, protocol.CompletionItem{
			Label: constant,
			Kind:  &constKind,
			Data:  "constant-" + constant,
		})
	}

	// 添加模块和模块内变量的补全
	filename := strings.ReplaceAll(params.TextDocument.URI, "file://", "")

	// 查找当前文件中的变量并添加到补全项
	var variables []string
	// 使用新的遍历器收集所有变量和导入语句
	scopeStack := newScopeStack()
	collector := NewCompletionCollector(&variables)
	collector.SetFilename(filename)

	// 读取文件内容
	content, err := os.ReadFile(filename)
	if err != nil {
		return items, nil
	}

	// 使用新的通用文件解析函数
	parsedFile, err := ParseFileContent(filename, string(content))
	if err != nil {
		return items, nil
	}

	traverser := NewTraverser(collector)
	traverser.TraverseStmts(parsedFile.Stmts, scopeStack)

	// 获取导入的模块
	importedModules := collector.GetImports()
	for moduleName, modulePath := range importedModules {
		// 添加模块名作为补全项
		moduleKind := protocol.CompletionItemKindModule
		items = append(items, protocol.CompletionItem{
			Label: moduleName,
			Kind:  &moduleKind,
			Data:  "module-" + moduleName,
		})

		// 查找模块中的导出变量
		exportedVars := findExportedVariables(modulePath)
		for _, variable := range exportedVars {
			// 添加module.variable形式的补全项
			label := moduleName + "." + variable
			detail := moduleName + " module variable"
			varKind := protocol.CompletionItemKindVariable
			items = append(items, protocol.CompletionItem{
				Label:  label,
				Kind:   &varKind,
				Detail: &detail,
				Data:   "module-variable-" + label,
			})
		}
	}

	varKind := protocol.CompletionItemKindVariable
	for _, variable := range variables {
		items = append(items, protocol.CompletionItem{
			Label: variable,
			Kind:  &varKind,
			Data:  "local-variable-" + variable,
		})
	}

	return items, nil
}

// findExportedVariables 查找模块中导出的变量
func findExportedVariables(modulePath string) []string {
	var variables []string

	// 读取模块文件
	content, err := os.ReadFile(modulePath)
	if err != nil {
		return variables
	}

	// 使用新的通用文件解析函数
	parsedFile, err := ParseFileContent(modulePath, string(content))
	if err != nil {
		return variables
	}

	// 使用CompletionCollector收集导出的变量
	collector := NewCompletionCollector(&variables)
	scopeStack := newScopeStack()
	traverser := NewTraverser(collector)
	traverser.TraverseStmts(parsedFile.Stmts, scopeStack)

	return variables
}
