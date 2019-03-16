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
			"type mismatch: Expected STRING or ARRAY. Got INTEGER",
		},
		{
			`true[0]`,
			"type error: BOOLEAN cannot be used as index expression",
		},
		{
			`[1][true]`,
			"type error: BOOLEAN cannot be used as index of ARRAY",
		},
		{
			`head()`,
			"type error: Expected 1 argument. Got 0",
		},
		{
			`head([1, 2], 1)`,
			"type error: Expected 1 argument. Got 2",
		},
		{
			`head(true)`,
			"type mismatch: Expected ARRAY. Got BOOLEAN",
		},
		{
			`foot()`,
			"type error: Expected 1 argument. Got 0",
		},
		{
			`foot([1, 2], 1)`,
			"type error: Expected 1 argument. Got 2",
		},
		{
			`foot(true)`,
			"type mismatch: Expected ARRAY. Got BOOLEAN",
		},
		{
			`tail()`,
			"type error: Expected 1 argument. Got 0",
		},
		{
			`tail([1, 2, 3], false)`,
			"type error: Expected 1 argument. Got 2",
		},
		{
			`tail(false)`,
			"type mismatch: Expected ARRAY. Got BOOLEAN",
		},
		{
			`push()`,
			"type error: Expected 2 arguments. Got 0",
		},
		{
			`push([])`,
			"type error: Expected 2 arguments. Got 1",
		},
		{
			`push([], 3, 5)`,
			"type error: Expected 2 arguments. Got 3",
		},
		{
			`push(2, true)`,
			"type mismatch: Expected ARRAY. Got INTEGER",
		},
		{
			`pop()`,
			"type error: Expected 1 argument. Got 0",
		},
		{
			`pop(true)`,
			"type mismatch: Expected ARRAY. Got BOOLEAN",
		},
		{
			`pop([], 2)`,
			"type error: Expected 1 argument. Got 2",
		},
		{
			`let dict = {fn(){}: 1}`,
			"value error: unhashable type as hash key: FUNCTION",
		},
		{
			`let dict = {{}: 1}`,
			"value error: unhashable type as hash key: HASH",
		},
		{
			`let null; let dict = {null: 1}`,
			"value error: unhashable type as hash key: NULL",
		},
		{
			`let dict = {[]: 1}`,
			"value error: unhashable type as hash key: ARRAY",
		},
		{
			`{}[{}]`,
			"value error: unhashable type as hash key: HASH",
		},
		{
			`{}[[]]`,
			"value error: unhashable type as hash key: ARRAY",
		},
		{
			`let null; {}[null]`,
			"value error: unhashable type as hash key: NULL",
		},
		{
			`{}[fn(){}]`,
			"value error: unhashable type as hash key: FUNCTION",
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

// BUILTINS

func TestLenBuiltinFunction(t *testing.T) {
	tests := []struct {
		code     string
		expected int
	}{
		{`len(" hello   there ")`, 15},
		{`len("")`, 0},
		{`len(fn (x) { "I ve " + x + " eyes"} ("2"))`, 11},
		{`len([])`, 0},
		{`len([1])`, 1},
		{`len([1, 2])`, 2},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.code)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestHeadBuiltinFunction(t *testing.T) {
	tests := []struct {
		code     string
		expected interface{}
	}{
		{
			`head([1, 2, 3])`,
			1,
		},
		{
			`head([true])`,
			true,
		},
		{
			`head([])`,
			nil,
		},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.code)
		switch expectedValue := test.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, expectedValue)
		case bool:
			testBooleanObject(t, evaluated, expectedValue)
		case nil:
			testNullObject(t, evaluated)
		}
	}
}

func TestFootBuiltinFunction(t *testing.T) {
	tests := []struct {
		code     string
		expected interface{}
	}{
		{
			`foot([1, 2, 3])`,
			3,
		},
		{
			`foot([true])`,
			true,
		},
		{
			`foot([])`,
			nil,
		},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.code)
		switch expectedValue := test.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, expectedValue)
		case bool:
			testBooleanObject(t, evaluated, expectedValue)
		case nil:
			testNullObject(t, evaluated)
		}
	}
}

