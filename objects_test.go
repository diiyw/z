package z_test

import (
	"testing"

	"github.com/diiyw/z"
	"github.com/diiyw/z/require"
	"github.com/diiyw/z/token"
)

func TestObject_TypeName(t *testing.T) {
	var o z.Object = &z.Int{}
	require.Equal(t, "int", o.TypeName())
	o = &z.Float{}
	require.Equal(t, "float", o.TypeName())
	o = &z.Char{}
	require.Equal(t, "char", o.TypeName())
	o = &z.String{}
	require.Equal(t, "string", o.TypeName())
	o = &z.Bool{}
	require.Equal(t, "bool", o.TypeName())
	o = &z.Array{}
	require.Equal(t, "array", o.TypeName())
	o = &z.Map{}
	require.Equal(t, "map", o.TypeName())
	o = &z.ArrayIterator{}
	require.Equal(t, "array-iterator", o.TypeName())
	o = &z.StringIterator{}
	require.Equal(t, "string-iterator", o.TypeName())
	o = &z.MapIterator{}
	require.Equal(t, "map-iterator", o.TypeName())
	o = &z.BuiltinFunction{Name: "fn"}
	require.Equal(t, "builtin-function:fn", o.TypeName())
	o = &z.UserFunction{Name: "fn"}
	require.Equal(t, "user-function:fn", o.TypeName())
	o = &z.CompiledFunction{}
	require.Equal(t, "compiled-function", o.TypeName())
	o = &z.Undefined{}
	require.Equal(t, "undefined", o.TypeName())
	o = &z.Error{}
	require.Equal(t, "error", o.TypeName())
	o = &z.Bytes{}
	require.Equal(t, "bytes", o.TypeName())
}

func TestObject_IsFalsy(t *testing.T) {
	var o z.Object = &z.Int{Value: 0}
	require.True(t, o.IsFalsy())
	o = &z.Int{Value: 1}
	require.False(t, o.IsFalsy())
	o = &z.Float{Value: 0}
	require.False(t, o.IsFalsy())
	o = &z.Float{Value: 1}
	require.False(t, o.IsFalsy())
	o = &z.Char{Value: ' '}
	require.False(t, o.IsFalsy())
	o = &z.Char{Value: 'T'}
	require.False(t, o.IsFalsy())
	o = &z.String{Value: ""}
	require.True(t, o.IsFalsy())
	o = &z.String{Value: " "}
	require.False(t, o.IsFalsy())
	o = &z.Array{Value: nil}
	require.True(t, o.IsFalsy())
	o = &z.Array{Value: []z.Object{nil}} // nil is not valid but still count as 1 element
	require.False(t, o.IsFalsy())
	o = &z.Map{Value: nil}
	require.True(t, o.IsFalsy())
	o = &z.Map{Value: map[string]z.Object{"a": nil}} // nil is not valid but still count as 1 element
	require.False(t, o.IsFalsy())
	o = &z.StringIterator{}
	require.True(t, o.IsFalsy())
	o = &z.ArrayIterator{}
	require.True(t, o.IsFalsy())
	o = &z.MapIterator{}
	require.True(t, o.IsFalsy())
	o = &z.BuiltinFunction{}
	require.False(t, o.IsFalsy())
	o = &z.CompiledFunction{}
	require.False(t, o.IsFalsy())
	o = &z.Undefined{}
	require.True(t, o.IsFalsy())
	o = &z.Error{}
	require.True(t, o.IsFalsy())
	o = &z.Bytes{}
	require.True(t, o.IsFalsy())
	o = &z.Bytes{Value: []byte{1, 2}}
	require.False(t, o.IsFalsy())
}

