package parser

import (
	"fmt"
	"pun/ast"
	"pun/error"
	"pun/lexer"
	"strings"
)

type Parser struct {
	lexer   *lexer.Lexer
	curTok  lexer.Token
	peekTok lexer.Token
	errors  []customError.SyntaxError
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.lexer.NextToken()

	for p.peekTok.Type == lexer.TOKEN_COMMENT {
		p.peekTok = p.lexer.NextToken()
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for p.peekTok.Type != lexer.TOKEN_EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
			// Biên dịch statement ngay sau khi parse

		} else {
			p.syncTo(lexer.TOKEN_SEMICOLON)
		}
	}

	return program // Trả về cả AST và compiler
}

// Hàm parseArguments GIỮ NGUYÊN như bản gốc
func (p *Parser) parseArguments() []ast.Expression {
	var args []ast.Expression

	if !p.expectCurrent(lexer.TOKEN_LPAREN) {
		return nil
	}

	p.nextToken()

	if p.curTok.Type == lexer.TOKEN_RPAREN {
		p.nextToken()
		return args
	}

	for p.curTok.Type != lexer.TOKEN_RPAREN && p.curTok.Type != lexer.TOKEN_EOF {
		arg := p.parseExpression(0)
		if arg == nil {
			return nil
		}
		args = append(args, arg)

		if p.curTok.Type == lexer.TOKEN_COMMA {
			p.nextToken()
		} else {
			break
		}
	}

	if !p.expectCurrent(lexer.TOKEN_RPAREN) {
		p.addError("Missing closing ')'", p.curTok.Line, p.curTok.Col)
		return nil
	}

	p.nextToken()
	return args
}

// Hàm addError dùng SyntaxError.Error()
func (p *Parser) addError(message string, line, col int) {
	err := customError.SyntaxError{
		PunError: customError.PunError{
			Message: message,
			Line:    line,
			Column:  col,
		},
		Context: fmt.Sprintf("Near token: %q (Type: %s)", p.curTok.Value, p.curTok.Type),
	}
	p.errors = append(p.errors, err)
}

func (p *Parser) expectPeek(t string) bool {
	if p.peekTok.Type == t {
		p.nextToken()
		return true
	}
	p.addError(
		fmt.Sprintf("Expected next token to be %s, got %s instead", t, p.peekTok.Type),
		p.peekTok.Line,
		p.peekTok.Col,
	)
	return false
}

func (p *Parser) expectCurrent(t string) bool {
	if p.curTok.Type == t {
		return true
	}
	p.addError(
		fmt.Sprintf("Expected current token to be %s, got %s instead", t, p.curTok.Type),
		p.curTok.Line,
		p.curTok.Col,
	)
	return false
}

// PrintErrors sẽ gọi SyntaxError.Error()
func (p *Parser) PrintErrors() {
	if !p.HasErrors() {
		return
	}

	fmt.Println("🚨 PARSER ERRORS:")
	for i, err := range p.errors {
		fmt.Printf("%d. %s\n", i+1, err.Error()) // Gọi phương thức Error() đã override
		fmt.Println(strings.Repeat("─", 60))
	}
}

func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0
}

func (p *Parser) syncTo(syncTokenType string) {
	for p.curTok.Type != syncTokenType && p.curTok.Type != lexer.TOKEN_EOF {
		p.nextToken()
	}
}
