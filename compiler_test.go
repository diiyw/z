package z_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/diiyw/z"
	"github.com/diiyw/z/parser"
	"github.com/diiyw/z/require"
	"github.com/diiyw/z/stdlib"
)

func TestCompiler_Compile(t *testing.T) {
	expectCompile(t, `1 + 2`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpBinaryOp, 11),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expectCompile(t, `1; 2`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expectCompile(t, `1 - 2`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpBinaryOp, 12),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expectCompile(t, `1 * 2`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpBinaryOp, 13),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expectCompile(t, `2 / 1`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpBinaryOp, 14),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(2),
				intObject(1))))

	expectCompile(t, `true`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpTrue),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray()))

	expectCompile(t, `false`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpFalse),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray()))

	expectCompile(t, `1 > 2`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpBinaryOp, 39),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expectCompile(t, `1 < 2`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpBinaryOp, 38),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expectCompile(t, `1 >= 2`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpBinaryOp, 44),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expectCompile(t, `1 <= 2`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpBinaryOp, 43),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expectCompile(t, `1 == 2`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpEqual),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expectCompile(t, `1 != 2`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpNotEqual),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expectCompile(t, `true == false`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpTrue),
				z.MakeInstruction(parser.OpFalse),
				z.MakeInstruction(parser.OpEqual),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray()))

	expectCompile(t, `true != false`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpTrue),
				z.MakeInstruction(parser.OpFalse),
				z.MakeInstruction(parser.OpNotEqual),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray()))

	expectCompile(t, `-1`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpMinus),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1))))

	expectCompile(t, `!true`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpTrue),
				z.MakeInstruction(parser.OpLNot),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray()))

	expectCompile(t, `if true { 10 }; 3333`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpTrue),          // 0000
				z.MakeInstruction(parser.OpJumpFalsy, 10), // 0001
				z.MakeInstruction(parser.OpConstant, 0),   // 0004
				z.MakeInstruction(parser.OpPop),           // 0007
				z.MakeInstruction(parser.OpConstant, 1),   // 0008
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)), // 0011
			objectsArray(
				intObject(10),
				intObject(3333))))

	expectCompile(t, `if (true) { 10 } else { 20 }; 3333;`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpTrue),          // 0000
				z.MakeInstruction(parser.OpJumpFalsy, 15), // 0001
				z.MakeInstruction(parser.OpConstant, 0),   // 0004
				z.MakeInstruction(parser.OpPop),           // 0007
				z.MakeInstruction(parser.OpJump, 19),      // 0008
				z.MakeInstruction(parser.OpConstant, 1),   // 0011
				z.MakeInstruction(parser.OpPop),           // 0014
				z.MakeInstruction(parser.OpConstant, 2),   // 0015
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)), // 0018
			objectsArray(
				intObject(10),
				intObject(20),
				intObject(3333))))

	expectCompile(t, `"kami"`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				stringObject("kami"))))

	expectCompile(t, `"ka" + "mi"`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpBinaryOp, 11),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				stringObject("ka"),
				stringObject("mi"))))

	expectCompile(t, `a := 1; b := 2; a += b`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpSetGlobal, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpSetGlobal, 1),
				z.MakeInstruction(parser.OpGetGlobal, 0),
				z.MakeInstruction(parser.OpGetGlobal, 1),
				z.MakeInstruction(parser.OpBinaryOp, 11),
				z.MakeInstruction(parser.OpSetGlobal, 0),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expectCompile(t, `a := 1; b := 2; a /= b`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpSetGlobal, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpSetGlobal, 1),
				z.MakeInstruction(parser.OpGetGlobal, 0),
				z.MakeInstruction(parser.OpGetGlobal, 1),
				z.MakeInstruction(parser.OpBinaryOp, 14),
				z.MakeInstruction(parser.OpSetGlobal, 0),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2))))

	expectCompile(t, `[]`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpArray, 0),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray()))

	expectCompile(t, `[1, 2, 3]`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpArray, 3),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2),
				intObject(3))))

	expectCompile(t, `[1 + 2, 3 - 4, 5 * 6]`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpBinaryOp, 11),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpConstant, 3),
				z.MakeInstruction(parser.OpBinaryOp, 12),
				z.MakeInstruction(parser.OpConstant, 4),
				z.MakeInstruction(parser.OpConstant, 5),
				z.MakeInstruction(parser.OpBinaryOp, 13),
				z.MakeInstruction(parser.OpArray, 3),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2),
				intObject(3),
				intObject(4),
				intObject(5),
				intObject(6))))

	expectCompile(t, `{}`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpMap, 0),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray()))

	expectCompile(t, `{a: 2, b: 4, c: 6}`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpConstant, 3),
				z.MakeInstruction(parser.OpConstant, 4),
				z.MakeInstruction(parser.OpConstant, 5),
				z.MakeInstruction(parser.OpMap, 6),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				stringObject("a"),
				intObject(2),
				stringObject("b"),
				intObject(4),
				stringObject("c"),
				intObject(6))))

	expectCompile(t, `{a: 2 + 3, b: 5 * 6}`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpBinaryOp, 11),
				z.MakeInstruction(parser.OpConstant, 3),
				z.MakeInstruction(parser.OpConstant, 4),
				z.MakeInstruction(parser.OpConstant, 5),
				z.MakeInstruction(parser.OpBinaryOp, 13),
				z.MakeInstruction(parser.OpMap, 4),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				stringObject("a"),
				intObject(2),
				intObject(3),
				stringObject("b"),
				intObject(5),
				intObject(6))))

	expectCompile(t, `[1, 2, 3][1 + 1]`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpArray, 3),
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpBinaryOp, 11),
				z.MakeInstruction(parser.OpIndex),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2),
				intObject(3))))

	expectCompile(t, `{a: 2}[2 - 1]`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpMap, 2),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpBinaryOp, 12),
				z.MakeInstruction(parser.OpIndex),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				stringObject("a"),
				intObject(2),
				intObject(1))))

	expectCompile(t, `[1, 2, 3][:]`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpArray, 3),
				z.MakeInstruction(parser.OpNull),
				z.MakeInstruction(parser.OpNull),
				z.MakeInstruction(parser.OpSliceIndex),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2),
				intObject(3))))

	expectCompile(t, `[1, 2, 3][0 : 2]`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpArray, 3),
				z.MakeInstruction(parser.OpConstant, 3),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpSliceIndex),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2),
				intObject(3),
				intObject(0))))

	expectCompile(t, `[1, 2, 3][:2]`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpArray, 3),
				z.MakeInstruction(parser.OpNull),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpSliceIndex),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2),
				intObject(3))))

	expectCompile(t, `[1, 2, 3][0:]`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpArray, 3),
				z.MakeInstruction(parser.OpConstant, 3),
				z.MakeInstruction(parser.OpNull),
				z.MakeInstruction(parser.OpSliceIndex),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2),
				intObject(3),
				intObject(0))))

	expectCompile(t, `f1 := func(a) { return a }; f1([1, 2]...);`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpSetGlobal, 0),
				z.MakeInstruction(parser.OpGetGlobal, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpArray, 2),
				z.MakeInstruction(parser.OpCall, 1, 1),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				compiledFunction(1, 1,
					z.MakeInstruction(parser.OpGetLocal, 0),
					z.MakeInstruction(parser.OpReturn, 1)),
				intObject(1),
				intObject(2))))

	expectCompile(t, `func() { return 5 + 10 }`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(5),
				intObject(10),
				compiledFunction(0, 0,
					z.MakeInstruction(parser.OpConstant, 0),
					z.MakeInstruction(parser.OpConstant, 1),
					z.MakeInstruction(parser.OpBinaryOp, 11),
					z.MakeInstruction(parser.OpReturn, 1)))))

	expectCompile(t, `func() { 5 + 10 }`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(5),
				intObject(10),
				compiledFunction(0, 0,
					z.MakeInstruction(parser.OpConstant, 0),
					z.MakeInstruction(parser.OpConstant, 1),
					z.MakeInstruction(parser.OpBinaryOp, 11),
					z.MakeInstruction(parser.OpPop),
					z.MakeInstruction(parser.OpReturn, 0)))))

	expectCompile(t, `func() { 1; 2 }`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2),
				compiledFunction(0, 0,
					z.MakeInstruction(parser.OpConstant, 0),
					z.MakeInstruction(parser.OpPop),
					z.MakeInstruction(parser.OpConstant, 1),
					z.MakeInstruction(parser.OpPop),
					z.MakeInstruction(parser.OpReturn, 0)))))

	expectCompile(t, `func() { 1; return 2 }`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2),
				compiledFunction(0, 0,
					z.MakeInstruction(parser.OpConstant, 0),
					z.MakeInstruction(parser.OpPop),
					z.MakeInstruction(parser.OpConstant, 1),
					z.MakeInstruction(parser.OpReturn, 1)))))

	expectCompile(t, `func() { if(true) { return 1 } else { return 2 } }`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2),
				compiledFunction(0, 0,
					z.MakeInstruction(parser.OpTrue),          // 0000
					z.MakeInstruction(parser.OpJumpFalsy, 11), // 0001
					z.MakeInstruction(parser.OpConstant, 0),   // 0004
					z.MakeInstruction(parser.OpReturn, 1),     // 0007
					z.MakeInstruction(parser.OpConstant, 1),   // 0009
					z.MakeInstruction(parser.OpReturn, 1)))))  // 0012

	expectCompile(t, `func() { 1; if(true) { 2 } else { 3 }; 4 }`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 4),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(1),
				intObject(2),
				intObject(3),
				intObject(4),
				compiledFunction(0, 0,
					z.MakeInstruction(parser.OpConstant, 0),   // 0000
					z.MakeInstruction(parser.OpPop),           // 0003
					z.MakeInstruction(parser.OpTrue),          // 0004
					z.MakeInstruction(parser.OpJumpFalsy, 19), // 0005
					z.MakeInstruction(parser.OpConstant, 1),   // 0008
					z.MakeInstruction(parser.OpPop),           // 0011
					z.MakeInstruction(parser.OpJump, 23),      // 0012
					z.MakeInstruction(parser.OpConstant, 2),   // 0015
					z.MakeInstruction(parser.OpPop),           // 0018
					z.MakeInstruction(parser.OpConstant, 3),   // 0019
					z.MakeInstruction(parser.OpPop),           // 0022
					z.MakeInstruction(parser.OpReturn, 0)))))  // 0023

	expectCompile(t, `func() { }`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				compiledFunction(0, 0,
					z.MakeInstruction(parser.OpReturn, 0)))))

	expectCompile(t, `func() { 24 }()`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpCall, 0, 0),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(24),
				compiledFunction(0, 0,
					z.MakeInstruction(parser.OpConstant, 0),
					z.MakeInstruction(parser.OpPop),
					z.MakeInstruction(parser.OpReturn, 0)))))

	expectCompile(t, `func() { return 24 }()`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpCall, 0, 0),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(24),
				compiledFunction(0, 0,
					z.MakeInstruction(parser.OpConstant, 0),
					z.MakeInstruction(parser.OpReturn, 1)))))

	expectCompile(t, `noArg := func() { 24 }; noArg();`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpSetGlobal, 0),
				z.MakeInstruction(parser.OpGetGlobal, 0),
				z.MakeInstruction(parser.OpCall, 0, 0),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(24),
				compiledFunction(0, 0,
					z.MakeInstruction(parser.OpConstant, 0),
					z.MakeInstruction(parser.OpPop),
					z.MakeInstruction(parser.OpReturn, 0)))))

	expectCompile(t, `noArg := func() { return 24 }; noArg();`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpSetGlobal, 0),
				z.MakeInstruction(parser.OpGetGlobal, 0),
				z.MakeInstruction(parser.OpCall, 0, 0),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(24),
				compiledFunction(0, 0,
					z.MakeInstruction(parser.OpConstant, 0),
					z.MakeInstruction(parser.OpReturn, 1)))))

	expectCompile(t, `n := 55; func() { n };`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpSetGlobal, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(55),
				compiledFunction(0, 0,
					z.MakeInstruction(parser.OpGetGlobal, 0),
					z.MakeInstruction(parser.OpPop),
					z.MakeInstruction(parser.OpReturn, 0)))))

	expectCompile(t, `func() { n := 55; return n }`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(55),
				compiledFunction(1, 0,
					z.MakeInstruction(parser.OpConstant, 0),
					z.MakeInstruction(parser.OpDefineLocal, 0),
					z.MakeInstruction(parser.OpGetLocal, 0),
					z.MakeInstruction(parser.OpReturn, 1)))))

	expectCompile(t, `func() { a := 55; b := 77; return a + b }`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(55),
				intObject(77),
				compiledFunction(2, 0,
					z.MakeInstruction(parser.OpConstant, 0),
					z.MakeInstruction(parser.OpDefineLocal, 0),
					z.MakeInstruction(parser.OpConstant, 1),
					z.MakeInstruction(parser.OpDefineLocal, 1),
					z.MakeInstruction(parser.OpGetLocal, 0),
					z.MakeInstruction(parser.OpGetLocal, 1),
					z.MakeInstruction(parser.OpBinaryOp, 11),
					z.MakeInstruction(parser.OpReturn, 1)))))

	expectCompile(t, `f1 := func(a) { return a }; f1(24);`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpSetGlobal, 0),
				z.MakeInstruction(parser.OpGetGlobal, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpCall, 1, 0),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				compiledFunction(1, 1,
					z.MakeInstruction(parser.OpGetLocal, 0),
					z.MakeInstruction(parser.OpReturn, 1)),
				intObject(24))))

	expectCompile(t, `varTest := func(...a) { return a }; varTest(1,2,3);`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpSetGlobal, 0),
				z.MakeInstruction(parser.OpGetGlobal, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpConstant, 3),
				z.MakeInstruction(parser.OpCall, 3, 0),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				compiledFunction(1, 1,
					z.MakeInstruction(parser.OpGetLocal, 0),
					z.MakeInstruction(parser.OpReturn, 1)),
				intObject(1), intObject(2), intObject(3))))

	expectCompile(t, `f1 := func(a, b, c) { a; b; return c; }; f1(24, 25, 26);`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpSetGlobal, 0),
				z.MakeInstruction(parser.OpGetGlobal, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpConstant, 3),
				z.MakeInstruction(parser.OpCall, 3, 0),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				compiledFunction(3, 3,
					z.MakeInstruction(parser.OpGetLocal, 0),
					z.MakeInstruction(parser.OpPop),
					z.MakeInstruction(parser.OpGetLocal, 1),
					z.MakeInstruction(parser.OpPop),
					z.MakeInstruction(parser.OpGetLocal, 2),
					z.MakeInstruction(parser.OpReturn, 1)),
				intObject(24),
				intObject(25),
				intObject(26))))

	expectCompile(t, `func() { n := 55; n = 23; return n }`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(55),
				intObject(23),
				compiledFunction(1, 0,
					z.MakeInstruction(parser.OpConstant, 0),
					z.MakeInstruction(parser.OpDefineLocal, 0),
					z.MakeInstruction(parser.OpConstant, 1),
					z.MakeInstruction(parser.OpSetLocal, 0),
					z.MakeInstruction(parser.OpGetLocal, 0),
					z.MakeInstruction(parser.OpReturn, 1)))))
	expectCompile(t, `len([]);`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpGetBuiltin, 0),
				z.MakeInstruction(parser.OpArray, 0),
				z.MakeInstruction(parser.OpCall, 1, 0),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray()))

	expectCompile(t, `func() { return len([]) }`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				compiledFunction(0, 0,
					z.MakeInstruction(parser.OpGetBuiltin, 0),
					z.MakeInstruction(parser.OpArray, 0),
					z.MakeInstruction(parser.OpCall, 1, 0),
					z.MakeInstruction(parser.OpReturn, 1)))))

	expectCompile(t, `func(a) { func(b) { return a + b } }`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				compiledFunction(1, 1,
					z.MakeInstruction(parser.OpGetFree, 0),
					z.MakeInstruction(parser.OpGetLocal, 0),
					z.MakeInstruction(parser.OpBinaryOp, 11),
					z.MakeInstruction(parser.OpReturn, 1)),
				compiledFunction(1, 1,
					z.MakeInstruction(parser.OpGetLocalPtr, 0),
					z.MakeInstruction(parser.OpClosure, 0, 1),
					z.MakeInstruction(parser.OpPop),
					z.MakeInstruction(parser.OpReturn, 0)))))

	expectCompile(t, `
func(a) {
	return func(b) {
		return func(c) {
			return a + b + c
		}
	}
}`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				compiledFunction(1, 1,
					z.MakeInstruction(parser.OpGetFree, 0),
					z.MakeInstruction(parser.OpGetFree, 1),
					z.MakeInstruction(parser.OpBinaryOp, 11),
					z.MakeInstruction(parser.OpGetLocal, 0),
					z.MakeInstruction(parser.OpBinaryOp, 11),
					z.MakeInstruction(parser.OpReturn, 1)),
				compiledFunction(1, 1,
					z.MakeInstruction(parser.OpGetFreePtr, 0),
					z.MakeInstruction(parser.OpGetLocalPtr, 0),
					z.MakeInstruction(parser.OpClosure, 0, 2),
					z.MakeInstruction(parser.OpReturn, 1)),
				compiledFunction(1, 1,
					z.MakeInstruction(parser.OpGetLocalPtr, 0),
					z.MakeInstruction(parser.OpClosure, 1, 1),
					z.MakeInstruction(parser.OpReturn, 1)))))

	expectCompile(t, `
g := 55;

func() {
	a := 66;

	return func() {
		b := 77;

		return func() {
			c := 88;

			return g + a + b + c;
		}
	}
}`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpSetGlobal, 0),
				z.MakeInstruction(parser.OpConstant, 6),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(55),
				intObject(66),
				intObject(77),
				intObject(88),
				compiledFunction(1, 0,
					z.MakeInstruction(parser.OpConstant, 3),
					z.MakeInstruction(parser.OpDefineLocal, 0),
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
					z.MakeInstruction(parser.OpDefineLocal, 0),
					z.MakeInstruction(parser.OpGetFreePtr, 0),
					z.MakeInstruction(parser.OpGetLocalPtr, 0),
					z.MakeInstruction(parser.OpClosure, 4, 2),
					z.MakeInstruction(parser.OpReturn, 1)),
				compiledFunction(1, 0,
					z.MakeInstruction(parser.OpConstant, 1),
					z.MakeInstruction(parser.OpDefineLocal, 0),
					z.MakeInstruction(parser.OpGetLocalPtr, 0),
					z.MakeInstruction(parser.OpClosure, 5, 1),
					z.MakeInstruction(parser.OpReturn, 1)))))

	expectCompile(t, `for i:=0; i<10; i++ {}`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpSetGlobal, 0),
				z.MakeInstruction(parser.OpGetGlobal, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpBinaryOp, 38),
				z.MakeInstruction(parser.OpJumpFalsy, 35),
				z.MakeInstruction(parser.OpGetGlobal, 0),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpBinaryOp, 11),
				z.MakeInstruction(parser.OpSetGlobal, 0),
				z.MakeInstruction(parser.OpJump, 6),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(0),
				intObject(10),
				intObject(1))))

	expectCompile(t, `m := {}; for k, v in m {}`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpMap, 0),
				z.MakeInstruction(parser.OpSetGlobal, 0),
				z.MakeInstruction(parser.OpGetGlobal, 0),
				z.MakeInstruction(parser.OpIteratorInit),
				z.MakeInstruction(parser.OpSetGlobal, 1),
				z.MakeInstruction(parser.OpGetGlobal, 1),
				z.MakeInstruction(parser.OpIteratorNext),
				z.MakeInstruction(parser.OpJumpFalsy, 41),
				z.MakeInstruction(parser.OpGetGlobal, 1),
				z.MakeInstruction(parser.OpIteratorKey),
				z.MakeInstruction(parser.OpSetGlobal, 2),
				z.MakeInstruction(parser.OpGetGlobal, 1),
				z.MakeInstruction(parser.OpIteratorValue),
				z.MakeInstruction(parser.OpSetGlobal, 3),
				z.MakeInstruction(parser.OpJump, 13),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray()))

	expectCompile(t, `a := 0; a == 0 && a != 1 || a < 1`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpSetGlobal, 0),
				z.MakeInstruction(parser.OpGetGlobal, 0),
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpEqual),
				z.MakeInstruction(parser.OpAndJump, 25),
				z.MakeInstruction(parser.OpGetGlobal, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpNotEqual),
				z.MakeInstruction(parser.OpOrJump, 38),
				z.MakeInstruction(parser.OpGetGlobal, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpBinaryOp, 38),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(0),
				intObject(1))))

	// unknown module name
	expectCompileError(t, `import("user1")`, "module 'user1' not found")

	// too many errors
	expectCompileError(t, `
r["x"] = {
    @a:1,
    @b:1,
    @c:1,
    @d:1,
    @e:1,
    @f:1,
    @g:1,
    @h:1,
    @i:1,
    @j:1,
    @k:1
}
`, "Parse Error: illegal character U+0040 '@'\n\tat test:3:5 (and 10 more errors)")

	expectCompileError(t, `import("")`, "empty module name")

	// https://github.com/d5/z/issues/314
	expectCompileError(t, `
(func() {
	fn := fn()
})()
`, "unresolved reference 'fn")
}

