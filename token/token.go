package token

const (
	// MISC
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	// Identifiers
	IDENTIFIER = "ident"

	// Literals
	INT    = "int"
	STRING = "string"
	TRUE   = "true"
	FALSE  = "false"

	// operators
	PLUS     = "+"
	MINUS    = "-"
	EQ       = "=="
	NOT_EQ   = "!="
	GT       = ">"
	LT       = "<"
	GTE      = ">="
	LTE      = "<="
	BANG     = "!"
	SLASH    = "/"
	ASTERISK = "*"
	PERCENT  = "%"
	POWER    = "^"

	//
	ASSIGNMENT = "="

	// keywords
	VAR    = "var"
	CONST  = "const"
	LET    = "let"
	FUNC   = "function"
	IF     = "if"
	ELSE   = "else"
	RETURN = "return"

	// Delimiters
	COMMA     = ","
	COLON     = ":"
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"
)

type TokenType string

var keywords = map[string]TokenType{
	"fn":       FUNC,
	"function": FUNC,
	"const":    CONST,
	"let":      LET,
	"var":      VAR,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"true":     TRUE,
	"false":    FALSE,
}

func LookupKeyword(literal string) TokenType {
	if tt, ok := keywords[literal]; ok {
		return tt
	}
	return IDENTIFIER
}

type Token struct {
	Type    TokenType
	Literal string
}

func New(tokenType TokenType, literal string) *Token {
	return &Token{Type: tokenType, Literal: literal}
}
