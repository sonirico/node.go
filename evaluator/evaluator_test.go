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
	return Eval(prg)
}

func testIntegerObjectEval(t *testing.T, obj object.Object, expected int64) bool {
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

func TestIntegerObjectEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"1", 1},
		{"0", 0},
		{"999", 999},
		//{"-2", -2},
	}
	for _, test := range tests {
		evaluated := testEval(t, test.input)
		if evaluated == nil {
			t.Fatal("evaluator.Eval returned nil")
		}
		testIntegerObjectEval(t, evaluated, test.expected)
	}
}