func TestCompilerErrorReport(t *testing.T) {
	expectCompileError(t, `import("user1")`,
		"Compile Error: module 'user1' not found\n\tat test:1:1")

	expectCompileError(t, `a = 1`,
		"Compile Error: unresolved reference 'a'\n\tat test:1:1")
	expectCompileError(t, `a := a`,
		"Compile Error: unresolved reference 'a'\n\tat test:1:6")
	expectCompileError(t, `a, b := 1, 2`,
		"Compile Error: tuple assignment not allowed\n\tat test:1:1")
	expectCompileError(t, `a.b := 1`,
		"not allowed with selector")
	expectCompileError(t, `a:=1; a:=3`,
		"Compile Error: 'a' redeclared in this block\n\tat test:1:7")

	expectCompileError(t, `return 5`,
		"Compile Error: return not allowed outside function\n\tat test:1:1")
	expectCompileError(t, `func() { break }`,
		"Compile Error: break not allowed outside loop\n\tat test:1:10")
	expectCompileError(t, `func() { continue }`,
		"Compile Error: continue not allowed outside loop\n\tat test:1:10")
	expectCompileError(t, `func() { export 5 }`,
		"Compile Error: export not allowed inside function\n\tat test:1:10")
}

func TestCompilerDeadCode(t *testing.T) {
	expectCompile(t, `
func() {
	a := 4
	return a

	b := 5 // dead code from here
	c := a
	return b
}`,
		bytecode(
			concatInsts(
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpSuspend)),
			objectsArray(
				intObject(4),
				intObject(5),
				compiledFunction(0, 0,
					z.MakeInstruction(parser.OpConstant, 0),
					z.MakeInstruction(parser.OpDefineLocal, 0),
					z.MakeInstruction(parser.OpGetLocal, 0),
					z.MakeInstruction(parser.OpReturn, 1)))))

	expectCompile(t, `
func() {
	if true {
		return 5
		a := 4  // dead code from here
		b := a
		return b
	} else {
		return 4
		c := 5  // dead code from here
		d := c
		return d
	}
}`, bytecode(
		concatInsts(
			z.MakeInstruction(parser.OpConstant, 2),
			z.MakeInstruction(parser.OpPop),
			z.MakeInstruction(parser.OpSuspend)),
		objectsArray(
			intObject(5),
			intObject(4),
			compiledFunction(0, 0,
				z.MakeInstruction(parser.OpTrue),
				z.MakeInstruction(parser.OpJumpFalsy, 11),
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpReturn, 1),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpReturn, 1)))))

	expectCompile(t, `
func() {
	a := 1
	for {
		if a == 5 {
			return 10
		}
		5 + 5
		return 20
		b := a
		return b
	}
}`, bytecode(
		concatInsts(
			z.MakeInstruction(parser.OpConstant, 4),
			z.MakeInstruction(parser.OpPop),
			z.MakeInstruction(parser.OpSuspend)),
		objectsArray(
			intObject(1),
			intObject(5),
			intObject(10),
			intObject(20),
			compiledFunction(0, 0,
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpDefineLocal, 0),
				z.MakeInstruction(parser.OpGetLocal, 0),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpEqual),
				z.MakeInstruction(parser.OpJumpFalsy, 21),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpReturn, 1),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpBinaryOp, 11),
				z.MakeInstruction(parser.OpPop),
				z.MakeInstruction(parser.OpConstant, 3),
				z.MakeInstruction(parser.OpReturn, 1)))))

	expectCompile(t, `
func() {
	if true {
		return 5
		a := 4  // dead code from here
		b := a
		return b
	} else {
		return 4
		c := 5  // dead code from here
		d := c
		return d
	}
}`, bytecode(
		concatInsts(
			z.MakeInstruction(parser.OpConstant, 2),
			z.MakeInstruction(parser.OpPop),
			z.MakeInstruction(parser.OpSuspend)),
		objectsArray(
			intObject(5),
			intObject(4),
			compiledFunction(0, 0,
				z.MakeInstruction(parser.OpTrue),
				z.MakeInstruction(parser.OpJumpFalsy, 11),
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpReturn, 1),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpReturn, 1)))))

	expectCompile(t, `
func() {
	if true {
		return
	}

    return

    return 123
}`, bytecode(
		concatInsts(
			z.MakeInstruction(parser.OpConstant, 1),
			z.MakeInstruction(parser.OpPop),
			z.MakeInstruction(parser.OpSuspend)),
		objectsArray(
			intObject(123),
			compiledFunction(0, 0,
				z.MakeInstruction(parser.OpTrue),
				z.MakeInstruction(parser.OpJumpFalsy, 8),
				z.MakeInstruction(parser.OpReturn, 0),
				z.MakeInstruction(parser.OpReturn, 0),
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpReturn, 1)))))
}

