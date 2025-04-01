package evaluator

import (
	"pun/ast"
)

// FunctionObject đại diện cho một function trong runtime
type FunctionObject struct {
	Name       string
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

type BuiltInFunction struct {
	Fn func(args ...interface{}) interface{}
}
