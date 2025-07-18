package z_test

import (
	"testing"

	"github.com/diiyw/z"
	"github.com/diiyw/z/require"
)

type VariableTest struct {
	Name        string
	Value       any
	ValueType   string
	IntValue    int
	Int64Value  int64
	FloatValue  float64
	CharValue   rune
	BoolValue   bool
	StringValue string
	Object      z.Object
	IsUndefined bool
}

func TestVariable(t *testing.T) {
	vars := []VariableTest{
		{
			Name:        "a",
			Value:       int64(1),
			ValueType:   "int",
			IntValue:    1,
			Int64Value:  1,
			FloatValue:  1.0,
			CharValue:   rune(1),
			BoolValue:   true,
			StringValue: "1",
			Object:      &z.Int{Value: 1},
		},
		{
			Name:        "b",
			Value:       "52.11",
			ValueType:   "string",
			FloatValue:  52.11,
			StringValue: "52.11",
			BoolValue:   true,
			Object:      &z.String{Value: "52.11"},
		},
		{
			Name:        "c",
			Value:       true,
			ValueType:   "bool",
			IntValue:    1,
			Int64Value:  1,
			FloatValue:  0,
			BoolValue:   true,
			StringValue: "true",
			Object:      z.TrueValue,
		},
		{
			Name:        "d",
			Value:       nil,
			ValueType:   "undefined",
			Object:      z.UndefinedValue,
			IsUndefined: true,
		},
	}

	for _, tc := range vars {
		v, err := z.NewVariable(tc.Name, tc.Value)
		require.NoError(t, err)
		require.Equal(t, tc.Value, v.Value(), "Name: %s", tc.Name)
		require.Equal(t, tc.ValueType, v.ValueType(), "Name: %s", tc.Name)
		require.Equal(t, tc.IntValue, v.Int(), "Name: %s", tc.Name)
		require.Equal(t, tc.Int64Value, v.Int64(), "Name: %s", tc.Name)
		require.Equal(t, tc.FloatValue, v.Float(), "Name: %s", tc.Name)
		require.Equal(t, tc.CharValue, v.Char(), "Name: %s", tc.Name)
		require.Equal(t, tc.BoolValue, v.Bool(), "Name: %s", tc.Name)
		require.Equal(t, tc.StringValue, v.String(), "Name: %s", tc.Name)
		require.Equal(t, tc.Object, v.Object(), "Name: %s", tc.Name)
		require.Equal(t, tc.IsUndefined, v.IsUndefined(), "Name: %s", tc.Name)
	}
}
