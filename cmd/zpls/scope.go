package main

// scopeStack 作用域栈，用于管理变量定义的作用域
type scopeStack struct {
	scopes [][]string // 每个作用域中的变量名
}

func newScopeStack() *scopeStack {
	return &scopeStack{
		scopes: make([][]string, 0),
	}
}

func (s *scopeStack) pushScope() {
	s.scopes = append(s.scopes, make([]string, 0))
}

func (s *scopeStack) popScope() {
	if len(s.scopes) > 0 {
		s.scopes = s.scopes[:len(s.scopes)-1]
	}
}

func (s *scopeStack) addVariable(name string) {
	if len(s.scopes) > 0 {
		s.scopes[len(s.scopes)-1] = append(s.scopes[len(s.scopes)-1], name)
	}
}

// isVariableInScope 检查变量是否在作用域中
func (s *scopeStack) isVariableInScope(name string) bool {
	// 从内到外检查作用域
	for i := len(s.scopes) - 1; i >= 0; i-- {
		scope := s.scopes[i]
		for _, variable := range scope {
			if variable == name {
				return true
			}
		}
	}
	return false
}