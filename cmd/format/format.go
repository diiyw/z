package format

import (
	"strings"

	"github.com/diiyw/z/parser"
)

type printer struct {
	level  int
	result string
}

func (p *printer) printStmts(stmts []parser.Stmt) {
	for _, stmt := range stmts {
		p.printStmt(stmt)
	}
}

func (p *printer) printStmt(stmt parser.Stmt) {
	switch s := stmt.(type) {
	case *parser.AssignStmt:
		p.printAssignStmt(s)
	case *parser.ExportStmt:
		p.printExportStmt(s)
	case *parser.BlockStmt:
		p.printBlockStmt(s)
	case *parser.ExprStmt:
		p.printExprStmt(s)
	case *parser.ForInStmt:
		p.printForInStmt(s)
	case *parser.ForStmt:
		p.printForStmt(s)
	case *parser.IfStmt:
		p.printIfStmt(s)
	case *parser.IncDecStmt:
		p.printIncDecStmt(s)
	case *parser.ReturnStmt:
		p.printReturnStmt(s)
	}
}

func (p *printer) printAssignStmt(s *parser.AssignStmt) {
	var lhs, rhs []string
	for _, e := range s.LHS {
		lhs = append(lhs, p.printExpr(e))
	}
	for _, e := range s.RHS {
		rhs = append(rhs, p.printExpr(e))
	}
	p.printLine(strings.Join(lhs, ", ")+" "+s.Token.String()+
		" "+strings.Join(rhs, ", "), true)
}

func (p *printer) printExportStmt(s *parser.ExportStmt) {
	result := p.printExpr(s.Result)
	p.printLine("export "+result, true)
}

func (p *printer) printBlockStmt(s *parser.BlockStmt) {
	p.printLine("{", false)
	if len(s.Stmts) > 0 {
		p.level++
		p.printStmts(s.Stmts)
		p.level--
	}
	p.printLine("}", true)
}

func (p *printer) printExprStmt(s *parser.ExprStmt) {
	p.printLine(p.printExpr(s.Expr), true)
}

func (p *printer) printForInStmt(s *parser.ForInStmt) {
	if s.Value != nil {
		p.print("for "+p.printExpr(s.Key)+", "+p.printExpr(s.Value)+
			" in "+p.printExpr(s.Iterable)+" ", false)
		p.printStmt(s.Body)
		return
	}
	p.printLine("for "+p.printExpr(s.Key)+" in "+p.printExpr(s.Iterable)+" ", true)
	p.printStmt(s.Body)
}

func (p *printer) printForStmt(s *parser.ForStmt) {
	var cond string
	if s.Cond != nil {
		cond = p.printExpr(s.Cond)
	}
	p.print("for ", true)
	if s.Init != nil && s.Post != nil {
		p.printStmt(s.Init)
		p.trim()
		p.printLine("; "+cond+"; ", false)
		p.trim()
		p.printStmt(s.Post)
		p.trim()
		p.print(" ", false)
	} else {
		if cond != "" {
			p.print(cond+" ", false)
		}
	}
	p.printBlockStmt(s.Body)
}

func (p *printer) printIfStmt(s *parser.IfStmt) {
	p.print("if ", true)
	if s.Init != nil {
		p.printStmt(s.Init)
		p.trim()
		p.print("; ", false)
	}
	p.print(p.printExpr(s.Cond)+" ", false)
	p.printStmt(s.Body)
	if s.Else != nil {
		p.trim()
		p.print(" else ", true)
		p.printStmt(s.Else)
	}
}

func (p *printer) printIncDecStmt(s *parser.IncDecStmt) {
	p.print(p.printExpr(s.Expr), true)
	p.printLine(s.Token.String(), true)
}

func (p *printer) printReturnStmt(s *parser.ReturnStmt) {
	if s.Result != nil {
		p.printLine("return "+p.printExpr(s.Result), true)
		return
	}
	p.printLine("return", true)
}

func (p *printer) printExpr(expr parser.Expr) string {
	switch e := expr.(type) {
	case *parser.ArrayLit:
		return p.printArrayLit(e)
	case *parser.BadExpr:
		return p.printBadExpr(e)
	case *parser.BinaryExpr:
		return p.printBinaryExpr(e)
	case *parser.BoolLit:
		return p.printBoolLit(e)
	case *parser.CallExpr:
		return p.printCallExpr(e)
	case *parser.CharLit:
		return p.printCharLit(e)
	case *parser.CondExpr:
		return p.printCondExpr(e)
	case *parser.ErrorExpr:
		return p.printErrorExpr(e)
	case *parser.FloatLit:
		return p.printFloatLit(e)
	case *parser.FuncLit:
		return p.printFuncLit(e)
	case *parser.FuncType:
		return p.printFuncType(e)
	case *parser.Ident:
		return p.printIdent(e)
	case *parser.ImmutableExpr:
		return p.printImmutableExpr(e)
	case *parser.ImportExpr:
		return p.printImportExpr(e)
	case *parser.IndexExpr:
		return p.printIndexExpr(e)
	case *parser.IntLit:
		return p.printIntLit(e)
	case *parser.MapElementLit:
		return p.printMapElementLit(e)
	case *parser.MapLit:
		return p.printMapLit(e)
	case *parser.ParenExpr:
		return p.printParenExpr(e)
	case *parser.SelectorExpr:
		return p.printSelectorExpr(e)
	case *parser.SliceExpr:
		return p.printSliceExpr(e)
	case *parser.StringLit:
		return p.printStringLit(e)
	case *parser.UnaryExpr:
		return p.printUnaryExpr(e)
	case *parser.UndefinedLit:
		return p.printUndefinedLit(e)
	}
	return ""
}

