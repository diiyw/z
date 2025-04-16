package stdlib

import (
	"os"
	"syscall"

	"github.com/diiyw/z"
)

func makeOSProcessState(state *os.ProcessState) *z.ImmutableMap {
	return &z.ImmutableMap{
		Value: map[string]z.Object{
			"exited": &z.UserFunction{
				Name:  "exited",
				Value: FuncARB(state.Exited),
			},
			"pid": &z.UserFunction{
				Name:  "pid",
				Value: FuncARI(state.Pid),
			},
			"string": &z.UserFunction{
				Name:  "string",
				Value: FuncARS(state.String),
			},
			"success": &z.UserFunction{
				Name:  "success",
				Value: FuncARB(state.Success),
			},
		},
	}
}

func makeOSProcess(proc *os.Process) *z.ImmutableMap {
	return &z.ImmutableMap{
		Value: map[string]z.Object{
			"kill": &z.UserFunction{
				Name:  "kill",
				Value: FuncARE(proc.Kill),
			},
			"release": &z.UserFunction{
				Name:  "release",
				Value: FuncARE(proc.Release),
			},
			"signal": &z.UserFunction{
				Name: "signal",
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
					return wrapError(proc.Signal(syscall.Signal(i1))), nil
				},
			},
			"wait": &z.UserFunction{
				Name: "wait",
				Value: func(args ...z.Object) (z.Object, error) {
					if len(args) != 0 {
						return nil, z.ErrWrongNumArguments
					}
					state, err := proc.Wait()
					if err != nil {
						return wrapError(err), nil
					}
					return makeOSProcessState(state), nil
				},
			},
		},
	}
}
