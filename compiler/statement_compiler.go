package compiler

import (
	"fmt"
	"pun/ast"
	"pun/bytecode"
	"strconv"
)

func (c *Compiler) compileStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		c.compileExpression(s.Expression)
	case *ast.AssignStatement:

		c.compileAssignStatement(s)
	case *ast.CompoundAssignStatement:
		c.compileCompoundAssignStatement(s)
	case *ast.IfStatement:
		c.compileIfStatement(s)
	case *ast.ForStatement:
		c.compileForStatement(s)
	case *ast.WhileStatement:
		c.compileWhileStatement(s)
	case *ast.FunctionDefinitionStatement:
		c.compileFunctionDefinitionStatement(s)
	case *ast.ReturnStatement:
		c.compileReturnStatement(s)
	default:
		c.addError(fmt.Sprintf("Unsupported statement type: %T", stmt), 1, 0, "compile statement")
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

func (c *Compiler) compileCompoundAssignStatement(s *ast.CompoundAssignStatement) {
	// Luôn compile giá trị bên phải trước
	c.compileExpression(s.Value)

	// Xử lý target assignment
	switch target := s.Name.(type) {
	case *ast.Identifier:
		name := target.Value
		// Global scope (không có sẵn thì báo lỗi)
		if len(c.Scopes) == 0 {
			idx, exists := c.GlobalSymbols[name]
			if !exists {
				c.addError("Undefined variable in compound assign statement", s.Line, 0, "compound assignment")
			}
			c.emit(bytecode.OP_STORE_GLOBAL, idx, s.Line)
			return
		}

		// Dùng resolveVariable để xử lí biến trong scope
		slot, _, isGlobal, exists := c.resolveVariable(name)

		if isGlobal {
			c.emit(bytecode.OP_STORE_GLOBAL, slot, s.Line) // Global override
		} else if exists {
			initDepth := c.getInitDepth(name)
			c.emit(bytecode.OP_STORE_LOCAL, &bytecode.LocalVar{Slot: slot, Depth: initDepth}, s.Line) // Local reassign
		} else {
			c.addError("Undefined variable in compound assign statement", s.Line, 0, "compound assignment")
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

func (c *Compiler) compileBlockStatement(s *ast.BlockStatement) {

	for _, stmt := range s.Statements {
		c.compileStatement(stmt)
	}

}

func (c *Compiler) compileIfBlockStatement(s *ast.BlockStatement) {
	//Vừa vào scope thì chưa biết có những biến nào nên tạm thời gán operand = 0
	c.enterScope()
	enterScopePos := len(c.Code)
	c.emit(bytecode.OP_ENTER_SCOPE, 0, s.Line)

	c.compileBlockStatement(s)

	//Gán lại operand cho lệnh enter scope là số lượng biến có trong scope
	c.Code[enterScopePos].Operand = len(c.CurrentScope)
	c.leaveScope()
	c.emit(bytecode.OP_LEAVE_SCOPE, nil, s.Line)
}

func (c *Compiler) compileIfStatement(s *ast.IfStatement) {
	// Compile điều kiện if
	c.compileExpression(s.Condition)

	// Jump nếu false -> else/end
	c.emitJumpToLabel(bytecode.OP_JUMP_IF_FALSE, "if_else", s.Line)

	// Compile if body (vào scope ở đây)
	c.compileIfBlockStatement(s.Body)

	// Jump qua phần else/elif (nếu có)
	c.emitJumpToLabel(bytecode.OP_JUMP, "if_end", s.Line)

	// Định nghĩa label "if_else" (bắt đầu elif/else)
	c.defineLabel("if_else")

	// Xử lý từng elif
	for idx, elif := range s.ElseIfs {
		elifEndLabel := "elif_end_" + strconv.Itoa(idx) //đánh dấu thứ tự elif để tránh trùng
		c.compileExpression(elif.Condition)
		c.emitJumpToLabel(bytecode.OP_JUMP_IF_FALSE, elifEndLabel, elif.Line)

		c.compileIfBlockStatement(elif.Body)

		c.emitJumpToLabel(bytecode.OP_JUMP, "if_end", elif.Line)

		c.defineLabel(elifEndLabel) // Kết thúc elif
	}

	// Xử lý else (nếu có)
	if s.ElseBlock != nil {
		c.compileIfBlockStatement(s.ElseBlock.Body)
	}

	// Định nghĩa label "if_end" (kết thúc toàn bộ if)
	c.defineLabel("if_end")

	// Resolve tất cả jumps sau khi biết vị trí label
	c.resolveJumps()

	//Reset lại labels và pending jumps đẻ tránh trùng
	c.resetLabels()
}

func (c *Compiler) compileForStatement(s *ast.ForStatement) {
	// 1. Vào scope trước khi khởi tạo biến lặp
	c.enterScope()
	enterScopePos := len(c.Code)
	c.emit(bytecode.OP_ENTER_SCOPE, 0, s.Line) // Operand tạm = 0

	// 2. Khởi tạo biến lặp
	c.compileStatement(s.Init)

	// 3. Định nghĩa label bắt đầu vòng lặp
	c.defineLabel("for_start")

	// 4. Compile điều kiện lặp
	c.compileExpression(s.Condition)
	c.emitJumpToLabel(bytecode.OP_JUMP_IF_FALSE, "for_end", s.Line)

	// 5. Compile thân vòng lặp
	c.compileBlockStatement(s.Body)

	// 6. Cập nhật biến lặp
	c.compileStatement(s.Update)

	// 7. Jump ngược về đầu vòng lặp
	c.emitJumpToLabel(bytecode.OP_JUMP, "for_start", s.Line)

	// 8. Định nghĩa label kết thúc
	c.defineLabel("for_end")

	// 9. Cập nhật operand ENTER_SCOPE với số lượng biến thực tế
	c.Code[enterScopePos].Operand = len(c.CurrentScope)

	//10.Resolve jumps
	c.resolveJumps()

	// 11. Reset jumps để tránh ảnh hưởng đến code sau
	c.resetLabels()

	// 12. Thoát scope
	c.leaveScope()
	c.emit(bytecode.OP_LEAVE_SCOPE, nil, s.Line)

}

func (c *Compiler) compileWhileStatement(s *ast.WhileStatement) {
	// 1. Vào scope trước khi khởi tạo biến lặp
	c.enterScope()
	enterScopePos := len(c.Code)
	c.emit(bytecode.OP_ENTER_SCOPE, 0, s.Line) // Operand tạm = 0

	// 2. Định nghĩa label bắt đầu vòng lặp
	c.defineLabel("while_start")

	// 3. Compile điều kiện lặp
	c.compileExpression(s.Condition)
	c.emitJumpToLabel(bytecode.OP_JUMP_IF_FALSE, "while_end", s.Line)

	// 4. Compile thân vòng lặp
	c.compileBlockStatement(s.Body)

	// 6. Jump ngược về đầu vòng lặp
	c.emitJumpToLabel(bytecode.OP_JUMP, "while_start", s.Line)

	// 7. Định nghĩa label kết thúc
	c.defineLabel("while_end")

	//8.Resolve jumps
	c.resolveJumps()

	// 9. Reset jumps để tránh ảnh hưởng đến code sau
	c.resetLabels()

	// 5. Cập nhật operand ENTER_SCOPE với số lượng biến thực tế
	c.Code[enterScopePos].Operand = len(c.CurrentScope)

	// 10. Thoát scope

	c.leaveScope()
	c.emit(bytecode.OP_LEAVE_SCOPE, nil, s.Line)
}

func (c *Compiler) compileReturnStatement(s *ast.ReturnStatement) {
	if !c.IsInsideFunction {
		c.addError("return statement outside of a function", s.Line, 0, "return")
	}
	//Có giá trị thì compile giá trị, không thì compile nothing
	if s.Value != nil {
		c.compileExpression(s.Value)
	} else {
		c.compileExpression(ast.NothingExpression{})
	}

	//emit lệnh return
	c.emit(bytecode.OP_RETURN, nil, s.Line)
}

func endsWithReturn(body *ast.BlockStatement) bool {
	if len(body.Statements) == 0 {
		return false
	}
	_, ok := body.Statements[len(body.Statements)-1].(*ast.ReturnStatement)
	return ok
}

func (c *Compiler) compileFunctionDefinitionStatement(s *ast.FunctionDefinitionStatement) {
	// 1. Kiểm tra global scope
	if len(c.Scopes) > 0 {
		c.addError("Function definitions are only allowed at the top-level (global scope)", s.Line, 0, "function definition")
		return
	}

	// 2. Kiểm tra tên hàm hợp lệ
	if !c.isValidVariableName(s.Name.Value, s.Line) {
		return
	}

	// 3. Tạo function object
	fn := &bytecode.Function{
		Name:      s.Name.Value, // Thêm tên hàm để debug
		Arity:     len(s.Parameters),
		StartPC:   0,                 // Sẽ cập nhật sau
		LocalSize: len(s.Parameters), // Số params ban đầu
	}

	// 4. Lưu hàm vào constants pool và emit code
	funcIndex := c.addConstant(fn)
	c.emit(bytecode.OP_LOAD_CONST, funcIndex, s.Line)
	c.emit(bytecode.OP_MAKE_FUNCTION, nil, s.Line)

	// 5. Gán hàm vào global scope
	idx := len(c.GlobalSymbols)
	c.GlobalSymbols[s.Name.Value] = idx
	c.emit(bytecode.OP_STORE_GLOBAL, idx, s.Line)

	// 6. Jump qua thân hàm
	jumpPos := len(c.Code)
	c.emit(bytecode.OP_JUMP, 0, s.Line)

	// 7. Cập nhật StartPC (vị trí bắt đầu thân hàm)
	fn.StartPC = len(c.Code)

	// 8. Vào scope hàm
	c.enterScope()
	c.emit(bytecode.OP_ENTER_SCOPE, len(s.Parameters), s.Line)

	// 9. Đăng ký params vào scope
	for i, param := range s.Parameters {
		c.CurrentScope[param.Value] = i // Slot = index của param
		c.LocalInitDepth[param.Value] = 1
	}

	// 10. Compile thân hàm với flag đang trong hàm
	prevInFunction := c.IsInsideFunction
	c.IsInsideFunction = true
	c.compileBlockStatement(s.Body)
	c.IsInsideFunction = prevInFunction

	// 11. Tự động thêm return nếu thân hàm không kết thúc bằng return
	if !endsWithReturn(s.Body) {
		c.emit(bytecode.OP_LOAD_NOTHING, nil, s.Line)
		c.emit(bytecode.OP_RETURN, nil, s.Line)
	}

	// 12. Cập nhật LocalSize (params + local vars)
	fn.LocalSize = len(c.CurrentScope)

	// 13. Thoát scope
	c.leaveScope()
	c.emit(bytecode.OP_LEAVE_SCOPE, nil, s.Line)

	// 14. Sửa jump offset (trừ 1 vì IP sẽ tự tăng sau khi đọc jump)
	c.Code[jumpPos].Operand = len(c.Code) - jumpPos - 1
}
