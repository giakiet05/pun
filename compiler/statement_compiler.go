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

		c.compileAssign(s)
	case *ast.IfStatement:
		c.compileIf(s)
	case *ast.ForStatement:
		c.compileFor(s)
	case *ast.WhileStatement:
		c.compileWhile(s)
	case *ast.FunctionDefinitionStatement:
		c.compileFuncDef(s)
	case *ast.ReturnStatement:
		c.compileReturn(s)
	case *ast.BreakStatement:
		c.compileBreak()
	case *ast.ContinueStatement:
		c.compileContinue()
	default:
		c.addError(fmt.Sprintf("Unsupported statement type: %T", stmt), 0, 0, "compile statement")
	}
}

func (c *Compiler) compileAssign(s *ast.AssignStatement) {
	// Luôn compile giá trị bên phải trước
	c.compileExpression(s.Value)

	// Xử lý target assignment
	switch target := s.Name.(type) {
	case *ast.Identifier:
		name := target.Value

		// Check tên biến hợp lệ (không trùng built-in)
		if !c.isValidVariableName(name) {
			return // Đã có error trong isValidVariableName
		}

		// Global scope (không có thì tạo mới, có thì cho operand = slot của cái đang có)
		if len(c.Scopes) == 0 {
			idx, exists := c.GlobalSymbols[name]
			if !exists {
				idx = len(c.GlobalSymbols)
				c.GlobalSymbols[name] = idx
			}
			c.emit(bytecode.OP_STORE_GLOBAL, idx)
			return
		}

		// Dùng resolveVariable để xử lí biến trong scope
		slot, depth, isGlobal, exists := c.resolveVariable(name)

		if isGlobal {
			c.emit(bytecode.OP_STORE_GLOBAL, slot) // Global override
		} else if exists {
			initDepth := c.LocalInitDepth[name] // lay depth cua scope ma bien nay duoc khoi tao (de gan cho bien do)
			operand := initDepth<<8 | slot
			c.emit(bytecode.OP_STORE_LOCAL, operand) // Local reassign
		} else {
			// Tạo local mới nếu biến chưa tồn tại anywhere
			newSlot := len(c.CurrentScope)
			c.CurrentScope[name] = newSlot
			c.LocalInitDepth[name] = depth
			operand := depth<<8 | newSlot
			c.emit(bytecode.OP_STORE_LOCAL, operand)
		}

	case *ast.ArrayIndexExpression:
		// Thêm check kiểu array trước khi gán
		c.compileExpression(target.Array)
		c.compileExpression(target.Index)
		c.emit(bytecode.OP_ARRAY_SET)

	default:
		c.addError(fmt.Sprintf("Unsupported assignment target: %T", target), 0, 0, "")
	}
}

func (c *Compiler) compileBlock(s *ast.BlockStatement) {

	for _, stmt := range s.Statements {
		c.compileStatement(stmt)
	}

}

func (c *Compiler) compileIfBlock(s *ast.BlockStatement) {
	// Start of scope
	c.enterScope()
	enterScopePos := c.emitWithPatch(bytecode.OP_ENTER_SCOPE)

	// Compile block contents
	c.compileBlock(s)

	// Patch the ENTER_SCOPE operand with actual local var count
	c.patchOperand(enterScopePos, len(c.CurrentScope))

	// Leave scope
	c.leaveScope()
	c.emit(bytecode.OP_LEAVE_SCOPE)
}

