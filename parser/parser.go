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
		} else {
			break
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

func (p *Parser) expectCurrent(t string) bool {
	if p.curTok.Type == t {
		return true
	}
	p.addError(fmt.Sprintf("Expected next token to be %s, got %s instead", t, p.peekTok.Type), p.peekTok.Line, p.peekTok.Col)
	return false
}

func (p *Parser) parseArguments() []ast.Expression {

	var args []ast.Expression

	if !p.expectCurrent(lexer.TOKEN_LPAREN) {
		return nil
	}

	p.nextToken()

	//If there is no args (peekToken is RPAREN), then return immediately
	if p.curTok.Type == lexer.TOKEN_RPAREN {
		p.nextToken()
		return args
	}
	if p.curTok.Type == lexer.TOKEN_EOF {
		p.addError("Missing ')'", p.curTok.Line, p.curTok.Col)
		return nil
	}

	// Allow zero or more arguments
	for p.curTok.Type != lexer.TOKEN_RPAREN && p.curTok.Type != lexer.TOKEN_EOF {
		arg := p.parseExpression(0)
		if arg == nil {
			return nil
		}
		args = append(args, arg)

		// Handle comma separation
		if p.curTok.Type == lexer.TOKEN_COMMA {
			p.nextToken() // Consume comma
		} else {
			break
		}
	}

	// Ensure closing parenthesis
	if !p.expectCurrent(lexer.TOKEN_RPAREN) {
		return nil
	}

	return args
}
