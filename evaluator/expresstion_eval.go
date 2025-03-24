package evaluator

import (
	"fmt"
	"math"
	"pun/ast"
)

func evalExpression(node ast.Expression, env *Environment) interface{} {
	switch node := node.(type) {
	case *ast.Identifier:
		val, ok := env.Get(node.Value)
		if !ok {
			fmt.Printf("Error: variable '%s' not found\n", node.Value)
			return nil
		}
		return val

	case *ast.NumberExpression:
		return node.Value

	case *ast.StringExpression:
		return node.Value
	case *ast.BooleanExpression:
		return node.Value

	case *ast.BinaryExpression:
		return evalBinaryExpression(node, env)

	default:
		return nil
	}
}

func evalBinaryExpression(node *ast.BinaryExpression, env *Environment) interface{} {
	left := evalExpression(node.Left, env)
	right := evalExpression(node.Right, env)

	leftStr, leftIsString := left.(string)
	rightStr, rightIsString := right.(string)

	// Nếu cả hai là chuỗi, gọi hàm tính toán chuỗi
	if leftIsString && rightIsString {
		return evalStringBinaryExpression(leftStr, rightStr, node.Operator)
	}
	if leftIsString || rightIsString {
		fmt.Println("Error: Cannot evaluate string and number")
		return nil
	}
	// Các toán tử khác: bắt buộc cả hai phải là số
	leftNum, leftIsNumber := left.(float64)
	rightNum, rightIsNumber := right.(float64)

	if !leftIsNumber || !rightIsNumber {
		fmt.Println("Error: Non-numeric value in arithmetic expression")
		return nil
	}

	return evalNumberBinaryExpression(leftNum, rightNum, node.Operator)
}

func evalNumberBinaryExpression(leftNum, rightNum float64, operator string) interface{} {
	switch operator {
	case "+":
		return leftNum + rightNum
	case "-":
		return leftNum - rightNum
	case "*":
		return leftNum * rightNum
	case "/":
		if rightNum == 0 {
			fmt.Println("Error: Division by zero")
			return nil
		}
		return leftNum / rightNum
	case "%":
		if rightNum == 0 {
			fmt.Println("Error: Division by zero")
			return nil
		}
		return math.Mod(leftNum, rightNum)
	case "==":
		return leftNum == rightNum
	case "!=":
		return leftNum != rightNum
	case "<":
		return leftNum < rightNum
	case ">":
		return leftNum > rightNum
	case "<=":
		return leftNum <= rightNum
	case ">=":
		return leftNum >= rightNum
	}
	fmt.Printf("Error: Unknown operator %s\n", operator)
	return nil
}

func evalStringBinaryExpression(leftStr, rightStr string, operator string) interface{} {
	switch operator {
	case "+":
		return leftStr + rightStr
	case "==":
		return leftStr == rightStr
	case "!=":
		return leftStr != rightStr
	case "<":
		return leftStr < rightStr
	case ">":
		return leftStr > rightStr
	case "<=":
		return leftStr <= rightStr
	case ">=":
		return leftStr >= rightStr
	}
	fmt.Printf("Error: Unknown operator %s\n", operator)
	return nil
}