func TestCompilerScopes(t *testing.T) {
	expectCompile(t, `
if a := 1; a {
    a = 2
	b := a
} else {
    a = 3
	b := a
}`, bytecode(
		concatInsts(
			z.MakeInstruction(parser.OpConstant, 0),
			z.MakeInstruction(parser.OpSetGlobal, 0),
			z.MakeInstruction(parser.OpGetGlobal, 0),
			z.MakeInstruction(parser.OpJumpFalsy, 31),
			z.MakeInstruction(parser.OpConstant, 1),
			z.MakeInstruction(parser.OpSetGlobal, 0),
			z.MakeInstruction(parser.OpGetGlobal, 0),
			z.MakeInstruction(parser.OpSetGlobal, 1),
			z.MakeInstruction(parser.OpJump, 43),
			z.MakeInstruction(parser.OpConstant, 2),
			z.MakeInstruction(parser.OpSetGlobal, 0),
			z.MakeInstruction(parser.OpGetGlobal, 0),
			z.MakeInstruction(parser.OpSetGlobal, 2),
			z.MakeInstruction(parser.OpSuspend)),
		objectsArray(
			intObject(1),
			intObject(2),
			intObject(3))))

	expectCompile(t, `
func() {
	if a := 1; a {
    	a = 2
		b := a
	} else {
    	a = 3
		b := a
	}
}`, bytecode(
		concatInsts(
			z.MakeInstruction(parser.OpConstant, 3),
			z.MakeInstruction(parser.OpPop),
			z.MakeInstruction(parser.OpSuspend)),
		objectsArray(
			intObject(1),
			intObject(2),
			intObject(3),
			compiledFunction(0, 0,
				z.MakeInstruction(parser.OpConstant, 0),
				z.MakeInstruction(parser.OpDefineLocal, 0),
				z.MakeInstruction(parser.OpGetLocal, 0),
				z.MakeInstruction(parser.OpJumpFalsy, 26),
				z.MakeInstruction(parser.OpConstant, 1),
				z.MakeInstruction(parser.OpSetLocal, 0),
				z.MakeInstruction(parser.OpGetLocal, 0),
				z.MakeInstruction(parser.OpDefineLocal, 1),
				z.MakeInstruction(parser.OpJump, 35),
				z.MakeInstruction(parser.OpConstant, 2),
				z.MakeInstruction(parser.OpSetLocal, 0),
				z.MakeInstruction(parser.OpGetLocal, 0),
				z.MakeInstruction(parser.OpDefineLocal, 1),
				z.MakeInstruction(parser.OpReturn, 0)))))
}

