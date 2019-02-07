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

func evalInfixIntegerExpression(
	operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	switch operator {
	case token.PLUS:
		return object.NewInteger(leftValue + rightValue)
	case token.ASTERISK:
		return object.NewInteger(leftValue * rightValue)
	case token.MINUS:
		return object.NewInteger(leftValue - rightValue)
	case token.SLASH:
		{
			if 0 == rightValue {
				return object.NULL
			}
			return object.NewInteger(leftValue / rightValue)
		}
	case token.EQ:
		return booleanToObject(leftValue == rightValue)
	case token.NOT_EQ:
		return booleanToObject(leftValue != rightValue)
	case token.LT:
		return booleanToObject(leftValue < rightValue)
	case token.GT:
		return booleanToObject(leftValue > rightValue)
	case token.LTE:
		return booleanToObject(leftValue <= rightValue)
	case token.GTE:
		return booleanToObject(leftValue >= rightValue)
	}
	return object.NULL
}

func evalInfixBooleanExpression(operator string, left object.Object, right object.Object) object.Object {
	switch operator {
	case token.EQ:
		return booleanToObject(left == right)
	case token.NOT_EQ:
		return booleanToObject(left != right)
	}
	return object.NULL
}

func evalInfixOperatorExpression(
	operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INT && right.Type() == object.INT:
		return evalInfixIntegerExpression(operator, left, right)
	case left.Type() == object.BOOL && right.Type() == object.BOOL:
		return evalInfixBooleanExpression(operator, left, right)
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
	case *ast.InfixExpression:
		{
			left := Eval(node.Left)
			right := Eval(node.Right)
			return evalInfixOperatorExpression(node.Operator, left, right)
		}
	case *ast.IntegerLiteral:
		return object.NewInteger(node.Value)
	case *ast.BooleanLiteral:
		return booleanToObject(node.Value)
	}

	return object.NULL
}
