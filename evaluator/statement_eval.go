package evaluator

import (
	"fmt"
	"pun/ast"
)

func Eval(node ast.Node, env *Environment) interface{} {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.MakeStatement:
		value := evalExpression(node.Value, env)
		env.Set(node.Name.Value, value)
		return value

	case *ast.AssignStatement:
		_, ok := env.Get(node.Name.Value)
		if !ok {
			fmt.Printf("Error: variable '%s' not found\n", node.Name.Value)
			return nil
		}
		value := evalExpression(node.Value, env)
		env.Set(node.Name.Value, value)
		return value

	case *ast.ShoutStatement:
		return evalShoutStatement(node, env)

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