func TestCompiler_custom_extension(t *testing.T) {
	pathFileSource := "./testdata/issue286/test.mshk"

	modules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)

	src, err := ioutil.ReadFile(pathFileSource)
	require.NoError(t, err)

	// Escape shegang
	if len(src) > 1 && string(src[:2]) == "#!" {
		copy(src, "//")
	}

	fileSet := parser.NewFileSet()
	srcFile := fileSet.AddFile(filepath.Base(pathFileSource), -1, len(src))

	p := parser.NewParser(srcFile, src, nil)
	file, err := p.ParseFile()
	require.NoError(t, err)

	c := z.NewCompiler(srcFile, nil, nil, modules, nil)
	c.EnableFileImport(true)
	c.SetImportDir(filepath.Dir(pathFileSource))

	// Search for "*.z" and ".mshk"(custom extension)
	c.SetImportFileExt(".z", ".mshk")

	err = c.Compile(file)
	require.NoError(t, err)
}

func TestCompilerNewCompiler_default_file_extension(t *testing.T) {
	modules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	input := "{}"
	fileSet := parser.NewFileSet()
	file := fileSet.AddFile("test", -1, len(input))

	c := z.NewCompiler(file, nil, nil, modules, nil)
	c.EnableFileImport(true)

	require.Equal(t, []string{".z"}, c.GetImportFileExt(),
		"newly created compiler object must contain the default extension")
}

