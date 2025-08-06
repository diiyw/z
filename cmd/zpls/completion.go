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

// VariableCollector 实现TraverserHandler接口，用于收集变量
type VariableCollector struct {
	variables *[]string
	varMap    map[string]bool
}

// NewVariableCollector 创建一个新的VariableCollector
func NewVariableCollector(variables *[]string) *VariableCollector {
	// 使用map去重
	varMap := make(map[string]bool)
	for _, v := range *variables {
		varMap[v] = true
	}

	return &VariableCollector{
		variables: variables,
		varMap:    varMap,
	}
}

func (vc *VariableCollector) HandleAssignStmt(stmt *parser.AssignStmt, scope *scopeStack) {
	// 添加左侧的变量到当前作用域
	for _, lh := range stmt.LHS {
		if ident, ok := lh.(*parser.Ident); ok {
			scope.addVariable(ident.Name)
			// 如果变量还没有加入到列表中，则添加
			if !vc.varMap[ident.Name] {
				*vc.variables = append(*vc.variables, ident.Name)
				vc.varMap[ident.Name] = true
			}
		}
	}
}

func (vc *VariableCollector) HandleBlockStmt(stmt *parser.BlockStmt, scope *scopeStack) {
	// 空实现
}

func (vc *VariableCollector) HandleFuncLit(stmt *parser.FuncLit, scope *scopeStack) {
	// 添加参数到作用域
	if stmt.Type.Params != nil {
		for _, param := range stmt.Type.Params.List {
			scope.addVariable(param.Name)
			// 如果变量还没有加入到列表中，则添加
			if !vc.varMap[param.Name] {
				*vc.variables = append(*vc.variables, param.Name)
				vc.varMap[param.Name] = true
			}
		}
	}
}

func (vc *VariableCollector) HandleForInStmt(stmt *parser.ForInStmt, scope *scopeStack) {
	// 添加key和value到作用域
	scope.addVariable(stmt.Key.Name)
	// 如果变量还没有加入到列表中，则添加
	if !vc.varMap[stmt.Key.Name] {
		*vc.variables = append(*vc.variables, stmt.Key.Name)
		vc.varMap[stmt.Key.Name] = true
	}

	scope.addVariable(stmt.Value.Name)
	// 如果变量还没有加入到列表中，则添加
	if !vc.varMap[stmt.Value.Name] {
		*vc.variables = append(*vc.variables, stmt.Value.Name)
		vc.varMap[stmt.Value.Name] = true
	}
}

func (vc *VariableCollector) HandleIdent(ident *parser.Ident, scope *scopeStack) {
	// 空实现
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

	// 查找当前文件中的import语句
	importedModules := findImportedModules(filename)
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

	// 查找当前文件中的变量并添加到补全项
	variables := findLocalVariables(filename)
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

// findImportedModules 查找文件中导入的模块
func findImportedModules(filename string) map[string]string {
	imports := make(map[string]string)

	// 读取文件内容
	content, err := os.ReadFile(filename)
	if err != nil {
		return imports
	}

	// 使用新的通用文件解析函数
	parsedFile, err := ParseFileContent(filename, string(content))
	if err != nil {
		return imports
	}

	// 遍历语法树查找import语句
	currentDir := filepath.Dir(filename)
	for _, stmt := range parsedFile.Stmts {
		if exprStmt, ok := stmt.(*parser.ExprStmt); ok {
			if importExpr, ok := exprStmt.Expr.(*parser.ImportExpr); ok {
				moduleName := importExpr.ModuleName
				modulePath := filepath.Join(currentDir, moduleName+".z")
				if _, err := os.Stat(modulePath); err == nil {
					imports[moduleName] = modulePath
				}
			}
		}
	}

	return imports
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

	// 遍历语法树查找export语句
	for _, stmt := range parsedFile.Stmts {
		if exportStmt, ok := stmt.(*parser.ExportStmt); ok {
			if m, ok := exportStmt.Result.(*parser.MapLit); ok {
				// 解析MapLit中的所有元素
				for _, element := range m.Elements {
					if ident, ok := element.Key.(*parser.Ident); ok {
						variables = append(variables, ident.Name)
					}
				}
			}
		}
	}

	return variables
}

// findLocalVariables 查找文件中的本地变量
func findLocalVariables(filename string) []string {
	var variables []string

	// 读取文件内容
	content, err := os.ReadFile(filename)
	if err != nil {
		return variables
	}

	// 使用新的通用文件解析函数
	parsedFile, err := ParseFileContent(filename, string(content))
	if err != nil {
		return variables
	}

	// 创建一个作用域栈来追踪变量定义
	scopeStack := newScopeStack()

	// 使用新的遍历器收集所有变量
	collector := NewVariableCollector(&variables)
	traverser := NewTraverser(collector)
	traverser.TraverseStmts(parsedFile.Stmts, scopeStack)

	return variables
}