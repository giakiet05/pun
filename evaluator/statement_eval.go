package evaluator

import (
	"fmt"
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
	env.Set(node.Name.Value, value)
	return value
}

func evalWhenStatement(node *ast.WhenStatement, env *Environment) interface{} {
	condition := evalExpression(node.Condition, env)
	if isTruthy(condition) {
		return evalBlock(node.Body, env)
	}

	for _, maybe := range node.ElseIfs {
		condition := evalExpression(maybe.Condition, env)
		if isTruthy(condition) {
			return evalBlock(maybe.Body, env)
		}
	}

	if node.ElseBlock != nil {
		return evalBlock(node.ElseBlock.Body, env)
	}
	return nil
}

func evalForStatement(node *ast.ForStatement, env *Environment) interface{} {
	Eval(node.Init, env)
	for {
		condition := evalExpression(node.Condition, env)
		if isTruthy(condition) {
			evalBlock(node.Body, env)
			Eval(node.Update, env)
		} else {
			break
		}
	}
	return nil
}

func evalWhileStatement(node *ast.WhileStatement, env *Environment) interface{} {
	for {
		condition := evalExpression(node.Condition, env)
		if isTruthy(condition) {
			evalBlock(node.Body, env)
		} else {
			break
		}
	}
	return nil
}

func evalUntilStatement(node *ast.UntilStatement, env *Environment) interface{} {
	for {
		condition := evalExpression(node.Condition, env)
		if !isTruthy(condition) {
			evalBlock(node.Body, env)
		} else {
			break
		}
	}
	return nil
}
