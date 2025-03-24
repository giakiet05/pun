package ast

// Node is the base interface for all AST nodes
type Node interface {
	TokenLiteral() string
}

// Program is the root node of our AST
// It contains a list of statements
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
