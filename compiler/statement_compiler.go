package compiler

import (
	"fmt"
	"pun/ast"
	"pun/bytecode"
)

func (c *Compiler) compileStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		c.compileExpression(s.Expression)
	case *ast.AssignStatement:

		c.compileAssignStatement(s)
	//case *ast.CompoundAssignStatement:
	//	c.compileCompoundAssignStatement(s)

	case *ast.IfStatement:
		c.compileIfStatement(s)

	default:
		c.addError(fmt.Sprintf("Unsupported statement type: %T", stmt), 1, 0, "Hello")
	}
}

func (c *Compiler) compileAssignStatement(s *ast.AssignStatement) {
	// Luôn compile giá trị bên phải trước
	c.compileExpression(s.Value)

	// Xử lý target assignment
	switch target := s.Name.(type) {
	case *ast.Identifier:
		name := target.Value

		// Check tên biến hợp lệ (không trùng built-in)
		if !c.isValidVariableName(name, s.Line) {
			return // Đã có error trong isValidVariableName
		}

		// Global scope (không có thì tạo mới, có thì cho operand = slot của cái đang có)
		if len(c.Scopes) == 0 {
			idx, exists := c.GlobalSymbols[name]
			if !exists {
				idx = len(c.GlobalSymbols)
				c.GlobalSymbols[name] = idx
			}
			c.emit(bytecode.OP_STORE_GLOBAL, idx, s.Line)
			return
		}

		// Dùng resolveVariable để xử lí biến trong scope
		slot, depth, isGlobal, exists := c.resolveVariable(name)

		if isGlobal {
			c.emit(bytecode.OP_STORE_GLOBAL, slot, s.Line) // Global override
		} else if exists {
			initDepth := c.getInitDepth(name)
			c.emit(bytecode.OP_STORE_LOCAL, &bytecode.LocalVar{Slot: slot, Depth: initDepth}, s.Line) // Local reassign
		} else {
			// Tạo local mới nếu biến chưa tồn tại anywhere
			newSlot := len(c.CurrentScope)
			c.CurrentScope[name] = newSlot
			c.LocalInitDepth[name] = depth
			c.emit(bytecode.OP_STORE_LOCAL, &bytecode.LocalVar{Slot: newSlot, Depth: depth}, s.Line)
		}

	case *ast.ArrayIndexExpression:
		// Thêm check kiểu array trước khi gán
		c.compileExpression(target.Array)
		c.compileExpression(target.Index)
		c.emit(bytecode.OP_ARRAY_SET, nil, s.Line)

	default:
		c.addError(fmt.Sprintf("Unsupported assignment target: %T", target), s.Line, 0, "")
	}
}

//func (c *Compiler) compileCompoundAssignStatement(s *ast.CompoundAssignStatement) {
//	// 1. Kiểm tra biến tồn tại (với identifier)
//	if ident, ok := s.Name.(*ast.Identifier); ok {
//		if _, exists := c.SymbolTable[ident.Value]; !exists {
//			c.addError("Undefined variable in compound assignment", s.Line, 0,
//				fmt.Sprintf("Variable: %s", ident.Value))
//		}
//	}
//
//	// 2. Compile binary expression (đã bao gồm load giá trị cũ + phép toán)
//	c.compileExpression(s.Value)
//
//	// 3. Store lại kết quả
//	switch target := s.Name.(type) {
//	case *ast.Identifier:
//		slot := c.SymbolTable[target.Value]
//		c.emit(bytecode.OP_STORE_VAR, slot, s.Line)
//
//	case *ast.ArrayIndexExpression:
//		c.compileExpression(target.Array)
//		c.compileExpression(target.Index)
//		c.emit(bytecode.OP_ARRAY_SET, nil, s.Line)
//	}
//
//}

func (c *Compiler) compileIfStatement(s *ast.IfStatement) {
	// Tạo slice lưu tất cả jump positions
	var pendingJumps []int

	//Comlile điều kiện
	c.compileExpression(s.Condition)

	//Emit jump if false (operand tạm thời bằng 0)
	jumpIfFalsePos := len(c.Code)
	c.emit(bytecode.OP_JUMP_IF_FALSE, 0, s.Line)

	// Compile body if
	c.compileBlockStatement(s.Body)

	//Sửa offset cho jump if false
	jumpIfFalseOffset := len(c.Code) - jumpIfFalsePos
	c.Code[jumpIfFalsePos].Operand = jumpIfFalseOffset

	// Thêm jump để nhảy qua các elif/else
	endJumpPos := len(c.Code)
	c.emit(bytecode.OP_JUMP, 0, s.Line)
	pendingJumps = append(pendingJumps, endJumpPos)

	// Xử lý từng elif
	for _, elif := range s.ElseIfs {
		// Compile condition
		c.compileExpression(elif.Condition)

		// Thêm jump mới
		jumpIfFalsePos = len(c.Code)
		c.emit(bytecode.OP_JUMP_IF_FALSE, 0, elif.Line)

		// Compile body
		c.compileBlockStatement(elif.Body)

		//Sửa offset cho jump if false
		jumpIfFalseOffset = len(c.Code) - jumpIfFalsePos
		c.Code[jumpIfFalsePos].Operand = jumpIfFalseOffset

		// Thêm jump cuối elif
		endJumpPos = len(c.Code)
		c.emit(bytecode.OP_JUMP, 0, elif.Line)
		pendingJumps = append(pendingJumps, endJumpPos)
	}
	if s.ElseBlock != nil {
		// Compile else body
		c.compileBlockStatement(s.ElseBlock.Body)

	}

	// Điền tất cả jump cuối (OP_JUMP) trỏ ra ngoài khối if
	for _, pos := range pendingJumps {
		c.Code[pos].Operand = len(c.Code) - pos
	}
}

func (c *Compiler) compileBlockStatement(s *ast.BlockStatement) {
	//Vừa vào scope thì chưa biết có những biến nào nên tạm thời gán operand = 0
	c.enterScope()
	enterScopePos := len(c.Code)
	c.emit(bytecode.OP_ENTER_SCOPE, 0, s.Line)

	for _, stmt := range s.Statements {
		c.compileStatement(stmt)
	}

	//Gán lại operand cho lệnh enter scope là số lượng biến có trong scope
	c.Code[enterScopePos].Operand = len(c.CurrentScope)
	c.leaveScope()
	c.emit(bytecode.OP_LEAVE_SCOPE, nil, s.Line)
}
