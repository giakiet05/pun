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
	Line  int
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Value }

// NumberExpression represents a numeric value
type NumberExpression struct {
	Value float64 // Đổi từ string -> float64
	Line  int
}

func (n *NumberExpression) expressionNode()      {}
func (n *NumberExpression) TokenLiteral() string { return fmt.Sprintf("%v", n.Value) }

// StringExpression represents a string value
type StringExpression struct {
	Value string
	Line  int
}

func (s *StringExpression) expressionNode()      {}
func (s *StringExpression) TokenLiteral() string { return s.Value }

type BooleanExpression struct {
	Value bool
	Line  int
}

func (b *BooleanExpression) expressionNode() {

}
func (b *BooleanExpression) TokenLiteral() string {
	return fmt.Sprintf("%v", b.Value)
}

type UnaryExpression struct {
	Operator string
	Value    Expression
	Line     int
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
	Line     int
}

func (b *BinaryExpression) expressionNode() {}
func (b *BinaryExpression) TokenLiteral() string {
	return b.Operator
}

type ArrayExpression struct {
	Elements []Expression
	Line     int
}

func (a ArrayExpression) TokenLiteral() string {
	return "array"
}

func (a ArrayExpression) expressionNode() {

}

type ArrayIndexExpression struct {
	Array Expression
	Index Expression
	Line  int
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
	Line      int
}

func (m MethodCallExpression) TokenLiteral() string {
	return "method"
}

func (m MethodCallExpression) expressionNode() {
}

type FunctionCallExpression struct {
	Function  Expression   // Hàm cần gọi (có thể là biến hoặc một biểu thức)
	Arguments []Expression // Danh sách tham số
	Line      int
}

func (f FunctionCallExpression) TokenLiteral() string {
	return "function"
}

func (f FunctionCallExpression) expressionNode() {
}

type NothingExpression struct {
	Line int
}

func (n NothingExpression) TokenLiteral() string {
	return "nothing"
}

func (n NothingExpression) expressionNode() {

}

type IncDecExpression struct {
	Operator string
	Value    Expression
	IsPrefix bool
	Line     int
}

func (i IncDecExpression) TokenLiteral() string {
	return "incdec"
}

func (i IncDecExpression) expressionNode() {

}
