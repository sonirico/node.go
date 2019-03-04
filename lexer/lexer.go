package lexer

import "node.go/token"

var WHITESPACES = map[byte]int{
	'\n': 1,
	'\r': 1,
	' ':  1,
	'\t': 1,
}

type Lexer struct {
	currentChar     byte
	currentPosition int64
	nextPosition    int64
	input           string
	inputLength     int64
}

func New(code string) *Lexer {
	lexer := Lexer{
		currentPosition: 0,
		nextPosition:    0,
		input:           code,
		inputLength:     int64(len(code)),
	}
	lexer.readChar()
	return &lexer
}

func (l *Lexer) readChar() {
	if l.nextPosition >= l.inputLength {
		l.currentChar = 0
	} else {
		l.currentChar = l.input[l.nextPosition]
	}

	l.currentPosition = l.nextPosition
	l.nextPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= l.inputLength {
		return 0
	}
	return l.input[l.nextPosition]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.consumeWhitespace()

	switch l.currentChar {
	// Delimiters
	case ',':
		tok = newToken(token.COMMA, l.currentChar)
	case '{':
		tok = newToken(token.LBRACE, l.currentChar)
		break
	case '}':
		tok = newToken(token.RBRACE, l.currentChar)
		break
	case '(':
		tok = newToken(token.LPAREN, l.currentChar)
		break
	case ')':
		tok = newToken(token.RPAREN, l.currentChar)
		break
	case ';':
		tok = newToken(token.SEMICOLON, l.currentChar)
		break
	case '[':
		tok = newToken(token.LBRACKET, l.currentChar)
		break
	case ']':
		tok = newToken(token.RBRACKET, l.currentChar)
		break
	case '"':
		l.readChar()
		tok.Type = token.STRING
		tok.Literal = l.readString()
		break
	// Operators
	case '<':
		{
			if l.peekChar() == '=' {
				ch := l.currentChar
				l.readChar()
				tok.Type = token.LTE
				tok.Literal = string(ch) + string(l.currentChar)
			} else {
				tok = newToken(token.LT, l.currentChar)
			}
			break
		}
	case '>':
		{
			if l.peekChar() == '=' {
				ch := l.currentChar
				l.readChar()
				tok.Type = token.GTE
				tok.Literal = string(ch) + string(l.currentChar)
			} else {
				tok = newToken(token.GT, l.currentChar)
			}
			break
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.currentChar
			l.readChar()
			tok.Type = token.NOT_EQ
			tok.Literal = string(ch) + string(l.currentChar)
		} else {
			tok = newToken(token.BANG, l.currentChar)
		}
		break
	case '/':
		tok = newToken(token.SLASH, l.currentChar)
		break
	case '+':
		tok = newToken(token.PLUS, l.currentChar)
		break
	case '-':
		tok = newToken(token.MINUS, l.currentChar)
		break
	case '*':
		tok = newToken(token.ASTERISK, l.currentChar)
		break
	case '%':
		tok = newToken(token.PERCENT, l.currentChar)
		break
	case '^':
		tok = newToken(token.POWER, l.currentChar)
		break
	case '=':
		{
			ch := l.currentChar
			switch l.peekChar() {
			case '=':
				l.readChar()
				tok.Type = token.EQ
				tok.Literal = string(ch) + string(l.currentChar)
				break
			default:
				tok = newToken(token.ASSIGNMENT, l.currentChar)
			}
		}
		break
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		break
	default:
		if isDigit(l.currentChar) {
			return token.Token{Type: token.INT, Literal: l.readNumber()}
		} else if isLetter(l.currentChar) {
			tokenLiteral := l.readWord()
			tokenType := token.LookupKeyword(tokenLiteral)
			return token.Token{Type: tokenType, Literal: tokenLiteral}
		}
		tok = newToken(token.ILLEGAL, l.currentChar)
	}
	l.readChar()
	return tok
}

func (l *Lexer) consumeWhitespace() {
	for isWhiteSpace(l.currentChar) {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	pos := l.currentPosition
	for isDigit(l.currentChar) {
		l.readChar()
	}
	return l.input[pos:l.currentPosition]
}

func (l *Lexer) readWord() string {
	pos := l.currentPosition
	for isLetter(l.currentChar) {
		l.readChar()
	}
	return l.input[pos:l.currentPosition]
}

func (l *Lexer) readString() string {
	pos := l.currentPosition

	for l.currentChar != '"' {
		l.readChar()
	}

	return l.input[pos:l.currentPosition]
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isLetter(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '_'
}

func isWhiteSpace(char byte) bool {
	if _, ok := WHITESPACES[char]; ok {
		return true
	}
	return false
}

func newToken(tokenType token.TokenType, literal byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(literal)}
}
