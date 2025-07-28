package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/diiyw/z"
	"github.com/diiyw/z/cmd/format"
	"github.com/diiyw/z/parser"
	"github.com/diiyw/z/stdlib"
)

const (
	sourceFileExt = ".z"
	replPrompt    = ">> "
)

var (
	compileOutput string
	showHelp      bool
	showVersion   bool
	resolvePath   bool // TODO Remove this flag at version 3
	version       = "dev"
)

func init() {
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.StringVar(&compileOutput, "o", "", "Compile output file")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.BoolVar(&resolvePath, "resolve", false,
		"Resolve relative import paths")
	flag.Parse()
}

func main() {
	if showHelp {
		doHelp()
		os.Exit(2)
	} else if showVersion {
		fmt.Println(version)
		return
	}

	// 新增 fmt 命令支持
	if flag.NArg() > 0 && flag.Arg(0) == "fmt" {
		handleFmtCommand()
		return
	}

	if flag.NArg() > 0 && flag.Arg(0) == "check" {
		handleCheckCommand()
		return
	}

	modules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	inputFile := flag.Arg(0)
	if inputFile == "" {
		// REPL
		RunREPL(modules, os.Stdin, os.Stdout)
		return
	}

	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr,
			"Error reading input file: %s\n", err.Error())
		os.Exit(1)
	}

	inputFile, err = filepath.Abs(inputFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error file path: %s\n", err)
		os.Exit(1)
	}

	if len(inputData) > 1 && string(inputData[:2]) == "#!" {
		copy(inputData, "//")
	}

	if compileOutput != "" {
		err := CompileOnly(modules, inputData, inputFile,
			compileOutput)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	} else if filepath.Ext(inputFile) == sourceFileExt {
		err := CompileAndRun(modules, inputData, inputFile)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	} else {
		if err := RunCompiled(modules, inputData); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}

// CompileOnly compiles the source code and writes the compiled binary into
// outputFile.
func CompileOnly(
	modules *z.ModuleMap,
	data []byte,
	inputFile, outputFile string,
) (err error) {
	bytecode, err := compileSrc(modules, data, inputFile)
	if err != nil {
		return
	}

	if outputFile == "" {
		outputFile = basename(inputFile) + ".out"
	}

	out, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = out.Close()
		} else {
			err = out.Close()
		}
	}()

	err = bytecode.Encode(out)
	if err != nil {
		return
	}
	fmt.Println(outputFile)
	return
}

// CompileAndRun compiles the source code and executes it.
func CompileAndRun(
	modules *z.ModuleMap,
	data []byte,
	inputFile string,
) (err error) {
	bytecode, err := compileSrc(modules, data, inputFile)
	if err != nil {
		return
	}

	machine := z.NewVM(bytecode, nil, -1)
	err = machine.Run()
	return
}

// RunCompiled reads the compiled binary from file and executes it.
func RunCompiled(modules *z.ModuleMap, data []byte) (err error) {
	bytecode := &z.Bytecode{}
	err = bytecode.Decode(bytes.NewReader(data), modules)
	if err != nil {
		return
	}

	machine := z.NewVM(bytecode, nil, -1)
	err = machine.Run()
	return
}

// RunREPL starts REPL.
func RunREPL(modules *z.ModuleMap, in io.Reader, out io.Writer) {
	stdin := bufio.NewScanner(in)
	fileSet := parser.NewFileSet()
	globals := make([]z.Object, z.GlobalsSize)
	symbolTable := z.NewSymbolTable()
	for idx, fn := range z.GetAllBuiltinFunctions() {
		symbolTable.DefineBuiltin(idx, fn.Name)
	}

	// embed println function
	symbol := symbolTable.Define("__repl_println__")
	globals[symbol.Index] = &z.UserFunction{
		Name: "println",
		Value: func(args ...z.Object) (ret z.Object, err error) {
			var printArgs []any
			for _, arg := range args {
				if _, isUndefined := arg.(*z.Undefined); isUndefined {
					printArgs = append(printArgs, "<undefined>")
				} else {
					s, _ := z.ToString(arg)
					printArgs = append(printArgs, s)
				}
			}
			printArgs = append(printArgs, "\n")
			_, _ = fmt.Print(printArgs...)
			return
		},
	}

	var constants []z.Object
	for {
		_, _ = fmt.Fprint(out, replPrompt)
		scanned := stdin.Scan()
		if !scanned {
			return
		}

		line := stdin.Text()
		srcFile := fileSet.AddFile("repl", -1, len(line))
		p := parser.NewParser(srcFile, []byte(line), nil)
		file, err := p.ParseFile()
		if err != nil {
			_, _ = fmt.Fprintln(out, err.Error())
			continue
		}

		file = addPrints(file)
		c := z.NewCompiler(srcFile, symbolTable, constants, modules, nil)
		if err := c.Compile(file); err != nil {
			_, _ = fmt.Fprintln(out, err.Error())
			continue
		}

		bytecode := c.Bytecode()
		machine := z.NewVM(bytecode, globals, -1)
		if err := machine.Run(); err != nil {
			_, _ = fmt.Fprintln(out, err.Error())
			continue
		}
		constants = bytecode.Constants
	}
}

