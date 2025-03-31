package evaluator

import (
	"fmt"
	"math"
	"pun/ast"
)

func Eval(node ast.Node, env *Environment) interface{} {
	switch node := node.(type) {
	case *ast.Program:

		return evalProgram(node, env)
	case *ast.AssignStatement:
		return evalAssignStatement(node, env)
	case *ast.ShoutStatement:
		return evalShoutStatement(node, env)
	case *ast.WhenStatement:
		return evalWhenStatement(node, env)
	case *ast.UntilStatement:
		return evalUntilStatement(node, env)
	case *ast.ForStatement:
		return evalForStatement(node, env)
	case *ast.WhileStatement:
		return evalWhileStatement(node, env)
	case *ast.ExpressionStatement:
		return evalExpressionStatement(node, env)
	case *ast.FunctionDefinitionStatement:
		return evalFunctionDefinitionStatement(node, env)
	case *ast.StopStatement:
		panic(&StopException{})

	case *ast.ContinueStatement:
		panic(&ContinueException{})

	case *ast.ReturnStatement:
		var returnValue interface{}
		if node.Value != nil {
			returnValue = evalExpression(node.Value, env)
		}
		panic(&ReturnException{Value: returnValue})
	default:
		return nil
	}
}

func evalProgram(prog *ast.Program, env *Environment) interface{} {
	var result interface{}
	for _, stmt := range prog.Statements {
		result = Eval(stmt, env)
	}
	return result
}

func evalBlock(block *ast.BlockStatement, env *Environment) interface{} {
	var result interface{}
	for _, stmt := range block.Statements {
		result = Eval(stmt, env)
	}
	return result
}

func evalShoutStatement(node *ast.ShoutStatement, env *Environment) interface{} {
	if len(node.Arguments) == 0 {
		fmt.Println() // No arguments -> print a newline
		return nil
	}

	for _, arg := range node.Arguments {
		value := evalExpression(arg, env)
		fmt.Print(value, " ") // Print each argument with a space
	}
	fmt.Println() // Newline after all arguments

	return nil
}

func evalAssignStatement(node *ast.AssignStatement, env *Environment) interface{} {
	value := evalExpression(node.Value, env)

	switch name := node.Name.(type) {
	case *ast.Identifier:
		// Gán giá trị cho biến thường
		env.Set(name.Value, value)

	case *ast.ArrayIndexExpression:
		// Đánh giá mảng và index
		arrayVal := evalExpression(name.Array, env)
		indexVal := evalExpression(name.Index, env)

		// Kiểm tra có phải mảng không
		arr, ok := arrayVal.([]interface{})
		if !ok {
			fmt.Println("Error: Cannot assign to non-array value:", arrayVal)
			return nil
		}

		// Kiểm tra index có hợp lệ không
		idxFloat, ok := indexVal.(float64)
		if !ok {
			fmt.Println("Error: Index must be a number:", indexVal)
			return nil
		}

		// Làm tròn index về số nguyên
		idx := int(math.Round(idxFloat))

		if idx < 0 || idx >= len(arr) {
			fmt.Println("Error: Array index out of bounds:", idx)
			return nil
		}

		// Gán giá trị vào phần tử mảng
		arr[idx] = value

	default:
		fmt.Println("Error: Invalid assignment target:", node.Name)
		return nil
	}

	return value
}

func evalWhenStatement(node *ast.WhenStatement, env *Environment) interface{} {
	whenEnv := NewEnclosedEnvironment(env)

	condition := evalExpression(node.Condition, whenEnv)
	if isTruthy(condition) {
		return evalBlock(node.Body, whenEnv)
	}

	for _, maybe := range node.ElseIfs {
		maybeEnv := NewEnclosedEnvironment(env)
		condition := evalExpression(maybe.Condition, maybeEnv)
		if isTruthy(condition) {
			return evalBlock(maybe.Body, maybeEnv)
		}
	}
	otherwiseEnv := NewEnclosedEnvironment(env)
	if node.ElseBlock != nil {
		return evalBlock(node.ElseBlock.Body, otherwiseEnv)
	}
	return nil
}

func evalForStatement(node *ast.ForStatement, env *Environment) interface{} {
	forEnv := NewEnclosedEnvironment(env)
	Eval(node.Init, forEnv)

	defer func() { // Bọc vòng lặp để bắt StopException
		if r := recover(); r != nil {
			if _, ok := r.(*StopException); ok {
				// Nếu là StopException -> Dừng vòng for luôn
				return
			}

			panic(r) // Nếu là lỗi khác -> Ném tiếp
		}
	}()

	for {
		condition := evalExpression(node.Condition, forEnv)
		if !isTruthy(condition) {
			break
		}

		func() { // IIFE để bọc `defer`
			defer func() {
				if r := recover(); r != nil {
					switch r.(type) {
					case *StopException:
						panic(r) // QUAN TRỌNG: Ném lại để vòng for bắt được
					case *ContinueException:
						Eval(node.Update, forEnv) // Chạy update trước khi tiếp tục
						return                    // Ném lại để vòng `for` tiếp tục
					}
				}
			}()

			evalBlock(node.Body, forEnv)
			Eval(node.Update, forEnv)
		}()
	}
	return nil
}

func evalWhileStatement(node *ast.WhileStatement, env *Environment) interface{} {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(*StopException); ok {
				return
			}
			panic(r) // Nếu là lỗi khác -> Ném tiếp
		}
	}()

	whileEnv := NewEnclosedEnvironment(env)

	for {
		condition := evalExpression(node.Condition, whileEnv) // ❗ Dùng env gốc để kiểm tra điều kiện
		if !isTruthy(condition) {
			break
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					switch r.(type) {
					case *StopException:
						panic(r)
					case *ContinueException:
						return
					}
				}
			}()

			evalBlock(node.Body, whileEnv) // ❗ Dùng whileEnv để tạo biến mới
		}()
	}
	return nil
}

func evalUntilStatement(node *ast.UntilStatement, env *Environment) interface{} {

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(*StopException); ok {
				return
			}
			panic(r) // Nếu là lỗi khác -> Ném tiếp
		}
	}()

	untilEnv := NewEnclosedEnvironment(env)

	for {
		condition := evalExpression(node.Condition, untilEnv)
		if isTruthy(condition) {
			break
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					switch r.(type) {
					case *StopException:
						panic(r)
					case *ContinueException:
						return
					}
				}
			}()

			evalBlock(node.Body, untilEnv)
		}()
	}
	return nil
}

func evalFunctionDefinitionStatement(node *ast.FunctionDefinitionStatement, env *Environment) interface{} {
	fn := &FunctionObject{
		Name:       node.Name.Value,
		Parameters: node.Parameters,
		Body:       node.Body,
		Env:        NewEnclosedEnvironment(env),
	}
	env.Set(node.Name.Value, fn)
	return fn
}

func evalExpressionStatement(node *ast.ExpressionStatement, env *Environment) interface{} {
	return evalExpression(node.Expression, env)
}
