package stdlib

import (
	"os/exec"

	"github.com/d5/tengo/v2"
)

func makeOSExecCommand(cmd *exec.Cmd) *z.ImmutableMap {
	return &z.ImmutableMap{
		Value: map[string]z.Object{
			// combined_output() => bytes/error
			"combined_output": &z.UserFunction{
				Name:  "combined_output",
				Value: FuncARYE(cmd.CombinedOutput),
			},
			// output() => bytes/error
			"output": &z.UserFunction{
				Name:  "output",
				Value: FuncARYE(cmd.Output),
			}, //
			// run() => error
			"run": &z.UserFunction{
				Name:  "run",
				Value: FuncARE(cmd.Run),
			}, //
			// start() => error
			"start": &z.UserFunction{
				Name:  "start",
				Value: FuncARE(cmd.Start),
			}, //
			// wait() => error
			"wait": &z.UserFunction{
				Name:  "wait",
				Value: FuncARE(cmd.Wait),
			}, //
			// set_path(path string)
			"set_path": &z.UserFunction{
				Name: "set_path",
				Value: func(args ...z.Object) (z.Object, error) {
					if len(args) != 1 {
						return nil, z.ErrWrongNumArguments
					}
					s1, ok := z.ToString(args[0])
					if !ok {
						return nil, z.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "string(compatible)",
							Found:    args[0].TypeName(),
						}
					}
					cmd.Path = s1
					return z.UndefinedValue, nil
				},
			},
			// set_dir(dir string)
			"set_dir": &z.UserFunction{
				Name: "set_dir",
				Value: func(args ...z.Object) (z.Object, error) {
					if len(args) != 1 {
						return nil, z.ErrWrongNumArguments
					}
					s1, ok := z.ToString(args[0])
					if !ok {
						return nil, z.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "string(compatible)",
							Found:    args[0].TypeName(),
						}
					}
					cmd.Dir = s1
					return z.UndefinedValue, nil
				},
			},
			// set_env(env array(string))
			"set_env": &z.UserFunction{
				Name: "set_env",
				Value: func(args ...z.Object) (z.Object, error) {
					if len(args) != 1 {
						return nil, z.ErrWrongNumArguments
					}

					var env []string
					var err error
					switch arg0 := args[0].(type) {
					case *z.Array:
						env, err = stringArray(arg0.Value, "first")
						if err != nil {
							return nil, err
						}
					case *z.ImmutableArray:
						env, err = stringArray(arg0.Value, "first")
						if err != nil {
							return nil, err
						}
					default:
						return nil, z.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "array",
							Found:    arg0.TypeName(),
						}
					}
					cmd.Env = env
					return z.UndefinedValue, nil
				},
			},
			// process() => imap(process)
			"process": &z.UserFunction{
				Name: "process",
				Value: func(args ...z.Object) (z.Object, error) {
					if len(args) != 0 {
						return nil, z.ErrWrongNumArguments
					}
					return makeOSProcess(cmd.Process), nil
				},
			},
		},
	}
}
