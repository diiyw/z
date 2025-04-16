package z_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/parser"
	"github.com/d5/tengo/v2/require"
)

type srcfile struct {
	name string
	size int
}

func TestBytecode(t *testing.T) {
	testBytecodeSerialization(t, bytecode(concatInsts(), objectsArray()))

	testBytecodeSerialization(t, bytecode(
		concatInsts(), objectsArray(
			&z.Char{Value: 'y'},
			&z.Float{Value: 93.11},
			compiledFunction(1, 0,
				z.MakeInstruction(parser.OpConstant, 3),
				z.MakeInstruction(parser.OpSetLocal, 0),
				z.MakeInstruction(parser.OpGetGlobal, 0),
				z.MakeInstruction(parser.OpGetFree, 0)),
			&z.Float{Value: 39.2},
			&z.Int{Value: 192},
			&z.String{Value: "bar"})))

	testBytecodeSerialization(t, bytecodeFileSet(
		concatInsts(
			z.MakeInstruction(parser.OpConstant, 0),
			z.MakeInstruction(parser.OpSetGlobal, 0),
			z.MakeInstruction(parser.OpConstant, 6),
			z.MakeInstruction(parser.OpPop)),
		objectsArray(
			&z.Int{Value: 55},
			&z.Int{Value: 66},
			&z.Int{Value: 77},
			&z.Int{Value: 88},
			&z.ImmutableMap{
				Value: map[string]z.Object{
					"array": &z.ImmutableArray{
						Value: []z.Object{
							&z.Int{Value: 1},
							&z.Int{Value: 2},
							&z.Int{Value: 3},
							z.TrueValue,
							z.FalseValue,
							z.UndefinedValue,
						},
					},
					"true":  z.TrueValue,
					"false": z.FalseValue,
					"bytes": &z.Bytes{Value: make([]byte, 16)},
					"char":  &z.Char{Value: 'Y'},
					"error": &z.Error{Value: &z.String{
						Value: "some error",
					}},
					"float": &z.Float{Value: -19.84},
					"immutable_array": &z.ImmutableArray{
						Value: []z.Object{
							&z.Int{Value: 1},
							&z.Int{Value: 2},
							&z.Int{Value: 3},
							z.TrueValue,
							z.FalseValue,
							z.UndefinedValue,
						},
					},
					"immutable_map": &z.ImmutableMap{
						Value: map[string]z.Object{
							"a": &z.Int{Value: 1},
							"b": &z.Int{Value: 2},
							"c": &z.Int{Value: 3},
							"d": z.TrueValue,
							"e": z.FalseValue,
							"f": z.UndefinedValue,
						},
					},
					"int": &z.Int{Value: 91},
					"map": &z.Map{
						Value: map[string]z.Object{
							"a": &z.Int{Value: 1},
							"b": &z.Int{Value: 2},
							"c": &z.Int{Value: 3},
							"d": z.TrueValue,
							"e": z.FalseValue,
							"f": z.UndefinedValue,
						},
					},
					"string":    &z.String{Value: "foo bar"},
					"time":      &z.Time{Value: time.Now()},
					"undefined": z.UndefinedValue,
				},
			},
			compiledFunction(1, 0,
				z.MakeInstruction(parser.OpConstant, 3),
				z.MakeInstruction(parser.OpSetLocal, 0),
				z.MakeInstruction(parser.OpGetGlobal, 0),
				z.MakeInstruction(parser.OpGetFree, 0),
				z.MakeInstruction(parser.OpBinaryOp, 11),
				z.MakeInstruction(parser.OpGetFree, 1),
				z.MakeInstruction(parser.OpBinaryOp, 11),
				z.MakeInstruction(parser.OpGetLocal, 0),
				z.MakeInstruction(parser.OpBinaryOp, 11),
				z.MakeInstruction(parser.OpReturn, 1)),
			compiledFunction(1, 0,
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpSetLocal, 0),
				z.MakeInstruction(parser.OpGetFree, 0),
				z.MakeInstruction(parser.OpGetLocal, 0),
				z.MakeInstruction(parser.OpClosure, 4, 2),
				z.MakeInstruction(parser.OpReturn, 1)),
			compiledFunction(1, 0,
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpSetLocal, 0),
				z.MakeInstruction(parser.OpGetLocal, 0),
				z.MakeInstruction(parser.OpClosure, 5, 1),
				z.MakeInstruction(parser.OpReturn, 1))),
		fileSet(srcfile{name: "file1", size: 100},
			srcfile{name: "file2", size: 200})))
}