func TestObject_String(t *testing.T) {
	var o z.Object = &z.Int{Value: 0}
	require.Equal(t, "0", o.String())
	o = &z.Int{Value: 1}
	require.Equal(t, "1", o.String())
	o = &z.Float{Value: 0}
	require.Equal(t, "0", o.String())
	o = &z.Float{Value: 1}
	require.Equal(t, "1", o.String())
	o = &z.Char{Value: ' '}
	require.Equal(t, " ", o.String())
	o = &z.Char{Value: 'T'}
	require.Equal(t, "T", o.String())
	o = &z.String{Value: ""}
	require.Equal(t, `""`, o.String())
	o = &z.String{Value: " "}
	require.Equal(t, `" "`, o.String())
	o = &z.Array{Value: nil}
	require.Equal(t, "[]", o.String())
	o = &z.Map{Value: nil}
	require.Equal(t, "{}", o.String())
	o = &z.Error{Value: nil}
	require.Equal(t, "error", o.String())
	o = &z.Error{Value: &z.String{Value: "error 1"}}
	require.Equal(t, `error: "error 1"`, o.String())
	o = &z.StringIterator{}
	require.Equal(t, "<string-iterator>", o.String())
	o = &z.ArrayIterator{}
	require.Equal(t, "<array-iterator>", o.String())
	o = &z.MapIterator{}
	require.Equal(t, "<map-iterator>", o.String())
	o = &z.Undefined{}
	require.Equal(t, "<undefined>", o.String())
	o = &z.Bytes{}
	require.Equal(t, "", o.String())
	o = &z.Bytes{Value: []byte("foo")}
	require.Equal(t, "foo", o.String())
}

func TestObject_BinaryOp(t *testing.T) {
	var o z.Object = &z.Char{}
	_, err := o.BinaryOp(token.Add, z.UndefinedValue)
	require.Error(t, err)
	o = &z.Bool{}
	_, err = o.BinaryOp(token.Add, z.UndefinedValue)
	require.Error(t, err)
	o = &z.Map{}
	_, err = o.BinaryOp(token.Add, z.UndefinedValue)
	require.Error(t, err)
	o = &z.ArrayIterator{}
	_, err = o.BinaryOp(token.Add, z.UndefinedValue)
	require.Error(t, err)
	o = &z.StringIterator{}
	_, err = o.BinaryOp(token.Add, z.UndefinedValue)
	require.Error(t, err)
	o = &z.MapIterator{}
	_, err = o.BinaryOp(token.Add, z.UndefinedValue)
	require.Error(t, err)
	o = &z.BuiltinFunction{}
	_, err = o.BinaryOp(token.Add, z.UndefinedValue)
	require.Error(t, err)
	o = &z.CompiledFunction{}
	_, err = o.BinaryOp(token.Add, z.UndefinedValue)
	require.Error(t, err)
	o = &z.Undefined{}
	_, err = o.BinaryOp(token.Add, z.UndefinedValue)
	require.Error(t, err)
	o = &z.Error{}
	_, err = o.BinaryOp(token.Add, z.UndefinedValue)
	require.Error(t, err)
}

func TestArray_BinaryOp(t *testing.T) {
	testBinaryOp(t, &z.Array{Value: nil}, token.Add,
		&z.Array{Value: nil}, &z.Array{Value: nil})
	testBinaryOp(t, &z.Array{Value: nil}, token.Add,
		&z.Array{Value: []z.Object{}}, &z.Array{Value: nil})
	testBinaryOp(t, &z.Array{Value: []z.Object{}}, token.Add,
		&z.Array{Value: nil}, &z.Array{Value: []z.Object{}})
	testBinaryOp(t, &z.Array{Value: []z.Object{}}, token.Add,
		&z.Array{Value: []z.Object{}},
		&z.Array{Value: []z.Object{}})
	testBinaryOp(t, &z.Array{Value: nil}, token.Add,
		&z.Array{Value: []z.Object{
			&z.Int{Value: 1},
		}}, &z.Array{Value: []z.Object{
			&z.Int{Value: 1},
		}})
	testBinaryOp(t, &z.Array{Value: nil}, token.Add,
		&z.Array{Value: []z.Object{
			&z.Int{Value: 1},
			&z.Int{Value: 2},
			&z.Int{Value: 3},
		}}, &z.Array{Value: []z.Object{
			&z.Int{Value: 1},
			&z.Int{Value: 2},
			&z.Int{Value: 3},
		}})
	testBinaryOp(t, &z.Array{Value: []z.Object{
		&z.Int{Value: 1},
		&z.Int{Value: 2},
		&z.Int{Value: 3},
	}}, token.Add, &z.Array{Value: nil},
		&z.Array{Value: []z.Object{
			&z.Int{Value: 1},
			&z.Int{Value: 2},
			&z.Int{Value: 3},
		}})
	testBinaryOp(t, &z.Array{Value: []z.Object{
		&z.Int{Value: 1},
		&z.Int{Value: 2},
		&z.Int{Value: 3},
	}}, token.Add, &z.Array{Value: []z.Object{
		&z.Int{Value: 4},
		&z.Int{Value: 5},
		&z.Int{Value: 6},
	}}, &z.Array{Value: []z.Object{
		&z.Int{Value: 1},
		&z.Int{Value: 2},
		&z.Int{Value: 3},
		&z.Int{Value: 4},
		&z.Int{Value: 5},
		&z.Int{Value: 6},
	}})
}

