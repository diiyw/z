package stdlib

import (
	"math/rand"

	"github.com/d5/tengo/v2"
)

var randModule = map[string]z.Object{
	"int": &z.UserFunction{
		Name:  "int",
		Value: FuncARI64(rand.Int63),
	},
	"float": &z.UserFunction{
		Name:  "float",
		Value: FuncARF(rand.Float64),
	},
	"intn": &z.UserFunction{
		Name:  "intn",
		Value: FuncAI64RI64(rand.Int63n),
	},
	"exp_float": &z.UserFunction{
		Name:  "exp_float",
		Value: FuncARF(rand.ExpFloat64),
	},
	"norm_float": &z.UserFunction{
		Name:  "norm_float",
		Value: FuncARF(rand.NormFloat64),
	},
	"perm": &z.UserFunction{
		Name:  "perm",
		Value: FuncAIRIs(rand.Perm),
	},
	"seed": &z.UserFunction{
		Name:  "seed",
		Value: FuncAI64R(rand.Seed),
	},
	"read": &z.UserFunction{
		Name: "read",
		Value: func(args ...z.Object) (ret z.Object, err error) {
			if len(args) != 1 {
				return nil, z.ErrWrongNumArguments
			}
			y1, ok := args[0].(*z.Bytes)
			if !ok {
				return nil, z.ErrInvalidArgumentType{
					Name:     "first",
					Expected: "bytes",
					Found:    args[0].TypeName(),
				}
			}
			res, err := rand.Read(y1.Value)
			if err != nil {
				ret = wrapError(err)
				return
			}
			return &z.Int{Value: int64(res)}, nil
		},
	},
	"rand": &z.UserFunction{
		Name: "rand",
		Value: func(args ...z.Object) (z.Object, error) {
			if len(args) != 1 {
				return nil, z.ErrWrongNumArguments
			}
			i1, ok := z.ToInt64(args[0])
			if !ok {
				return nil, z.ErrInvalidArgumentType{
					Name:     "first",
					Expected: "int(compatible)",
					Found:    args[0].TypeName(),
				}
			}
			src := rand.NewSource(i1)
			return randRand(rand.New(src)), nil
		},
	},
}

func randRand(r *rand.Rand) *z.ImmutableMap {
	return &z.ImmutableMap{
		Value: map[string]z.Object{
			"int": &z.UserFunction{
				Name:  "int",
				Value: FuncARI64(r.Int63),
			},
			"float": &z.UserFunction{
				Name:  "float",
				Value: FuncARF(r.Float64),
			},
			"intn": &z.UserFunction{
				Name:  "intn",
				Value: FuncAI64RI64(r.Int63n),
			},
			"exp_float": &z.UserFunction{
				Name:  "exp_float",
				Value: FuncARF(r.ExpFloat64),
			},
			"norm_float": &z.UserFunction{
				Name:  "norm_float",
				Value: FuncARF(r.NormFloat64),
			},
			"perm": &z.UserFunction{
				Name:  "perm",
				Value: FuncAIRIs(r.Perm),
			},
			"seed": &z.UserFunction{
				Name:  "seed",
				Value: FuncAI64R(r.Seed),
			},
			"read": &z.UserFunction{
				Name: "read",
				Value: func(args ...z.Object) (
					ret z.Object,
					err error,
				) {
					if len(args) != 1 {
						return nil, z.ErrWrongNumArguments
					}
					y1, ok := args[0].(*z.Bytes)
					if !ok {
						return nil, z.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "bytes",
							Found:    args[0].TypeName(),
						}
					}
					res, err := r.Read(y1.Value)
					if err != nil {
						ret = wrapError(err)
						return
					}
					return &z.Int{Value: int64(res)}, nil
				},
			},
		},
	}
}
