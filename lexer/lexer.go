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
	l.nextChar() // Initialize first character
	return l
}

// NextToken extracts the next token from the input
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	startCol := l.col

	switch l.ch {
	case '.':
		l.nextChar()
		return Token{Type: TOKEN_DOT, Value: ".", Line: l.line, Col: startCol}
	case ',':
		l.nextChar()
		return Token{Type: TOKEN_COMMA, Value: ",", Line: l.line, Col: startCol}
	case '(':
		l.nextChar()
		return Token{Type: TOKEN_LPAREN, Value: "(", Line: l.line, Col: startCol}
	case ')':
		l.nextChar()
		return Token{Type: TOKEN_RPAREN, Value: ")", Line: l.line, Col: startCol}
	case '{':
		l.nextChar()
		return Token{Type: TOKEN_LCURLY, Value: "{", Line: l.line, Col: startCol}
	case '}':
		l.nextChar()
		return Token{Type: TOKEN_RCURLY, Value: "}", Line: l.line, Col: startCol}
	case '[':
		l.nextChar()
		return Token{Type: TOKEN_LSQUARE, Value: "[", Line: l.line, Col: startCol}
	case ']':
		l.nextChar()
		return Token{Type: TOKEN_RSQUARE, Value: "]", Line: l.line, Col: startCol}
	case ';':
		l.nextChar()
		return Token{Type: TOKEN_SEMICOLON, Value: ";", Line: l.line, Col: startCol}
	case '"':
		return Token{Type: TOKEN_STRING, Value: l.readString(), Line: l.line, Col: startCol}
	case '/':
		if l.peekChar() == '/' { // Line comment (//)
			return l.readLineComment()
		} else if l.peekChar() == '*' { // Block comment (/* */)
			return l.readBlockComment()
		}
		// Nếu không phải comment, xử lý như toán tử /
		l.nextChar()
		return Token{Type: TOKEN_ARITHMETIC, Value: "/", Line: l.line, Col: startCol}
	case 0:
		return Token{Type: TOKEN_EOF, Value: "", Line: l.line, Col: startCol}
	default:
		if unicode.IsLetter(l.ch) {
			return l.readKeyword()
		}
		if unicode.IsDigit(l.ch) {
			return l.readNumber()
		}
		return l.readOperator()
	}
}

// matchTwoCharToken handles tokens like ==, !=, >=, <=
func (l *Lexer) matchTwoCharToken(expectedNext rune, singleType, doubleType string, startCol int) Token {
	tok := Token{Type: singleType, Value: string(l.ch), Line: l.line, Col: startCol}
	if l.peekChar() == expectedNext {
		tok = Token{Type: doubleType, Value: string(l.ch) + string(expectedNext), Line: l.line, Col: startCol}
		l.nextChar()
	}
	l.nextChar()
	return tok
}

// nextChar advances the position in the input
func (l *Lexer) nextChar() {
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
		l.nextChar()
	}

	ident := l.input[start:l.position]
	tokType := TOKEN_IDENTIFIER
	if kwType, ok := keywords[ident]; ok {
		tokType = kwType
	}

	return Token{Type: tokType, Value: ident, Line: l.line, Col: startCol}
}

func (l *Lexer) readOperator() Token {

	startCol := l.col

	op := string(l.ch)

	switch l.ch {
	case '+':
		if l.peekChar() == '=' || l.peekChar() == '+' {
			l.nextChar()
			op += string(l.ch)
		}
	case '-':
		if l.peekChar() == '=' || l.peekChar() == '-' {
			l.nextChar()
			op += string(l.ch)
		}
	case '*':
		if l.peekChar() == '=' || l.peekChar() == '*' {
			l.nextChar()
			op += string(l.ch)
		}

	case '/', '%', '=', '!':
		if l.peekChar() == '=' {
			l.nextChar()
			op += string(l.ch)
		}
	case '<':
		if l.peekChar() == '<' || l.peekChar() == '=' {
			l.nextChar()
			op += string(l.ch)
		}
	case '>':
		if l.peekChar() == '>' || l.peekChar() == '=' {
			l.nextChar()
			op += string(l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			l.nextChar()
			op += string(l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			l.nextChar()
			op += string(l.ch)
		}
	}

	l.nextChar()

	tokType, ok := operators[op]
	if !ok {
		return Token{Type: TOKEN_UNKNOWN, Value: op, Line: l.line, Col: startCol}
	}

	return Token{Type: tokType, Value: op, Line: l.line, Col: startCol}
}

// skipWhitespace skips spaces and tabs
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.nextChar()
	}
}

// readNumber reads a number (supports decimals)
func (l *Lexer) readNumber() Token {
	start := l.position
	startCol := l.col

	// Read the integer part
	for unicode.IsDigit(l.ch) {
		l.nextChar()
	}

	// Check for a decimal point
	if l.ch == '.' {
		l.nextChar()

		// Ensure there’s at least one digit after the decimal
		if !unicode.IsDigit(l.ch) {
			return Token{Type: TOKEN_NUMBER, Value: l.input[start:l.position], Line: l.line, Col: startCol}
		}

		// Read the fractional part
		for unicode.IsDigit(l.ch) {
			l.nextChar()
		}
	}

	return Token{Type: TOKEN_NUMBER, Value: l.input[start:l.position], Line: l.line, Col: startCol}
}

// readString reads a string enclosed in quotes and supports escape sequences
func (l *Lexer) readString() string {
	l.nextChar() // Skip opening quote
	//start := l.position
	var strBuilder []rune

	for l.ch != '"' && l.ch != 0 {
		// Xử lý escape sequence
		if l.ch == '\\' {
			l.nextChar()
			switch l.ch {
			case 'n':
				strBuilder = append(strBuilder, '\n')
			case 't':
				strBuilder = append(strBuilder, '\t')
			case '"':
				strBuilder = append(strBuilder, '"')
			case '\\':
				strBuilder = append(strBuilder, '\\')
			default:
				strBuilder = append(strBuilder, '\\', l.ch) // Giữ nguyên nếu escape không hợp lệ
			}
		} else {
			strBuilder = append(strBuilder, l.ch)
		}
		l.nextChar()
	}

	l.nextChar() // Skip closing quote
	return string(strBuilder)
}

// Đọc block comment (/* */)
func (l *Lexer) readBlockComment() Token {
	startPos := l.position
	startLine := l.line
	startCol := l.col

	l.nextChar() // Bỏ qua '*'
	l.nextChar() // Bỏ qua '/' ban đầu

	for {
		if l.ch == 0 { // EOF trước khi đóng comment
			return Token{
				Type:  TOKEN_COMMENT,
				Value: l.input[startPos:l.position],
				Line:  startLine,
				Col:   startCol,
			}
		}
		if l.ch == '*' && l.peekChar() == '/' {
			l.nextChar() // Bỏ qua '*'
			l.nextChar() // Bỏ qua '/'
			break
		}
		l.nextChar()
	}

	return Token{
		Type:  TOKEN_COMMENT,
		Value: l.input[startPos:l.position],
		Line:  startLine,
		Col:   startCol,
	}
}

// Đọc line comment (// đến hết dòng)
func (l *Lexer) readLineComment() Token {
	startPos := l.position
	startLine := l.line
	startCol := l.col

	for l.ch != '\n' && l.ch != 0 {
		l.nextChar()
	}

	return Token{
		Type:  TOKEN_COMMENT,
		Value: l.input[startPos:l.position],
		Line:  startLine,
		Col:   startCol,
	}
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return rune(l.input[l.readPosition])
}
