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

type AskExpression struct {
	Prompt Expression
}

func (a AskExpression) TokenLiteral() string {
	return "ask"
}

func (a AskExpression) expressionNode() {

}

type ArrayExpression struct {
	Elements []Expression
}

func (a ArrayExpression) TokenLiteral() string {
	return "array"
}

func (a ArrayExpression) expressionNode() {

}

type ArrayIndexExpression struct {
	Array Expression
	Index Expression
}

func (a ArrayIndexExpression) TokenLiteral() string {
	return "[]"
}

func (a ArrayIndexExpression) expressionNode() {
}

type MethodCallExpression struct {
	Caller    Expression   // Thằng gọi method (ví dụ: array trong array.inject())
	Method    string       // Tên method ("inject" hoặc "vomit")
	Arguments []Expression // Danh sách đối số (nếu có)
}

func (m MethodCallExpression) TokenLiteral() string {
	return "method"
}

func (m MethodCallExpression) expressionNode() {
}

type FunctionCallExpression struct {
	Function  Expression   // Hàm cần gọi (có thể là biến hoặc một biểu thức)
	Arguments []Expression // Danh sách tham số
}

func (f FunctionCallExpression) TokenLiteral() string {
	return "function"
}

func (f FunctionCallExpression) expressionNode() {
}