func TestCompilerSetImportExt_extension_name_validation(t *testing.T) {
	c := new(z.Compiler) // Instantiate a new compiler object with no initialization

	// Test of empty arg
	err := c.SetImportFileExt()

	require.Error(t, err, "empty arg should return an error")

	// Test of various arg types
	for _, test := range []struct {
		extensions []string
		expect     []string
		requireErr bool
		msgFail    string
	}{
		{[]string{".z"}, []string{".z"}, false,
			"well-formed extension should not return an error"},
		{[]string{""}, []string{".z"}, true,
			"empty extension name should return an error"},
		{[]string{"foo"}, []string{".z"}, true,
			"name without dot prefix should return an error"},
		{[]string{"foo.bar"}, []string{".z"}, true,
			"malformed extension should return an error"},
		{[]string{"foo."}, []string{".z"}, true,
			"malformed extension should return an error"},
		{[]string{".mshk"}, []string{".mshk"}, false,
			"name with dot prefix should be added"},
		{[]string{".foo", ".bar"}, []string{".foo", ".bar"}, false,
			"it should replace instead of appending"},
	} {
		err := c.SetImportFileExt(test.extensions...)
		if test.requireErr {
			require.Error(t, err, test.msgFail)
		}

		expect := test.expect
		actual := c.GetImportFileExt()
		require.Equal(t, expect, actual, test.msgFail)
	}
}