func TestError_Equals(t *testing.T) {
	err1 := &z.Error{Value: &z.String{Value: "some error"}}
	err2 := err1
	require.True(t, err1.Equals(err2))
	require.True(t, err2.Equals(err1))

	err2 = &z.Error{Value: &z.String{Value: "some error"}}
	require.False(t, err1.Equals(err2))
	require.False(t, err2.Equals(err1))
}

func TestFloat_BinaryOp(t *testing.T) {
	// float + float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &z.Float{Value: l}, token.Add,
				&z.Float{Value: r}, &z.Float{Value: l + r})
		}
	}

	// float - float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &z.Float{Value: l}, token.Sub,
				&z.Float{Value: r}, &z.Float{Value: l - r})
		}
	}

	// float * float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &z.Float{Value: l}, token.Mul,
				&z.Float{Value: r}, &z.Float{Value: l * r})
		}
	}

	// float / float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			if r != 0 {
				testBinaryOp(t, &z.Float{Value: l}, token.Quo,
					&z.Float{Value: r}, &z.Float{Value: l / r})
			}
		}
	}

	// float < float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &z.Float{Value: l}, token.Less,
				&z.Float{Value: r}, boolValue(l < r))
		}
	}

	// float > float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &z.Float{Value: l}, token.Greater,
				&z.Float{Value: r}, boolValue(l > r))
		}
	}

	// float <= float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &z.Float{Value: l}, token.LessEq,
				&z.Float{Value: r}, boolValue(l <= r))
		}
	}

	// float >= float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &z.Float{Value: l}, token.GreaterEq,
				&z.Float{Value: r}, boolValue(l >= r))
		}
	}

	// float + int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &z.Float{Value: l}, token.Add,
				&z.Int{Value: r}, &z.Float{Value: l + float64(r)})
		}
	}

	// float - int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &z.Float{Value: l}, token.Sub,
				&z.Int{Value: r}, &z.Float{Value: l - float64(r)})
		}
	}

	// float * int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &z.Float{Value: l}, token.Mul,
				&z.Int{Value: r}, &z.Float{Value: l * float64(r)})
		}
	}

	// float / int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			if r != 0 {
				testBinaryOp(t, &z.Float{Value: l}, token.Quo,
					&z.Int{Value: r},
					&z.Float{Value: l / float64(r)})
			}
		}
	}

	// float < int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &z.Float{Value: l}, token.Less,
				&z.Int{Value: r}, boolValue(l < float64(r)))
		}
	}

	// float > int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &z.Float{Value: l}, token.Greater,
				&z.Int{Value: r}, boolValue(l > float64(r)))
		}
	}

	// float <= int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &z.Float{Value: l}, token.LessEq,
				&z.Int{Value: r}, boolValue(l <= float64(r)))
		}
	}

	// float >= int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &z.Float{Value: l}, token.GreaterEq,
				&z.Int{Value: r}, boolValue(l >= float64(r)))
		}
	}
}

