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
	case *ast.UnaryExpression:
		return evalUnaryExpression(node, env)
	case *ast.BinaryExpression:
		return evalBinaryExpression(node, env)

	default:
		return nil
	}
}

func evalUnaryExpression(node *ast.UnaryExpression, env *Environment) interface{} {
	value := evalExpression(node.Value, env)

	switch node.Operator {
	case "-":
		if numVal, isNum := value.(float64); isNum {
			return -numVal
		} else {
			fmt.Println("Error: The unary minus operator (-) can only be used for numbers. Want to use it for other datatypes? Find another language!")
			return nil
		}
	case "!":
		if boolVal, isBool := value.(bool); isBool {
			return !boolVal
		} else {
			fmt.Println("Error: The logical NOT operator (!) can only be used for booleans. Trying to invert a non-boolean? What are you smoking?")
			return nil
		}
	default:
		fmt.Println("What are you writing???")
		return nil
	}
}

func evalBinaryExpression(node *ast.BinaryExpression, env *Environment) interface{} {
	left := evalExpression(node.Left, env)
	right := evalExpression(node.Right, env)

	switch l := left.(type) {
	case string:
		if r, ok := right.(string); ok {
			return evalStringBinaryExpression(l, r, node.Operator)
		}
	case bool:
		if r, ok := right.(bool); ok {
			return evalBooleanBinaryExpression(l, r, node.Operator)
		}
	case float64:
		if r, ok := right.(float64); ok {
			return evalNumberBinaryExpression(l, r, node.Operator)
		}
	}
	fmt.Println("Error: Cannot evaluate different data types")
	return nil
}

func evalNumberBinaryExpression(left, right float64, operator string) interface{} {
	switch operator {
	case "+":
		return left + right
	case "-":
		return left - right
	case "*":
		return left * right
	case "/":
		if right == 0 {
			fmt.Println("Error: Division by zero")
			return nil
		}
		return left / right
	case "%":
		if right == 0 {
			fmt.Println("Error: Division by zero")
			return nil
		}
		return math.Mod(left, right)
	case "==":
		return left == right
	case "!=":
		return left != right
	case "<":
		return left < right
	case ">":
		return left > right
	case "<=":
		return left <= right
	case ">=":
		return left >= right
	}
	fmt.Printf("Error: Unknown operator %s\n", operator)
	return nil
}

func evalStringBinaryExpression(left, right string, operator string) interface{} {
	switch operator {
	case "+":
		return left + right
	case "==":
		return left == right
	case "!=":
		return left != right
	case "<":
		return left < right
	case ">":
		return left > right
	case "<=":
		return left <= right
	case ">=":
		return left >= right
	}
	fmt.Printf("Error: Unknown operator %s\n", operator)
	return nil
}

func evalBooleanBinaryExpression(left, right bool, operator string) interface{} {
	switch operator {
	case "&&":
		return left && right
	case "||":
		return left || right
	}
	fmt.Printf("Error: Unknown operator %s\n", operator)
	return nil
}
