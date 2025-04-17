package stdlib

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/diiyw/z"
)

var osModule = map[string]z.Object{
	"platform":            &z.String{Value: runtime.GOOS},
	"arch":                &z.String{Value: runtime.GOARCH},
	"o_rdonly":            &z.Int{Value: int64(os.O_RDONLY)},
	"o_wronly":            &z.Int{Value: int64(os.O_WRONLY)},
	"o_rdwr":              &z.Int{Value: int64(os.O_RDWR)},
	"o_append":            &z.Int{Value: int64(os.O_APPEND)},
	"o_create":            &z.Int{Value: int64(os.O_CREATE)},
	"o_excl":              &z.Int{Value: int64(os.O_EXCL)},
	"o_sync":              &z.Int{Value: int64(os.O_SYNC)},
	"o_trunc":             &z.Int{Value: int64(os.O_TRUNC)},
	"mode_dir":            &z.Int{Value: int64(os.ModeDir)},
	"mode_append":         &z.Int{Value: int64(os.ModeAppend)},
	"mode_exclusive":      &z.Int{Value: int64(os.ModeExclusive)},
	"mode_temporary":      &z.Int{Value: int64(os.ModeTemporary)},
	"mode_symlink":        &z.Int{Value: int64(os.ModeSymlink)},
	"mode_device":         &z.Int{Value: int64(os.ModeDevice)},
	"mode_named_pipe":     &z.Int{Value: int64(os.ModeNamedPipe)},
	"mode_socket":         &z.Int{Value: int64(os.ModeSocket)},
	"mode_setuid":         &z.Int{Value: int64(os.ModeSetuid)},
	"mode_setgui":         &z.Int{Value: int64(os.ModeSetgid)},
	"mode_char_device":    &z.Int{Value: int64(os.ModeCharDevice)},
	"mode_sticky":         &z.Int{Value: int64(os.ModeSticky)},
	"mode_type":           &z.Int{Value: int64(os.ModeType)},
	"mode_perm":           &z.Int{Value: int64(os.ModePerm)},
	"path_separator":      &z.Char{Value: os.PathSeparator},
	"path_list_separator": &z.Char{Value: os.PathListSeparator},
	"dev_null":            &z.String{Value: os.DevNull},
	"seek_set":            &z.Int{Value: int64(io.SeekStart)},
	"seek_cur":            &z.Int{Value: int64(io.SeekCurrent)},
	"seek_end":            &z.Int{Value: int64(io.SeekEnd)},
	"args": &z.UserFunction{
		Name:  "args",
		Value: osArgs,
	}, // args() => array(string)
	"chdir": &z.UserFunction{
		Name:  "chdir",
		Value: FuncASRE(os.Chdir),
	}, // chdir(dir string) => error
	"chmod": osFuncASFmRE("chmod", os.Chmod), // chmod(name string, mode int) => error
	"chown": &z.UserFunction{
		Name:  "chown",
		Value: FuncASIIRE(os.Chown),
	}, // chown(name string, uid int, gid int) => error
	"clearenv": &z.UserFunction{
		Name:  "clearenv",
		Value: FuncAR(os.Clearenv),
	}, // clearenv()
	"environ": &z.UserFunction{
		Name:  "environ",
		Value: FuncARSs(os.Environ),
	}, // environ() => array(string)
	"exit": &z.UserFunction{
		Name:  "exit",
		Value: FuncAIR(os.Exit),
	}, // exit(code int)
	"expand_env": &z.UserFunction{
		Name:  "expand_env",
		Value: osExpandEnv,
	}, // expand_env(s string) => string
	"getegid": &z.UserFunction{
		Name:  "getegid",
		Value: FuncARI(os.Getegid),
	}, // getegid() => int
	"getenv": &z.UserFunction{
		Name:  "getenv",
		Value: FuncASRS(os.Getenv),
	}, // getenv(s string) => string
	"geteuid": &z.UserFunction{
		Name:  "geteuid",
		Value: FuncARI(os.Geteuid),
	}, // geteuid() => int
	"getgid": &z.UserFunction{
		Name:  "getgid",
		Value: FuncARI(os.Getgid),
	}, // getgid() => int
	"getgroups": &z.UserFunction{
		Name:  "getgroups",
		Value: FuncARIsE(os.Getgroups),
	}, // getgroups() => array(string)/error
	"getpagesize": &z.UserFunction{
		Name:  "getpagesize",
		Value: FuncARI(os.Getpagesize),
	}, // getpagesize() => int
	"getpid": &z.UserFunction{
		Name:  "getpid",
		Value: FuncARI(os.Getpid),
	}, // getpid() => int
	"getppid": &z.UserFunction{
		Name:  "getppid",
		Value: FuncARI(os.Getppid),
	}, // getppid() => int
	"getuid": &z.UserFunction{
		Name:  "getuid",
		Value: FuncARI(os.Getuid),
	}, // getuid() => int
	"getwd": &z.UserFunction{
		Name:  "getwd",
		Value: FuncARSE(os.Getwd),
	}, // getwd() => string/error
	"hostname": &z.UserFunction{
		Name:  "hostname",
		Value: FuncARSE(os.Hostname),
	}, // hostname() => string/error
	"lchown": &z.UserFunction{
		Name:  "lchown",
		Value: FuncASIIRE(os.Lchown),
	}, // lchown(name string, uid int, gid int) => error
	"link": &z.UserFunction{
		Name:  "link",
		Value: FuncASSRE(os.Link),
	}, // link(oldname string, newname string) => error
	"lookup_env": &z.UserFunction{
		Name:  "lookup_env",
		Value: osLookupEnv,
	}, // lookup_env(key string) => string/false
	"mkdir":     osFuncASFmRE("mkdir", os.Mkdir),        // mkdir(name string, perm int) => error
	"mkdir_all": osFuncASFmRE("mkdir_all", os.MkdirAll), // mkdir_all(name string, perm int) => error
	"readlink": &z.UserFunction{
		Name:  "readlink",
		Value: FuncASRSE(os.Readlink),
	}, // readlink(name string) => string/error
	"remove": &z.UserFunction{
		Name:  "remove",
		Value: FuncASRE(os.Remove),
	}, // remove(name string) => error
	"remove_all": &z.UserFunction{
		Name:  "remove_all",
		Value: FuncASRE(os.RemoveAll),
	}, // remove_all(name string) => error
	"rename": &z.UserFunction{
		Name:  "rename",
		Value: FuncASSRE(os.Rename),
	}, // rename(oldpath string, newpath string) => error
	"setenv": &z.UserFunction{
		Name:  "setenv",
		Value: FuncASSRE(os.Setenv),
	}, // setenv(key string, value string) => error
	"symlink": &z.UserFunction{
		Name:  "symlink",
		Value: FuncASSRE(os.Symlink),
	}, // symlink(oldname string newname string) => error
	"temp_dir": &z.UserFunction{
		Name:  "temp_dir",
		Value: FuncARS(os.TempDir),
	}, // temp_dir() => string
	"truncate": &z.UserFunction{
		Name:  "truncate",
		Value: FuncASI64RE(os.Truncate),
	}, // truncate(name string, size int) => error
	"unsetenv": &z.UserFunction{
		Name:  "unsetenv",
		Value: FuncASRE(os.Unsetenv),
	}, // unsetenv(key string) => error
	"create": &z.UserFunction{
		Name:  "create",
		Value: osCreate,
	}, // create(name string) => imap(file)/error
	"open": &z.UserFunction{
		Name:  "open",
		Value: osOpen,
	}, // open(name string) => imap(file)/error
	"open_file": &z.UserFunction{
		Name:  "open_file",
		Value: osOpenFile,
	}, // open_file(name string, flag int, perm int) => imap(file)/error
	"find_process": &z.UserFunction{
		Name:  "find_process",
		Value: osFindProcess,
	}, // find_process(pid int) => imap(process)/error
	"start_process": &z.UserFunction{
		Name:  "start_process",
		Value: osStartProcess,
	}, // start_process(name string, argv array(string), dir string, env array(string)) => imap(process)/error
	"exec_look_path": &z.UserFunction{
		Name:  "exec_look_path",
		Value: FuncASRSE(exec.LookPath),
	}, // exec_look_path(file) => string/error
	"exec": &z.UserFunction{
		Name:  "exec",
		Value: osExec,
	}, // exec(name, args...) => command
	"stat": &z.UserFunction{
		Name:  "stat",
		Value: osStat,
	}, // stat(name) => imap(fileinfo)/error
	"read_file": &z.UserFunction{
		Name:  "read_file",
		Value: osReadFile,
	}, // readfile(name) => array(byte)/error
}