func compileSrc(
	modules *z.ModuleMap,
	src []byte,
	inputFile string,
) (*z.Bytecode, error) {
	fileSet := parser.NewFileSet()
	srcFile := fileSet.AddFile(filepath.Base(inputFile), -1, len(src))

	p := parser.NewParser(srcFile, src, nil)
	file, err := p.ParseFile()
	if err != nil {
		return nil, err
	}

	c := z.NewCompiler(srcFile, nil, nil, modules, nil)
	c.EnableFileImport(true)
	if resolvePath {
		c.SetImportDir(filepath.Dir(inputFile))
	}

	if err := c.Compile(file); err != nil {
		return nil, err
	}

	bytecode := c.Bytecode()
	bytecode.RemoveDuplicates()
	return bytecode, nil
}

func doHelp() {
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("	z [flags] {input-file}")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println()
	fmt.Println("	-o        compile output file")
	fmt.Println("	-version  show version")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println()
	fmt.Println("	z")
	fmt.Println()
	fmt.Println("	          Start Z REPL")
	fmt.Println("	z fmt myapp.z")
	fmt.Println()
	fmt.Println("	          Format z code (myapp.z)")
	fmt.Println("	z fmt -text")
	fmt.Println()
	fmt.Println("	          Format z code from stdin")
	fmt.Println("	z check")
	fmt.Println()
	fmt.Println("	          Check z code from stdin")
	fmt.Println()
	fmt.Println("	z myapp.z")
	fmt.Println()
	fmt.Println("	          Compile and run source file (myapp.z)")
	fmt.Println("	          Source file must have .z extension")
	fmt.Println()
	fmt.Println("	z -o myapp myapp.z")
	fmt.Println()
	fmt.Println("	          Compile source file (myapp.z) into bytecode file (myapp)")
	fmt.Println()
	fmt.Println("	z myapp")
	fmt.Println()
	fmt.Println("	          Run bytecode file (myapp)")
	fmt.Println()
	fmt.Println()
}

func addPrints(file *parser.File) *parser.File {
	var stmts []parser.Stmt
	for _, s := range file.Stmts {
		switch s := s.(type) {
		case *parser.ExprStmt:
			stmts = append(stmts, &parser.ExprStmt{
				Expr: &parser.CallExpr{
					Func: &parser.Ident{Name: "__repl_println__"},
					Args: []parser.Expr{s.Expr},
				},
			})
		case *parser.AssignStmt:
			stmts = append(stmts, s)

			stmts = append(stmts, &parser.ExprStmt{
				Expr: &parser.CallExpr{
					Func: &parser.Ident{
						Name: "__repl_println__",
					},
					Args: s.LHS,
				},
			})
		default:
			stmts = append(stmts, s)
		}
	}
	return &parser.File{
		InputFile: file.InputFile,
		Stmts:     stmts,
	}
}

func basename(s string) string {
	s = filepath.Base(s)
	n := strings.LastIndexByte(s, '.')
	if n > 0 {
		return s[:n]
	}
	return s
}

// handleFmtCommand handles the 'fmt' command for formatting Z code
func handleFmtCommand() {
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
			fmt.Fprintln(os.Stderr, "Usage: z fmt <file>")
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

func handleCheckCommand() {
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
