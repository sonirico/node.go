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

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	stringObj, ok := obj.(*object.String)
	if !ok {
		t.Errorf("Object is not String. Got %s", obj.Type())
		if obj.Type() == object.ERROR {
			t.Errorf("The error is %s", obj.Inspect())
		}
		return false
	}
	if expected != stringObj.Value {
		t.Errorf("String.Value is not %s. Got %s", expected, stringObj.Value)
		return false
	}
	return true
}

func testFunctionObject(t *testing.T, evaluated object.Object, arguments []string, body string) bool {
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Errorf("function is not an object.Function. Got %T(%+v)", evaluated, evaluated)
		return false
	}
	if len(fn.Parameters) != len(arguments) {
		t.Errorf("function expected to have %d params. Got %d", len(arguments), len(fn.Parameters))
		return false
	}
	for index, argument := range arguments {
		if fn.Parameters[index].String() != argument {
			t.Errorf("parameter expected to be '%s'. Got '%s'", argument, fn.Parameters[index].String())
			return false
		}
	}
	if fn.Body.String() != body {
		t.Errorf("function expected to have body equal to '%s'. Got '%s'",
			body, fn.Body.String())
		return false
	}
	return true
}

func testArrayObject(t *testing.T, evaluated object.Object, assertions func(array *object.Array)) {
	arr, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("expected Array. Got %T(%+v)", evaluated, evaluated)
	}
	assertions(arr)
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

func TestEvalStringObject(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"I am a spoon";`, "I am a spoon"},
		{`"0"`, "0"},
	}
	for _, test := range tests {
		evaluated := testEval(t, test.input)
		testStringObject(t, evaluated, test.expected)
	}
}

func TestEvalStringConcatenation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let a = "Hello";let b = " world!"; let c = a + b; c;`, "Hello world!"},
	}
	for _, test := range tests {
		evaluated := testEval(t, test.input)
		testStringObject(t, evaluated, test.expected)
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
		{
			`"hello" - "world"`,
			"unknown operator: STRING - STRING",
		},
		{
			`len()`,
			"Type error: Expected 1 argument. Got 0",
		},
		{
			`len(1, 2)`,
			"Type error: Expected 1 argument. Got 2",
		},
		{
			`len(1)`,
			"Type mismatch: Expected STRING. Got INTEGER",
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
	expectedParams := []string{"x"}
	expectedBody := "{(x + 2)}"
	testFunctionObject(t, evaluated, expectedParams, expectedBody)
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

func TestLenBuiltinFunction(t *testing.T) {
	tests := []struct {
		code     string
		expected int
	}{
		{`len(" hello   there ")`, 15},
		{`len("")`, 0},
		{`len(fn (x) { "I ve " + x + " eyes"} ("2"))`, 11},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.code)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestArrayLiteralEvaluation(t *testing.T) {
	tests := []struct {
		code       string
		assertions func(*object.Array)
	}{
		{
			`[true, 1, "hey", fn(){}]`,
			func(arr *object.Array) {
				if len(arr.Items) != 4 {
					t.Fatalf("expected Array to have %d items. Got %d", 4, len(arr.Items))
				}
				testBooleanObject(t, arr.Items[0], true)
				testIntegerObject(t, arr.Items[1], 1)
				testStringObject(t, arr.Items[2], "hey")
				testFunctionObject(t, arr.Items[3], []string{}, "{}")
			},
		},
		{
			`[]`,
			func(arr *object.Array) {},
		},
		{
			`let f = fn(z){z * z}; [f(-2), fn(x, y){x + y}("a", "b")]`,
			func(arr *object.Array) {
				if len(arr.Items) != 2 {
					t.Fatalf("expected Array to have %d items. Got %d", 4, len(arr.Items))
				}
				testIntegerObject(t, arr.Items[0], 4)
				testStringObject(t, arr.Items[1], "ab")
			},
		},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.code)
		testArrayObject(t, evaluated, test.assertions)
	}
}
