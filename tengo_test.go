package z_test

import (
	"strings"
	"testing"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/parser"
	"github.com/d5/tengo/v2/require"
)

func TestInstructions_String(t *testing.T) {
	assertInstructionString(t,
		[][]byte{
			z.MakeInstruction(parser.OpConstant, 1),
			z.MakeInstruction(parser.OpConstant, 2),
			z.MakeInstruction(parser.OpConstant, 65535),
		},
		`0000 CONST   1    
0003 CONST   2    
0006 CONST   65535`)

	assertInstructionString(t,
		[][]byte{
			z.MakeInstruction(parser.OpBinaryOp, 11),
			z.MakeInstruction(parser.OpConstant, 2),
			z.MakeInstruction(parser.OpConstant, 65535),
		},
		`0000 BINARYOP 11   
0002 CONST   2    
0005 CONST   65535`)

	assertInstructionString(t,
		[][]byte{
			z.MakeInstruction(parser.OpBinaryOp, 11),
			z.MakeInstruction(parser.OpGetLocal, 1),
			z.MakeInstruction(parser.OpConstant, 2),
			z.MakeInstruction(parser.OpConstant, 65535),
		},
		`0000 BINARYOP 11   
0002 GETL    1    
0004 CONST   2    
0007 CONST   65535`)
}

func TestMakeInstruction(t *testing.T) {
	makeInstruction(t, []byte{parser.OpConstant, 0, 0},
		parser.OpConstant, 0)
	makeInstruction(t, []byte{parser.OpConstant, 0, 1},
		parser.OpConstant, 1)
	makeInstruction(t, []byte{parser.OpConstant, 255, 254},
		parser.OpConstant, 65534)
	makeInstruction(t, []byte{parser.OpPop}, parser.OpPop)
	makeInstruction(t, []byte{parser.OpTrue}, parser.OpTrue)
	makeInstruction(t, []byte{parser.OpFalse}, parser.OpFalse)
}

func TestNumObjects(t *testing.T) {
	testCountObjects(t, &z.Array{}, 1)
	testCountObjects(t, &z.Array{Value: []z.Object{
		&z.Int{Value: 1},
		&z.Int{Value: 2},
		&z.Array{Value: []z.Object{
			&z.Int{Value: 3},
			&z.Int{Value: 4},
			&z.Int{Value: 5},
		}},
	}}, 7)
	testCountObjects(t, z.TrueValue, 1)
	testCountObjects(t, z.FalseValue, 1)
	testCountObjects(t, &z.BuiltinFunction{}, 1)
	testCountObjects(t, &z.Bytes{Value: []byte("foobar")}, 1)
	testCountObjects(t, &z.Char{Value: '가'}, 1)
	testCountObjects(t, &z.CompiledFunction{}, 1)
	testCountObjects(t, &z.Error{Value: &z.Int{Value: 5}}, 2)
	testCountObjects(t, &z.Float{Value: 19.84}, 1)
	testCountObjects(t, &z.ImmutableArray{Value: []z.Object{
		&z.Int{Value: 1},
		&z.Int{Value: 2},
		&z.ImmutableArray{Value: []z.Object{
			&z.Int{Value: 3},
			&z.Int{Value: 4},
			&z.Int{Value: 5},
		}},
	}}, 7)
	testCountObjects(t, &z.ImmutableMap{
		Value: map[string]z.Object{
			"k1": &z.Int{Value: 1},
			"k2": &z.Int{Value: 2},
			"k3": &z.Array{Value: []z.Object{
				&z.Int{Value: 3},
				&z.Int{Value: 4},
				&z.Int{Value: 5},
			}},
		}}, 7)
	testCountObjects(t, &z.Int{Value: 1984}, 1)
	testCountObjects(t, &z.Map{Value: map[string]z.Object{
		"k1": &z.Int{Value: 1},
		"k2": &z.Int{Value: 2},
		"k3": &z.Array{Value: []z.Object{
			&z.Int{Value: 3},
			&z.Int{Value: 4},
			&z.Int{Value: 5},
		}},
	}}, 7)
	testCountObjects(t, &z.String{Value: "foo bar"}, 1)
	testCountObjects(t, &z.Time{Value: time.Now()}, 1)
	testCountObjects(t, z.UndefinedValue, 1)
}

func testCountObjects(t *testing.T, o z.Object, expected int) {
	require.Equal(t, expected, z.CountObjects(o))
}

func assertInstructionString(
	t *testing.T,
	instructions [][]byte,
	expected string,
) {
	concatted := make([]byte, 0)
	for _, e := range instructions {
		concatted = append(concatted, e...)
	}
	require.Equal(t, expected, strings.Join(
		z.FormatInstructions(concatted, 0), "\n"))
}

func makeInstruction(
	t *testing.T,
	expected []byte,
	opcode parser.Opcode,
	operands ...int,
) {
	inst := z.MakeInstruction(opcode, operands...)
	require.Equal(t, expected, inst)
}
