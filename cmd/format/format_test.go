package format

import "testing"

func TestFormat(t *testing.T) {
	var tests = []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple assignment",
			input:    "x= 1 + 2",
			expected: "x = 1 + 2",
		},
		{
			name:     "simple assignment with spaces",
			input:    "print(1,2,3)",
			expected: "print(1, 2, 3)",
		},
		{
			name:     "multiple assignments",
			input:    "x=2*1+2",
			expected: "x = 2 * 1 + 2",
		},
		{
			name:     "multiple assignments",
			input:    "x=2*(1+2)",
			expected: "x = 2 * (1 + 2)",
		},
		{
			name:     "assignment with spaces",
			input:    "fmt:=import(\"fmt\")",
			expected: "fmt := import(\"fmt\")",
		},
		{
			name:     "assignment with spaces and export",
			input:    "fmt:=import(\"fmt\")\nexport{}",
			expected: "fmt := import(\"fmt\")\nexport {}",
		},
		{
			name: "assignment with spaces and export with function",
			input: `fmt:=import("fmt")
export{fn:func(){}}`,
			expected: `fmt := import("fmt")
export {
	fn: func() {
	}
}`,
		},
		{
			name: "assignment with spaces and export with function and value",
			input: `fmt:=import("fmt")
export{a:1,fn:func(){}}`,
			expected: `fmt := import("fmt")
export {
	a: 1,
	fn: func() {
	}
}`,
		},
		{
			name: "assignment with spaces and export with function and value and print",
			input: `fmt:=import("fmt")
export{a:1,fn:func(){print(1)}}`,
			expected: `fmt := import("fmt")
export {
	a: 1,
	fn: func() {
		print(1)
	}
}`,
		},
		{
			name:  "for statement",
			input: `for {}`,
			expected: `for {
}`,
		},
		{
			name:  "for statement with cond",
			input: `for true{}`,
			expected: `for true {
}`,
		},
		{
			name:  "for statement full",
			input: `for i:=0;i<=10;i++{print(i)}`,
			expected: `for i := 0; i <= 10; i++ {
	print(i)
}`,
		},
		{
			name:  "if statement with cond",
			input: `if x>1{}`,
			expected: `if x > 1 {
}`,
		},
		{
			name:  "if statement with init and cond",
			input: `if x:=1;x>1{}`,
			expected: `if x := 1; x > 1 {
}`,
		},
		{
			name:  "for-in statement",
			input: `for v in vv{}`,
			expected: `for _, v in vv {
}`,
		},
		{
			name:  "for-in statement with key-value",
			input: `for k,v in vv{}`,
			expected: `for k, v in vv {
}`,
		},
		{
			name:  "multi level 1",
			input: `for{for{for true{ m:={a:1}}}}`,
			expected: `for {
	for {
		for true {
			m := {
				a: 1
			}
		}
	}
}`,
		},
		{
			name: "multi level 2",
			input: `for{for{fmt:=import("fmt")
			if true{print(1)}}}`,
			expected: `for {
	for {
		fmt := import("fmt")
		if true {
			print(1)
		}
	}
}`,
		},
		{
			name:     "array",
			input:    `[x>1,1,2]`,
			expected: `[x > 1, 1, 2]`,
		},
		{
			name:  "map",
			input: `{a:1,b:2,c:3}`,
			expected: `{
	a: 1,
	b: 2,
	c: 3
}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Format([]byte(test.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != test.expected {
				t.Errorf("expected:\n%s\n-----\ngot:\n%s", test.expected, result)
			}
		})
	}
}
