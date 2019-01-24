package parser

import (
	"fmt"
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

// INTEGER literal testing
func testIntegerLiteralExpression(t *testing.T, e ast.Expression, expectedValue int64) {
	integerLiteralExp, ok := e.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt is not *ast.IntegerLiteral. Got '%q' instead.", integerLiteralExp)
	}
	if integerLiteralExp.Value != expectedValue {
		t.Errorf("IntegerLiteral.Value is not %d. Got %s", expectedValue, integerLiteralExp)
	}
	if integerLiteralExp.TokenLiteral() != fmt.Sprintf("%d", expectedValue) {
		t.Fatalf("IntegerLiteral.TokenLiteral is not '%d'. Got '%s' instead",
			expectedValue, integerLiteralExp.TokenLiteral())
	}
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

func TestIntegerLiteralExpression(t *testing.T) {
	code := "2;"
	lex := lexer.New(code)
	par := New(lex)
	program := par.ParseProgram()
	checkParserErrors(t, par)
	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("Got and unexpected number of statements: %d': 1. Expected %d",
			len(program.Statements), 1)
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. Got '%q' instead.", stmt)
	}
	integerLiteralExp, ok := stmt.Expression.(*ast.IntegerLiteral)
	testIntegerLiteralExpression(t, integerLiteralExp, 2)
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		input            string
		expectedOperator string
		expectedValue    int64
	}{
		{"!1", "!", 1},
		{"-4", "-", 4},
	}

	for _, test := range tests {
		lex := lexer.New(test.input)
		par := New(lex)
		program := par.ParseProgram()
		checkParserErrors(t, par)
		if program == nil {
			t.Fatalf("ParseProgram returned nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("Got and unexpected number of statements: %d': 1. Expected %d",
				len(program.Statements), 1)
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not *ast.ExpressionStatement. Got '%q' instead.", stmt)
		}
		prefixExpression, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("prefixExpression is not *ast.PrefixExpression. Got '%q' instead.", stmt)
		}
		if prefixExpression.Operator != test.expectedOperator {
			t.Fatalf("Expected PrefixExpression.Operator to be '%s'. Got '%s' instead",
				test.expectedOperator, prefixExpression.Operator)
		}
		testIntegerLiteralExpression(t, prefixExpression.Right, test.expectedValue)
	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input            string
		leftValue        int64
		expectedOperator string
		rightValue       int64
	}{
		{"1 + 2", 1, "+", 2},
		{"1 - 2", 1, "-", 2},
		{"1 * 2", 1, "*", 2},
		{"1 / 2", 1, "/", 2},
		{"1 % 2", 1, "%", 2},
		{"1 ^ 2", 1, "^", 2},
		{"1 < 2", 1, "<", 2},
		{"1 > 2", 1, ">", 2},
		{"1 >= 2", 1, ">=", 2},
		{"1 <= 2", 1, "<=", 2},
		{"1 == 2", 1, "==", 2},
		{"1 != 2", 1, "!=", 2},
	}

	for _, test := range tests {
		lex := lexer.New(test.input)
		par := New(lex)
		program := par.ParseProgram()
		checkParserErrors(t, par)
		if program == nil {
			t.Fatalf("ParseProgram returned nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("Got and unexpected number of statements: %d': 1. Expected %d",
				len(program.Statements), 1)
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not *ast.ExpressionStatement. Got '%q' instead.", stmt)
		}
		infixExpression, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("infixExpression is not *ast.InfixExpression. Got '%q' instead.", stmt)
		}
		testIntegerLiteralExpression(t, infixExpression.Left, test.leftValue)
		if infixExpression.Operator != test.expectedOperator {
			t.Fatalf("Expected InfixExpression.Operator to be '%s'. Got '%s' instead",
				test.expectedOperator, infixExpression.Operator)
		}
		testIntegerLiteralExpression(t, infixExpression.Right, test.rightValue)
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"1 + 2 + 3", "((1 + 2) + 3)"},
		{"1 + 2 % 1 * 3 / 2 ^ 6", "(1 + (2 % ((1 * 3) / (2 ^ 6))))"},
		{"1 > 2 >= 3 < 4 <= 5", "((((1 > 2) >= 3) < 4) <= 5)"},
	}

	for _, test := range tests {
		lex := lexer.New(test.input)
		par := New(lex)
		prog := par.ParseProgram()
		if prog == nil {
			t.Fatalf("ParseProgram returned nil")
		}
		checkParserErrors(t, par)
		if test.expectedOutput != prog.String() {
			t.Fatalf("Error when parsing precedence. Expected '%s'. Got '%s'",
				test.expectedOutput, prog.String())
		}
	}
}