func (c *Compiler) compileIf(s *ast.IfStatement) {
	// Compile condition
	c.compileExpression(s.Condition)

	// Emit jump-if-false with temporary operand
	jumpToElsePos := c.emitWithPatch(bytecode.OP_JUMP_IF_FALSE)

	// Compile if body
	c.compileIfBlock(s.Body)

	// Emit jump-to-end with temporary operand
	jumpToEndPos := c.emitWithPatch(bytecode.OP_JUMP)

	// Record position where else/elif starts
	elsePos := len(c.Code)
	// Patch the jump-if-false to jump here
	c.patchOperand(jumpToElsePos, elsePos)

	// Handle elseifs
	var endJumps []int
	for _, elif := range s.ElseIfs {
		c.compileExpression(elif.Condition)
		jumpToNextPos := c.emitWithPatch(bytecode.OP_JUMP_IF_FALSE)

		c.compileIfBlock(elif.Body)

		// Add jump to end
		endJumps = append(endJumps, c.emitWithPatch(bytecode.OP_JUMP))

		// Patch the jump-if-false
		nextPos := len(c.Code)
		c.patchOperand(jumpToNextPos, nextPos)
	}

	// Handle else block
	if s.ElseBlock != nil {
		c.compileIfBlock(s.ElseBlock.Body)
	}

	// Record end position
	endPos := len(c.Code)

	// Patch all jumps to end
	c.patchOperand(jumpToEndPos, endPos)
	for _, pos := range endJumps {
		c.patchOperand(pos, endPos)
	}
}

func (c *Compiler) compileFor(s *ast.ForStatement) {
	// Save current break positions stack
	oldBreakPositions := c.breakPositions
	c.breakPositions = make([]int, 0)
	// Save current continue positions stack
	oldContinuePositions := c.continuePositions
	c.continuePositions = make([]int, 0)
	// 1. Create new scope for loop variables
	c.enterScope()
	// Save position for ENTER_SCOPE instruction - will patch with final local var count
	enterScopePos := c.emitWithPatch(bytecode.OP_ENTER_SCOPE)

	// 2. Compile initialization statement (runs once before loop)
	c.compileStatement(s.Init)

	// 3. Save position where loop condition check begins
	startPos := len(c.Code)

	// 4. Compile loop condition
	c.compileExpression(s.Condition)

	// 5. Emit conditional jump to end (if condition is false)
	// Save position to patch later with end position
	endJumpPos := c.emitWithPatch(bytecode.OP_JUMP_IF_FALSE)

	// 6. Compile loop body
	c.compileBlock(s.Body)

	// 7. Compile update statement (runs after each iteration)
	c.compileStatement(s.Update)

	// 8. Jump back to condition check
	c.emit(bytecode.OP_JUMP, startPos)

	// 9. Record end position and patch the conditional jump
	endPos := len(c.Code)
	c.patchOperand(endJumpPos, endPos)

	// Patch all break jumps to end position
	for _, pos := range c.breakPositions {
		c.patchOperand(pos, endPos)
	}

	// Restore previous break positions
	c.breakPositions = oldBreakPositions

	// Patch all continue jumps to point to start of condition check
	for _, pos := range c.continuePositions {
		c.patchOperand(pos, startPos)
	}

	// Restore previous continue positions
	c.continuePositions = oldContinuePositions
	// 10. Patch ENTER_SCOPE with final local variable count
	c.patchOperand(enterScopePos, len(c.CurrentScope))

	// 11. Clean up scope
	c.leaveScope()
	c.emit(bytecode.OP_LEAVE_SCOPE)
}

func (c *Compiler) compileWhile(s *ast.WhileStatement) {
	// Save current break positions stack
	oldBreakPositions := c.breakPositions
	c.breakPositions = make([]int, 0)
	// Save current continue positions stack
	oldContinuePositions := c.continuePositions
	c.continuePositions = make([]int, 0)
	// 1. Create new scope for loop variables
	c.enterScope()
	// Save position for ENTER_SCOPE instruction - will patch with final local var count
	enterScopePos := c.emitWithPatch(bytecode.OP_ENTER_SCOPE)

	// 2. Mark start of loop for continue statements
	startPos := len(c.Code)

	// 3. Compile condition
	c.compileExpression(s.Condition)

	// 4. Emit conditional jump to end with temporary operand
	endJumpPos := c.emitWithPatch(bytecode.OP_JUMP_IF_FALSE)

	// 5. Compile loop body
	c.compileBlock(s.Body)

	// 6. Jump back to condition check
	c.emit(bytecode.OP_JUMP, startPos)

	// 7. Record end position and patch the conditional jump
	endPos := len(c.Code)
	c.patchOperand(endJumpPos, endPos)

	// Patch all break jumps to end position
	for _, pos := range c.breakPositions {
		c.patchOperand(pos, endPos)
	}

	// Restore previous break positions
	c.breakPositions = oldBreakPositions

	// Patch all continue jumps to point to start of condition check
	for _, pos := range c.continuePositions {
		c.patchOperand(pos, startPos)
	}

	// Restore previous continue positions
	c.continuePositions = oldContinuePositions
	// 8. Patch ENTER_SCOPE with final local variable count
	c.patchOperand(enterScopePos, len(c.CurrentScope))

	// 9. Clean up scope
	c.leaveScope()
	c.emit(bytecode.OP_LEAVE_SCOPE)
}

