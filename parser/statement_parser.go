package parser

import (
	"fmt"
	"pun/ast"
	"pun/lexer"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curTok.Value {
	case "if":
		return p.parseIfStatement()
	case "for":
		return p.parseForStatement()
	case "while":
		return p.parseWhileStatement()
	case "until":
		return p.parseUntilStatement()
	case "break":
		return p.parseBreakStatement()
	case "continue":
		return p.parseContinueStatement()
	case "return":
		return p.parseReturnStatement()
	case "func":
		return p.parseFunctionDefinitionStatement()
	case "++", "--":
		return p.parseIncDecStatement()
	default:
		//If the current token is an identifier and the next token is =, then we know this is an assignment
		if p.curTok.Type == lexer.TOKEN_IDENTIFIER {
			// Parse expression cơ bản trước
			expr := p.parseExpression(0)
			if expr == nil {
				return nil
			}

			// Xử lý theo token tiếp theo
			switch p.curTok.Type {
			case lexer.TOKEN_ASSIGN:
				return p.parseAssignStatement(expr)
			default: //Các trường hợp còn lại
				return &ast.ExpressionStatement{Expression: expr, Line: p.curTok.Line}
			}
		}
	}

	p.addError(fmt.Sprintf("Unexpected statement: %s", p.curTok.Value), p.curTok.Line, p.curTok.Col)
	return nil
}

func (p *Parser) parseAssignStatement(expr ast.Expression) ast.Statement {
	stmt := &ast.AssignStatement{Line: p.curTok.Line}

	// Parse left-hand side
	left := expr
	if left == nil || !p.isValidAssignmentTarget(left) {
		p.addError("Invalid assignment target", p.curTok.Line, p.curTok.Col)
		return nil
	}

	stmt.Name = left

	// Check and consume '='
	if !p.expectCurrent(lexer.TOKEN_ASSIGN) {
		return nil
	}
	p.nextToken()

	// Parse right-hand side
	stmt.Value = p.parseExpression(0)
	if stmt.Value == nil {
		p.addError("Invalid value in assignment", p.curTok.Line, p.curTok.Col)
		return nil
	}

	return stmt
}

// Helper method to check valid assignment targets
func (p *Parser) isValidAssignmentTarget(expr ast.Expression) bool {
	switch expr.(type) {
	case *ast.Identifier, *ast.ArrayIndexExpression:
		return true
	default:
		return false
	}
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	ifStmt := &ast.IfStatement{Line: p.curTok.Line}
	p.nextToken()
	condition := p.parseExpression(0)

	if condition == nil {
		p.addError("Invalid condition in 'when' statement", p.curTok.Line, p.curTok.Col)
		return nil
	}

	ifStmt.Condition = condition

	if !p.expectCurrent(lexer.TOKEN_LCURLY) {
		return nil
	}

	ifStmt.Body = p.parseBlockStatement()

	if !p.expectCurrent(lexer.TOKEN_RCURLY) {
		return nil
	}

	p.nextToken()

	ifStmt.ElseIfs = []*ast.ElifStatement{}

	for p.curTok.Value == "elif" {
		elifStmt := p.parseElifStatement()
		if elifStmt != nil {
			ifStmt.ElseIfs = append(ifStmt.ElseIfs, elifStmt)
		}
	}

	if p.curTok.Value == "else" {
		ifStmt.ElseBlock = p.parseElseStatement()
	}

	return ifStmt
}

func (p *Parser) parseElifStatement() *ast.ElifStatement {
	elifStmt := &ast.ElifStatement{Line: p.curTok.Line}
	p.nextToken()
	condition := p.parseExpression(0)

	if condition == nil {
		p.addError("Invalid condition in 'elif' statement", p.curTok.Line, p.curTok.Col)
		return nil
	}

	elifStmt.Condition = condition

	if !p.expectCurrent(lexer.TOKEN_LCURLY) {
		return nil
	}

	elifStmt.Body = p.parseBlockStatement()

	if !p.expectCurrent(lexer.TOKEN_RCURLY) {
		return nil
	}

	p.nextToken()
	return elifStmt
}

func (p *Parser) parseElseStatement() *ast.ElseStatement {
	elseStmt := &ast.ElseStatement{Line: p.curTok.Line}

	p.nextToken()
	if !p.expectCurrent(lexer.TOKEN_LCURLY) {
		return nil
	}

	elseStmt.Body = p.parseBlockStatement()

	if !p.expectCurrent(lexer.TOKEN_RCURLY) {
		return nil
	}

	p.nextToken()
	return elseStmt

}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Line: p.curTok.Line}
	p.nextToken()
	for p.curTok.Type != lexer.TOKEN_RCURLY && p.curTok.Type != lexer.TOKEN_EOF {
		stmt := p.parseStatement()
		if stmt == nil {
			p.addError("Invalid statement", p.curTok.Line, p.curTok.Col)
			p.nextToken() // ⚠️ Quan trọng: Phải consume token lỗi để tránh infinite loop
			continue
		}
		block.Statements = append(block.Statements, stmt)
	}
	return block
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	forStmt := &ast.ForStatement{Line: p.curTok.Line}
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
	whileStmt := &ast.WhileStatement{Line: p.curTok.Line}
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
	untilStmt := &ast.UntilStatement{Line: p.curTok.Line}
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
	stmt := &ast.FunctionDefinitionStatement{Line: p.curTok.Line}
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
		param := &ast.Identifier{Value: p.curTok.Value, Line: p.curTok.Line}
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

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	p.nextToken() // Bỏ qua "stop"
	return &ast.BreakStatement{Line: p.curTok.Line}
}

func (p *Parser) parseContinueStatement() *ast.ContinueStatement {
	p.nextToken() // Bỏ qua "continue"
	return &ast.ContinueStatement{Line: p.curTok.Line}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Line: p.curTok.Line}
	p.nextToken() // Bỏ qua "return"

	// Nếu có giá trị return thì parse nó
	if p.curTok.Type != lexer.TOKEN_SEMICOLON && p.curTok.Type != lexer.TOKEN_RCURLY {
		stmt.Value = p.parseExpression(0)
	}

	return stmt
}

func (p *Parser) parseIncDecStatement() *ast.ExpressionStatement {
	expr := p.parseExpression(0)
	p.nextToken()
	return &ast.ExpressionStatement{Expression: expr, Line: p.curTok.Line}
}
