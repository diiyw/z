package z_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/diiyw/z"
)

func Test_builtinDelete(t *testing.T) {
	var builtinDelete func(args ...z.Object) (z.Object, error)
	for _, f := range z.GetAllBuiltinFunctions() {
		if f.Name == "delete" {
			builtinDelete = f.Value
			break
		}
	}
	if builtinDelete == nil {
		t.Fatal("builtin delete not found")
	}
	type args struct {
		args []z.Object
	}
	tests := []struct {
		name      string
		args      args
		want      z.Object
		wantErr   bool
		wantedErr error
		target    interface{}
	}{
		{name: "invalid-arg", args: args{[]z.Object{&z.String{},
			&z.String{}}}, wantErr: true,
			wantedErr: z.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "map",
				Found:    "string"},
		},
		{name: "no-args",
			wantErr: true, wantedErr: z.ErrWrongNumArguments},
		{name: "empty-args", args: args{[]z.Object{}}, wantErr: true,
			wantedErr: z.ErrWrongNumArguments,
		},
		{name: "3-args", args: args{[]z.Object{
			(*z.Map)(nil), (*z.String)(nil), (*z.String)(nil)}},
			wantErr: true, wantedErr: z.ErrWrongNumArguments,
		},
		{name: "nil-map-empty-key",
			args: args{[]z.Object{&z.Map{}, &z.String{}}},
			want: z.UndefinedValue,
		},
		{name: "nil-map-nonstr-key",
			args: args{[]z.Object{
				&z.Map{}, &z.Int{}}}, wantErr: true,
			wantedErr: z.ErrInvalidArgumentType{
				Name: "second", Expected: "string", Found: "int"},
		},
		{name: "nil-map-no-key",
			args: args{[]z.Object{&z.Map{}}}, wantErr: true,
			wantedErr: z.ErrWrongNumArguments,
		},
		{name: "map-missing-key",
			args: args{
				[]z.Object{
					&z.Map{Value: map[string]z.Object{
						"key": &z.String{Value: "value"},
					}},
					&z.String{Value: "key1"}}},
			want: z.UndefinedValue,
			target: &z.Map{
				Value: map[string]z.Object{
					"key": &z.String{
						Value: "value"}}},
		},
		{name: "map-emptied",
			args: args{
				[]z.Object{
					&z.Map{Value: map[string]z.Object{
						"key": &z.String{Value: "value"},
					}},
					&z.String{Value: "key"}}},
			want:   z.UndefinedValue,
			target: &z.Map{Value: map[string]z.Object{}},
		},
		{name: "map-multi-keys",
			args: args{
				[]z.Object{
					&z.Map{Value: map[string]z.Object{
						"key1": &z.String{Value: "value1"},
						"key2": &z.Int{Value: 10},
					}},
					&z.String{Value: "key1"}}},
			want: z.UndefinedValue,
			target: &z.Map{Value: map[string]z.Object{
				"key2": &z.Int{Value: 10}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := builtinDelete(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("builtinDelete() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if tt.wantErr && !errors.Is(err, tt.wantedErr) {
				if err.Error() != tt.wantedErr.Error() {
					t.Errorf("builtinDelete() error = %v, wantedErr %v",
						err, tt.wantedErr)
					return
				}
			}
			if got != tt.want {
				t.Errorf("builtinDelete() = %v, want %v", got, tt.want)
				return
			}
			if !tt.wantErr && tt.target != nil {
				switch v := tt.args.args[0].(type) {
				case *z.Map, *z.Array:
					if !reflect.DeepEqual(tt.target, tt.args.args[0]) {
						t.Errorf("builtinDelete() objects are not equal "+
							"got: %+v, want: %+v", tt.args.args[0], tt.target)
					}
				default:
					t.Errorf("builtinDelete() unsuporrted arg[0] type %s",
						v.TypeName())
					return
				}
			}
		})
	}
}

func Test_builtinSplice(t *testing.T) {
	var builtinSplice func(args ...z.Object) (z.Object, error)
	for _, f := range z.GetAllBuiltinFunctions() {
		if f.Name == "splice" {
			builtinSplice = f.Value
			break
		}
	}
	if builtinSplice == nil {
		t.Fatal("builtin splice not found")
	}
	tests := []struct {
		name      string
		args      []z.Object
		deleted   z.Object
		Array     *z.Array
		wantErr   bool
		wantedErr error
	}{
		{name: "no args", args: []z.Object{}, wantErr: true,
			wantedErr: z.ErrWrongNumArguments,
		},
		{name: "invalid args", args: []z.Object{&z.Map{}},
			wantErr: true,
			wantedErr: z.ErrInvalidArgumentType{
				Name: "first", Expected: "array", Found: "map"},
		},
		{name: "invalid args",
			args:    []z.Object{&z.Array{}, &z.String{}},
			wantErr: true,
			wantedErr: z.ErrInvalidArgumentType{
				Name: "second", Expected: "int", Found: "string"},
		},
		{name: "negative index",
			args:      []z.Object{&z.Array{}, &z.Int{Value: -1}},
			wantErr:   true,
			wantedErr: z.ErrIndexOutOfBounds},
		{name: "non int count",
			args: []z.Object{
				&z.Array{}, &z.Int{Value: 0},
				&z.String{Value: ""}},
			wantErr: true,
			wantedErr: z.ErrInvalidArgumentType{
				Name: "third", Expected: "int", Found: "string"},
		},
		{name: "negative count",
			args: []z.Object{
				&z.Array{Value: []z.Object{
					&z.Int{Value: 0},
					&z.Int{Value: 1},
					&z.Int{Value: 2}}},
				&z.Int{Value: 0},
				&z.Int{Value: -1}},
			wantErr:   true,
			wantedErr: z.ErrIndexOutOfBounds,
		},
		{name: "insert with zero count",
			args: []z.Object{
				&z.Array{Value: []z.Object{
					&z.Int{Value: 0},
					&z.Int{Value: 1},
					&z.Int{Value: 2}}},
				&z.Int{Value: 0},
				&z.Int{Value: 0},
				&z.String{Value: "b"}},
			deleted: &z.Array{Value: []z.Object{}},
			Array: &z.Array{Value: []z.Object{
				&z.String{Value: "b"},
				&z.Int{Value: 0},
				&z.Int{Value: 1},
				&z.Int{Value: 2}}},
		},
		{name: "insert",
			args: []z.Object{
				&z.Array{Value: []z.Object{
					&z.Int{Value: 0},
					&z.Int{Value: 1},
					&z.Int{Value: 2}}},
				&z.Int{Value: 1},
				&z.Int{Value: 0},
				&z.String{Value: "c"},
				&z.String{Value: "d"}},
			deleted: &z.Array{Value: []z.Object{}},
			Array: &z.Array{Value: []z.Object{
				&z.Int{Value: 0},
				&z.String{Value: "c"},
				&z.String{Value: "d"},
				&z.Int{Value: 1},
				&z.Int{Value: 2}}},
		},
		{name: "insert with zero count",
			args: []z.Object{
				&z.Array{Value: []z.Object{
					&z.Int{Value: 0},
					&z.Int{Value: 1},
					&z.Int{Value: 2}}},
				&z.Int{Value: 1},
				&z.Int{Value: 0},
				&z.String{Value: "c"},
				&z.String{Value: "d"}},
			deleted: &z.Array{Value: []z.Object{}},
			Array: &z.Array{Value: []z.Object{
				&z.Int{Value: 0},
				&z.String{Value: "c"},
				&z.String{Value: "d"},
				&z.Int{Value: 1},
				&z.Int{Value: 2}}},
		},
		{name: "insert with delete",
			args: []z.Object{
				&z.Array{Value: []z.Object{
					&z.Int{Value: 0},
					&z.Int{Value: 1},
					&z.Int{Value: 2}}},
				&z.Int{Value: 1},
				&z.Int{Value: 1},
				&z.String{Value: "c"},
				&z.String{Value: "d"}},
			deleted: &z.Array{
				Value: []z.Object{&z.Int{Value: 1}}},
			Array: &z.Array{Value: []z.Object{
				&z.Int{Value: 0},
				&z.String{Value: "c"},
				&z.String{Value: "d"},
				&z.Int{Value: 2}}},
		},
		{name: "insert with delete multi",
			args: []z.Object{
				&z.Array{Value: []z.Object{
					&z.Int{Value: 0},
					&z.Int{Value: 1},
					&z.Int{Value: 2}}},
				&z.Int{Value: 1},
				&z.Int{Value: 2},
				&z.String{Value: "c"},
				&z.String{Value: "d"}},
			deleted: &z.Array{Value: []z.Object{
				&z.Int{Value: 1},
				&z.Int{Value: 2}}},
			Array: &z.Array{
				Value: []z.Object{
					&z.Int{Value: 0},
					&z.String{Value: "c"},
					&z.String{Value: "d"}}},
		},
		{name: "delete all with positive count",
			args: []z.Object{
				&z.Array{Value: []z.Object{
					&z.Int{Value: 0},
					&z.Int{Value: 1},
					&z.Int{Value: 2}}},
				&z.Int{Value: 0},
				&z.Int{Value: 3}},
			deleted: &z.Array{Value: []z.Object{
				&z.Int{Value: 0},
				&z.Int{Value: 1},
				&z.Int{Value: 2}}},
			Array: &z.Array{Value: []z.Object{}},
		},
		{name: "delete all with big count",
			args: []z.Object{
				&z.Array{Value: []z.Object{
					&z.Int{Value: 0},
					&z.Int{Value: 1},
					&z.Int{Value: 2}}},
				&z.Int{Value: 0},
				&z.Int{Value: 5}},
			deleted: &z.Array{Value: []z.Object{
				&z.Int{Value: 0},
				&z.Int{Value: 1},
				&z.Int{Value: 2}}},
			Array: &z.Array{Value: []z.Object{}},
		},
		{name: "nothing2",
			args: []z.Object{
				&z.Array{Value: []z.Object{
					&z.Int{Value: 0},
					&z.Int{Value: 1},
					&z.Int{Value: 2}}}},
			Array: &z.Array{Value: []z.Object{}},
			deleted: &z.Array{Value: []z.Object{
				&z.Int{Value: 0},
				&z.Int{Value: 1},
				&z.Int{Value: 2}}},
		},
		{name: "pop without count",
			args: []z.Object{
				&z.Array{Value: []z.Object{
					&z.Int{Value: 0},
					&z.Int{Value: 1},
					&z.Int{Value: 2}}},
				&z.Int{Value: 2}},
			deleted: &z.Array{Value: []z.Object{&z.Int{Value: 2}}},
			Array: &z.Array{Value: []z.Object{
				&z.Int{Value: 0}, &z.Int{Value: 1}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := builtinSplice(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("builtinSplice() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.deleted) {
				t.Errorf("builtinSplice() = %v, want %v", got, tt.deleted)
			}
			if tt.wantErr && tt.wantedErr.Error() != err.Error() {
				t.Errorf("builtinSplice() error = %v, wantedErr %v",
					err, tt.wantedErr)
			}
			if tt.Array != nil && !reflect.DeepEqual(tt.Array, tt.args[0]) {
				t.Errorf("builtinSplice() arrays are not equal expected"+
					" %s, got %s", tt.Array, tt.args[0].(*z.Array))
			}
		})
	}
}

func Test_builtinRange(t *testing.T) {
	var builtinRange func(args ...z.Object) (z.Object, error)
	for _, f := range z.GetAllBuiltinFunctions() {
		if f.Name == "range" {
			builtinRange = f.Value
			break
		}
	}
	if builtinRange == nil {
		t.Fatal("builtin range not found")
	}
	tests := []struct {
		name      string
		args      []z.Object
		result    *z.Array
		wantErr   bool
		wantedErr error
	}{
		{name: "no args", args: []z.Object{}, wantErr: true,
			wantedErr: z.ErrWrongNumArguments,
		},
		{name: "single args", args: []z.Object{&z.Map{}},
			wantErr:   true,
			wantedErr: z.ErrWrongNumArguments,
		},
		{name: "4 args", args: []z.Object{&z.Map{}, &z.String{}, &z.String{}, &z.String{}},
			wantErr:   true,
			wantedErr: z.ErrWrongNumArguments,
		},
		{name: "invalid start",
			args:    []z.Object{&z.String{}, &z.String{}},
			wantErr: true,
			wantedErr: z.ErrInvalidArgumentType{
				Name: "start", Expected: "int", Found: "string"},
		},
		{name: "invalid stop",
			args:    []z.Object{&z.Int{}, &z.String{}},
			wantErr: true,
			wantedErr: z.ErrInvalidArgumentType{
				Name: "stop", Expected: "int", Found: "string"},
		},
		{name: "invalid step",
			args:    []z.Object{&z.Int{}, &z.Int{}, &z.String{}},
			wantErr: true,
			wantedErr: z.ErrInvalidArgumentType{
				Name: "step", Expected: "int", Found: "string"},
		},
		{name: "zero step",
			args:      []z.Object{&z.Int{}, &z.Int{}, &z.Int{}}, //must greate than 0
			wantErr:   true,
			wantedErr: z.ErrInvalidRangeStep,
		},
		{name: "negative step",
			args:      []z.Object{&z.Int{}, &z.Int{}, intObject(-2)}, //must greate than 0
			wantErr:   true,
			wantedErr: z.ErrInvalidRangeStep,
		},
		{name: "same bound",
			args:    []z.Object{&z.Int{}, &z.Int{}},
			wantErr: false,
			result: &z.Array{
				Value: nil,
			},
		},
		{name: "positive range",
			args:    []z.Object{&z.Int{}, &z.Int{Value: 5}},
			wantErr: false,
			result: &z.Array{
				Value: []z.Object{
					intObject(0),
					intObject(1),
					intObject(2),
					intObject(3),
					intObject(4),
				},
			},
		},
		{name: "negative range",
			args:    []z.Object{&z.Int{}, &z.Int{Value: -5}},
			wantErr: false,
			result: &z.Array{
				Value: []z.Object{
					intObject(0),
					intObject(-1),
					intObject(-2),
					intObject(-3),
					intObject(-4),
				},
			},
		},

		{name: "positive with step",
			args:    []z.Object{&z.Int{}, &z.Int{Value: 5}, &z.Int{Value: 2}},
			wantErr: false,
			result: &z.Array{
				Value: []z.Object{
					intObject(0),
					intObject(2),
					intObject(4),
				},
			},
		},

		{name: "negative with step",
			args:    []z.Object{&z.Int{}, &z.Int{Value: -10}, &z.Int{Value: 2}},
			wantErr: false,
			result: &z.Array{
				Value: []z.Object{
					intObject(0),
					intObject(-2),
					intObject(-4),
					intObject(-6),
					intObject(-8),
				},
			},
		},

		{name: "large range",
			args:    []z.Object{intObject(-10), intObject(10), &z.Int{Value: 3}},
			wantErr: false,
			result: &z.Array{
				Value: []z.Object{
					intObject(-10),
					intObject(-7),
					intObject(-4),
					intObject(-1),
					intObject(2),
					intObject(5),
					intObject(8),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := builtinRange(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("builtinRange() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.wantedErr.Error() != err.Error() {
				t.Errorf("builtinRange() error = %v, wantedErr %v",
					err, tt.wantedErr)
			}
			if tt.result != nil && !reflect.DeepEqual(tt.result, got) {
				t.Errorf("builtinRange() arrays are not equal expected"+
					" %s, got %s", tt.result, got.(*z.Array))
			}
		})
	}
}