func concatInsts(instructions ...[]byte) []byte {
	var concat []byte
	for _, i := range instructions {
		concat = append(concat, i...)
	}
	return concat
}

func bytecode(
	instructions []byte,
	constants []z.Object,
) *z.Bytecode {
	return &z.Bytecode{
		FileSet:      parser.NewFileSet(),
		MainFunction: &z.CompiledFunction{Instructions: instructions},
		Constants:    constants,
	}
}

func expectCompile(
	t *testing.T,
	input string,
	expected *z.Bytecode,
) {
	actual, trace, err := traceCompile(input, nil)

	var ok bool
	defer func() {
		if !ok {
			for _, tr := range trace {
				t.Log(tr)
			}
		}
	}()

	require.NoError(t, err)
	equalBytecode(t, expected, actual)
	ok = true
}

func expectCompileError(t *testing.T, input, expected string) {
	_, trace, err := traceCompile(input, nil)

	var ok bool
	defer func() {
		if !ok {
			for _, tr := range trace {
				t.Log(tr)
			}
		}
	}()

	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), expected),
		"expected error string: %s, got: %s", expected, err.Error())
	ok = true
}

func equalBytecode(t *testing.T, expected, actual *z.Bytecode) {
	require.Equal(t, expected.MainFunction, actual.MainFunction)
	equalConstants(t, expected.Constants, actual.Constants)
}

