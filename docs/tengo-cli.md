# Z CLI Tool

Z is designed as an embedding script language for Go, but, it can also be
compiled and executed as native binary using `z` CLI tool.

## Installing Z CLI

To install `z` tool, run:

```bash
go get github.com/d5/z/cmd/z
```

Or, you can download the precompiled binaries from
[here](https://github.com/d5/z/releases/latest).

## Compiling and Executing Z Code

You can directly execute the Z source code by running `z` tool with
your Z source file (`*.z`).

```bash
z myapp.z
```

Or, you can compile the code into a binary file and execute it later.

```bash
z -o myapp myapp.z   # compile 'myapp.z' into binary file 'myapp'
z myapp                  # execute the compiled binary `myapp`
```

Or, you can make z source file executable

```bash
# copy z executable to a dir where PATH environment variable includes
cp z /usr/local/bin/

# add shebang line to source file
cat > myapp.z << EOF
#!/usr/local/bin/z
fmt := import("fmt")
fmt.println("Hello World!")
EOF

# make myapp.z file executable
chmod +x myapp.z

# run your script
./myapp.z
```

**Note: Your source file must have `.z` extension.**

## Resolving Relative Import Paths

If there are z source module files which are imported with relative import
paths, CLI has `-resolve` flag. Flag enables to import a module relative to
importing file. This behavior will be default at version 3.

## Z REPL

You can run Z [REPL](https://en.wikipedia.org/wiki/Read–eval–print_loop)
if you run `z` with no arguments.

```bash
z
```
