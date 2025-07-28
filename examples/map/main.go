package main

import (
	"context"

	"github.com/diiyw/z"
	"github.com/diiyw/z/stdlib"
)

func main() {
	// Z script code
	src := `fmt:= import("fmt");
m := {"*":1,"()":2,a:3}
fmt.println(m["*"],m["()"])`

	// create a new Script instance
	script := z.NewScript([]byte(src))
	moduleMap := z.NewModuleMap()
	for name, m := range stdlib.BuiltinModules {
		moduleMap.AddBuiltinModule(name, m)
	}
	script.SetImports(moduleMap)
	// run the script
	_, err := script.RunContext(context.Background())
	if err != nil {
		panic(err)
	}
	// Output:
	// 1 2
}
