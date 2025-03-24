package ast

// Statement represents a statement in Pun (like make, shout, etc.)
type Statement interface {
	Node
	statementNode()
}

// MakeStatement represents variable declaration (like: make x = 5)
type MakeStatement struct {
	Name  *Identifier
	Value Expression
}

func (ms *MakeStatement) statementNode()       {}
func (ms *MakeStatement) TokenLiteral() string { return "make" }

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
