package compiler

import (
	"fmt"
	"math"
	"pun/ast"
	"pun/bytecode"
	"pun/error"
	"strings"
)

type Compiler struct {
	Constants        []interface{}          // Pool hằng số
	Code             []bytecode.Instruction // Chương trình bytecode
	GlobalSymbols    map[string]int         // Chỉ cho biến global
	CurrentScope     map[string]int         //Scope hiện tại
	Scopes           []map[string]int       // Chỉ cho local scopes (không chứa global)
	LocalInitDepth   map[string]int         //Lưu depth của scope mà biến local lần đầu được tạo (dùng cho nested scope)
	Labels           map[string]int         //Lưu các label
	PendingJumps     map[string][]int       // Lưu các vị trí jump ứng với lỗi label
	BuiltinFuncs     map[string]bool        //Lưu tên các hàm built-in
	BuiltinConstants map[string]int         //Lưu tên hằng số và index trong constants pool
	IsInsideFunction bool                   //Kiểm tra xem có đang trong hàm không (quản lí return)
	Errors           []customError.CompilationError
}

// Dùng để lưu biến local cùng depth của scope chứa nó (giúp vm xác định đúng)

func NewCompiler() *Compiler {
	c := &Compiler{
		BuiltinFuncs:     make(map[string]bool),
		BuiltinConstants: make(map[string]int),
		GlobalSymbols:    make(map[string]int),
		Scopes:           make([]map[string]int, 0), // Bắt đầu với empty stack
		LocalInitDepth:   make(map[string]int),
		Labels:           make(map[string]int),
		PendingJumps:     make(map[string][]int),
		IsInsideFunction: false,
	}
	//Thêm hàm builtin
	c.registerBuiltinFunc("print")
	c.registerBuiltinFunc("ask")

	//Thêm hằng số
	c.registerBuiltinConstant("PI", math.Pi)
	c.registerBuiltinConstant("E", math.E)

	return c

}

func (c *Compiler) registerBuiltinConstant(name string, value interface{}) {
	index := c.addConstant(value)
	c.BuiltinConstants[name] = index
}

func (c *Compiler) registerBuiltinFunc(name string) {
	c.BuiltinFuncs[name] = true
}

func (c *Compiler) CompileProgram(program *ast.Program) {
	for _, stmt := range program.Statements {
		c.compileStatement(stmt)
	}
}

// Thêm hằng số vào pool
func (c *Compiler) addConstant(value interface{}) int {
	//Nếu hằng số đã có sẵn trong constant pool thì khỏi cần thêm nữa (tối ưu bộ nhớ)
	for i, val := range c.Constants {
		if val == value {
			return i
		}
	}
	c.Constants = append(c.Constants, value)
	return len(c.Constants) - 1
}

// Tạo instruction mới
func (c *Compiler) emit(op string, operand interface{}, line int) {
	c.Code = append(c.Code, bytecode.Instruction{
		Op:      op,
		Operand: operand,
		Line:    line,
	})
}

// Định nghĩa label tại vị trí hiện tại
func (c *Compiler) defineLabel(name string) {
	c.Labels[name] = len(c.Code) // Lưu PC hiện tại
}

// Emit jump đến label (chưa biết PC)
func (c *Compiler) emitJumpToLabel(op string, label string, line int) {
	c.emit(op, label, line) // Operand là tên label (tạm thời)
	c.PendingJumps[label] = append(c.PendingJumps[label], len(c.Code)-1)
}

// Resolve tất cả jumps sau khi biết vị trí label (chuyển label thành offset trong instruction jump)
func (c *Compiler) resolveJumps() {
	for label, jumps := range c.PendingJumps {
		targetPC, ok := c.Labels[label]
		if !ok {
			c.addError(fmt.Sprintf("undefined label: %s", label), 0, 0, "jump resolution")
			continue
		}
		for _, pc := range jumps { //pc là vị trí của lệnh jump
			offset := targetPC - pc - 1 // Tính offset
			c.Code[pc].Operand = offset // Sửa operand
		}
	}
}

func (c *Compiler) resetLabels() {
	c.Labels = make(map[string]int)
	c.PendingJumps = make(map[string][]int)
}

func (c *Compiler) enterScope() {
	newScope := make(map[string]int)
	c.Scopes = append(c.Scopes, newScope)
	c.CurrentScope = newScope
}

func (c *Compiler) leaveScope() {
	if len(c.Scopes) > 0 {
		c.Scopes = c.Scopes[:len(c.Scopes)-1]
		if len(c.Scopes) > 0 {
			c.CurrentScope = c.Scopes[len(c.Scopes)-1]
		} else {
			c.CurrentScope = nil
		}
	}

}

func (c *Compiler) resolveVariable(name string) (slot int, depth int, isGlobal bool, exists bool) {
	// 1. Tìm trong các scope local (từ trong ra ngoài)
	//Depth bắt đầu từ 1 (0 là depth của global)
	for i := len(c.Scopes) - 1; i >= 0; i-- {
		if idx, ok := c.Scopes[i][name]; ok {
			// Slot + Depth tính từ current scope
			return idx, len(c.Scopes) - i, false, true
		}
	}

	// 2. Tìm trong global
	if idx, ok := c.GlobalSymbols[name]; ok {
		return idx, 0, true, true
	}

	//Không tìm thấy ở cả local và global => Biến chưa được tạo ở scope hiện tại
	// => Depth của scope hiện tại = độ dài mảng scope (do depth bắt đầu từ 1)
	return 0, len(c.Scopes), false, false
}

func (c *Compiler) addError(message string, line, col int, context string) {
	err := customError.CompilationError{
		PunError: customError.PunError{
			Message: message,
			Line:    line,
			Column:  col,
		},
		Context: context,
	}
	c.Errors = append(c.Errors, err)
}
func (c *Compiler) isValidVariableName(name string, line int) bool {
	if c.BuiltinFuncs[name] || c.BuiltinConstants[name] != 0 {
		c.addError("Cannot redeclare built-in name", line, 0, name)
		return false
	}
	return true
}
func (c *Compiler) HasErrors() bool {
	return len(c.Errors) > 0
}

func (c *Compiler) PrintErrors() {
	if !c.HasErrors() {
		return
	}

	fmt.Println("🚨 COMPILATION ERRORS:")
	for i, err := range c.Errors {
		fmt.Printf("%d. %s\n", i+1, err.Error())
		fmt.Println(strings.Repeat("─", 60))
	}
}

func (c *Compiler) getInitDepth(name string) int {
	return c.LocalInitDepth[name]
}
