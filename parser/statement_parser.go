package parser

import (
	"fmt"
	"pun/ast"
	"pun/lexer"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curTok.Value {
	case "make":
		return p.parseMakeStatement()
	case "shout":
		return p.parseShoutStatement()
	default:
		//If the current token is an identifier and the next token is =, then we know this is an assignment
		if p.curTok.Type == lexer.TOKEN_IDENTIFIER && p.peekTok.Type == lexer.TOKEN_ASSIGN {
			return p.parseAssignStatement()
		}

		p.addError(fmt.Sprintf("Unexpected statement: %s", p.curTok.Value), p.curTok.Line, p.curTok.Col)
		return nil
	}
}

func (p *Parser) parseMakeStatement() *ast.MakeStatement {
	stmt := &ast.MakeStatement{}

	if !p.expectPeek(lexer.TOKEN_IDENTIFIER) {
		return nil
	}

	p.nextToken()

	stmt.Name = &ast.Identifier{Value: p.curTok.Value}

	if !p.expectPeek(lexer.TOKEN_ASSIGN) {
		return nil
	}

	p.nextToken()
	p.nextToken()

	stmt.Value = p.parseExpression(0) // Start with lowest precedence
	if stmt.Value == nil {
		p.addError("Invalid value in assignment", p.curTok.Line, p.curTok.Col)
		return nil
	}

	return stmt
}

func (p *Parser) parseAssignStatement() ast.Statement {
	stmt := &ast.AssignStatement{}

	stmt.Name = &ast.Identifier{Value: p.curTok.Value}

	if !p.expectPeek(lexer.TOKEN_ASSIGN) {
		return nil
	}
	p.nextToken()
	p.nextToken()
	stmt.Value = p.parseExpression(0)
	if stmt.Value == nil {
		p.addError("Invalid value in assignment", p.curTok.Line, p.curTok.Col)
		return nil
	}
	return stmt
}

func (p *Parser) parseShoutStatement() *ast.ShoutStatement {
	stmt := &ast.ShoutStatement{}
	stmt.Arguments = []ast.Expression{}

	if !p.expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}
	p.nextToken()

	//If there is no args (peekToken is RPAREN), then return immediately
	if p.peekTok.Type == lexer.TOKEN_RPAREN {
		return stmt
	}

	p.nextToken()

	// Allow zero or more arguments
	for p.curTok.Type != lexer.TOKEN_RPAREN && p.curTok.Type != lexer.TOKEN_EOF {
		arg := p.parseExpression(0)
		if arg == nil {
			return nil
		}
		stmt.Arguments = append(stmt.Arguments, arg)

		// Handle comma separation
		if p.curTok.Type == lexer.TOKEN_COMMA {
			p.nextToken() // Consume comma
		} else {
			break
		}
	}

	// Ensure closing parenthesis
	if p.curTok.Type != lexer.TOKEN_RPAREN {
		return nil
	}
	p.nextToken()
	return stmt
}
