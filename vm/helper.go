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
	return val
}
