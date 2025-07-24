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
			expected: "fmt := import(\"fmt\")\nexport {\n}",
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
