package ast

import (
	"node.go/token"
	"testing"
)

func TestAstString(t *testing.T) {
	prog := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{
					Type:    token.LET,
					Literal: "let",
				},
				Name: &Identifier{
					Token: token.Token{
						Type:    token.IDENTIFIER,
						Literal: "variable",
					},
					Value: "variable",
				},
				Value: &Identifier{
					Token: token.Token{
						Type:    token.INT,
						Literal: "1",
					},
					Value: "1",
				},
			},
		},
	}
	expected := "let variable = 1;"
	if expected != prog.String() {
		t.Errorf("Unexpected program.String(). Expected '%s'. Got '%s'.",
			expected, prog.String())
	}
}
