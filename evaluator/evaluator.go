package evaluator

import (
	"node.go/ast"
	"node.go/object"
	"node.go/token"
)

func booleanToObject(value bool) object.Object {
	if value {
		return object.TRUE
	} else {
		return object.FALSE
	}
}

func evalBooleanLiteral(obj object.Object) object.Object {
	switch obj {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	}
	return object.NULL
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)
	}
	return result
}

func evalMinusOperatorExpression(obj object.Object) object.Object {
	if obj.Type() != object.INT {
		return object.NULL
	}
	intObj, _ := obj.(*object.Integer)
	return &object.Integer{Value: -intObj.Value}
}

func evalIntegerToBoolean(value int64) *object.Boolean {
	if value == 0 {
		return object.TRUE
	} else {
		return object.FALSE
	}
}

func evalBangOperatorExpression(obj object.Object) object.Object {
	switch obj.Type() {
	case object.BOOL:
		return evalBooleanLiteral(obj)
	case object.INT:
		{
			intObj, _ := obj.(*object.Integer)
			return evalIntegerToBoolean(intObj.Value)
		}
	case object.NULL_TYPE:
		return object.TRUE
	}
	return object.NULL
}

func evalPrefixExpression(operator string, obj object.Object) object.Object {
	switch operator {
	case token.MINUS:
		return evalMinusOperatorExpression(obj)
	case token.BANG:
		return evalBangOperatorExpression(obj)
	}
	return object.NULL
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		{
			operator := node.Operator
			right := Eval(node.Right)
			return evalPrefixExpression(operator, right)
		}
	case *ast.IntegerLiteral:
		return object.NewInteger(node.Value)
	case *ast.BooleanLiteral:
		return booleanToObject(node.Value)
	}

	return object.NULL
}
