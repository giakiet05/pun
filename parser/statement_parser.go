package parser

import (
	"fmt"
	"pun/ast"
	"pun/lexer"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curTok.Value {
	case "shout":
		return p.parseShoutStatement()
	case "when":
		return p.parseWhenStatement()
	case "for":
		return p.parseForStatement()
	case "while":
		return p.parseWhileStatement()
	case "until":
		return p.parseUntilStatement()
	default:
		//If the current token is an identifier and the next token is =, then we know this is an assignment
		if p.curTok.Type == lexer.TOKEN_IDENTIFIER && p.peekTok.Type == lexer.TOKEN_ASSIGN {
			return p.parseAssignStatement()
		}

		p.addError(fmt.Sprintf("Unexpected statement: %s", p.curTok.Value), p.curTok.Line, p.curTok.Col)
		return nil
	}
}

func (p *Parser) parseAssignStatement() ast.Statement {
	stmt := &ast.AssignStatement{}

	stmt.Name = &ast.Identifier{Value: p.curTok.Value}

	if !p.expectPeek(lexer.TOKEN_ASSIGN) {
		p.nextToken()
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
		p.nextToken()
		return nil
	}
	p.nextToken()

	//If there is no args (peekToken is RPAREN), then return immediately
	if p.peekTok.Type == lexer.TOKEN_RPAREN {
		p.nextToken()
		return stmt
	}
	if p.peekTok.Type == lexer.TOKEN_EOF {
		p.addError("Unterminated 'shout' statement, missing ')'", p.curTok.Line, p.curTok.Col)
		return nil
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
		p.addError("Unterminated 'shout' statement, missing ')'", p.curTok.Line, p.curTok.Col)
		return nil
	}
	p.nextToken()
	return stmt
}

func (p *Parser) parseWhenStatement() *ast.WhenStatement {
	whenStmt := &ast.WhenStatement{}
	p.nextToken()
	condition := p.parseExpression(0)

	if condition == nil {
		p.addError("Invalid condition in 'when' statement", p.curTok.Line, p.curTok.Col)
		return nil
	}

	whenStmt.Condition = condition

	if !p.expectCurrent(lexer.TOKEN_LCURLY) {
		return nil
	}

	whenStmt.Body = p.parseBlockStatement()

	if !p.expectCurrent(lexer.TOKEN_RCURLY) {
		return nil
	}

	p.nextToken()

	whenStmt.ElseIfs = []ast.MaybeStatement{}

	for p.curTok.Value == "maybe" {
		maybeStatement := p.parseMaybeStatement()
		if maybeStatement != nil {
			whenStmt.ElseIfs = append(whenStmt.ElseIfs, *maybeStatement)
		}
	}

	if p.curTok.Value == "otherwise" {
		whenStmt.ElseBlock = p.parseOtherwiseStatement()
	}

	p.nextToken()

	return whenStmt
}

func (p *Parser) parseMaybeStatement() *ast.MaybeStatement {
	maybeStmt := &ast.MaybeStatement{}
	p.nextToken()
	condition := p.parseExpression(0)

	if condition == nil {
		p.addError("Invalid condition in 'maybe' statement", p.curTok.Line, p.curTok.Col)
		return nil
	}

	maybeStmt.Condition = condition

	if !p.expectCurrent(lexer.TOKEN_LCURLY) {
		return nil
	}

	maybeStmt.Body = p.parseBlockStatement()

	if !p.expectCurrent(lexer.TOKEN_RCURLY) {
		return nil
	}

	p.nextToken()
	return maybeStmt
}

func (p *Parser) parseOtherwiseStatement() *ast.OtherwiseStatement {
	otherwiseStmt := &ast.OtherwiseStatement{}

	p.nextToken()
	if !p.expectCurrent(lexer.TOKEN_LCURLY) {
		return nil
	}

	otherwiseStmt.Body = p.parseBlockStatement()

	if !p.expectCurrent(lexer.TOKEN_RCURLY) {
		return nil
	}

	p.nextToken()
	return otherwiseStmt

}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{}
	p.nextToken()
	for p.curTok.Type != lexer.TOKEN_RCURLY && p.curTok.Type != lexer.TOKEN_EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}
	return block
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	forStmt := &ast.ForStatement{}
	p.nextToken()

	init := p.parseStatement()

	if init == nil {
		p.addError("Wrong init statement!!!", p.curTok.Line, p.curTok.Col)
		return nil
	}

	if !p.expectCurrent(lexer.TOKEN_SEMICOLON) {
		return nil
	}

	p.nextToken()

	condition := p.parseExpression(0)
	if condition == nil {
		p.addError("Invalid condition in 'for' statement", p.curTok.Line, p.curTok.Col)
		return nil
	}

	if !p.expectCurrent(lexer.TOKEN_SEMICOLON) {
		return nil
	}

	p.nextToken()

	update := p.parseStatement()

	if update == nil {
		p.addError("Wrong update statement!!!", p.curTok.Line, p.curTok.Col)
		return nil
	}

	forStmt.Init = init
	forStmt.Condition = condition
	forStmt.Update = update

	if !p.expectCurrent(lexer.TOKEN_LCURLY) {
		return nil
	}

	forStmt.Body = p.parseBlockStatement()

	if !p.expectCurrent(lexer.TOKEN_RCURLY) {
		return nil
	}

	p.nextToken()
	return forStmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	whileStmt := &ast.WhileStatement{}
	p.nextToken()
	condition := p.parseExpression(0)

	if condition == nil {
		p.addError("Invalid condition in 'while' statement", p.curTok.Line, p.curTok.Col)
		return nil
	}

	whileStmt.Condition = condition

	if !p.expectCurrent(lexer.TOKEN_LCURLY) {
		return nil
	}

	whileStmt.Body = p.parseBlockStatement()

	if !p.expectCurrent(lexer.TOKEN_RCURLY) {
		return nil
	}

	p.nextToken()
	return whileStmt
}

func (p *Parser) parseUntilStatement() *ast.UntilStatement {
	untilStmt := &ast.UntilStatement{}
	p.nextToken()
	condition := p.parseExpression(0)

	if condition == nil {
		p.addError("Invalid condition in 'until' statement", p.curTok.Line, p.curTok.Col)
		return nil
	}

	untilStmt.Condition = condition

	if !p.expectCurrent(lexer.TOKEN_LCURLY) {
		return nil
	}

	untilStmt.Body = p.parseBlockStatement()

	if !p.expectCurrent(lexer.TOKEN_RCURLY) {
		return nil
	}

	p.nextToken()
	return untilStmt
}
