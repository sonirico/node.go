package parser

import (
	"fmt"
	"node.go/ast"
	"node.go/lexer"
	"testing"
)

func ParseTesting(t *testing.T, code string) *ast.Program {
	lex := lexer.New(code)
	par := New(lex)
	program := par.ParseProgram()
	checkParserErrors(t, par)
	return program
}

func checkProgram(t *testing.T, p *ast.Program) {
	if p == nil {
		t.Fatalf("Parser.ParseProgram returned nil")
	}
}

func checkProgramStatements(t *testing.T, p *ast.Program, expectedStatementsLength int) {
	checkProgram(t, p)
	if expectedStatementsLength != len(p.Statements) {
		t.Fatalf("Expected program to have %d statements. Got %d.",
			expectedStatementsLength, len(p.Statements))
	}
}

func expectAnyParserErrors(t *testing.T, p *Parser) {
	if len(p.Errors()) < 1 {
		t.Fatalf("Expected parser to have errors.")
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	if len(p.Errors()) < 1 {
		return
	}

	for _, errorMessage := range p.Errors() {
		t.Errorf("Parser error: %s", errorMessage)
	}

	t.FailNow()
}

func testExpressionStatement(t *testing.T, stmt ast.Statement) *ast.ExpressionStatement {
	expStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. Got '%q' instead.", stmt)
	}

	return expStmt
}

// FUNCTION expression

func testExpressionIsFunctionLiteral(t *testing.T, exp ast.Expression) *ast.FunctionLiteral {
	funcExp, ok := exp.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("Expression is not FunctionLiteral. Got %q", exp)
	}
	return funcExp
}

func testFunctionExpressionParameters(t *testing.T, funcExp *ast.FunctionLiteral, idents []string) {
	if len(idents) != len(funcExp.Parameters) {
		t.Fatalf("Expected parameters cardinal to be %d, Got %d",
			len(idents), len(funcExp.Parameters))
	}
	for index, identName := range idents {
		if identName != funcExp.Parameters[index].TokenLiteral() {
			t.Fatalf("FunctionLiteral does not have param %s at position %d. Got %s",
				identName, index, funcExp.Parameters[index].TokenLiteral())
		}
	}
}

func testFunctionExpressionBody(
	t *testing.T,
	funcExp *ast.FunctionLiteral,
	expectedBodyStatementsLength int) {
	if len(funcExp.Body.Statements) != expectedBodyStatementsLength {
		t.Fatalf("Expected body statements cardinal to be %d, Got %d",
			expectedBodyStatementsLength, len(funcExp.Parameters))
	}
}

func testFunctionLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	idents []string,
	bodyStmts int) *ast.FunctionLiteral {

	funcExp := testExpressionIsFunctionLiteral(t, exp)
	testFunctionExpressionParameters(t, funcExp, idents)
	testFunctionExpressionBody(t, funcExp, bodyStmts)
	return funcExp
}

