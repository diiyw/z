package stdlib

import (
	"os"

	"github.com/d5/tengo/v2"
)

func makeOSFile(file *os.File) *z.ImmutableMap {
	return &z.ImmutableMap{
		Value: map[string]z.Object{
			// chdir() => true/error
			"chdir": &z.UserFunction{
				Name:  "chdir",
				Value: FuncARE(file.Chdir),
			}, //
			// chown(uid int, gid int) => true/error
			"chown": &z.UserFunction{
				Name:  "chown",
				Value: FuncAIIRE(file.Chown),
			}, //
			// close() => error
			"close": &z.UserFunction{
				Name:  "close",
				Value: FuncARE(file.Close),
			}, //
			// name() => string
			"name": &z.UserFunction{
				Name:  "name",
				Value: FuncARS(file.Name),
			}, //
			// readdirnames(n int) => array(string)/error
			"readdirnames": &z.UserFunction{
				Name:  "readdirnames",
				Value: FuncAIRSsE(file.Readdirnames),
			}, //
			// sync() => error
			"sync": &z.UserFunction{
				Name:  "sync",
				Value: FuncARE(file.Sync),
			}, //
			// write(bytes) => int/error
			"write": &z.UserFunction{
				Name:  "write",
				Value: FuncAYRIE(file.Write),
			}, //
			// write(string) => int/error
			"write_string": &z.UserFunction{
				Name:  "write_string",
				Value: FuncASRIE(file.WriteString),
			}, //
			// read(bytes) => int/error
			"read": &z.UserFunction{
				Name:  "read",
				Value: FuncAYRIE(file.Read),
			}, //
			// chmod(mode int) => error
			"chmod": &z.UserFunction{
				Name: "chmod",
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
					return wrapError(file.Chmod(os.FileMode(i1))), nil
				},
			},
			// seek(offset int, whence int) => int/error
			"seek": &z.UserFunction{
				Name: "seek",
				Value: func(args ...z.Object) (z.Object, error) {
					if len(args) != 2 {
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
					i2, ok := z.ToInt(args[1])
					if !ok {
						return nil, z.ErrInvalidArgumentType{
							Name:     "second",
							Expected: "int(compatible)",
							Found:    args[1].TypeName(),
						}
					}
					res, err := file.Seek(i1, i2)
					if err != nil {
						return wrapError(err), nil
					}
					return &z.Int{Value: res}, nil
				},
			},
			// stat() => imap(fileinfo)/error
			"stat": &z.UserFunction{
				Name: "stat",
				Value: func(args ...z.Object) (z.Object, error) {
					if len(args) != 0 {
						return nil, z.ErrWrongNumArguments
					}
					return osStat(&z.String{Value: file.Name()})
				},
			},
		},
	}
}