func TestTailBuiltinFunction(t *testing.T) {
	tests := []struct {
		code     string
		expected interface{}
	}{
		{
			`tail([1, 2, 3])`,
			[]int{2, 3},
		},
		{
			`tail([true])`,
			[]int{},
		},
		{
			`tail([])`,
			nil,
		},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.code)
		switch expectedValue := test.expected.(type) {
		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Fatalf("expected array")
			}
			for index, expectedIntValue := range expectedValue {
				testIntegerObject(t, array.Items[index], expectedIntValue)
			}
		case nil:
			testNullObject(t, evaluated)
		}
	}
}

func TestPushBuiltinFunction(t *testing.T) {
	tests := []struct {
		code     string
		expected []int
	}{
		{
			`push([1, 2, 3], 4)`,
			[]int{1, 2, 3, 4},
		},
		{
			`push([1], 2)`,
			[]int{1, 2},
		},
		{
			`push([], 1)`,
			[]int{1},
		},
		{
			`let arr = [1]; push(arr, 2); arr;`,
			[]int{1},
		},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.code)
		array, ok := evaluated.(*object.Array)
		if !ok {
			t.Fatalf("expected object.Array. Got %T(%+v)", evaluated, evaluated)
		}
		for index, expectedIntValue := range test.expected {
			testIntegerObject(t, array.Items[index], expectedIntValue)
		}
	}
}

func TestPopBuiltinFunction(t *testing.T) {
	tests := []struct {
		code     string
		expected interface{}
	}{
		{
			`pop([1, 2, 3])`,
			3,
		},
		{
			`pop([1])`,
			1,
		},
		{
			`pop([])`,
			nil,
		},
		{
			`let arr = [1]; pop(arr); arr;`,
			[]int{},
		},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.code)
		switch expected := test.expected.(type) {
		case nil:
			testNullObject(t, evaluated)
		case int:
			testIntegerObject(t, evaluated, expected)
		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Fatalf("expected object.Array. Got %T(%+v)", evaluated, evaluated)
			}
			if len(expected) != len(array.Items) {
				t.Fatalf("expected object.Array to have length of %d. Got %d",
					len(expected), len(array.Items))
			}
			for index, expectedIntValue := range expected {
				testIntegerObject(t, array.Items[index], expectedIntValue)
			}
		}
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

func TestIndexExpressionEvaluation(t *testing.T) {
	tests := []struct {
		code     string
		expected interface{}
	}{
		{
			`[1, 2 * 3, fn(n){n - 1}(1)][2]`,
			0,
		},
		{
			`let array = [1, 2 * 3, fn(n){n - 1}(1)]; array[2]`,
			0,
		},
		{
			`let getArray = fn () {
				return [1, 2 * 3, fn(n){n - 1}(1)]
			};
			getArray()[1];
			`,
			6,
		},
		{
			"[1][1]",
			nil,
		},
	}
	for _, test := range tests {
		evaluated := testEval(t, test.code)
		switch exp := test.expected.(type) {
		case nil:
			testNullObject(t, evaluated)
		case int:
			if !testIntegerObject(t, evaluated, exp) {
				t.Fatalf("expected IndexExpression to return %d. Got %s",
					test.expected, evaluated.Type())
			}
		}
	}
}

func TestStringHashKey(t *testing.T) {
	s1 := object.NewString("hola que tal amigo")
	s2 := object.NewString("hola que tal amigo")
	if s1.HashKey().Value != s2.HashKey().Value {
		t.Fatalf("StringLiterarl.HashKey differ: %d != %d",
			s1.HashKey().Value, s2.HashKey().Value)
	}

	s1 = object.NewString("vamos a subirle el IVA a los chuches tambien!")
	s2 = object.NewString("tenemos que construir maquinas...")
	if s1.HashKey().Value == s2.HashKey().Value {
		t.Fatalf("StringLiterarl.HashKey should differ: %d != %d",
			s1.HashKey().Value, s2.HashKey().Value)
	}
}

func TestIntegersHashKey(t *testing.T) {
	s1 := object.NewInteger(189423879)
	s2 := object.NewInteger(189423879)
	if s1.HashKey().Value != s2.HashKey().Value {
		t.Fatalf("Integer.HashKey differ: %d != %d",
			s1.HashKey().Value, s2.HashKey().Value)
	}

	s1 = object.NewInteger(189423879)
	s2 = object.NewInteger(289423879)
	if s1.HashKey().Value == s2.HashKey().Value {
		t.Fatalf("Integer.HashKey should differ: %d != %d",
			s1.HashKey().Value, s2.HashKey().Value)
	}
}

