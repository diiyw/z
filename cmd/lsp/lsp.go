package lsp

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/diiyw/z/cmd/format"
	"github.com/diiyw/z/parser"
)

// OnDocumentFormatting handles the 'fmt' command for formatting Z code
func OnDocumentFormatting() {
	// Create a new flag set for the fmt command
	fmtFlags := flag.NewFlagSet("fmt", flag.ExitOnError)
	textFlag := fmtFlags.Bool("text", false, "Format Z code from text input")

	// Parse flags starting from the second argument (after "fmt")
	err := fmtFlags.Parse(flag.Args()[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %s\n", err)
		os.Exit(1)
	}

	if *textFlag {
		// Format from stdin text input
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err)
			os.Exit(1)
		}

		result, err := format.Format(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Format error: %s\n", err)
			os.Exit(1)
		}
		fmt.Print(result)
	} else {
		// Format from file (existing behavior)
		if fmtFlags.NArg() < 1 {
			fmt.Fprintln(os.Stderr, "Usage: z formatting <file>")
			os.Exit(1)
		}
		file := fmtFlags.Arg(0)
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
			os.Exit(1)
		}
		result, err := format.Format(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Format error: %s\n", err)
			os.Exit(1)
		}
		fmt.Print(result)
	}
}

func OnDiagnostics() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err)
		os.Exit(1)
	}
	fileSet := parser.NewFileSet()
	sourceFile := fileSet.AddFile("fmt.z", -1, len(input))
	p := parser.NewParser(sourceFile, input, nil)
	_, err = p.ParseFile()
	if err != nil {
		var errList parser.ErrorList
		errors.As(err, &errList)
		data, _ := json.Marshal(errList)
		fmt.Println(string(data))
	}
}

type DefinitionItem struct {
	Name  string `json:"name"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

type Definition struct {
	Import  DefinitionItem   `json:"import"`
	Targets []DefinitionItem `json:"targets"`
}

func OnDefinition(input []byte) {
	var err error
	if input == nil {
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err)
			os.Exit(1)
		}
	}
	var m = make(map[string]any)
	if err = json.Unmarshal(input, &m); err != nil {
		//fmt.Fprintf(os.Stderr, "Error parsing JSON: %s\n", err)
		os.Exit(0)
	}
	code := m["code"].(string)
	offset := int(m["offset"].(float64))
	fileSet := parser.NewFileSet()
	sourceFile := fileSet.AddFile("definition.z", -1, len(code))
	p := parser.NewParser(sourceFile, []byte(code), nil)
	parsedFile, err := p.ParseFile()
	if err != nil {
		// fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err)
		os.Exit(0)
	}
	var def = &Definition{
		Targets: make([]DefinitionItem, 0),
	}
	var expr = findNode(parsedFile.Stmts, offset)
	if e, ok := expr.(*parser.ImportExpr); ok {
		def.Import = DefinitionItem{
			Name:  e.ModuleName + ".z",
			Start: int(e.Pos()) + 7,
			End:   int(e.End()) - 3,
		}
	}
	data, _ := json.Marshal(def)
	os.Stdout.Write(data)
}

func findNode(stmts []parser.Stmt, offset int) parser.Expr {
	for _, stmt := range stmts {
		if node := findExprNode(stmt, offset); node != nil {
			return node
		}
	}
	return nil
}

func findExprNode(stmt parser.Stmt, start int) parser.Expr {
	ss, se := stmt.Pos(), stmt.End()
	if int(ss) <= start && start < int(se) {
		switch stmt := stmt.(type) {
		case *parser.AssignStmt:
			for _, expr := range stmt.LHS {
				if int(expr.Pos()) <= start && start < int(expr.End()) {
					return expr
				}
			}
			for _, expr := range stmt.RHS {
				if int(expr.Pos()) <= start && start < int(expr.End()) {
					return expr
				}
			}
		case *parser.ExportStmt:
			expr := stmt.Result
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
		case *parser.BlockStmt:
			for _, expr := range stmt.Stmts {
				if int(expr.Pos()) <= start && start < int(expr.End()) {
					return findExprNode(expr, start)
				}
			}
		case *parser.ExprStmt:
			expr := stmt.Expr
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
		case *parser.ForInStmt:
			expr := stmt.Iterable
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
			expr = stmt.Key
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
			expr = stmt.Value
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
			return findExprNode(stmt.Body, start)
		case *parser.ForStmt:
			expr := stmt.Cond
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
			if node := findExprNode(stmt.Init, start); node != nil {
				return node
			}
			if node := findExprNode(stmt.Post, start); node != nil {
				return node
			}
			return findExprNode(stmt.Body, start)
		case *parser.IfStmt:
			expr := stmt.Cond
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
			if node := findExprNode(stmt.Init, start); node != nil {
				return node
			}
			if node := findExprNode(stmt.Else, start); node != nil {
				return node
			}
			return findExprNode(stmt.Body, start)
		case *parser.IncDecStmt:
			expr := stmt.Expr
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
		case *parser.ReturnStmt:
			expr := stmt.Result
			if int(expr.Pos()) <= start && start < int(expr.End()) {
				return expr
			}
		}
	}
	return nil
}