func equalConstants(t *testing.T, expected, actual []z.Object) {
	require.Equal(t, len(expected), len(actual))
	for i := 0; i < len(expected); i++ {
		require.Equal(t, expected[i], actual[i])
	}
}

type compileTracer struct {
	Out []string
}

func (o *compileTracer) Write(p []byte) (n int, err error) {
	o.Out = append(o.Out, string(p))
	return len(p), nil
}

func traceCompile(
	input string,
	symbols map[string]z.Object,
) (res *z.Bytecode, trace []string, err error) {
	fileSet := parser.NewFileSet()
	file := fileSet.AddFile("test", -1, len(input))

	p := parser.NewParser(file, []byte(input), nil)

	symTable := z.NewSymbolTable()
	for name := range symbols {
		symTable.Define(name)
	}
	for idx, fn := range z.GetAllBuiltinFunctions() {
		symTable.DefineBuiltin(idx, fn.Name)
	}

	tr := &compileTracer{}
	c := z.NewCompiler(file, symTable, nil, nil, tr)
	parsed, err := p.ParseFile()
	if err != nil {
		return
	}

	err = c.Compile(parsed)
	res = c.Bytecode()
	res.RemoveDuplicates()
	{
		trace = append(trace, fmt.Sprintf("Compiler Trace:\n%s",
			strings.Join(tr.Out, "")))
		trace = append(trace, fmt.Sprintf("Compiled Constants:\n%s",
			strings.Join(res.FormatConstants(), "\n")))
		trace = append(trace, fmt.Sprintf("Compiled Instructions:\n%s\n",
			strings.Join(res.FormatInstructions(), "\n")))
	}
	if err != nil {
		return
	}
	return
}

func objectsArray(o ...z.Object) []z.Object {
	return o
}

func intObject(v int64) *z.Int {
	return &z.Int{Value: v}
}

func stringObject(v string) *z.String {
	return &z.String{Value: v}
}

func compiledFunction(
	numLocals, numParams int,
	insts ...[]byte,
) *z.CompiledFunction {
	return &z.CompiledFunction{
		Instructions:  concatInsts(insts...),
		NumLocals:     numLocals,
		NumParameters: numParams,
	}
}