func TestBytecode_RemoveDuplicates(t *testing.T) {
	testBytecodeRemoveDuplicates(t,
		bytecode(
			concatInsts(), objectsArray(
				&z.Char{Value: 'y'},
				&z.Float{Value: 93.11},
				compiledFunction(1, 0,
					z.MakeInstruction(parser.OpConstant, 3),
					z.MakeInstruction(parser.OpSetLocal, 0),
					z.MakeInstruction(parser.OpGetGlobal, 0),
					z.MakeInstruction(parser.OpGetFree, 0)),
				&z.Float{Value: 39.2},
				&z.Int{Value: 192},
				&z.String{Value: "bar"})),
		bytecode(
			concatInsts(), objectsArray(
				&z.Char{Value: 'y'},
				&z.Float{Value: 93.11},
				compiledFunction(1, 0,
					z.MakeInstruction(parser.OpConstant, 3),
					z.MakeInstruction(parser.OpSetLocal, 0),
					z.MakeInstruction(parser.OpGetGlobal, 0),
					z.MakeInstruction(parser.OpGetFree, 0)),
				&z.Float{Value: 39.2},
				&z.Int{Value: 192},
				&z.String{Value: "bar"})))

	testBytecodeRemoveDuplicates(t,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpConstant, 3),
				z.MakeInstruction(parser.OpConstant, 4),
				z.MakeInstruction(parser.OpConstant, 5),
				z.MakeInstruction(parser.OpConstant, 6),
				z.MakeInstruction(parser.OpConstant, 7),
				z.MakeInstruction(parser.OpConstant, 8),
				z.MakeInstruction(parser.OpClosure, 4, 1)),
			objectsArray(
				&z.Int{Value: 1},
				&z.Float{Value: 2.0},
				&z.Char{Value: '3'},
				&z.String{Value: "four"},
				compiledFunction(1, 0,
					z.MakeInstruction(parser.OpConstant, 3),
					z.MakeInstruction(parser.OpConstant, 7),
					z.MakeInstruction(parser.OpSetLocal, 0),
					z.MakeInstruction(parser.OpGetGlobal, 0),
					z.MakeInstruction(parser.OpGetFree, 0)),
				&z.Int{Value: 1},
				&z.Float{Value: 2.0},
				&z.Char{Value: '3'},
				&z.String{Value: "four"})),
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpConstant, 3),
				z.MakeInstruction(parser.OpConstant, 4),
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpConstant, 3),
				z.MakeInstruction(parser.OpClosure, 4, 1)),
			objectsArray(
				&z.Int{Value: 1},
				&z.Float{Value: 2.0},
				&z.Char{Value: '3'},
				&z.String{Value: "four"},
				compiledFunction(1, 0,
					z.MakeInstruction(parser.OpConstant, 3),
					z.MakeInstruction(parser.OpConstant, 2),
					z.MakeInstruction(parser.OpSetLocal, 0),
					z.MakeInstruction(parser.OpGetGlobal, 0),
					z.MakeInstruction(parser.OpGetFree, 0)))))

	testBytecodeRemoveDuplicates(t,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpConstant, 3),
				z.MakeInstruction(parser.OpConstant, 4)),
			objectsArray(
				&z.Int{Value: 1},
				&z.Int{Value: 2},
				&z.Int{Value: 3},
				&z.Int{Value: 1},
				&z.Int{Value: 3})),
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 2)),
			objectsArray(
				&z.Int{Value: 1},
				&z.Int{Value: 2},
				&z.Int{Value: 3})))
}

func TestBytecode_CountObjects(t *testing.T) {
	b := bytecode(
		concatInsts(),
		objectsArray(
			&z.Int{Value: 55},
			&z.Int{Value: 66},
			&z.Int{Value: 77},
			&z.Int{Value: 88},
			compiledFunction(1, 0,
				z.MakeInstruction(parser.OpConstant, 3),
				z.MakeInstruction(parser.OpReturn, 1)),
			compiledFunction(1, 0,
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpReturn, 1)),
			compiledFunction(1, 0,
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpReturn, 1))))
	require.Equal(t, 7, b.CountObjects())
}

func fileSet(files ...srcfile) *parser.SourceFileSet {
	fileSet := parser.NewFileSet()
	for _, f := range files {
		fileSet.AddFile(f.name, -1, f.size)
	}
	return fileSet
}

func bytecodeFileSet(
	instructions []byte,
	constants []z.Object,
	fileSet *parser.SourceFileSet,
) *z.Bytecode {
	return &z.Bytecode{
		FileSet:      fileSet,
		MainFunction: &z.CompiledFunction{Instructions: instructions},
		Constants:    constants,
	}
}

func testBytecodeRemoveDuplicates(
	t *testing.T,
	input, expected *z.Bytecode,
) {
	input.RemoveDuplicates()

	require.Equal(t, expected.FileSet, input.FileSet)
	require.Equal(t, expected.MainFunction, input.MainFunction)
	require.Equal(t, expected.Constants, input.Constants)
}

func testBytecodeSerialization(t *testing.T, b *z.Bytecode) {
	var buf bytes.Buffer
	err := b.Encode(&buf)
	require.NoError(t, err)

	r := &z.Bytecode{}
	err = r.Decode(bytes.NewReader(buf.Bytes()), nil)
	require.NoError(t, err)

	require.Equal(t, b.FileSet, r.FileSet)
	require.Equal(t, b.MainFunction, r.MainFunction)
	require.Equal(t, b.Constants, r.Constants)
}
