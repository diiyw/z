package stdlib_test

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/require"
	"github.com/d5/tengo/v2/stdlib"
)

func TestFuncAIR(t *testing.T) {
	uf := stdlib.FuncAIR(func(int) {})
	ret, err := funcCall(uf, &z.Int{Value: 10})
	require.NoError(t, err)
	require.Equal(t, z.UndefinedValue, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncAR(t *testing.T) {
	uf := stdlib.FuncAR(func() {})
	ret, err := funcCall(uf)
	require.NoError(t, err)
	require.Equal(t, z.UndefinedValue, ret)
	_, err = funcCall(uf, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncARI(t *testing.T) {
	uf := stdlib.FuncARI(func() int { return 10 })
	ret, err := funcCall(uf)
	require.NoError(t, err)
	require.Equal(t, &z.Int{Value: 10}, ret)
	_, err = funcCall(uf, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncARE(t *testing.T) {
	uf := stdlib.FuncARE(func() error { return nil })
	ret, err := funcCall(uf)
	require.NoError(t, err)
	require.Equal(t, z.TrueValue, ret)
	uf = stdlib.FuncARE(func() error { return errors.New("some error") })
	ret, err = funcCall(uf)
	require.NoError(t, err)
	require.Equal(t, &z.Error{
		Value: &z.String{Value: "some error"},
	}, ret)
	_, err = funcCall(uf, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncARIsE(t *testing.T) {
	uf := stdlib.FuncARIsE(func() ([]int, error) {
		return []int{1, 2, 3}, nil
	})
	ret, err := funcCall(uf)
	require.NoError(t, err)
	require.Equal(t, array(&z.Int{Value: 1},
		&z.Int{Value: 2}, &z.Int{Value: 3}), ret)
	uf = stdlib.FuncARIsE(func() ([]int, error) {
		return nil, errors.New("some error")
	})
	ret, err = funcCall(uf)
	require.NoError(t, err)
	require.Equal(t, &z.Error{
		Value: &z.String{Value: "some error"},
	}, ret)
	_, err = funcCall(uf, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncARS(t *testing.T) {
	uf := stdlib.FuncARS(func() string { return "foo" })
	ret, err := funcCall(uf)
	require.NoError(t, err)
	require.Equal(t, &z.String{Value: "foo"}, ret)
	_, err = funcCall(uf, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncARSE(t *testing.T) {
	uf := stdlib.FuncARSE(func() (string, error) { return "foo", nil })
	ret, err := funcCall(uf)
	require.NoError(t, err)
	require.Equal(t, &z.String{Value: "foo"}, ret)
	uf = stdlib.FuncARSE(func() (string, error) {
		return "", errors.New("some error")
	})
	ret, err = funcCall(uf)
	require.NoError(t, err)
	require.Equal(t, &z.Error{
		Value: &z.String{Value: "some error"},
	}, ret)
	_, err = funcCall(uf, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncARSs(t *testing.T) {
	uf := stdlib.FuncARSs(func() []string { return []string{"foo", "bar"} })
	ret, err := funcCall(uf)
	require.NoError(t, err)
	require.Equal(t, array(&z.String{Value: "foo"},
		&z.String{Value: "bar"}), ret)
	_, err = funcCall(uf, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncASRE(t *testing.T) {
	uf := stdlib.FuncASRE(func(a string) error { return nil })
	ret, err := funcCall(uf, &z.String{Value: "foo"})
	require.NoError(t, err)
	require.Equal(t, z.TrueValue, ret)
	uf = stdlib.FuncASRE(func(a string) error {
		return errors.New("some error")
	})
	ret, err = funcCall(uf, &z.String{Value: "foo"})
	require.NoError(t, err)
	require.Equal(t, &z.Error{
		Value: &z.String{Value: "some error"},
	}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncASRS(t *testing.T) {
	uf := stdlib.FuncASRS(func(a string) string { return a })
	ret, err := funcCall(uf, &z.String{Value: "foo"})
	require.NoError(t, err)
	require.Equal(t, &z.String{Value: "foo"}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncASRSs(t *testing.T) {
	uf := stdlib.FuncASRSs(func(a string) []string { return []string{a} })
	ret, err := funcCall(uf, &z.String{Value: "foo"})
	require.NoError(t, err)
	require.Equal(t, array(&z.String{Value: "foo"}), ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncASI64RE(t *testing.T) {
	uf := stdlib.FuncASI64RE(func(a string, b int64) error { return nil })
	ret, err := funcCall(uf, &z.String{Value: "foo"}, &z.Int{Value: 5})
	require.NoError(t, err)
	require.Equal(t, z.TrueValue, ret)
	uf = stdlib.FuncASI64RE(func(a string, b int64) error {
		return errors.New("some error")
	})
	ret, err = funcCall(uf, &z.String{Value: "foo"}, &z.Int{Value: 5})
	require.NoError(t, err)
	require.Equal(t,
		&z.Error{Value: &z.String{Value: "some error"}}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncAIIRE(t *testing.T) {
	uf := stdlib.FuncAIIRE(func(a, b int) error { return nil })
	ret, err := funcCall(uf, &z.Int{Value: 5}, &z.Int{Value: 7})
	require.NoError(t, err)
	require.Equal(t, z.TrueValue, ret)
	uf = stdlib.FuncAIIRE(func(a, b int) error {
		return errors.New("some error")
	})
	ret, err = funcCall(uf, &z.Int{Value: 5}, &z.Int{Value: 7})
	require.NoError(t, err)
	require.Equal(t,
		&z.Error{Value: &z.String{Value: "some error"}}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncASIIRE(t *testing.T) {
	uf := stdlib.FuncASIIRE(func(a string, b, c int) error { return nil })
	ret, err := funcCall(uf, &z.String{Value: "foo"}, &z.Int{Value: 5},
		&z.Int{Value: 7})
	require.NoError(t, err)
	require.Equal(t, z.TrueValue, ret)
	uf = stdlib.FuncASIIRE(func(a string, b, c int) error {
		return errors.New("some error")
	})
	ret, err = funcCall(uf, &z.String{Value: "foo"}, &z.Int{Value: 5},
		&z.Int{Value: 7})
	require.NoError(t, err)
	require.Equal(t,
		&z.Error{Value: &z.String{Value: "some error"}}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncASRSE(t *testing.T) {
	uf := stdlib.FuncASRSE(func(a string) (string, error) { return a, nil })
	ret, err := funcCall(uf, &z.String{Value: "foo"})
	require.NoError(t, err)
	require.Equal(t, &z.String{Value: "foo"}, ret)
	uf = stdlib.FuncASRSE(func(a string) (string, error) {
		return a, errors.New("some error")
	})
	ret, err = funcCall(uf, &z.String{Value: "foo"})
	require.NoError(t, err)
	require.Equal(t,
		&z.Error{Value: &z.String{Value: "some error"}}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncASSRE(t *testing.T) {
	uf := stdlib.FuncASSRE(func(a, b string) error { return nil })
	ret, err := funcCall(uf, &z.String{Value: "foo"},
		&z.String{Value: "bar"})
	require.NoError(t, err)
	require.Equal(t, z.TrueValue, ret)
	uf = stdlib.FuncASSRE(func(a, b string) error {
		return errors.New("some error")
	})
	ret, err = funcCall(uf, &z.String{Value: "foo"},
		&z.String{Value: "bar"})
	require.NoError(t, err)
	require.Equal(t,
		&z.Error{Value: &z.String{Value: "some error"}}, ret)
	_, err = funcCall(uf, &z.String{Value: "foo"})
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncASsRS(t *testing.T) {
	uf := stdlib.FuncASsSRS(func(a []string, b string) string {
		return strings.Join(a, b)
	})
	ret, err := funcCall(uf, array(&z.String{Value: "foo"},
		&z.String{Value: "bar"}), &z.String{Value: " "})
	require.NoError(t, err)
	require.Equal(t, &z.String{Value: "foo bar"}, ret)
	_, err = funcCall(uf, &z.String{Value: "foo"})
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncARF(t *testing.T) {
	uf := stdlib.FuncARF(func() float64 { return 10.0 })
	ret, err := funcCall(uf)
	require.NoError(t, err)
	require.Equal(t, &z.Float{Value: 10.0}, ret)
	_, err = funcCall(uf, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncAFRF(t *testing.T) {
	uf := stdlib.FuncAFRF(func(a float64) float64 { return a })
	ret, err := funcCall(uf, &z.Float{Value: 10.0})
	require.NoError(t, err)
	require.Equal(t, &z.Float{Value: 10.0}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
	_, err = funcCall(uf, z.TrueValue, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncAIRF(t *testing.T) {
	uf := stdlib.FuncAIRF(func(a int) float64 {
		return float64(a)
	})
	ret, err := funcCall(uf, &z.Int{Value: 10.0})
	require.NoError(t, err)
	require.Equal(t, &z.Float{Value: 10.0}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
	_, err = funcCall(uf, z.TrueValue, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncAFRI(t *testing.T) {
	uf := stdlib.FuncAFRI(func(a float64) int {
		return int(a)
	})
	ret, err := funcCall(uf, &z.Float{Value: 10.5})
	require.NoError(t, err)
	require.Equal(t, &z.Int{Value: 10}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
	_, err = funcCall(uf, z.TrueValue, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncAFRB(t *testing.T) {
	uf := stdlib.FuncAFRB(func(a float64) bool {
		return a > 0.0
	})
	ret, err := funcCall(uf, &z.Float{Value: 0.1})
	require.NoError(t, err)
	require.Equal(t, z.TrueValue, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
	_, err = funcCall(uf, z.TrueValue, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncAFFRF(t *testing.T) {
	uf := stdlib.FuncAFFRF(func(a, b float64) float64 {
		return a + b
	})
	ret, err := funcCall(uf, &z.Float{Value: 10.0},
		&z.Float{Value: 20.0})
	require.NoError(t, err)
	require.Equal(t, &z.Float{Value: 30.0}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
	_, err = funcCall(uf, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncASIRS(t *testing.T) {
	uf := stdlib.FuncASIRS(func(a string, b int) string {
		return strings.Repeat(a, b)
	})
	ret, err := funcCall(uf, &z.String{Value: "ab"}, &z.Int{Value: 2})
	require.NoError(t, err)
	require.Equal(t, &z.String{Value: "abab"}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
	_, err = funcCall(uf, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncAIFRF(t *testing.T) {
	uf := stdlib.FuncAIFRF(func(a int, b float64) float64 {
		return float64(a) + b
	})
	ret, err := funcCall(uf, &z.Int{Value: 10}, &z.Float{Value: 20.0})
	require.NoError(t, err)
	require.Equal(t, &z.Float{Value: 30.0}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
	_, err = funcCall(uf, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncAFIRF(t *testing.T) {
	uf := stdlib.FuncAFIRF(func(a float64, b int) float64 {
		return a + float64(b)
	})
	ret, err := funcCall(uf, &z.Float{Value: 10.0}, &z.Int{Value: 20})
	require.NoError(t, err)
	require.Equal(t, &z.Float{Value: 30.0}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
	_, err = funcCall(uf, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncAFIRB(t *testing.T) {
	uf := stdlib.FuncAFIRB(func(a float64, b int) bool {
		return a < float64(b)
	})
	ret, err := funcCall(uf, &z.Float{Value: 10.0}, &z.Int{Value: 20})
	require.NoError(t, err)
	require.Equal(t, z.TrueValue, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
	_, err = funcCall(uf, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncAIRSsE(t *testing.T) {
	uf := stdlib.FuncAIRSsE(func(a int) ([]string, error) {
		return []string{"foo", "bar"}, nil
	})
	ret, err := funcCall(uf, &z.Int{Value: 10})
	require.NoError(t, err)
	require.Equal(t, array(&z.String{Value: "foo"},
		&z.String{Value: "bar"}), ret)
	uf = stdlib.FuncAIRSsE(func(a int) ([]string, error) {
		return nil, errors.New("some error")
	})
	ret, err = funcCall(uf, &z.Int{Value: 10})
	require.NoError(t, err)
	require.Equal(t,
		&z.Error{Value: &z.String{Value: "some error"}}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncASSRSs(t *testing.T) {
	uf := stdlib.FuncASSRSs(func(a, b string) []string {
		return []string{a, b}
	})
	ret, err := funcCall(uf, &z.String{Value: "foo"},
		&z.String{Value: "bar"})
	require.NoError(t, err)
	require.Equal(t, array(&z.String{Value: "foo"},
		&z.String{Value: "bar"}), ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncASSIRSs(t *testing.T) {
	uf := stdlib.FuncASSIRSs(func(a, b string, c int) []string {
		return []string{a, b, strconv.Itoa(c)}
	})
	ret, err := funcCall(uf, &z.String{Value: "foo"},
		&z.String{Value: "bar"}, &z.Int{Value: 5})
	require.NoError(t, err)
	require.Equal(t, array(&z.String{Value: "foo"},
		&z.String{Value: "bar"}, &z.String{Value: "5"}), ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncARB(t *testing.T) {
	uf := stdlib.FuncARB(func() bool { return true })
	ret, err := funcCall(uf)
	require.NoError(t, err)
	require.Equal(t, z.TrueValue, ret)
	_, err = funcCall(uf, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncARYE(t *testing.T) {
	uf := stdlib.FuncARYE(func() ([]byte, error) {
		return []byte("foo bar"), nil
	})
	ret, err := funcCall(uf)
	require.NoError(t, err)
	require.Equal(t, &z.Bytes{Value: []byte("foo bar")}, ret)
	uf = stdlib.FuncARYE(func() ([]byte, error) {
		return nil, errors.New("some error")
	})
	ret, err = funcCall(uf)
	require.NoError(t, err)
	require.Equal(t,
		&z.Error{Value: &z.String{Value: "some error"}}, ret)
	_, err = funcCall(uf, z.TrueValue)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncASRIE(t *testing.T) {
	uf := stdlib.FuncASRIE(func(a string) (int, error) { return 5, nil })
	ret, err := funcCall(uf, &z.String{Value: "foo"})
	require.NoError(t, err)
	require.Equal(t, &z.Int{Value: 5}, ret)
	uf = stdlib.FuncASRIE(func(a string) (int, error) {
		return 0, errors.New("some error")
	})
	ret, err = funcCall(uf, &z.String{Value: "foo"})
	require.NoError(t, err)
	require.Equal(t,
		&z.Error{Value: &z.String{Value: "some error"}}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncAYRIE(t *testing.T) {
	uf := stdlib.FuncAYRIE(func(a []byte) (int, error) { return 5, nil })
	ret, err := funcCall(uf, &z.Bytes{Value: []byte("foo")})
	require.NoError(t, err)
	require.Equal(t, &z.Int{Value: 5}, ret)
	uf = stdlib.FuncAYRIE(func(a []byte) (int, error) {
		return 0, errors.New("some error")
	})
	ret, err = funcCall(uf, &z.Bytes{Value: []byte("foo")})
	require.NoError(t, err)
	require.Equal(t,
		&z.Error{Value: &z.String{Value: "some error"}}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncASSRI(t *testing.T) {
	uf := stdlib.FuncASSRI(func(a, b string) int { return len(a) + len(b) })
	ret, err := funcCall(uf,
		&z.String{Value: "foo"}, &z.String{Value: "bar"})
	require.NoError(t, err)
	require.Equal(t, &z.Int{Value: 6}, ret)
	_, err = funcCall(uf, &z.String{Value: "foo"})
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncASSRS(t *testing.T) {
	uf := stdlib.FuncASSRS(func(a, b string) string { return a + b })
	ret, err := funcCall(uf,
		&z.String{Value: "foo"}, &z.String{Value: "bar"})
	require.NoError(t, err)
	require.Equal(t, &z.String{Value: "foobar"}, ret)
	_, err = funcCall(uf, &z.String{Value: "foo"})
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncASSRB(t *testing.T) {
	uf := stdlib.FuncASSRB(func(a, b string) bool { return len(a) > len(b) })
	ret, err := funcCall(uf,
		&z.String{Value: "123"}, &z.String{Value: "12"})
	require.NoError(t, err)
	require.Equal(t, z.TrueValue, ret)
	_, err = funcCall(uf, &z.String{Value: "foo"})
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncAIRS(t *testing.T) {
	uf := stdlib.FuncAIRS(func(a int) string { return strconv.Itoa(a) })
	ret, err := funcCall(uf, &z.Int{Value: 55})
	require.NoError(t, err)
	require.Equal(t, &z.String{Value: "55"}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncAIRIs(t *testing.T) {
	uf := stdlib.FuncAIRIs(func(a int) []int { return []int{a, a} })
	ret, err := funcCall(uf, &z.Int{Value: 55})
	require.NoError(t, err)
	require.Equal(t, array(&z.Int{Value: 55}, &z.Int{Value: 55}), ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncAI64R(t *testing.T) {
	uf := stdlib.FuncAIR(func(a int) {})
	ret, err := funcCall(uf, &z.Int{Value: 55})
	require.NoError(t, err)
	require.Equal(t, z.UndefinedValue, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncARI64(t *testing.T) {
	uf := stdlib.FuncARI64(func() int64 { return 55 })
	ret, err := funcCall(uf)
	require.NoError(t, err)
	require.Equal(t, &z.Int{Value: 55}, ret)
	_, err = funcCall(uf, &z.Int{Value: 55})
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncASsSRS(t *testing.T) {
	uf := stdlib.FuncASsSRS(func(a []string, b string) string {
		return strings.Join(a, b)
	})
	ret, err := funcCall(uf,
		array(&z.String{Value: "abc"}, &z.String{Value: "def"}),
		&z.String{Value: "-"})
	require.NoError(t, err)
	require.Equal(t, &z.String{Value: "abc-def"}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func TestFuncAI64RI64(t *testing.T) {
	uf := stdlib.FuncAI64RI64(func(a int64) int64 { return a * 2 })
	ret, err := funcCall(uf, &z.Int{Value: 55})
	require.NoError(t, err)
	require.Equal(t, &z.Int{Value: 110}, ret)
	_, err = funcCall(uf)
	require.Equal(t, z.ErrWrongNumArguments, err)
}

func funcCall(
	fn z.CallableFunc,
	args ...z.Object,
) (z.Object, error) {
	userFunc := &z.UserFunction{Value: fn}
	return userFunc.Call(args...)
}

func array(elements ...z.Object) *z.Array {
	return &z.Array{Value: elements}
}
