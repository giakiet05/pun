package evaluator

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"pun/ast"
	"strings"
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
	case *ast.AskExpression:
		return evalAskExpression(node, env)
	case *ast.ArrayExpression:
		return evalArrayExpression(node, env)
	case *ast.ArrayIndexExpression:
		return evalArrayIndexExpression(node, env)
	case *ast.FunctionCallExpression:
		return evalFunctionCallExpression(node, env)
	default:
		fmt.Printf("Error: Unsupported expression type: %T\n", node)
		return nil
	}
}

func evalArrayIndexExpression(node *ast.ArrayIndexExpression, env *Environment) interface{} {
	ident := evalExpression(node.Array, env)

	if ident == nil {
		fmt.Println("Error: Array doesn't exist")
		return nil
	}

	arr, isArray := ident.([]interface{})
	if !isArray {
		fmt.Println("Error: Value is not an array")
		return nil
	}

	val := evalExpression(node.Index, env)
	if val == nil {
		fmt.Println("Error: Incorrect index")
		return nil
	}

	// Ép kiểu về float64
	indexFloat, isFloat := val.(float64)
	if !isFloat {
		fmt.Println("Error: Index must be a number")
		return nil
	}

	// Kiểm tra số nguyên
	if indexFloat != math.Floor(indexFloat) {
		fmt.Println("Error: Index must be an integer")
		return nil
	}

	index := int(indexFloat) // Ép float64 -> int

	// Kiểm tra index hợp lệ
	if index < 0 || index >= len(arr) {
		fmt.Println("Error: Index out of range")
		return nil
	}

	return arr[index]
}

func evalArrayExpression(node *ast.ArrayExpression, env *Environment) []interface{} {
	var values []interface{}

	for _, element := range node.Elements {
		value := evalExpression(element, env)
		if value == nil {
			fmt.Println("Error: Failed to evaluate array element")
			return nil
		}
		values = append(values, value)
	}
	return values
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

func evalAskExpression(node *ast.AskExpression, env *Environment) interface{} {
	if node.Prompt != nil {
		fmt.Print(evalExpression(node.Prompt, env)) // ✅ In Prompt
	}

	reader := bufio.NewReader(os.Stdin)
	value, _ := reader.ReadString('\n') // ✅ Đọc nguyên dòng
	value = strings.TrimSpace(value)    // ✅ Xóa dấu xuống dòng
	return value
}

func evalFunctionCallExpression(node *ast.FunctionCallExpression, env *Environment) interface{} {
	funcObj, ok := env.Get(node.Function.(*ast.Identifier).Value)
	if !ok {
		fmt.Println("Error: function not found")
		return nil
	}

	fn, ok := funcObj.(*FunctionObject)
	if !ok {
		fmt.Println("Error: Not a function")
		return nil
	}

	fnEnv := NewEnclosedEnvironment(fn.Env)

	if len(node.Arguments) != len(fn.Parameters) {
		fmt.Println("Error: argument count mismatch")
		return nil
	}

	for i, param := range fn.Parameters {
		argVal := evalExpression(node.Arguments[i], env)
		fnEnv.Set(param.Value, argVal)
	}

	result := evalBlock(fn.Body, fnEnv)

	return result
}