func TestInt_BinaryOp(t *testing.T) {
	// int + int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &z.Int{Value: l}, token.Add,
				&z.Int{Value: r}, &z.Int{Value: l + r})
		}
	}

	// int - int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &z.Int{Value: l}, token.Sub,
				&z.Int{Value: r}, &z.Int{Value: l - r})
		}
	}

	// int * int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &z.Int{Value: l}, token.Mul,
				&z.Int{Value: r}, &z.Int{Value: l * r})
		}
	}

	// int / int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			if r != 0 {
				testBinaryOp(t, &z.Int{Value: l}, token.Quo,
					&z.Int{Value: r}, &z.Int{Value: l / r})
			}
		}
	}

	// int % int
	for l := int64(-4); l <= 4; l++ {
		for r := -int64(-4); r <= 4; r++ {
			if r == 0 {
				testBinaryOp(t, &z.Int{Value: l}, token.Rem,
					&z.Int{Value: r}, &z.Int{Value: l % r})
			}
		}
	}

	// int & int
	testBinaryOp(t,
		&z.Int{Value: 0}, token.And, &z.Int{Value: 0},
		&z.Int{Value: int64(0)})
	testBinaryOp(t,
		&z.Int{Value: 1}, token.And, &z.Int{Value: 0},
		&z.Int{Value: int64(1) & int64(0)})
	testBinaryOp(t,
		&z.Int{Value: 0}, token.And, &z.Int{Value: 1},
		&z.Int{Value: int64(0) & int64(1)})
	testBinaryOp(t,
		&z.Int{Value: 1}, token.And, &z.Int{Value: 1},
		&z.Int{Value: int64(1)})
	testBinaryOp(t,
		&z.Int{Value: 0}, token.And, &z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(0) & int64(0xffffffff)})
	testBinaryOp(t,
		&z.Int{Value: 1}, token.And, &z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(1) & int64(0xffffffff)})
	testBinaryOp(t,
		&z.Int{Value: int64(0xffffffff)}, token.And,
		&z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(0xffffffff)})
	testBinaryOp(t,
		&z.Int{Value: 1984}, token.And,
		&z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(1984) & int64(0xffffffff)})
	testBinaryOp(t, &z.Int{Value: -1984}, token.And,
		&z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(-1984) & int64(0xffffffff)})

	// int | int
	testBinaryOp(t,
		&z.Int{Value: 0}, token.Or, &z.Int{Value: 0},
		&z.Int{Value: int64(0)})
	testBinaryOp(t,
		&z.Int{Value: 1}, token.Or, &z.Int{Value: 0},
		&z.Int{Value: int64(1) | int64(0)})
	testBinaryOp(t,
		&z.Int{Value: 0}, token.Or, &z.Int{Value: 1},
		&z.Int{Value: int64(0) | int64(1)})
	testBinaryOp(t,
		&z.Int{Value: 1}, token.Or, &z.Int{Value: 1},
		&z.Int{Value: int64(1)})
	testBinaryOp(t,
		&z.Int{Value: 0}, token.Or, &z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(0) | int64(0xffffffff)})
	testBinaryOp(t,
		&z.Int{Value: 1}, token.Or, &z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(1) | int64(0xffffffff)})
	testBinaryOp(t,
		&z.Int{Value: int64(0xffffffff)}, token.Or,
		&z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(0xffffffff)})
	testBinaryOp(t,
		&z.Int{Value: 1984}, token.Or,
		&z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(1984) | int64(0xffffffff)})
	testBinaryOp(t,
		&z.Int{Value: -1984}, token.Or,
		&z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(-1984) | int64(0xffffffff)})

	// int ^ int
	testBinaryOp(t,
		&z.Int{Value: 0}, token.Xor, &z.Int{Value: 0},
		&z.Int{Value: int64(0)})
	testBinaryOp(t,
		&z.Int{Value: 1}, token.Xor, &z.Int{Value: 0},
		&z.Int{Value: int64(1) ^ int64(0)})
	testBinaryOp(t,
		&z.Int{Value: 0}, token.Xor, &z.Int{Value: 1},
		&z.Int{Value: int64(0) ^ int64(1)})
	testBinaryOp(t,
		&z.Int{Value: 1}, token.Xor, &z.Int{Value: 1},
		&z.Int{Value: int64(0)})
	testBinaryOp(t,
		&z.Int{Value: 0}, token.Xor, &z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(0) ^ int64(0xffffffff)})
	testBinaryOp(t,
		&z.Int{Value: 1}, token.Xor, &z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(1) ^ int64(0xffffffff)})
	testBinaryOp(t,
		&z.Int{Value: int64(0xffffffff)}, token.Xor,
		&z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(0)})
	testBinaryOp(t,
		&z.Int{Value: 1984}, token.Xor,
		&z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(1984) ^ int64(0xffffffff)})
	testBinaryOp(t,
		&z.Int{Value: -1984}, token.Xor,
		&z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(-1984) ^ int64(0xffffffff)})

	// int &^ int
	testBinaryOp(t,
		&z.Int{Value: 0}, token.AndNot, &z.Int{Value: 0},
		&z.Int{Value: int64(0)})
	testBinaryOp(t,
		&z.Int{Value: 1}, token.AndNot, &z.Int{Value: 0},
		&z.Int{Value: int64(1) &^ int64(0)})
	testBinaryOp(t,
		&z.Int{Value: 0}, token.AndNot,
		&z.Int{Value: 1}, &z.Int{Value: int64(0) &^ int64(1)})
	testBinaryOp(t,
		&z.Int{Value: 1}, token.AndNot, &z.Int{Value: 1},
		&z.Int{Value: int64(0)})
	testBinaryOp(t,
		&z.Int{Value: 0}, token.AndNot,
		&z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(0) &^ int64(0xffffffff)})
	testBinaryOp(t,
		&z.Int{Value: 1}, token.AndNot,
		&z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(1) &^ int64(0xffffffff)})
	testBinaryOp(t,
		&z.Int{Value: int64(0xffffffff)}, token.AndNot,
		&z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(0)})
	testBinaryOp(t,
		&z.Int{Value: 1984}, token.AndNot,
		&z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(1984) &^ int64(0xffffffff)})
	testBinaryOp(t,
		&z.Int{Value: -1984}, token.AndNot,
		&z.Int{Value: int64(0xffffffff)},
		&z.Int{Value: int64(-1984) &^ int64(0xffffffff)})

	// int << int
	for s := int64(0); s < 64; s++ {
		testBinaryOp(t,
			&z.Int{Value: 0}, token.Shl, &z.Int{Value: s},
			&z.Int{Value: int64(0) << uint(s)})
		testBinaryOp(t,
			&z.Int{Value: 1}, token.Shl, &z.Int{Value: s},
			&z.Int{Value: int64(1) << uint(s)})
		testBinaryOp(t,
			&z.Int{Value: 2}, token.Shl, &z.Int{Value: s},
			&z.Int{Value: int64(2) << uint(s)})
		testBinaryOp(t,
			&z.Int{Value: -1}, token.Shl, &z.Int{Value: s},
			&z.Int{Value: int64(-1) << uint(s)})
		testBinaryOp(t,
			&z.Int{Value: -2}, token.Shl, &z.Int{Value: s},
			&z.Int{Value: int64(-2) << uint(s)})
		testBinaryOp(t,
			&z.Int{Value: int64(0xffffffff)}, token.Shl,
			&z.Int{Value: s},
			&z.Int{Value: int64(0xffffffff) << uint(s)})
	}

	// int >> int
	for s := int64(0); s < 64; s++ {
		testBinaryOp(t,
			&z.Int{Value: 0}, token.Shr, &z.Int{Value: s},
			&z.Int{Value: int64(0) >> uint(s)})
		testBinaryOp(t,
			&z.Int{Value: 1}, token.Shr, &z.Int{Value: s},
			&z.Int{Value: int64(1) >> uint(s)})
		testBinaryOp(t,
			&z.Int{Value: 2}, token.Shr, &z.Int{Value: s},
			&z.Int{Value: int64(2) >> uint(s)})
		testBinaryOp(t,
			&z.Int{Value: -1}, token.Shr, &z.Int{Value: s},
			&z.Int{Value: int64(-1) >> uint(s)})
		testBinaryOp(t,
			&z.Int{Value: -2}, token.Shr, &z.Int{Value: s},
			&z.Int{Value: int64(-2) >> uint(s)})
		testBinaryOp(t,
			&z.Int{Value: int64(0xffffffff)}, token.Shr,
			&z.Int{Value: s},
			&z.Int{Value: int64(0xffffffff) >> uint(s)})
	}

	// int < int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &z.Int{Value: l}, token.Less,
				&z.Int{Value: r}, boolValue(l < r))
		}
	}

	// int > int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &z.Int{Value: l}, token.Greater,
				&z.Int{Value: r}, boolValue(l > r))
		}
	}

	// int <= int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &z.Int{Value: l}, token.LessEq,
				&z.Int{Value: r}, boolValue(l <= r))
		}
	}

	// int >= int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &z.Int{Value: l}, token.GreaterEq,
				&z.Int{Value: r}, boolValue(l >= r))
		}
	}

	// int + float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &z.Int{Value: l}, token.Add,
				&z.Float{Value: r},
				&z.Float{Value: float64(l) + r})
		}
	}

	// int - float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &z.Int{Value: l}, token.Sub,
				&z.Float{Value: r},
				&z.Float{Value: float64(l) - r})
		}
	}

	// int * float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &z.Int{Value: l}, token.Mul,
				&z.Float{Value: r},
				&z.Float{Value: float64(l) * r})
		}
	}

	// int / float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			if r != 0 {
				testBinaryOp(t, &z.Int{Value: l}, token.Quo,
					&z.Float{Value: r},
					&z.Float{Value: float64(l) / r})
			}
		}
	}

	// int < float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &z.Int{Value: l}, token.Less,
				&z.Float{Value: r}, boolValue(float64(l) < r))
		}
	}

	// int > float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &z.Int{Value: l}, token.Greater,
				&z.Float{Value: r}, boolValue(float64(l) > r))
		}
	}

	// int <= float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &z.Int{Value: l}, token.LessEq,
				&z.Float{Value: r}, boolValue(float64(l) <= r))
		}
	}

	// int >= float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &z.Int{Value: l}, token.GreaterEq,
				&z.Float{Value: r}, boolValue(float64(l) >= r))
		}
	}
}

