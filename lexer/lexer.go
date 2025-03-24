package lexer

import "unicode"

// Lexer structure
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune
	line         int
	col          int
}

// NewLexer creates a new lexer
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1, col: 0}
	l.readChar() // Initialize first character
	return l
}

// NextToken extracts the next token from the input
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	startCol := l.col

	switch l.ch {
	case '=':
		return l.matchTwoCharToken('=', TOKEN_ASSIGN, TOKEN_OPERATOR, startCol)
	case '!', '<', '>':
		return l.matchTwoCharToken('=', TOKEN_OPERATOR, TOKEN_OPERATOR, startCol)
	case ',':
		l.readChar()
		return Token{Type: TOKEN_COMMA, Value: ",", Line: l.line, Col: startCol}
	case '(':
		l.readChar()
		return Token{Type: TOKEN_LPAREN, Value: "(", Line: l.line, Col: startCol}
	case ')':
		l.readChar()
		return Token{Type: TOKEN_RPAREN, Value: ")", Line: l.line, Col: startCol}
	case '"':
		return Token{Type: TOKEN_STRING, Value: l.readString(), Line: l.line, Col: startCol}
	case 0:
		return Token{Type: TOKEN_EOF, Value: "", Line: l.line, Col: startCol}
	default:
		if unicode.IsLetter(l.ch) {
			return l.readKeyword()
		}
		if unicode.IsDigit(l.ch) {
			return l.readNumber()
		}
		ch := l.ch
		l.readChar()
		return Token{Type: TOKEN_OPERATOR, Value: string(ch), Line: l.line, Col: startCol}
	}
}

// matchTwoCharToken handles tokens like ==, !=, >=, <=
func (l *Lexer) matchTwoCharToken(expectedNext rune, singleType, doubleType string, startCol int) Token {
	tok := Token{Type: singleType, Value: string(l.ch), Line: l.line, Col: startCol}
	if l.readPosition < len(l.input) && rune(l.input[l.readPosition]) == expectedNext {
		tok = Token{Type: doubleType, Value: string(l.ch) + string(expectedNext), Line: l.line, Col: startCol}
		l.readChar()
	}
	l.readChar()
	return tok
}

// readChar advances the position in the input
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		l.ch = rune(l.input[l.readPosition])
	}

	if l.ch == '\n' {
		l.line++
		l.col = 1 // Reset column for new line
	} else {
		l.col++ // Move column forward
	}

	l.position = l.readPosition
	l.readPosition++
}

// readKeyword reads an identifier or keyword or boolean
func (l *Lexer) readKeyword() Token {
	start := l.position
	startCol := l.col

	for unicode.IsLetter(l.ch) {
		l.readChar()
	}

	ident := l.input[start:l.position]
	tokType := TOKEN_IDENTIFIER
	if kwType, ok := keywords[ident]; ok {
		tokType = kwType
	}

	return Token{Type: tokType, Value: ident, Line: l.line, Col: startCol}
}

// skipWhitespace skips spaces and tabs
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

// readNumber reads a number (supports decimals)
func (l *Lexer) readNumber() Token {
	start := l.position
	startCol := l.col

	// Read the integer part
	for unicode.IsDigit(l.ch) {
		l.readChar()
	}

	// Check for a decimal point
	if l.ch == '.' {
		l.readChar()

		// Ensure there’s at least one digit after the decimal
		if !unicode.IsDigit(l.ch) {
			return Token{Type: TOKEN_NUMBER, Value: l.input[start:l.position], Line: l.line, Col: startCol}
		}

		// Read the fractional part
		for unicode.IsDigit(l.ch) {
			l.readChar()
		}
	}

	return Token{Type: TOKEN_NUMBER, Value: l.input[start:l.position], Line: l.line, Col: startCol}
}

// readString reads a string enclosed in quotes
func (l *Lexer) readString() string {
	l.readChar() // Skip opening quote
	start := l.position

	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}

	str := l.input[start:l.position]
	l.readChar() // Skip closing quote
	return str
}
