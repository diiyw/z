package stdlib

import (
	"fmt"

	"github.com/d5/tengo/v2"
)

var fmtModule = map[string]z.Object{
	"print":   &z.UserFunction{Name: "print", Value: fmtPrint},
	"printf":  &z.UserFunction{Name: "printf", Value: fmtPrintf},
	"println": &z.UserFunction{Name: "println", Value: fmtPrintln},
	"sprintf": &z.UserFunction{Name: "sprintf", Value: fmtSprintf},
}

func fmtPrint(args ...z.Object) (ret z.Object, err error) {
	printArgs, err := getPrintArgs(args...)
	if err != nil {
		return nil, err
	}
	_, _ = fmt.Print(printArgs...)
	return nil, nil
}

func fmtPrintf(args ...z.Object) (ret z.Object, err error) {
	numArgs := len(args)
	if numArgs == 0 {
		return nil, z.ErrWrongNumArguments
	}

	format, ok := args[0].(*z.String)
	if !ok {
		return nil, z.ErrInvalidArgumentType{
			Name:     "format",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}
	if numArgs == 1 {
		fmt.Print(format)
		return nil, nil
	}

	s, err := z.Format(format.Value, args[1:]...)
	if err != nil {
		return nil, err
	}
	fmt.Print(s)
	return nil, nil
}

func fmtPrintln(args ...z.Object) (ret z.Object, err error) {
	printArgs, err := getPrintArgs(args...)
	if err != nil {
		return nil, err
	}
	printArgs = append(printArgs, "\n")
	_, _ = fmt.Print(printArgs...)
	return nil, nil
}

func fmtSprintf(args ...z.Object) (ret z.Object, err error) {
	numArgs := len(args)
	if numArgs == 0 {
		return nil, z.ErrWrongNumArguments
	}

	format, ok := args[0].(*z.String)
	if !ok {
		return nil, z.ErrInvalidArgumentType{
			Name:     "format",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}
	if numArgs == 1 {
		// okay to return 'format' directly as String is immutable
		return format, nil
	}
	s, err := z.Format(format.Value, args[1:]...)
	if err != nil {
		return nil, err
	}
	return &z.String{Value: s}, nil
}

func getPrintArgs(args ...z.Object) ([]interface{}, error) {
	var printArgs []interface{}
	l := 0
	for _, arg := range args {
		s, _ := z.ToString(arg)
		slen := len(s)
		// make sure length does not exceed the limit
		if l+slen > z.MaxStringLen {
			return nil, z.ErrStringLimit
		}
		l += slen
		printArgs = append(printArgs, s)
	}
	return printArgs, nil
}
