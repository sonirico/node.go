package parser

import (
	"node.go/ast"
	"testing"
)

func TestHashLiteralStringKeys(t *testing.T) {
	tests := []struct {
		code     string
		expected map[string]int64
	}{
		{
			`{"key": 0, "hello": 3, "negative": 1}`,
			map[string]int64{
				"key":      0,
				"hello":    3,
				"negative": 1,
			},
		},
	}

	for _, test := range tests {
		program := ParseTesting(t, test.code)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ExpressionStatement. Got %T(%+v)",
				program.Statements[0], program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.HashLiteral)
		if !ok {
			t.Fatalf("exp is not HashLiteral. Got %T(%+v)",
				stmt.Expression, stmt.Expression)
		}
		if len(test.expected) != len(exp.Pairs) {
			t.Fatalf("HashLiteral has an unexpected number of items. Expected %d, Got %d",
				len(test.expected), len(exp.Pairs))
		}
		for k, v := range exp.Pairs {
			literal, ok := k.(*ast.StringLiteral)
			if !ok {
				t.Fatalf("HashLiteral key is not StringLiteral'")
			}
			stringL := literal.String()
			expected := test.expected[stringL]
			testIntegerLiteralExpression(t, v, expected)
		}
	}
}

func TestParseEmptyHashLiteralExpression(t *testing.T) {
	payload := `{}`
	program := ParseTesting(t, payload)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt expected to be ExpressionStatement. Got %T(%+v)",
			program.Statements[0], program.Statements[0])
	}
	hashLiteral, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp expected to be HashLiteral. Got %T(%+v)",
			stmt.Expression, stmt.Expression)
	}
	if len(hashLiteral.Pairs) != 0 {
		t.Fatalf("hashLiteral expected to have %d items. Got %d",
			0, len(hashLiteral.Pairs))
	}
}

func TestParseHashLiteralWithInfixExpressions(t *testing.T) {
	// {"keys": 1 + 1}
}
