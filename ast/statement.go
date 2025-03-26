package ast

// Statement represents a statement in Pun (like make, shout, etc.)
type Statement interface {
	Node
	statementNode()
}

// Represents the shout() func
type ShoutStatement struct {
	Arguments []Expression // Can be numbers, strings, or identifiers
}

func (ss *ShoutStatement) statementNode() {}
func (ss *ShoutStatement) TokenLiteral() string {
	return "shout"
}

type AssignStatement struct {
	Name  *Identifier
	Value Expression
}

func (as *AssignStatement) statementNode() {}
func (as *AssignStatement) TokenLiteral() string {
	return "="
}

type BlockStatement struct {
	Statements []Statement
}

func (b BlockStatement) TokenLiteral() string {
	return "block"
}

func (b BlockStatement) statementNode() {
}

type WhenStatement struct {
	Condition Expression
	Body      *BlockStatement
	ElseIfs   []MaybeStatement
	ElseBlock *OtherwiseStatement
}

func (w WhenStatement) TokenLiteral() string {
	return "when"
}

func (w WhenStatement) statementNode() {
}

type MaybeStatement struct {
	Condition Expression
	Body      *BlockStatement
}

func (m MaybeStatement) TokenLiteral() string {
	return "maybe"
}

func (m MaybeStatement) statementNode() {
}

type OtherwiseStatement struct {
	Body *BlockStatement
}

func (o OtherwiseStatement) TokenLiteral() string {
	return "otherwise"
}

func (o OtherwiseStatement) statementNode() {
}

type ForStatement struct {
	Init      Statement  // Khởi tạo biến (i = 0)
	Condition Expression // Điều kiện (i < 10)
	Update    Statement  // Cập nhật biến (i = i + 1)
	Body      *BlockStatement
}

func (f *ForStatement) statementNode()       {}
func (f *ForStatement) TokenLiteral() string { return "for" }

type WhileStatement struct {
	Condition Expression
	Body      *BlockStatement
}

func (w *WhileStatement) statementNode()       {}
func (w *WhileStatement) TokenLiteral() string { return "while" }

type UntilStatement struct {
	Condition Expression
	Body      *BlockStatement
}

func (u *UntilStatement) statementNode()       {}
func (u *UntilStatement) TokenLiteral() string { return "until" }
