package main

import (
	"fmt"
	"os"

	"github.com/diiyw/z/parser"
)

func main() {
	code, err := Fmt(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Stdout.WriteString(code)
}

// Fmt 将AST打印为格式化的代码
func Fmt(code string) (string, error) {
	fileSet := parser.NewFileSet()
	sf := fileSet.AddFile("fmt", -1, len(code))
	p := parser.NewParser(sf, []byte(code), nil)
	fi, err := p.ParseFile()
	if err != nil {
		return "", err
	}
	return fi.String(), nil
}
