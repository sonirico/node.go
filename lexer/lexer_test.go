package lexer

import (
	"node.go/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
		let a_b = 5 + 10;
		!/*%^
		(1 + 2) * 3 == 9;
		0 != 9
		4 > 3 >= 3 < 2 <= 2;
		if (true != false) {return 1}
		fn(x, y, z){};
		sum(1, 2)
		""
		"I am a teapot"
		[1, 3, true]
		{"key": 0}
	`
	expected := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "a_b"},
		{token.ASSIGNMENT, "="},
		{token.INT, "5"},
		{token.PLUS, "+"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.PERCENT, "%"},
		{token.POWER, "^"},
		{token.LPAREN, "("},
		{token.INT, "1"},
		{token.PLUS, "+"},
		{token.INT, "2"},
		{token.RPAREN, ")"},
		{token.ASTERISK, "*"},
		{token.INT, "3"},
		{token.EQ, "=="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.INT, "0"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.INT, "4"},
		{token.GT, ">"},
		{token.INT, "3"},
		{token.GTE, ">="},
		{token.INT, "3"},
		{token.LT, "<"},
		{token.INT, "2"},
		{token.LTE, "<="},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.TRUE, "true"},
		{token.NOT_EQ, "!="},
		{token.FALSE, "false"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.INT, "1"},
		{token.RBRACE, "}"},
		{token.FUNC, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "z"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.IDENTIFIER, "sum"},
		{token.LPAREN, "("},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RPAREN, ")"},
		{token.STRING, ""},
		{token.STRING, "I am a teapot"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "3"},
		{token.COMMA, ","},
		{token.TRUE, "true"},
		{token.RBRACKET, "]"},
		{token.LBRACE, "{"},
		{token.STRING, "key"},
		{token.COLON, ":"},
		{token.INT, "0"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for _, expectedToken := range expected {
		actualToken := lexer.NextToken()
		if expectedToken.expectedType != actualToken.Type {
			t.Fatalf("Expected TokenType to be %q, got %q", expectedToken.expectedType, actualToken.Type)
		}
		if expectedToken.expectedLiteral != actualToken.Literal {
			t.Fatalf("Expected TokenType to be %q, got %q", expectedToken.expectedLiteral, actualToken.Literal)
		}
	}
}
