package ast

import (
	"node.go/token"
	"testing"
)

func testProgramString(t *testing.T, actual *Program, expectedRpr string) {
	if expectedRpr != actual.String() {
		t.Errorf("Unexpected program.String(). Expected '%s'. Got '%s'.",
			expectedRpr, actual.String())
	}
}

func TestLetStatement_String(t *testing.T) {
	expected := "let variable = 1;"
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
	testProgramString(t, prog, expected)
}

func TestIfExpression_String(t *testing.T) {
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
					Consequence: &BlockStatement{
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

func TestFunctionLiteral_String(t *testing.T) {
	expected := "fn() {}"
	prg := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Token: token.Token{
					Type:    token.FUNC,
					Literal: "fn",
				},
				Expression: &FunctionLiteral{
					Token: token.Token{
						Type:    token.FUNC,
						Literal: "fn",
					},
					Parameters: []*Identifier{},
					Body: &BlockStatement{
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
	testProgramString(t, prg, expected)
}

func TestCallExpression_String(t *testing.T) {
	expected := "sum(a, 2)"
	prg := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Token: token.Token{
					Type:    token.FUNC,
					Literal: "sum",
				},
				Expression: &CallExpression{
					Token: token.Token{
						Type:    token.FUNC,
						Literal: "sum",
					},
					Function: &Identifier{
						Token: token.Token{
							Type:    token.IDENTIFIER,
							Literal: "sum",
						},
						Value: "sum",
					},
					Arguments: []Expression{
						&Identifier{
							Token: token.Token{
								Type:    token.IDENTIFIER,
								Literal: "a",
							},
							Value: "a",
						},
						&Identifier{
							Token: token.Token{
								Type:    token.IDENTIFIER,
								Literal: "2",
							},
							Value: "2",
						},
					},
				},
			},
		},
	}
	testProgramString(t, prg, expected)
}
