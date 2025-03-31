package lexer

// Token types
const (
	TOKEN_EOF        = "EOF"
	TOKEN_IDENTIFIER = "IDENTIFIER"
	TOKEN_NUMBER     = "NUMBER"
	TOKEN_STRING     = "STRING"
	TOKEN_KEYWORD    = "KEYWORD"
	TOKEN_OPERATOR   = "OPERATOR"
	TOKEN_ASSIGN     = "ASSIGN"
	TOKEN_LPAREN     = "LPAREN"
	TOKEN_RPAREN     = "RPAREN"
	TOKEN_COMMA      = "COMMA"
	TOKEN_DOT        = "DOT"
	TOKEN_BOOLEAN    = "BOOLEAN"
	TOKEN_LOGICAL    = "LOGICAL"
	TOKEN_BITWISE    = "BITWISE"
	TOKEN_LCURLY     = "LCURLY"
	TOKEN_RCURLY     = "RCURLY"
	TOKEN_LSQUARE    = "LSQUARE"
	TOKEN_RSQUARE    = "RSQUARE"
	TOKEN_SEMICOLON  = "SEMICOLON"
)

// Keywords in Pun
var keywords = map[string]string{
	"ask":       TOKEN_KEYWORD,
	"shout":     TOKEN_KEYWORD,
	"when":      TOKEN_KEYWORD,
	"maybe":     TOKEN_KEYWORD,
	"otherwise": TOKEN_KEYWORD,
	"stop":      TOKEN_KEYWORD,
	"continue":  TOKEN_KEYWORD,
	"return":    TOKEN_KEYWORD,
	"for":       TOKEN_KEYWORD,
	"while":     TOKEN_KEYWORD,
	"until":     TOKEN_KEYWORD,
	"nothing":   TOKEN_KEYWORD,
	"cook":      TOKEN_KEYWORD,
	"true":      TOKEN_BOOLEAN,
	"false":     TOKEN_BOOLEAN,
}

// Token structure
type Token struct {
	Type  string
	Value string
	Line  int
	Col   int
}
