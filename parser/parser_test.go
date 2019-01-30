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

// LITERAL expression
func testLiteralExpression(t *testing.T, actual ast.Expression, expected interface{}) bool {
	switch value := expected.(type) {
	case int64:
		return testIntegerLiteralExpression(t, actual, value)
	case int:
		return testIntegerLiteralExpression(t, actual, int64(value))
	case string:
		return testIdentifier(t, actual, value)
	case bool:
		return testBooleanLiteralExpression(t, actual, value)
	}
	t.Errorf("There is no literal test check for expression of type %q", expected)
	return false
}

// INFIX expression
func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{}, operator string, right interface{}) {

	infixExpression, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("InfixExpression is not *ast.InfixExpression. Got '%q' instead.",
			exp)
	}
	if ok := testLiteralExpression(t, infixExpression.Left, left); !ok {
		t.Fatalf("InfixExpression.Left is not '%s'. Got '%s'",
			left, infixExpression.Left.String())
	}
	if infixExpression.Operator != operator {
		t.Fatalf("InfixExpression.Operator to be '%s'. Got '%s' instead",
			operator, infixExpression.Operator)
	}
	if ok := testLiteralExpression(t, infixExpression.Right, right); !ok {
		t.Fatalf("InfixExpression.Right is not '%q'. Got '%q'",
			right, infixExpression.Right)
	}
}

// INDENT expression
func testIdentifier(t *testing.T, expression ast.Expression, expectedName string) bool {
	identifier, ok := expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression is not *ast.Identifier. Got '%q' instead.", expression)
		return false
	}
	if identifier.Value != expectedName {
		t.Errorf("identifier.Value is not '%s'. Got '%s'", expectedName,
			identifier.Value)
		return false
	}
	if identifier.TokenLiteral() != expectedName {
		t.Fatalf("identifier.TokenLiteral is not '%s'. Got '%s' instead",
			expectedName, identifier.TokenLiteral())
		return false
	}

	return true
}

// BOOLEAN literal
func testBooleanLiteralExpression(t *testing.T, exp ast.Expression, expectedValue bool) bool {
	boolExp, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("expression is not *ast.BooleanLiteral. Got '%q' instead.", boolExp)
		return false
	}
	if boolExp.Value != expectedValue {
		t.Errorf("BooleanLiteral.Value is not %t. Got %t", expectedValue, boolExp.Value)
		return false
	}
	expectedStringLiteral := fmt.Sprintf("%t", expectedValue)
	if boolExp.TokenLiteral() != expectedStringLiteral {
		t.Fatalf("BooleanLiteral.TokenLiteral is not '%s'. Got '%s' instead",
			expectedStringLiteral, boolExp.TokenLiteral())
		return false
	}
	return true
}

// INTEGER literal
func testIntegerLiteralExpression(t *testing.T, e ast.Expression, expectedValue int64) bool {
	integerLiteralExp, ok := e.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression is not *ast.IntegerLiteral. Got '%q' instead.", integerLiteralExp)
		return false
	}
	if integerLiteralExp.Value != expectedValue {
		t.Errorf("IntegerLiteral.Value is not %d. Got %d", expectedValue, integerLiteralExp.Value)
		return false
	}
	if integerLiteralExp.TokenLiteral() != fmt.Sprintf("%d", expectedValue) {
		t.Fatalf("IntegerLiteral.TokenLiteral is not '%d'. Got '%s' instead",
			expectedValue, integerLiteralExp.TokenLiteral())
		return false
	}
	return true
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

func TestIdentifierLiteralExpression(t *testing.T) {
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
	testIdentifier(t, expressionStmt.Expression, "fizzbuzz")
}

func TestBooleanLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
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
		testBooleanLiteralExpression(t, stmt.Expression, test.expected)
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
	testIntegerLiteralExpression(t, stmt.Expression, 2)
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		input            string
		expectedOperator string
		expectedValue    interface{}
	}{
		{"!true", "!", true},
		{"!false", "!", false},
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
		testLiteralExpression(t, prefixExpression.Right, test.expectedValue)
	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input            string
		leftValue        interface{}
		expectedOperator string
		rightValue       interface{}
	}{
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false != true", false, "!=", true},
		{"false == false", false, "==", false},
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
		testInfixExpression(t, stmt.Expression, test.leftValue, test.expectedOperator, test.rightValue)
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"!1 == 2", "((!1) == 2)"},
		{"!1 ^ 2", "((!1) ^ 2)"},
		{"1 + 2 + 3", "((1 + 2) + 3)"},
		{"1 + 2 % 1 * 3 / 2 ^ 6", "(1 + (2 % ((1 * 3) / (2 ^ 6))))"},
		{"1 > 2 >= 3 < 4 <= 5", "((((1 > 2) >= 3) < 4) <= 5)"},
		{"1 + 2 * 3", "(1 + (2 * 3))"},
		{"(1 + 2) * 3", "((1 + 2) * 3)"},
		{"(1 > 2) == false", "((1 > 2) == false)"},
		{"(1 > 2) ^ (2 > 3)", "((1 > 2) ^ (2 > 3))"},
		{"true == (2 == 2)", "(true == (2 == 2))"},
		{"!2 / (1 + 1) > 1", "(((!2) / (1 + 1)) > 1)"},
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

func TestIfExpression(t *testing.T) {
	code := `
		if (z > 1) {z}
	`
	lex := lexer.New(code)
	par := New(lex)
	prg := par.ParseProgram()
	if prg == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	checkParserErrors(t, par)
	if len(prg.Statements) != 1 {
		t.Fatalf("Got and unexpected number of statements: %d': 1. Expected %d",
			len(prg.Statements), 1)
	}
	stmt, ok := prg.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. Got '%q' instead.", stmt)
	}
	ifExp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expression is not *ast.IfExpression. Got '%q' instead.", ifExp)
	}
	// The Condition!
	testInfixExpression(t, ifExp.Condition, "z", ">", 1)

	// The consequence!
	consequence := ifExp.Consequence
	if !ok {
		t.Fatalf("consequence is not *ast.BlockStatement. Got '%q' instead.", consequence)
	}
	if len(consequence.Statements) != 1 {
		t.Fatalf("Expected BlockStatement to have %d statements. Got %d",
			1, len(consequence.Statements))
	}
	consequenceStmt, ok := consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence.Statements[0] is not *ast.ExpressionStatement. Got '%q' instead.",
			consequenceStmt)
	}
	if !testLiteralExpression(t, consequenceStmt.Expression, "z") {
		return
	}
	if ifExp.Alternative != nil {
		t.Fatalf("IfExp.Alternative is not nil. Got %q", ifExp.Alternative)
	}
}

func TestIfWithElseExpression(t *testing.T) {
	code := `
		if (true) {z} else {1}'
	`
	lex := lexer.New(code)
	par := New(lex)
	prg := par.ParseProgram()
	if prg == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	checkParserErrors(t, par)
	if len(prg.Statements) != 1 {
		t.Fatalf("Got and unexpected number of statements: %d': 1. Expected %d",
			len(prg.Statements), 1)
	}
	stmt, ok := prg.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. Got '%q' instead.", stmt)
	}
	ifExp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expression is not *ast.IfExpression. Got '%q' instead.", ifExp)
	}
	// The Condition!
	testLiteralExpression(t, ifExp.Condition, true)
	// The consequence!
	consequence := ifExp.Consequence
	if len(consequence.Statements) != 1 {
		t.Fatalf("Expected BlockStatement to have %d statements. Got %d",
			1, len(consequence.Statements))
	}
	consequenceStmt, ok := consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence.Statements[0] is not *ast.ExpressionStatement. Got '%q' instead.",
			consequenceStmt)
	}
	if !testLiteralExpression(t, consequenceStmt.Expression, "z") {
		return
	}
	if ifExp.Alternative == nil {
		t.Fatalf("IfExpression.Alternative is nil. Got '%q' instead.",
			ifExp.Alternative)
	}
	if len(ifExp.Alternative.Statements) != 1 {
		t.Fatalf("IfExpression.Alternative got %d statements. Expected %d",
			len(ifExp.Alternative.Statements), 1)
	}
	elseFirstStmt, ok := ifExp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("IfExpression.Alternative.Statements[0] is not *ast.ExpressionStatement. Got '%q' instead.",
			elseFirstStmt)
	}
	testLiteralExpression(t, elseFirstStmt.Expression, 1)
}
