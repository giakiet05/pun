package compiler

import (
	"pun/ast"
	"pun/bytecode"
)

func (c *Compiler) compileExpression(expr ast.Expression) {
	switch e := expr.(type) {
	case *ast.NumberExpression:
		constIndex := c.addConstant(e.Value)
		c.emit(bytecode.OP_LOAD_CONST, constIndex)

	case *ast.StringExpression:
		constIndex := c.addConstant(e.Value)
		c.emit(bytecode.OP_LOAD_CONST, constIndex)

	case *ast.BooleanExpression:
		constIndex := c.addConstant(e.Value)
		c.emit(bytecode.OP_LOAD_CONST, constIndex)

	case *ast.NothingExpression:
		c.emit(bytecode.OP_LOAD_NOTHING)
	case *ast.Identifier:
		// 1. Kiểm tra nếu là built-in constant
		if index, ok := c.BuiltinConstants[e.Value]; ok {
			c.emit(bytecode.OP_LOAD_CONST, index)
			return
		}

		// 2. Kiểm tra nếu là built-in function
		if c.BuiltinFuncs[e.Value] {
			constIndex := c.addConstant(e.Value) // Lưu tên hàm như string constant
			c.emit(bytecode.OP_LOAD_CONST, constIndex)
			return
		}

		// 3. Xử lý biến thông thường
		if slot, depth, isGlobal, exists := c.resolveVariable(e.Value); exists {
			if isGlobal {
				c.emit(bytecode.OP_LOAD_GLOBAL, slot)
			} else {
				//initDepth := c.getInitDepth(e.Value)
				operand := (depth << 8) | slot

				c.emit(bytecode.OP_LOAD_LOCAL, operand)
			}
		} else {
			c.addError("undefined variable", 0, 0, e.Value)
		}

	case *ast.UnaryExpression:
		c.compileExpression(e.Value)
		switch e.Operator {
		case "-":
			c.emit(bytecode.OP_NEG)
		case "!":
			c.emit(bytecode.OP_NOT)
		}

	case *ast.ArrayExpression:
		// Compile từng element
		for _, elem := range e.Elements {
			c.compileExpression(elem)
		}
		// Tạo array với số lượng element
		c.emit(bytecode.OP_MAKE_ARRAY, len(e.Elements))

	case *ast.ArrayIndexExpression:
		c.compileExpression(e.Array)
		c.compileExpression(e.Index)
		c.emit(bytecode.OP_ARRAY_GET)

	case *ast.FunctionCallExpression:
		// Compile từng argument
		for _, arg := range e.Arguments {
			c.compileExpression(arg)
		}
		// Compile function expression
		c.compileExpression(e.Function)
		// Gọi function với số argument
		c.emit(bytecode.OP_CALL, len(e.Arguments))

	case *ast.BinaryExpression:
		c.compileExpression(e.Left)
		c.compileExpression(e.Right)
		switch e.Operator {
		case "+":
			c.emit(bytecode.OP_ADD)
		case "-":
			c.emit(bytecode.OP_SUB)
		case "*":
			c.emit(bytecode.OP_MUL)
		case "/":
			c.emit(bytecode.OP_DIV)
		case "%":
			c.emit(bytecode.OP_MOD)
		case "**":
			c.emit(bytecode.OP_POW)
		case "==":
			c.emit(bytecode.OP_EQ)
		case "!=":
			c.emit(bytecode.OP_NEQ)
		case "<":
			c.emit(bytecode.OP_LT)
		case ">":
			c.emit(bytecode.OP_GT)
		case "<=":
			c.emit(bytecode.OP_LTE)
		case ">=":
			c.emit(bytecode.OP_GTE)
		case "&&":
			c.emit(bytecode.OP_AND)
		case "||":
			c.emit(bytecode.OP_OR)
		}
	}
}
