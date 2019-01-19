package lexer

import (
	"node.go/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
		let a_b 5 + 1
	`
	expected := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "a_b"},
		{token.INT, "5"},
		{token.PLUS, "+"},
		{token.INT, "1"},
	}

	lexer := New(input)

	for _, expectedToken := range expected {
		var actualToken = lexer.NextToken()
		if expectedToken.expectedType != actualToken.Type {
			t.Fatalf("Expected TokenType to be %q, got %q", expectedToken.expectedType, actualToken.Type)
		}
		if expectedToken.expectedLiteral != actualToken.Literal {
			t.Fatalf("Expected TokenType to be %q, got %q", expectedToken.expectedLiteral, actualToken.Literal)
		}
	}
}
