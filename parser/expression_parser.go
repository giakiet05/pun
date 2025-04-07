package parser

import (
	"fmt"
	"pun/ast"
	"pun/lexer"
	"strconv"
)

var precedences = map[string]int{
	"!":  7, // Unary có mức ưu tiên cao nhất
	"**": 6,
	"*":  5, "/": 5, "%": 5,
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
		lit := &ast.NumberExpression{Value: value, Line: p.curTok.Line} // 🛠 Đổi từ string -> float64
		p.nextToken()
		return lit

	case lexer.TOKEN_STRING:
		lit := &ast.StringExpression{Value: p.curTok.Value, Line: p.curTok.Line}
		p.nextToken()
		return lit

	case lexer.TOKEN_IDENTIFIER:
		ident := &ast.Identifier{Value: p.curTok.Value, Line: p.curTok.Line}
		p.nextToken()
		// Nếu có dấu `[` => Đây là truy xuất mảng
		if p.curTok.Type == lexer.TOKEN_LSQUARE {
			return p.parseArrayIndexExpression(ident) // Truyền array vào
		}
		if p.curTok.Type == lexer.TOKEN_LPAREN {
			return p.parseFunctionCallExpression(ident)
		}
		//Nếu có dấu . phía sau thì là method
		if p.curTok.Type == lexer.TOKEN_DOT {
			return p.parseMethodCallExpression(ident)
		}
		return ident

	case lexer.TOKEN_BOOLEAN:
		booleanVal := p.curTok.Value == "true"
		lit := &ast.BooleanExpression{Value: booleanVal, Line: p.curTok.Line}
		p.nextToken()
		return lit
	case lexer.TOKEN_NOTHING:
		p.nextToken()
		return &ast.NothingExpression{Line: p.curTok.Line}
	case lexer.TOKEN_LPAREN:
		p.nextToken()
		expr := p.parseExpression(0)

		if p.curTok.Type != lexer.TOKEN_RPAREN {
			p.addError("Expected closing ')'", p.curTok.Line, p.curTok.Col)
			return nil
		}
		p.nextToken() // Ăn dấu ')'
		return expr
	case lexer.TOKEN_LSQUARE:
		return p.parseArrayExpression()
	case lexer.TOKEN_ARITHMETIC:
		if p.curTok.Value == "-" {
			operator := p.curTok.Value
			p.nextToken()
			value := p.parseExpression(p.getMaxPrec())
			expr := &ast.UnaryExpression{Operator: operator, Value: value, Line: p.curTok.Line}
			return expr
		}
		return nil
	case lexer.TOKEN_LOGICAL:
		if p.curTok.Value == "!" {
			operator := p.curTok.Value
			p.nextToken()
			value := p.parseExpression(p.getMaxPrec())
			expr := &ast.UnaryExpression{Operator: operator, Value: value, Line: p.curTok.Line}
			return expr
		}
		return nil
	default:
		p.addError(fmt.Sprintf("Unexpected token: %s", p.curTok.Value), p.curTok.Line, p.curTok.Col)
		return nil
	}
}

func (p *Parser) parseIncDecExpression(ident *ast.Identifier, op string, isPrefix bool) ast.Expression {
	expr := &ast.IncDecExpression{
		Operator: op,
		Value:    ident,
		IsPrefix: isPrefix,
		Line:     p.curTok.Line,
	}
	p.nextToken()
	return expr
}

func (p *Parser) parseArrayExpression() ast.Expression {
	array := &ast.ArrayExpression{Line: p.curTok.Line}

	p.nextToken() //skip "["

	if p.curTok.Type == lexer.TOKEN_RSQUARE {
		p.nextToken()
		return array
	}

	for p.curTok.Type != lexer.TOKEN_RSQUARE && p.curTok.Type != lexer.TOKEN_EOF {
		element := p.parseExpression(0)
		if element != nil {
			array.Elements = append(array.Elements, element)
		} else {
			break
		}
		if p.curTok.Type == lexer.TOKEN_COMMA {
			if p.peekTok.Type == lexer.TOKEN_RSQUARE {
				p.addError("Trailling comma in array is not allowed", p.curTok.Line, p.curTok.Col)
				return nil
			}
			p.nextToken()
		}
	}

	if !p.expectCurrent(lexer.TOKEN_RSQUARE) {
		return nil
	}

	p.nextToken()

	return array

}

func (p *Parser) parseArrayIndexExpression(array ast.Expression) ast.Expression {
	expr := &ast.ArrayIndexExpression{Array: array, Line: p.curTok.Line}

	p.nextToken() // Bỏ qua '['

	index := p.parseExpression(0)

	if index == nil {
		return nil
	}

	expr.Index = index

	if !p.expectCurrent(lexer.TOKEN_RSQUARE) {
		return nil
	}

	p.nextToken()

	return expr
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

func (p *Parser) parseMethodCallExpression(caller ast.Expression) ast.Expression {
	expr := &ast.MethodCallExpression{Caller: caller, Line: p.curTok.Line}

	p.nextToken()

	if !p.expectCurrent(lexer.TOKEN_IDENTIFIER) {
		return nil
	}

	expr.Method = p.curTok.Value

	p.nextToken()
	expr.Arguments = p.parseArguments()
	p.nextToken()

	return expr
}

func (p *Parser) parseFunctionCallExpression(function *ast.Identifier) ast.Expression {
	expr := &ast.FunctionCallExpression{Function: function, Line: p.curTok.Line}

	expr.Arguments = p.parseArguments()

	return expr
}
