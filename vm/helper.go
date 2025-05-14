package vm

import (
	"fmt"
	"pun/error"
	"strings"
)

// Thêm lỗi vào danh sách
func (v *VM) addError(message string, line, col int, context string) {
	err := customError.RuntimeError{
		PunError: customError.PunError{
			Message: message,
			Line:    line,
			Column:  col,
		},
		Context: context,
	}
	v.Errors = append(v.Errors, err)
}

// Kiểm tra có lỗi hay không
func (v *VM) HasErrors() bool {
	return len(v.Errors) > 0
}

// In tất cả lỗi
func (v *VM) PrintErrors() {
	if !v.HasErrors() {
		return
	}

	fmt.Println("🚨 RUNTIME ERRORS:")
	for i, err := range v.Errors {
		fmt.Printf("%d. %s\n", i+1, err.Error())
		fmt.Println(strings.Repeat("─", 60))
	}
}

// Helper methods
func (v *VM) push(val interface{}) {
	v.Sp++
	if v.Sp == len(v.Stack) {
		v.Stack = append(v.Stack, val)
	} else {
		v.Stack[v.Sp] = val
	}
}

func (v *VM) pop() interface{} {
	val := v.Stack[v.Sp]
	v.Sp--
	v.Stack = v.Stack[:len(v.Stack)-1]
	return val
}

func (v *VM) pushScope(localSize int) {
	scope := &Scope{Locals: make([]interface{}, localSize), Parent: v.CurrentScope}
	v.ScopeStack = append(v.ScopeStack, scope)
	v.CurrentScope = scope
}

func (v *VM) popScope() {
	if len(v.ScopeStack) > 1 { // Giữ lại global scope
		v.ScopeStack = v.ScopeStack[:len(v.ScopeStack)-1]
		v.CurrentScope = v.ScopeStack[len(v.ScopeStack)-1]
	}
}
