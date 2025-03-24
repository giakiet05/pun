package parser

import (
	"fmt"
	"pun/ast"
	"pun/lexer"
	"strconv"
)

var precedences = map[string]int{
	"+":  2,
	"-":  2,
	"*":  3,
	"/":  3,
	"%":  3,
	"==": 1,
	"!=": 1,
	">":  1,
	"<":  1,
	">=": 1,
	"<=": 1,
}

func (p *Parser) isOperator(op string) bool {
	switch op {
	case "+", "-", "*", "/", "%", "==", "!=", ">", "<", ">=", "<=":
		return true
	default:
		return false
	}
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

	default:
		p.addError(fmt.Sprintf("Unexpected token: %s", p.curTok.Value), p.curTok.Line, p.curTok.Col)
		return nil
	}
}
