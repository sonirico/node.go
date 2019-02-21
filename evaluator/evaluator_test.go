package evaluator

import (
	"node.go/lexer"
	"node.go/object"
	"node.go/parser"
	"testing"
)

func checkParserErrors(t *testing.T, p *parser.Parser) {
	if len(p.Errors()) < 1 {
		return
	}

	for _, errorMessage := range p.Errors() {
		t.Errorf("Parser error: %s", errorMessage)
	}

	t.FailNow()
}
func testEval(t *testing.T, code string) object.Object {
	lex := lexer.New(code)
	par := parser.New(lex)
	prg := par.ParseProgram()
	checkParserErrors(t, par)
	if prg == nil {
		t.Fatalf("Parser.ParseProgram returned nil")
	}
	if len(prg.Statements) < 1 {
		t.Fatalf("ast.Program has got no nodes")
	}
	environment := object.NewEnvironment()
	evaluated := Eval(prg, environment)
	if evaluated == nil {
		t.Fatalf("Eval returned nil")
	}
	return evaluated
}

func testErrorObject(t *testing.T, obj object.Object, message string) bool {
	errorObj, ok := obj.(*object.Error)
	if !ok {
		t.Errorf("Object is not Error. Got %T(%+v)", obj, obj)
		return false
	}
	if errorObj.Message != message {
		t.Errorf("Error.Message is not '%s'. Got '%s'",
			message, errorObj.Message)
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	_, ok := obj.(*object.Null)
	if !ok {
		t.Errorf("Object is not Null. Got %q", obj)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	booleanObj, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Object is not Boolean. Got %s<%s>", obj.Inspect(), obj.Type())
		return false
	}
	if booleanObj.Value != expected {
		t.Errorf("Boolean.Value is not %t. Got %t", expected, booleanObj.Value)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int) bool {
	integer, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Object is not Integer. Got %s", obj.Type())
		if obj.Type() == object.ERROR {
			t.Errorf("The error is %s", obj.Inspect())
		}
		return false
	}
	if expected != int(integer.Value) {
		t.Errorf("Integer.Value is not %d. Got %d", expected, integer.Value)
		return false
	}
	return true
}

func TestEvalIntegerObject(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"1", 1},
		{"0", 0},
		{"999", 999},
		{"-2", -2},
		{"1 + 1", 2},
		{"1 - 3", -2},
		{"1 - -3", 4},
		{"2 * (1 + 3)", 8},
		{"1 + 2 * 3 / 2 - 4", 0},
	}
	for _, test := range tests {
		evaluated := testEval(t, test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestEvalBooleanObject(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!1", false},
		{"!!1", true},
		{"true == true", true},
		{"false == true", false},
		{"false != true", true},
		{"false != false", false},
		{"1 > 1", false},
		{"1 > -1", true},
		{"1 == 1", true},
		{"0 < 1", true},
		{"2 < 1", false},
		{"1 != 1", false},
		{"-1 != 1", true},
		{"(2 > 0) == true", true},
		{"!(2 > 0) == true", false},
	}
	for _, test := range tests {
		evaluated := testEval(t, test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func TestEvalNullObject(t *testing.T) {
	tests := []string{"1 / 0"}
	for _, test := range tests {
		evaluated := testEval(t, test)
		testNullObject(t, evaluated)
	}
}

func TestIfConditionalEval(t *testing.T) {
	tests := []struct {
		code     string
		expected interface{}
	}{
		{"if (true) {1;}", 1},
		{"if (false) {1}", nil},
		{"if (2 > 0) {1}", 1},
		{"if (false) {1} else {2}", 2},
		{"if (false == (1 == 1)) {1} else {2}", 2},
		{"if (0) {} else {2}", 2},
		{"if (true) {return 1;} else {2}", 1},
		{"if (true) {return;} else {2}", nil},
	}
	for _, test := range tests {
		expected := test.expected
		evaluated := testEval(t, test.code)
		if expected == nil {
			testNullObject(t, evaluated)
		} else {
			intExp, _ := expected.(int)
			testIntegerObject(t, evaluated, intExp)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"return;", nil},
		{"9; return 1; 5", 1},
		{"9; return; 5", nil},
		{"return 1; 5", 1},
		{"9; 0; return 2;", 2},
		{
			`if (1 > 0) {
					if (1 > 0) {
						return 2;
				    }
					return 0;
               }`,
			2,
		},
		{
			"if (1 > 0) {return 2;}; return 0;", 2,
		},
	}

	for _, test := range tests {
		expected := test.expected
		evaluated := testEval(t, test.input)
		if nil == expected {
			testNullObject(t, evaluated)
		} else {
			expectedNumber := expected.(int)
			testIntegerObject(t, evaluated, expectedNumber)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input                string
		expectedErrorMessage string
	}{
		{
			"1 == true",
			"type mismatch: INTEGER == BOOLEAN",
		},
		{
			"true > false",
			"unknown operator: BOOLEAN > BOOLEAN",
		},
		{
			"true + false",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"true - false",
			"unknown operator: BOOLEAN - BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"1 > (false == 2)",
			"type mismatch: BOOLEAN == INTEGER",
		},
		{
			"!(true * true)",
			"unknown operator: BOOLEAN * BOOLEAN",
		},
		{
			`if (true) {
				if (1 != false) {
					return 4
				}
			}`,
			"type mismatch: INTEGER != BOOLEAN",
		},
		{
			"false <= 1; return 2;",
			"type mismatch: BOOLEAN <= INTEGER",
		},
		{
			"a + 1;",
			"reference error: a is not defined",
		},
		{
			"let b = a * 3",
			"reference error: a is not defined",
		},
		{
			"let f = fn(x, y) {a + y} (1, 2)",
			"reference error: a is not defined",
		},
	}
	for _, test := range tests {
		evaluated := testEval(t, test.input)
		testErrorObject(t, evaluated, test.expectedErrorMessage)
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{
			"let a = 1; a;",
			1,
		},
		{
			"let a = 1; let b = a - 2; b;",
			-1,
		},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	code := "fn (x) { x + 2;}"
	evaluated := testEval(t, code)
	evalFunc, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("evalFunc is not an object.Function. Got %T(%+v)", evaluated, evaluated)
	}
	if len(evalFunc.Parameters) != 1 {
		t.Fatalf("function expected to have %d params. Got %d", 2, len(evalFunc.Parameters))
	}
	if evalFunc.Parameters[0].String() != "x" {
		t.Fatalf("parameter expected to be '%s'. Got '%s'", "x", evalFunc.Parameters[0].String())
	}
	expectedBody := "{(x + 2)}"
	if expectedBody != evalFunc.Body.String() {
		t.Fatalf("function expected to have body equal to '%s'. Got '%s'",
			expectedBody, evalFunc.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		code     string
		expected int
	}{
		{
			`let identity = fn (x) {x}; identity(1);`,
			1,
		},
		{
			`let double = fn (x) {x * 2}; double(2); 4`,
			4,
		},
		{
			`let add = fn (a, b) {a + b}; add(add(1, 3), add(-1, -3))`,
			0,
		},
		{
			`fn(a, b){a + b}(1, 2);`,
			3,
		},
		{
			`
			let a = fn (x) {
				let b = fn (y) {
					return y * 2
				}
				return x + b(x)
			}
			a(2);
			`,
			6,
		},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.code)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestFunctionClosures(t *testing.T) {
	tests := []struct {
		code     string
		expected int
	}{
		{
			`
			let sumGenerator = fn (left) {
				return fn (y) { left + y }
			}
		
			let addTwo = sumGenerator(2);
			addTwo(4)
			`,
			6,
		},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.code)
		testIntegerObject(t, evaluated, test.expected)
	}
}
