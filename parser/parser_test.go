package parser

import (
	"node.go/ast"
	"node.go/lexer"
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	if len(p.Errors()) < 1 {
		return
	}

	for _, errorMessage := range p.Errors() {
		t.Errorf("Parser error: %s", errorMessage)
	}

	t.FailNow()
}

// LET testing
func testLetStatement(t *testing.T, actualStmt ast.Statement, expectedIdentName string) bool {
	if actualStmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral wasn't 'let'. Got %q", actualStmt.TokenLiteral())
		return false
	}
	letStmt, ok := actualStmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt wasn't LetStatement, got %s", actualStmt)
		return false
	}
	if letStmt.Name.Value != expectedIdentName {
		t.Errorf("LetStmt.Name.Value wasn't '%s'. Got '%s'",
			expectedIdentName, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != expectedIdentName {
		t.Errorf("LetStmt.Name.TokenLiteral wasn't '%s'. Got '%s'",
			expectedIdentName, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}

func TestLetStatements(t *testing.T) {
	code := `
		let foo;
		let bar = foo;
		let foobar;
	`
	lex := lexer.New(code)
	par := New(lex)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Expected statements: 3. Got %d", len(program.Statements))
	}
	expected := []struct {
		Name string
	}{
		{"foo"},
		{"bar"},
		{"foobar"},
	}
	for index, expectedIdentifier := range expected {
		stmt := program.Statements[index]
		testLetStatement(t, stmt, expectedIdentifier.Name)
	}
}

func TestReturnStatement(t *testing.T) {
	code := `
		return 1;
		return 3 * 4;
		return function (x) {x};
	`
	lex := lexer.New(code)
	par := New(lex)
	prog := par.ParseProgram()
	checkParserErrors(t, par)
	if prog == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(prog.Statements) != 3 {
		t.Fatalf("Expected 'ReturnStatements': 3. Got %d", len(prog.Statements))
	}
	for _, stmt := range prog.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt is not an *ast.ReturnStatement. Got '%q' instead", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStatment.TokenLiteral is not 'return'. Got '%s' instead",
				returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "fizzbuzz;"
	lex := lexer.New(input)
	par := New(lex)
	prog := par.ParseProgram()
	checkParserErrors(t, par)
	if prog == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(prog.Statements) != 1 {
		t.Fatalf("Expected 'Identifiers': 1. Got %d", len(prog.Statements))
	}
	stmt := prog.Statements[0]
	expressionStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. Got '%q' instead.", stmt)
	}
	identifier, ok := expressionStmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt is not *ast.Identifier. Got '%q' instead.", expressionStmt)
	}
	if identifier.Value != "fizzbuzz" {
		t.Errorf("identifier.Value is not %s. Got %s", "fizzbuzz", identifier.Value)
	}
	if identifier.TokenLiteral() != "fizzbuzz" {
		t.Fatalf("expressionStmt.TokenLiteral is not 'fizzbuzz'. Got '%s' instead",
			identifier.TokenLiteral())
	}
}
