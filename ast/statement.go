package ast

// Statement represents a statement in Pun (like make, shout, etc.)
type Statement interface {
	Node
	statementNode()
}

type AssignStatement struct {
	Name  Expression
	Value Expression
	Line  int
}

func (as *AssignStatement) statementNode() {}
func (as *AssignStatement) TokenLiteral() string {
	return "="
}

type CompoundAssignStatement struct {
	Name  Expression
	Value Expression
	Line  int
}

func (as *CompoundAssignStatement) statementNode() {}
func (as *CompoundAssignStatement) TokenLiteral() string {
	return "="
}

type BlockStatement struct {
	Statements []Statement
	Line       int
}

func (b BlockStatement) TokenLiteral() string {
	return "block"
}

func (b BlockStatement) statementNode() {
}

type IfStatement struct {
	Condition Expression
	Body      *BlockStatement
	ElseIfs   []*ElifStatement
	ElseBlock *ElseStatement
	Line      int
}

func (w IfStatement) TokenLiteral() string {
	return "if"
}

func (w IfStatement) statementNode() {
}

type ElifStatement struct {
	Condition Expression
	Body      *BlockStatement
	Line      int
}

func (m ElifStatement) TokenLiteral() string {
	return "elif"
}

func (m ElifStatement) statementNode() {
}

type ElseStatement struct {
	Body *BlockStatement
	Line int
}

func (o ElseStatement) TokenLiteral() string {
	return "else"
}

func (o ElseStatement) statementNode() {
}

type ForStatement struct {
	Init      Statement  // Khởi tạo biến (i = 0)
	Condition Expression // Điều kiện (i < 10)
	Update    Statement  // Cập nhật biến (i = i + 1)
	Body      *BlockStatement
	Line      int
}

func (f *ForStatement) statementNode()       {}
func (f *ForStatement) TokenLiteral() string { return "for" }

type WhileStatement struct {
	Condition Expression
	Body      *BlockStatement
	Line      int
}

func (w *WhileStatement) statementNode()       {}
func (w *WhileStatement) TokenLiteral() string { return "while" }

type UntilStatement struct {
	Condition Expression
	Body      *BlockStatement
	Line      int
}

func (u *UntilStatement) statementNode()       {}
func (u *UntilStatement) TokenLiteral() string { return "until" }

type ExpressionStatement struct {
	Expression Expression
	Line       int
}

func (e ExpressionStatement) TokenLiteral() string {
	return "expression statement"
}

func (e ExpressionStatement) statementNode() {

}

type FunctionDefinitionStatement struct {
	Name       *Identifier     // Tên function
	Parameters []*Identifier   // Danh sách tham số
	Body       *BlockStatement // Thân hàm
	Line       int
}

func (f FunctionDefinitionStatement) TokenLiteral() string {
	return "func"
}

func (f FunctionDefinitionStatement) statementNode() {

}

type MethodDefinitionStatement struct {
	Receiver   *Identifier   // Tên của object (ví dụ: String)
	Name       *Identifier   // Tên method (ví dụ: uppercase)
	Parameters []*Identifier // Danh sách tham số
	Body       *BlockStatement
	Line       int
}

func (m MethodDefinitionStatement) TokenLiteral() string {
	return "method definition"
}

func (m MethodDefinitionStatement) statementNode() {

}

type BreakStatement struct{ Line int }

func (s BreakStatement) TokenLiteral() string {
	return "break"
}

func (s BreakStatement) statementNode() {

}

type ContinueStatement struct{ Line int }

func (c ContinueStatement) TokenLiteral() string {
	return "continue"
}

func (c ContinueStatement) statementNode() {

}

type ReturnStatement struct {
	Value Expression
	Line  int
}

func (r ReturnStatement) TokenLiteral() string {
	return "return"
}

func (r ReturnStatement) statementNode() {

}