// LITERAL expression
func testLiteralExpression(t *testing.T, actual ast.Expression, expected interface{}) bool {
	switch value := expected.(type) {
	case int64:
		return testIntegerLiteralExpression(t, actual, value)
	case int:
		return testIntegerLiteralExpression(t, actual, int64(value))
	case string:
		{
			switch actual.(type) {
			case *ast.StringLiteral:
				return testStringLiteral(t, actual, value)
			default:
				return testIdentifier(t, actual, value)
			}
		}
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

func testLetStatement(t *testing.T, actualStmt ast.Statement, expected struct {
	Name  string
	Value interface{}
}) bool {
	if actualStmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral wasn't 'let'. Got %q", actualStmt.TokenLiteral())
		return false
	}
	letStmt, ok := actualStmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt wasn't LetStatement, got %s", actualStmt)
		return false
	}
	if letStmt.Name.Value != expected.Name {
		t.Errorf("LetStmt.Name.Value wasn't '%s'. Got '%s'",
			expected.Name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != expected.Name {
		t.Errorf("LetStmt.Name.TokenLiteral wasn't '%s'. Got '%s'",
			expected.Name, letStmt.Name.TokenLiteral())
		return false
	}
	if expected.Value != nil && !testLiteralExpression(t, letStmt.Value, expected.Value) {
		t.Errorf("LetStatement.Value is wrong. Expected %q. Got %q",
			expected.Value, letStmt.Value)
	}
	return true
}

func TestLetStatements(t *testing.T) {
	code := `
		let foo;
		let bar = true;
		let foobar = bar;
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
		Name  string
		Value interface{}
	}{
		{"foo", nil},
		{"bar", true},
		{"foobar", "bar"},
	}
	for index, expectedIdentifier := range expected {
		stmt := program.Statements[index]
		testLetStatement(t, stmt, expectedIdentifier)
	}
}

func testReturnStatement(t *testing.T, stmt ast.Statement, expected interface{}) {
	returnStmt, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("stmt is not an *ast.ReturnStatement. Got '%q' instead", stmt)
	}
	if returnStmt.TokenLiteral() != "return" {
		t.Fatalf("ReturnStatment.TokenLiteral is not 'return'. Got '%s' instead",
			returnStmt.TokenLiteral())
	}
	if expected == nil {
		if returnStmt.ReturnValue != nil {
			t.Fatalf("ReturnStatement.ReturnValue is not nil. Got '%q'",
				returnStmt.ReturnValue)
		}
	} else if !testLiteralExpression(t, returnStmt.ReturnValue, expected) {
		t.Fatalf("ReturnStatement.ReturnValue is wrong. Expected %q. Got %q",
			expected, returnStmt.ReturnValue)
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral interface{}
	}{
		{"return;", nil},
		{"return 3", 3},
		{"return theadventofcode", "theadventofcode"},
	}
	for _, test := range tests {
		lex := lexer.New(test.input)
		par := New(lex)
		prg := par.ParseProgram()
		checkParserErrors(t, par)
		checkProgramStatements(t, prg, 1)
		stmt := prg.Statements[0]
		testReturnStatement(t, stmt, test.expectedLiteral)
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

func testStringLiteral(t *testing.T, exp ast.Expression, expected string) bool {
	stringExp, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Errorf("Expected expression to be StringLiteral, got %q instead", exp)
		return false
	}
	if stringExp.Value != expected {
		t.Errorf("StringLiteral.Value is not '%s'. Got '%s'", expected, stringExp.Value)
		return false
	}
	return true
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"I am a fork";1;`
	program := ParseTesting(t, input)
	if len(program.Statements) < 1 {
		t.Fatalf("Parse string got no statements, expected %d", 1)
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. Got '%q' instead.", stmt)
	}
	testStringLiteral(t, stmt.Expression, "I am a fork")
}

func testPrefixExpression(t *testing.T, expression ast.Expression, expected interface{}, operator string) bool {
	prefixExpression, ok := expression.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("prefixExpression is not *ast.PrefixExpression. Got '%q' instead.", expression)
	}
	if prefixExpression.Operator != operator {
		t.Fatalf("Expected PrefixExpression.Operator to be '%s'. Got '%s' instead",
			operator, prefixExpression.Operator)
	}
	return testLiteralExpression(t, prefixExpression.Right, expected)
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
		testPrefixExpression(t, stmt.Expression, test.expectedValue, test.expectedOperator)
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
		{`"hello " + "world"`, "hello ", "+", "world"},
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
		{
			"sum(1 + 2, 3 * 4 ^ 5, fn(){}, sub(1, 0))",
			"sum((1 + 2), (3 * (4 ^ 5)), fn() {}, sub(1, 0))",
		},
		{
			"2 ^ add(2, 4) * 8",
			"((2 ^ add(2, 4)) * 8)",
		},
		{
			"2 ^ add(2, 4) * 8",
			"((2 ^ add(2, 4)) * 8)",
		},
		{
			"!isTrue(1 > 2)",
			"(!isTrue((1 > 2)))",
		},
		{
			"1 + [1, 2, 3][0] - 2",
			"((1 + ([1, 2, 3][0])) - 2)",
		},
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
		if (true) {z} else {1}
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

func TestFunctionLiteralExpression(t *testing.T) {
	input := "fn(){}"
	lex := lexer.New(input)
	par := New(lex)
	prg := par.ParseProgram()
	checkParserErrors(t, par)
	if len(prg.Statements) != 1 {
		t.Fatalf("FunctionLiteral got and unexpected number of statements: %d': 1. Expected %d",
			len(prg.Statements), 1)
	}
	stmt := testExpressionStatement(t, prg.Statements[0])
	testFunctionLiteralExpression(t, stmt.Expression, []string{}, 0)
}

func TestFunctionParameters_Zero(t *testing.T) {
	input := "fn(){}"
	lex := lexer.New(input)
	par := New(lex)
	prg := par.ParseProgram()
	checkParserErrors(t, par)
	if len(prg.Statements) != 1 {
		t.Fatalf("FunctionLiteral got and unexpected number of statements: %d': 1. Expected %d",
			len(prg.Statements), 1)
	}
	stmt := testExpressionStatement(t, prg.Statements[0])
	testFunctionLiteralExpression(t, stmt.Expression, []string{}, 0)
}

func TestFunctionParameters_One(t *testing.T) {
	input := "fn(x){}"
	lex := lexer.New(input)
	par := New(lex)
	prg := par.ParseProgram()
	checkParserErrors(t, par)
	if len(prg.Statements) != 1 {
		t.Fatalf("FunctionLiteral got and unexpected number of statements: %d': 1. Expected %d",
			len(prg.Statements), 1)
	}
	stmt := testExpressionStatement(t, prg.Statements[0])
	testFunctionLiteralExpression(t, stmt.Expression, []string{"x"}, 0)
}

func TestFunctionParameters_Several(t *testing.T) {
	input := "fn(x, y, z){}"
	lex := lexer.New(input)
	par := New(lex)
	prg := par.ParseProgram()
	checkParserErrors(t, par)
	if len(prg.Statements) != 1 {
		t.Fatalf("FunctionLiteral got and unexpected number of statements: %d': 1. Expected %d",
			len(prg.Statements), 1)
	}
	stmt := testExpressionStatement(t, prg.Statements[0])
	testFunctionLiteralExpression(t, stmt.Expression, []string{"x", "y", "z"}, 0)
}

func TestFunctionParameters_SyntaxError(t *testing.T) {
	input := "fn(x, y,){}" // There is a missing identifier
	lex := lexer.New(input)
	par := New(lex)
	par.ParseProgram()
	expectAnyParserErrors(t, par)
}

func TestCallExpression(t *testing.T) {
	input := `sum(1, 2)`
	lex := lexer.New(input)
	par := New(lex)
	prg := par.ParseProgram()
	checkProgram(t, prg)
	stmt := testExpressionStatement(t, prg.Statements[0])
	callExp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expression expected to be CallExpression. Got %q", stmt.Expression)
	}

	testLiteralExpression(t, callExp.Function, "sum")

	expectedParameters := []int{1, 2}

	if len(expectedParameters) != len(callExp.Arguments) {
		t.Fatalf("CallExpression expected to have %d parameters. Got %d",
			len(expectedParameters), len(callExp.Arguments))
	}

	for index, expParam := range expectedParameters {
		actualParam := callExp.Arguments[index]
		testLiteralExpression(t, actualParam, expParam)
	}

}

func TestCallExpressionNoArgs(t *testing.T) {
	input := `sum()`
	lex := lexer.New(input)
	par := New(lex)
	prg := par.ParseProgram()
	checkProgram(t, prg)
	stmt := testExpressionStatement(t, prg.Statements[0])
	callExp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expression expected to be CallExpression. Got %q", stmt.Expression)
	}

	testLiteralExpression(t, callExp.Function, "sum")

	expectedParameters := make([]int, 0)

	if len(expectedParameters) != len(callExp.Arguments) {
		t.Fatalf("CallExpression expected to have %d parameters. Got %d",
			len(expectedParameters), len(callExp.Arguments))
	}

	for index, expParam := range expectedParameters {
		actualParam := callExp.Arguments[index]
		testLiteralExpression(t, actualParam, expParam)
	}
}