func osReadFile(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		return nil, z.ErrWrongNumArguments
	}
	fname, ok := z.ToString(args[0])
	if !ok {
		return nil, z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}
	bytes, err := os.ReadFile(fname)
	if err != nil {
		return wrapError(err), nil
	}
	if len(bytes) > z.MaxBytesLen {
		return nil, z.ErrBytesLimit
	}
	return &z.Bytes{Value: bytes}, nil
}

func osStat(args ...z.Object) (ret z.Object, err error) {
	if len(args) != 1 {
		return nil, z.ErrWrongNumArguments
	}
	fname, ok := z.ToString(args[0])
	if !ok {
		return nil, z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}
	stat, err := os.Stat(fname)
	if err != nil {
		return wrapError(err), nil
	}
	fstat := &z.ImmutableMap{
		Value: map[string]z.Object{
			"name":  &z.String{Value: stat.Name()},
			"mtime": &z.Time{Value: stat.ModTime()},
			"size":  &z.Int{Value: stat.Size()},
			"mode":  &z.Int{Value: int64(stat.Mode())},
		},
	}
	if stat.IsDir() {
		fstat.Value["directory"] = z.TrueValue
	} else {
		fstat.Value["directory"] = z.FalseValue
	}
	return fstat, nil
}

func osCreate(args ...z.Object) (z.Object, error) {
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
	res, err := os.Create(s1)
	if err != nil {
		return wrapError(err), nil
	}
	return makeOSFile(res), nil
}

