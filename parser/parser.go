package parser

import (
	"fmt"
	"pun/ast"
	"pun/errors"
	"pun/lexer"
)

type Parser struct {
	l       *lexer.Lexer
	curTok  lexer.Token
	peekTok lexer.Token
	errors  []errors.PunError
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for p.peekTok.Type != lexer.TOKEN_EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}
	return program
}

// Use expectPeek when we want to make sure that the next token must be the type we want
func (p *Parser) expectPeek(t string) bool {
	if p.peekTok.Type == t {
		return true
	}
	p.addError(fmt.Sprintf("Expected next token to be %s, got %s instead", t, p.peekTok.Type), p.peekTok.Line, p.peekTok.Col)
	return false
}
