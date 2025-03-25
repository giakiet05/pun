package ast

import "fmt"

// Expression represents an expression (like math operations, function calls)
type Expression interface {
	Node
	expressionNode()
}

// Identifier represents variable names
type Identifier struct {
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Value }

// NumberExpression represents a numeric value
type NumberExpression struct {
	Value float64 // Đổi từ string -> float64
}

func (n *NumberExpression) expressionNode()      {}
func (n *NumberExpression) TokenLiteral() string { return fmt.Sprintf("%v", n.Value) }

// StringExpression represents a string value
type StringExpression struct {
	Value string
}

func (s *StringExpression) expressionNode()      {}
func (s *StringExpression) TokenLiteral() string { return s.Value }

type BooleanExpression struct {
	Value bool
}

func (b *BooleanExpression) expressionNode() {

}
func (b *BooleanExpression) TokenLiteral() string {
	return fmt.Sprintf("%v", b.Value)
}

type UnaryExpression struct {
	Operator string
	Value    Expression
}

func (u UnaryExpression) TokenLiteral() string {
	return u.Operator
}

func (u UnaryExpression) expressionNode() {
}

// Represent arithmetic expressions
type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (b *BinaryExpression) expressionNode() {}
func (b *BinaryExpression) TokenLiteral() string {
	return b.Operator
}
