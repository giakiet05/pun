package parser

import (
	"fmt"
	"pun/ast"
	"pun/lexer"
	"strconv"
)

var precedences = map[string]int{
	"!": 6, // Unary có mức ưu tiên cao nhất
	"*": 5, "/": 5, "%": 5,
	"+": 4, "-": 4,
	"==": 3, "!=": 3, ">": 3, "<": 3, ">=": 3, "<=": 3,
	"&&": 2, // AND cao hơn OR
	"||": 1, // OR thấp nhất nhưng vẫn bắt đầu từ 1
}

func (p *Parser) getPrecedence(op string) int {
	if prec, ok := precedences[op]; ok {
		return prec
	}
	return 0
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	left := p.parsePrimaryExpression()
	if left == nil {
		return nil
	}

	for p.getPrecedence(p.curTok.Value) > precedence {
		op := p.curTok.Value
		p.nextToken()
		right := p.parseExpression(precedences[op])
		if right == nil {
			return nil
		}
		left = &ast.BinaryExpression{Left: left, Operator: op, Right: right}
	}
	return left
}

func (p *Parser) parsePrimaryExpression() ast.Expression {
	switch p.curTok.Type {
	case lexer.TOKEN_NUMBER:
		value, err := strconv.ParseFloat(p.curTok.Value, 64)
		if err != nil {
			p.addError(fmt.Sprintf("Invalid number: %s", p.curTok.Value), p.curTok.Line, p.curTok.Col)
			return nil
		}
		lit := &ast.NumberExpression{Value: value} // 🛠 Đổi từ string -> float64
		p.nextToken()
		return lit

	case lexer.TOKEN_STRING:
		lit := &ast.StringExpression{Value: p.curTok.Value}
		p.nextToken()
		return lit

	case lexer.TOKEN_IDENTIFIER:
		ident := &ast.Identifier{Value: p.curTok.Value}
		p.nextToken()
		return ident

	case lexer.TOKEN_BOOLEAN:
		booleanVal := p.curTok.Value == "true"
		lit := &ast.BooleanExpression{Value: booleanVal}
		p.nextToken()
		return lit

	case lexer.TOKEN_LPAREN:
		p.nextToken()
		expr := p.parseExpression(0)

		if p.curTok.Type != lexer.TOKEN_RPAREN {
			p.addError("Expected closing ')'", p.curTok.Line, p.curTok.Col)
			return nil
		}
		p.nextToken() // Ăn dấu ')'
		return expr
	case lexer.TOKEN_OPERATOR:
		if p.curTok.Value == "-" {
			operator := p.curTok.Value
			p.nextToken()
			value := p.parseExpression(p.getMaxPrec())
			expr := &ast.UnaryExpression{Operator: operator, Value: value}
			return expr
		}
		return nil
	case lexer.TOKEN_LOGICAL:
		if p.curTok.Value == "!" {
			operator := p.curTok.Value
			p.nextToken()
			value := p.parseExpression(p.getMaxPrec())
			expr := &ast.UnaryExpression{Operator: operator, Value: value}
			return expr
		}
		return nil
	case lexer.TOKEN_KEYWORD: // ✅ Thêm xử lý ask()
		if p.curTok.Value == "ask" {
			return p.parseAskExpression()
		}
		return nil

	default:
		p.addError(fmt.Sprintf("Unexpected token: %s", p.curTok.Value), p.curTok.Line, p.curTok.Col)
		return nil
	}
}

func (p *Parser) parseAskExpression() ast.Expression {
	askExpr := &ast.AskExpression{}

	p.nextToken() // Bỏ qua 'ask'

	if !p.expectCurrent(lexer.TOKEN_LPAREN) {
		return nil
	}
	p.nextToken() // Qua nội dung trong ngoặc

	// Nếu có nội dung trong ask("...")
	if p.curTok.Type != lexer.TOKEN_RPAREN {
		askExpr.Prompt = p.parseExpression(0)
	}

	if !p.expectCurrent(lexer.TOKEN_RPAREN) {
		return nil
	}
	p.nextToken() // Bỏ qua ')'

	return askExpr
}

func (p *Parser) getMaxPrec() int {
	maxPrec := 0
	for _, prec := range precedences {
		if prec > maxPrec {
			maxPrec = prec
		}
	}
	return maxPrec
}