func TestBooleanHashKey(t *testing.T) {
	s1 := object.NewBoolean(true)
	s2 := object.NewBoolean(true)
	if s1.HashKey().Value != s2.HashKey().Value {
		t.Fatalf("Boolean.HashKey differ: %d != %d",
			s1.HashKey().Value, s2.HashKey().Value)
	}
	s1 = object.NewBoolean(true)
	s2 = object.NewBoolean(false)
	if s1.HashKey().Value == s2.HashKey().Value {
		t.Fatalf("Boolean.HashKey should differ: %d != %d",
			s1.HashKey().Value, s2.HashKey().Value)
	}
}

func TestHashEvaluation(t *testing.T) {
	payload := `let k = 6; {
		!false: -0,
		!!!true: 2 / 2,
		(4 / 2) - (-3 / 3): 2 > 1,
		"for" + "th": -(-3),
		fn(){"fifth"}(): fn(x){x}(4),
		k: 5
	}`
	expected := map[object.HashKey]object.HashPair{
		object.TRUE.HashKey():               {Key: object.TRUE, Value: object.NewInteger(0)},
		object.FALSE.HashKey():              {Key: object.FALSE, Value: object.NewInteger(1)},
		object.NewInteger(3).HashKey():      {Key: object.NewInteger(3), Value: object.TRUE},
		object.NewString("forth").HashKey(): {Key: object.NewString("forth"), Value: object.NewInteger(3)},
		object.NewString("fifth").HashKey(): {Key: object.NewString("fifth"), Value: object.NewInteger(4)},
		object.NewInteger(6).HashKey():      {Key: object.NewInteger(6), Value: object.NewInteger(5)},
	}
	evaluated := testEval(t, payload)
	hash, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("hash evaluation. expected object to be %s. Got %s",
			object.HASH, evaluated.Type())
	}
	if len(expected) != len(hash.Pairs) {
		t.Fatalf("Hash got a wrong number of pairs. Expected %d. Got %d",
			len(expected), len(hash.Pairs))
	}
	for expectedHashKey, expectedHashPair := range expected {
		actualHashPair, ok := hash.Pairs[expectedHashKey]
		if !ok {
			t.Fatalf("Hash did not contain HashKey %d for Key %s",
				expectedHashKey.Value, expectedHashPair.Key.Inspect())
		}
		actualKey := actualHashPair.Key
		actualValue := actualHashPair.Value
		expectedKey := expectedHashPair.Key
		expectedValue := expectedHashPair.Value
		if actualKey.Type() != expectedKey.Type() {
			t.Fatalf("Hash contains an unexpected hash key %s. Expected %s",
				actualKey.Type(), expectedKey.Type())
		}
		if actualKey.Inspect() != expectedKey.Inspect() {
			t.Fatalf("Hash key was expected to be %s. Got %s",
				actualKey.Inspect(), expectedKey.Inspect())
		}
		if actualValue.Type() != expectedValue.Type() {
			t.Fatalf("Hash contains an unexpected hash value %s. Expected %s",
				actualValue.Type(), expectedValue.Type())
		}
		if actualValue.Inspect() != expectedValue.Inspect() {
			t.Fatalf("Hash value was expected to be %s. Got %s",
				actualValue.Inspect(), expectedValue.Inspect())
		}
	}
}

func TestHashIndexExpressionEvaluation(t *testing.T) {
	tests := []struct {
		code     string
		expected object.Object
	}{
		{`{0: 1}[0]`, object.NewInteger(1)},
		{`{"two": 2}["one"]`, object.NULL},
		{`{"two": 2}["two"]`, object.NewInteger(2)},
		{`{false: false}[false]`, object.FALSE},
		{`{"sum": fn(a, b){a + b}}["sum"](1, 2)`, object.NewInteger(3)},
		{`let regKey = "registry"; {regKey: {}}[regKey]`, object.NewHash()},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.code)
		if test.expected.Type() != evaluated.Type() {
			t.Fatalf("Expected object.Type to be %s. Got %s",
				test.expected.Type(), evaluated.Type())
		}
		if test.expected.Inspect() != evaluated.Inspect() {
			t.Fatalf("Expected object.Inspect to be %s. Got %s",
				test.expected.Inspect(), evaluated.Inspect())
		}
	}
}
