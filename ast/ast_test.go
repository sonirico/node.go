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

func TestIfExpressionToString(t *testing.T) {
	prg := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Token: token.Token{
					Type:    token.IF,
					Literal: "if",
				},
				Expression: &IfExpression{
					Token: token.Token{
						Type:    token.IF,
						Literal: "if",
					},
					Condition: &BooleanLiteral{
						Token: token.Token{
							Type:    token.TRUE,
							Literal: "true",
						},
						Value: true,
					},
					Consequence: BlockStatement{
						Token: token.Token{
							Type:    token.LBRACE,
							Literal: "{",
						},
						Statements: []Statement{
							&ExpressionStatement{
								Token: token.Token{
									Type:    token.FALSE,
									Literal: "false",
								},
								Expression: &BooleanLiteral{
									Token: token.Token{
										Type:    token.FALSE,
										Literal: "false",
									},
									Value: false,
								},
							},
						},
					},
					Alternative: &BlockStatement{
						Token: token.Token{
							Type:    token.LBRACE,
							Literal: "{",
						},
						Statements: []Statement{},
					},
				},
			},
		},
	}
	// Do note that empty alternatives (else) result in them being totally ignored
	expected := "if (true) {false}"
	if expected != prg.String() {
		t.Errorf("Unexpected program.String(). Expected '%s'. Got '%s'.",
			expected, prg.String())
	}
}