func TestMap_Index(t *testing.T) {
	m := &z.Map{Value: make(map[string]z.Object)}
	k := &z.Int{Value: 1}
	v := &z.String{Value: "abcdef"}
	err := m.IndexSet(k, v)

	require.NoError(t, err)

	res, err := m.IndexGet(k)
	require.NoError(t, err)
	require.Equal(t, v, res)
}

func TestString_BinaryOp(t *testing.T) {
	lstr := "abcde"
	rstr := "01234"
	for l := 0; l < len(lstr); l++ {
		for r := 0; r < len(rstr); r++ {
			ls := lstr[l:]
			rs := rstr[r:]
			testBinaryOp(t, &z.String{Value: ls}, token.Add,
				&z.String{Value: rs},
				&z.String{Value: ls + rs})

			rc := []rune(rstr)[r]
			testBinaryOp(t, &z.String{Value: ls}, token.Add,
				&z.Char{Value: rc},
				&z.String{Value: ls + string(rc)})
		}
	}
}

func testBinaryOp(
	t *testing.T,
	lhs z.Object,
	op token.Token,
	rhs z.Object,
	expected z.Object,
) {
	t.Helper()
	actual, err := lhs.BinaryOp(op, rhs)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func boolValue(b bool) z.Object {
	if b {
		return z.TrueValue
	}
	return z.FalseValue
}
