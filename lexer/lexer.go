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
	if l.currentPosition >= l.inputLength {
		l.currentChar = 0
	} else {
		l.currentChar = l.input[l.nextPosition]
	}

	l.currentPosition = l.nextPosition
	l.nextPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.consumeWhitespace()
	switch l.currentChar {
	case '+':
		tok = newToken(token.PLUS, l.currentChar)
		break
	case '-':
		tok = newToken(token.MINUS, l.currentChar)
		break
	case '=':
		tok = newToken(token.ASSIGNMENT, l.currentChar)
		break
	default:
		if isDigit(l.currentChar) {
			return token.Token{Type: token.INT, Literal: l.readNumber()}
		}
		if isLetter(l.currentChar) {
			tokenLiteral := l.readWord()
			tokenType := token.LookupKeyword(tokenLiteral)
			return token.Token{Type: tokenType, Literal: tokenLiteral}
		}
		tok = newToken(token.INT, l.currentChar)
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
