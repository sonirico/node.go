package parser

import (
	"node.go/ast"
	"node.go/lexer"
	"testing"
)

// LET testing
func testLetStatement(t *testing.T, actualStmt ast.Statement, expectedIdentName string) bool {
	if actualStmt.Literal() != "let" {
		t.Errorf("stmt.Literal wasn't 'let'. Got %q", actualStmt.Literal())
		return false
	}
	letStmt, ok := actualStmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt wasn't LetStatement, got %s", actualStmt)
		return false
	}
	if letStmt.Name.Value != expectedIdentName {
		t.Errorf("LetStmt.Name.Value wasn't '%s'. Got '%s'", expectedIdentName, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.Literal() != expectedIdentName {
		t.Errorf("LetStmt.Name.Literal wasn't '%s'. Got '%s'", expectedIdentName, letStmt.Name.Literal())
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
