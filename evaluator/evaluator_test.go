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
		t.Errorf("Object is not Boolean. Got %q", obj)
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
	}
	for _, test := range tests {
		evaluated := testEval(t, test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func TestEvalNullObject(t *testing.T) {
	evaluated := testEval(t, "null")
	testNullObject(t, evaluated)
}