func (p *printer) printArrayLit(e *parser.ArrayLit) string {
	var elements []string
	for _, m := range e.Elements {
		elements = append(elements, p.printExpr(m))
	}
	return "[" + strings.Join(elements, ", ") + "]"
}

func (p *printer) printBadExpr(e *parser.BadExpr) string {
	return ""
}

func (p *printer) printBinaryExpr(e *parser.BinaryExpr) string {
	return p.printExpr(e.LHS) + " " + e.Token.String() +
		" " + p.printExpr(e.RHS)
}

func (p *printer) printBoolLit(e *parser.BoolLit) string {
	return e.String()
}

func (p *printer) printCallExpr(e *parser.CallExpr) string {
	var args []string
	for _, e := range e.Args {
		args = append(args, strings.TrimSuffix(p.printExpr(e), "\n"))
	}
	if len(args) > 0 && e.Ellipsis.IsValid() {
		args[len(args)-1] = args[len(args)-1] + "..."
	}
	return e.Func.String() + "(" + strings.Join(args, ", ") + ")"
}

func (p *printer) printCharLit(e *parser.CharLit) string {
	return e.String()
}

func (p *printer) printCondExpr(e *parser.CondExpr) string {
	return "(" + p.printExpr(e.Cond) + " ? " + p.printExpr(e.True) +
		" : " + p.printExpr(e.False) + ")"
}

func (p *printer) printErrorExpr(e *parser.ErrorExpr) string {
	return "error(" + p.printExpr(e.Expr) + ")"
}

func (p *printer) printFloatLit(e *parser.FloatLit) string {
	return e.String()
}

func (p *printer) printFuncLit(e *parser.FuncLit) string {
	var np = &printer{level: p.level}
	np.printBlockStmt(e.Body)
	return "func" + e.Type.Params.String() + " " + np.result
}

func (p *printer) printFuncType(e *parser.FuncType) string {
	return e.String()
}

func (p *printer) printIdent(e *parser.Ident) string {
	return e.String()
}

func (p *printer) printImmutableExpr(e *parser.ImmutableExpr) string {
	return "immutable(" + p.printExpr(e.Expr) + ")"
}

func (p *printer) printImportExpr(e *parser.ImportExpr) string {
	return e.String()
}

func (p *printer) printIndexExpr(e *parser.IndexExpr) string {
	var index string
	if e.Index != nil {
		index = p.printExpr(e.Index)
	}
	return p.printExpr(e.Expr) + "[" + index + "]"
}

func (p *printer) printIntLit(e *parser.IntLit) string {
	return e.String()
}

func (p *printer) printMapElementLit(e *parser.MapElementLit) string {
	return p.printExpr(e.Key) + ": " + p.printExpr(e.Value)
}

func (p *printer) printMapLit(e *parser.MapLit) string {
	var elements []string
	p.level++
	subSpaces := strings.Repeat("\t", p.level)
	for _, m := range e.Elements {
		elements = append(elements, strings.TrimSuffix(subSpaces+p.printExpr(m), "\n"))
	}
	p.level--
	if len(elements) == 0 {
		return "{}"
	}
	spaces := strings.Repeat("\t", p.level)
	return "{\n" + strings.Join(elements, ",\n") + "\n" + spaces + "}"
}

func (p *printer) printParenExpr(e *parser.ParenExpr) string {
	return "(" + p.printExpr(e.Expr) + ")"
}

func (p *printer) printSelectorExpr(e *parser.SelectorExpr) string {
	return p.printExpr(e.Expr) + "." + p.printExpr(e.Sel)
}

func (p *printer) printSliceExpr(e *parser.SliceExpr) string {
	var low, high string
	if e.Low != nil {
		low = p.printExpr(e.Low)
	}
	if e.High != nil {
		high = p.printExpr(e.High)
	}
	return p.printExpr(e.Expr) + "[" + low + ":" + high + "]"
}

func (p *printer) printStringLit(e *parser.StringLit) string {
	return e.String()
}

func (p *printer) printUnaryExpr(e *parser.UnaryExpr) string {
	return e.Token.String() + p.printExpr(e.Expr)
}

func (p *printer) printUndefinedLit(e *parser.UndefinedLit) string {
	return e.String()
}

func (p *printer) trim() {
	p.result = strings.TrimSuffix(p.result, "\n")
}

func (p *printer) print(v string, prefix bool) {
	if prefix {
		p.result += strings.Repeat("\t", p.level)
	}
	p.result += v
}

func (p *printer) printLine(v string, prefix bool) {
	if prefix {
		p.result += strings.Repeat("\t", p.level)
	}
	p.result += v + "\n"
}

func Format(src []byte) (string, error) {
	fileSet := parser.NewFileSet()
	sourceFile := fileSet.AddFile("fmt.z", -1, len(src))
	p := parser.NewParser(sourceFile, src, nil)
	parsedFile, err := p.ParseFile()
	if err != nil {
		return "", err
	}
	var pt = &printer{}
	pt.printStmts(parsedFile.Stmts)
	pt.trim()
	return pt.result, nil
}