func osOpen(args ...z.Object) (z.Object, error) {
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
	res, err := os.Open(s1)
	if err != nil {
		return wrapError(err), nil
	}
	return makeOSFile(res), nil
}

func osOpenFile(args ...z.Object) (z.Object, error) {
	if len(args) != 3 {
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
	i2, ok := z.ToInt(args[1])
	if !ok {
		return nil, z.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
	}
	i3, ok := z.ToInt(args[2])
	if !ok {
		return nil, z.ErrInvalidArgumentType{
			Name:     "third",
			Expected: "int(compatible)",
			Found:    args[2].TypeName(),
		}
	}
	res, err := os.OpenFile(s1, i2, os.FileMode(i3))
	if err != nil {
		return wrapError(err), nil
	}
	return makeOSFile(res), nil
}

func osArgs(args ...z.Object) (z.Object, error) {
	if len(args) != 0 {
		return nil, z.ErrWrongNumArguments
	}
	arr := &z.Array{}
	for _, osArg := range os.Args {
		if len(osArg) > z.MaxStringLen {
			return nil, z.ErrStringLimit
		}
		arr.Value = append(arr.Value, &z.String{Value: osArg})
	}
	return arr, nil
}

func osFuncASFmRE(
	name string,
	fn func(string, os.FileMode) error,
) *z.UserFunction {
	return &z.UserFunction{
		Name: name,
		Value: func(args ...z.Object) (z.Object, error) {
			if len(args) != 2 {
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
			i2, ok := z.ToInt64(args[1])
			if !ok {
				return nil, z.ErrInvalidArgumentType{
					Name:     "second",
					Expected: "int(compatible)",
					Found:    args[1].TypeName(),
				}
			}
			return wrapError(fn(s1, os.FileMode(i2))), nil
		},
	}
}

func osLookupEnv(args ...z.Object) (z.Object, error) {
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
	res, ok := os.LookupEnv(s1)
	if !ok {
		return z.FalseValue, nil
	}
	if len(res) > z.MaxStringLen {
		return nil, z.ErrStringLimit
	}
	return &z.String{Value: res}, nil
}

func osExpandEnv(args ...z.Object) (z.Object, error) {
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
	var vlen int
	var failed bool
	s := os.Expand(s1, func(k string) string {
		if failed {
			return ""
		}
		v := os.Getenv(k)

		// this does not count the other texts that are not being replaced
		// but the code checks the final length at the end
		vlen += len(v)
		if vlen > z.MaxStringLen {
			failed = true
			return ""
		}
		return v
	})
	if failed || len(s) > z.MaxStringLen {
		return nil, z.ErrStringLimit
	}
	return &z.String{Value: s}, nil
}

func osExec(args ...z.Object) (z.Object, error) {
	if len(args) == 0 {
		return nil, z.ErrWrongNumArguments
	}
	name, ok := z.ToString(args[0])
	if !ok {
		return nil, z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}
	var execArgs []string
	for idx, arg := range args[1:] {
		execArg, ok := z.ToString(arg)
		if !ok {
			return nil, z.ErrInvalidArgumentType{
				Name:     fmt.Sprintf("args[%d]", idx),
				Expected: "string(compatible)",
				Found:    args[1+idx].TypeName(),
			}
		}
		execArgs = append(execArgs, execArg)
	}
	return makeOSExecCommand(exec.Command(name, execArgs...)), nil
}

func osFindProcess(args ...z.Object) (z.Object, error) {
	if len(args) != 1 {
		return nil, z.ErrWrongNumArguments
	}
	i1, ok := z.ToInt(args[0])
	if !ok {
		return nil, z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
	}
	proc, err := os.FindProcess(i1)
	if err != nil {
		return wrapError(err), nil
	}
	return makeOSProcess(proc), nil
}

func osStartProcess(args ...z.Object) (z.Object, error) {
	if len(args) != 4 {
		return nil, z.ErrWrongNumArguments
	}
	name, ok := z.ToString(args[0])
	if !ok {
		return nil, z.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}
	var argv []string
	var err error
	switch arg1 := args[1].(type) {
	case *z.Array:
		argv, err = stringArray(arg1.Value, "second")
		if err != nil {
			return nil, err
		}
	case *z.ImmutableArray:
		argv, err = stringArray(arg1.Value, "second")
		if err != nil {
			return nil, err
		}
	default:
		return nil, z.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "array",
			Found:    arg1.TypeName(),
		}
	}

	dir, ok := z.ToString(args[2])
	if !ok {
		return nil, z.ErrInvalidArgumentType{
			Name:     "third",
			Expected: "string(compatible)",
			Found:    args[2].TypeName(),
		}
	}

	var env []string
	switch arg3 := args[3].(type) {
	case *z.Array:
		env, err = stringArray(arg3.Value, "fourth")
		if err != nil {
			return nil, err
		}
	case *z.ImmutableArray:
		env, err = stringArray(arg3.Value, "fourth")
		if err != nil {
			return nil, err
		}
	default:
		return nil, z.ErrInvalidArgumentType{
			Name:     "fourth",
			Expected: "array",
			Found:    arg3.TypeName(),
		}
	}

	proc, err := os.StartProcess(name, argv, &os.ProcAttr{
		Dir: dir,
		Env: env,
	})
	if err != nil {
		return wrapError(err), nil
	}
	return makeOSProcess(proc), nil
}

func stringArray(arr []z.Object, argName string) ([]string, error) {
	var sarr []string
	for idx, elem := range arr {
		str, ok := elem.(*z.String)
		if !ok {
			return nil, z.ErrInvalidArgumentType{
				Name:     fmt.Sprintf("%s[%d]", argName, idx),
				Expected: "string",
				Found:    elem.TypeName(),
			}
		}
		sarr = append(sarr, str.Value)
	}
	return sarr, nil
}
