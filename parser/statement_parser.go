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
	case "stop":
		return p.parseStopStatement()
	case "continue":
		return p.parseContinueStatement()
	case "return":
		return p.parseReturnStatement()
	case "cook":
		return p.parseFunctionDefinitionStatement()
	default:
		//If the current token is an identifier and the next token is =, then we know this is an assignment
		if p.curTok.Type == lexer.TOKEN_IDENTIFIER {

			if p.peekTok.Type == lexer.TOKEN_DOT || p.peekTok.Type == lexer.TOKEN_LPAREN {
				expr := p.parseExpression(0)
				if expr != nil {
					return &ast.ExpressionStatement{Expression: expr}
				}
			}

			return p.parseAssignStatement()
		}

		p.addError(fmt.Sprintf("Unexpected statement: %s", p.curTok.Value), p.curTok.Line, p.curTok.Col)
		return nil
	}
}

func (p *Parser) parseAssignStatement() ast.Statement {
	stmt := &ast.AssignStatement{}

	// Parse bên trái dấu '=' (có thể là biến hoặc phần tử mảng)
	left := p.parseExpression(0)

	if left == nil {
		p.addError("Invalid assignment target", p.curTok.Line, p.curTok.Col)
		return nil
	}

	// Kiểm tra nếu không phải Identifier hay ArrayIndexExpression thì báo lỗi
	if _, ok := left.(*ast.Identifier); !ok {
		if _, ok := left.(*ast.ArrayIndexExpression); !ok {
			p.addError("Invalid assignment target", p.curTok.Line, p.curTok.Col)
			return nil
		}
	}

	stmt.Name = left // Gán biến hoặc ArrayIndexExpression vào Name

	if !p.expectCurrent(lexer.TOKEN_ASSIGN) {
		return nil
	}
	p.nextToken() // Bỏ qua '='

	stmt.Value = p.parseExpression(0) // Parse giá trị bên phải
	if stmt.Value == nil {
		p.addError("Invalid value in assignment", p.curTok.Line, p.curTok.Col)
		return nil
	}

	return stmt
}

func (p *Parser) parseShoutStatement() *ast.ShoutStatement {
	stmt := &ast.ShoutStatement{}

	p.nextToken()

	stmt.Arguments = p.parseArguments()

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

func (p *Parser) parseFunctionDefinitionStatement() *ast.FunctionDefinitionStatement {
	stmt := &ast.FunctionDefinitionStatement{}
	p.nextToken()

	if !p.expectCurrent(lexer.TOKEN_IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Value: p.curTok.Value}

	p.nextToken()

	if !p.expectCurrent(lexer.TOKEN_LPAREN) {
		return nil
	}

	p.nextToken()

	// Parse danh sách tham số
	stmt.Parameters = []*ast.Identifier{}

	for p.curTok.Type != lexer.TOKEN_RPAREN && p.curTok.Type != lexer.TOKEN_EOF {
		if !p.expectCurrent(lexer.TOKEN_IDENTIFIER) {
			return nil
		}
		param := &ast.Identifier{Value: p.curTok.Value}
		stmt.Parameters = append(stmt.Parameters, param)

		p.nextToken()

		if p.curTok.Type == lexer.TOKEN_COMMA {
			p.nextToken()
		}
	}
	if !p.expectCurrent(lexer.TOKEN_RPAREN) {
		return nil
	}
	p.nextToken()

	if !p.expectCurrent(lexer.TOKEN_LCURLY) {
		return nil
	}
	stmt.Body = p.parseBlockStatement()

	if !p.expectCurrent(lexer.TOKEN_RCURLY) {
		return nil
	}

	p.nextToken()

	return stmt

}

func (p *Parser) parseStopStatement() ast.Statement {
	p.nextToken() // Bỏ qua "stop"
	return &ast.StopStatement{}
}

func (p *Parser) parseContinueStatement() ast.Statement {
	p.nextToken() // Bỏ qua "continue"
	return &ast.ContinueStatement{}
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{}
	p.nextToken() // Bỏ qua "return"

	// Nếu có giá trị return thì parse nó
	if p.curTok.Type != lexer.TOKEN_SEMICOLON && p.curTok.Type != lexer.TOKEN_RCURLY {
		stmt.Value = p.parseExpression(0)
	}

	return stmt
}