func (c *Compiler) compileBreak() {
	// Save position of the break jump instruction to patch later
	breakPos := c.emitWithPatch(bytecode.OP_JUMP)

	// Track this break position to patch when we know the end of the loop
	c.breakPositions = append(c.breakPositions, breakPos)
}

func (c *Compiler) compileContinue() {
	// Save position of the continue jump instruction to patch later
	continuePos := c.emitWithPatch(bytecode.OP_JUMP)

	// Track this continue position to patch when we know the update/start position
	c.continuePositions = append(c.continuePositions, continuePos)
}

func (c *Compiler) compileReturn(s *ast.ReturnStatement) {
	if !c.IsInsideFunction {
		c.addError("return statement outside of a function", 0, 0, "return")
	}
	//Có giá trị thì compile giá trị, không thì compile nothing
	if s.Value != nil {
		c.compileExpression(s.Value)
	} else {
		c.compileExpression(ast.NothingExpression{})
	}

	//emit lệnh return
	c.emit(bytecode.OP_RETURN)
}

func (c *Compiler) compileFuncDef(s *ast.FunctionDefinitionStatement) {
	// 1. Kiểm tra global scope
	if len(c.Scopes) > 0 {
		c.addError("Function definitions are only allowed at the top-level (global scope)", 0, 0, "function definition")
		return
	}

	// 2. Kiểm tra tên hàm hợp lệ
	if !c.isValidVariableName(s.Name.Value) {
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
	c.emit(bytecode.OP_LOAD_CONST, funcIndex)
	c.emit(bytecode.OP_MAKE_FUNCTION)

	// 5. Gán hàm vào global scope
	idx := len(c.GlobalSymbols)
	c.GlobalSymbols[s.Name.Value] = idx
	c.emit(bytecode.OP_STORE_GLOBAL, idx)

	// 6. Jump qua thân hàm
	jumpPos := c.emitWithPatch(bytecode.OP_JUMP)

	// 7. Cập nhật StartPC (vị trí bắt đầu thân hàm)
	fn.StartPC = len(c.Code)

	// 8. Vào scope hàm
	c.enterScope()
	c.emit(bytecode.OP_ENTER_SCOPE, len(s.Parameters))

	// 9. Đăng ký params vào scope
	for i, param := range s.Parameters {
		c.CurrentScope[param.Value] = i // Slot = index của param
		c.LocalInitDepth[param.Value] = 1
	}

	// 10. Compile thân hàm với flag đang trong hàm
	prevInFunction := c.IsInsideFunction
	c.IsInsideFunction = true
	c.compileBlock(s.Body)
	c.IsInsideFunction = prevInFunction

	// 11. Tự động thêm return nếu thân hàm không kết thúc bằng return
	if !endsWithReturn(s.Body) {
		c.emit(bytecode.OP_LOAD_NOTHING)
		c.emit(bytecode.OP_RETURN)
	}

	// 12. Cập nhật LocalSize (params + local vars)
	fn.LocalSize = len(c.CurrentScope)

	// 13. Thoát scope
	c.leaveScope()
	c.emit(bytecode.OP_LEAVE_SCOPE)

	// 14. Sửa jump offset (trừ 1 vì IP sẽ tự tăng sau khi đọc jump)
	c.patchOperand(jumpPos, len(c.Code)-jumpPos-1)

}

func endsWithReturn(body *ast.BlockStatement) bool {
	if len(body.Statements) == 0 {
		return false
	}
	_, ok := body.Statements[len(body.Statements)-1].(*ast.ReturnStatement)
	return ok
}
