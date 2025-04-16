package stdlib

import (
	"bytes"
	gojson "encoding/json"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib/json"
)

var jsonModule = map[string]z.Object{
	"decode": &z.UserFunction{
		Name:  "decode",
		Value: jsonDecode,
	},
	"encode": &z.UserFunction{
		Name:  "encode",
		Value: jsonEncode,
	},
	"indent": &z.UserFunction{
		Name:  "encode",
		Value: jsonIndent,
	},
	"html_escape": &z.UserFunction{
		Name:  "html_escape",
		Value: jsonHTMLEscape,
	},
}

func jsonDecode(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		return nil, z.ErrWrongNumArguments
	}

	switch o := args[0].(type) {
	case *z.Bytes:
		v, err := json.Decode(o.Value)
		if err != nil {
			return &z.Error{
				Value: &z.String{Value: err.Error()},
			}, nil
		}
		return v, nil
	case *z.String:
		v, err := json.Decode([]byte(o.Value))
		if err != nil {
			return &z.Error{
				Value: &z.String{Value: err.Error()},
			}, nil
		}
		return v, nil
	default:
		return nil, z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "bytes/string",
			Found:    args[0].TypeName(),
		}
	}
}

func jsonEncode(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		return nil, z.ErrWrongNumArguments
	}

	b, err := json.Encode(args[0])
	if err != nil {
		return &z.Error{Value: &z.String{Value: err.Error()}}, nil
	}

	return &z.Bytes{Value: b}, nil
}

func jsonIndent(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 3 {
		return nil, z.ErrWrongNumArguments
	}

	prefix, ok := z.ToString(args[1])
	if !ok {
		return nil, z.ErrInvalidArgumentType{
			Name:     "prefix",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	indent, ok := z.ToString(args[2])
	if !ok {
		return nil, z.ErrInvalidArgumentType{
			Name:     "indent",
			Expected: "string(compatible)",
			Found:    args[2].TypeName(),
		}
	}

	switch o := args[0].(type) {
	case *z.Bytes:
		var dst bytes.Buffer
		err := gojson.Indent(&dst, o.Value, prefix, indent)
		if err != nil {
			return &z.Error{
				Value: &z.String{Value: err.Error()},
			}, nil
		}
		return &z.Bytes{Value: dst.Bytes()}, nil
	case *z.String:
		var dst bytes.Buffer
		err := gojson.Indent(&dst, []byte(o.Value), prefix, indent)
		if err != nil {
			return &z.Error{
				Value: &z.String{Value: err.Error()},
			}, nil
		}
		return &z.Bytes{Value: dst.Bytes()}, nil
	default:
		return nil, z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "bytes/string",
			Found:    args[0].TypeName(),
		}
	}
}

func jsonHTMLEscape(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		return nil, z.ErrWrongNumArguments
	}

	switch o := args[0].(type) {
	case *z.Bytes:
		var dst bytes.Buffer
		gojson.HTMLEscape(&dst, o.Value)
		return &z.Bytes{Value: dst.Bytes()}, nil
	case *z.String:
		var dst bytes.Buffer
		gojson.HTMLEscape(&dst, []byte(o.Value))
		return &z.Bytes{Value: dst.Bytes()}, nil
	default:
		return nil, z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "bytes/string",
			Found:    args[0].TypeName(),
		}
	}
}