func testArrayLiteral(t *testing.T, actual ast.Expression, expected []int64) bool {
	arr, ok := actual.(*ast.ArrayLiteral)
	if !ok {
		t.Errorf("Expression is not ArrayLiteral. Got %T(%+v)", actual, actual)
		return false
	}
	if len(expected) != len(arr.Items) {
		t.Errorf("ArrayLiteral expected to have %d elements. Got %d",
			len(expected), len(arr.Items))
		return false
	}
	for index, expected := range expected {
		if !testIntegerLiteralExpression(t, arr.Items[index], expected) {
			return false
		}
	}
	return true
}

func TestArrayLiteralExpression(t *testing.T) {
	tests := []struct {
		code     string
		expected []int64
	}{
		{"[0, 1, 4]", []int64{0, 1, 4}},
		{"[]", []int64{}},
	}

	for _, test := range tests {
		program := ParseTesting(t, test.code)
		stmt := testExpressionStatement(t, program.Statements[0])
		testArrayLiteral(t, stmt.Expression, test.expected)
	}
}

func TestIndexExpressionParsing(t *testing.T) {
	// TODO: Add more!
	code := `[1, 2, 3][1]`
	program := ParseTesting(t, code)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement. Got %T", program.Statements[0])
	}
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("expected IndexExpression. Got %T", stmt.Expression)
	}
	if indexExp.Container.String() != "[1, 2, 3]" {
		t.Fatalf("expected IndexExpression.Container to be %s. Got %s",
			"[1, 2, 3]", indexExp.Container.String())
	}
	if !testLiteralExpression(t, indexExp.Index, 1) {
		t.Fatalf("expected IndexExpression.Index to be %s. Got %s", "1", indexExp.Index.String())
	}
}
