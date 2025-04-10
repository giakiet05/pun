package compiler

import (
	"pun/ast"
	"pun/bytecode"
)

func (c *Compiler) compileExpression(expr ast.Expression) {
	switch e := expr.(type) {
	case *ast.NumberExpression:
		constIndex := c.addConstant(e.Value)
		c.emit(bytecode.OP_LOAD_CONST, constIndex, e.Line)

	case *ast.StringExpression:
		constIndex := c.addConstant(e.Value)
		c.emit(bytecode.OP_LOAD_CONST, constIndex, e.Line)

	case *ast.BooleanExpression:
		constIndex := c.addConstant(e.Value)
		c.emit(bytecode.OP_LOAD_CONST, constIndex, e.Line)

	case *ast.NothingExpression:
		constIndex := c.addConstant(nil)
		c.emit(bytecode.OP_LOAD_CONST, constIndex, e.Line)
	case *ast.Identifier:
		// 1. Kiểm tra nếu là built-in constant
		if index, ok := c.BuiltinConstants[e.Value]; ok {
			c.emit(bytecode.OP_LOAD_CONST, index, e.Line)
			return
		}

		// 2. Kiểm tra nếu là built-in function
		if c.BuiltinFuncs[e.Value] {
			constIndex := c.addConstant(e.Value) // Lưu tên hàm như string constant
			c.emit(bytecode.OP_LOAD_CONST, constIndex, e.Line)
			return
		}

		// 3. Xử lý biến thông thường
		if slot, _, isGlobal, exists := c.resolveVariable(e.Value); exists {
			if isGlobal {
				c.emit(bytecode.OP_LOAD_GLOBAL, slot, e.Line)
			} else {
				initDepth := c.getInitDepth(e.Value)
				c.emit(bytecode.OP_LOAD_LOCAL, &bytecode.LocalVar{Slot: slot, Depth: initDepth}, e.Line)
			}
		} else {
			c.addError("undefined variable", e.Line, 0, e.Value)
		}

	case *ast.UnaryExpression:
		c.compileExpression(e.Value)
		switch e.Operator {
		case "-":
			c.emit(bytecode.OP_NEG, nil, e.Line)
		case "!":
			c.emit(bytecode.OP_NOT, nil, e.Line)
		}

	case *ast.ArrayExpression:
		// Compile từng element
		for _, elem := range e.Elements {
			c.compileExpression(elem)
		}
		// Tạo array với số lượng element
		c.emit(bytecode.OP_MAKE_ARRAY, len(e.Elements), e.Line)

	case *ast.ArrayIndexExpression:
		c.compileExpression(e.Array)
		c.compileExpression(e.Index)
		c.emit(bytecode.OP_ARRAY_GET, nil, e.Line)

	case *ast.FunctionCallExpression:
		// Compile từng argument
		for _, arg := range e.Arguments {
			c.compileExpression(arg)
		}
		// Compile function expression
		c.compileExpression(e.Function)
		// Gọi function với số argument
		c.emit(bytecode.OP_CALL, len(e.Arguments), e.Line)

	case *ast.BinaryExpression:
		c.compileExpression(e.Left)
		c.compileExpression(e.Right)
		switch e.Operator {
		case "+":
			c.emit(bytecode.OP_ADD, nil, e.Line)
		case "-":
			c.emit(bytecode.OP_SUB, nil, e.Line)
		case "*":
			c.emit(bytecode.OP_MUL, nil, e.Line)
		case "/":
			c.emit(bytecode.OP_DIV, nil, e.Line)
		case "%":
			c.emit(bytecode.OP_MOD, nil, e.Line)
		case "**":
			c.emit(bytecode.OP_POW, nil, e.Line)
		case "==":
			c.emit(bytecode.OP_EQ, nil, e.Line)
		case "!=":
			c.emit(bytecode.OP_NEQ, nil, e.Line)
		case "<":
			c.emit(bytecode.OP_LT, nil, e.Line)
		case ">":
			c.emit(bytecode.OP_GT, nil, e.Line)
		case "<=":
			c.emit(bytecode.OP_LTE, nil, e.Line)
		case ">=":
			c.emit(bytecode.OP_GTE, nil, e.Line)
		case "&&":
			c.emit(bytecode.OP_AND, nil, e.Line)
		case "||":
			c.emit(bytecode.OP_OR, nil, e.Line)
		}
	}
}
