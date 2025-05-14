package lexer

// Token types
const (
	TOKEN_EOF        = "EOF"
	TOKEN_IDENTIFIER = "IDENTIFIER"
	TOKEN_NUMBER     = "NUMBER"
	TOKEN_STRING     = "STRING"
	TOKEN_KEYWORD    = "KEYWORD"
	TOKEN_ASSIGN     = "ASSIGN"
	TOKEN_LPAREN     = "LPAREN"
	TOKEN_RPAREN     = "RPAREN"
	TOKEN_COMMA      = "COMMA"
	TOKEN_DOT        = "DOT"
	TOKEN_BOOLEAN    = "BOOLEAN"
	TOKEN_LCURLY     = "LCURLY"
	TOKEN_RCURLY     = "RCURLY"
	TOKEN_LSQUARE    = "LSQUARE"
	TOKEN_RSQUARE    = "RSQUARE"
	TOKEN_SEMICOLON  = "SEMICOLON"
	TOKEN_COMMENT    = "COMMENT"
	TOKEN_NOTHING    = "NOTHING"
	TOKEN_ARITHMETIC = "ARITHMETIC" // + - * / % **
	TOKEN_COMPARISON = "COMPARISON" // == != > < >= <=
	TOKEN_LOGICAL    = "LOGICAL"    // && || !
	TOKEN_BITWISE    = "BITWISE"    // & | ^ ~ << >>
	TOKEN_UNKNOWN    = "UNKNOWN"
)

// Keywords in Pun
var keywords = map[string]string{
	"if":       TOKEN_KEYWORD,
	"elif":     TOKEN_KEYWORD,
	"else":     TOKEN_KEYWORD,
	"break":    TOKEN_KEYWORD,
	"continue": TOKEN_KEYWORD,
	"return":   TOKEN_KEYWORD,
	"for":      TOKEN_KEYWORD,
	"while":    TOKEN_KEYWORD,
	"until":    TOKEN_KEYWORD,
	"func":     TOKEN_KEYWORD,
	"true":     TOKEN_BOOLEAN,
	"false":    TOKEN_BOOLEAN,
	"nothing":  TOKEN_NOTHING,
}

var operators = map[string]string{
	//Gán
	"=": TOKEN_ASSIGN,

	// Số học
	"+":  TOKEN_ARITHMETIC,
	"-":  TOKEN_ARITHMETIC,
	"*":  TOKEN_ARITHMETIC,
	"/":  TOKEN_ARITHMETIC,
	"%":  TOKEN_ARITHMETIC,
	"**": TOKEN_ARITHMETIC,

	// So sánh
	"==": TOKEN_COMPARISON,
	"!=": TOKEN_COMPARISON,
	">":  TOKEN_COMPARISON,
	"<":  TOKEN_COMPARISON,
	">=": TOKEN_COMPARISON,
	"<=": TOKEN_COMPARISON,

	// Logic
	"&&": TOKEN_LOGICAL,
	"||": TOKEN_LOGICAL,
	"!":  TOKEN_LOGICAL,

	// Bitwise (nếu cần)
	"&":  TOKEN_BITWISE,
	"|":  TOKEN_BITWISE,
	"^":  TOKEN_BITWISE,
	"~":  TOKEN_BITWISE,
	"<<": TOKEN_BITWISE,
	">>": TOKEN_BITWISE,
}

// Token structure
type Token struct {
	Type  string
	Value string
	Line  int
	Col   int
}
