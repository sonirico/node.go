package evaluator

import (
	"node.go/lexer"
	"node.go/object"
	"node.go/parser"
	"testing"
)

func testEval(t *testing.T, code string) object.Object {
	lex := lexer.New(code)
	par := parser.New(lex)
	prg := par.ParseProgram()
	if prg == nil {
		t.Fatalf("Parser.ParseProgram returned nil")
	}
	if len(prg.Statements) < 1 {
		t.Fatalf("ast.Program has got no nodes")
	}
	evaluated := Eval(prg)
	if evaluated == nil {
		t.Fatalf("Eval returned nil")
	}
	return evaluated
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

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	integer, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Object is not Integer. Got %s", obj.Type())
		return false
	}
	if expected != integer.Value {
		t.Errorf("Integer.Value is not %d. Got %d", expected, integer.Value)
		return false
	}
	return true
}

func TestEvalIntegerObject(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
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
	tests := []string{"null", "1 / 0"}
	for _, test := range tests {
		evaluated := testEval(t, test)
		testNullObject(t, evaluated)
	}
}
