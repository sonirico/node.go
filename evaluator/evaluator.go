package evaluator

import (
	"fmt"
	"node.go/ast"
	"node.go/object"
	"node.go/token"
)

func isError(obj object.Object) bool {
	return obj.Type() == object.ERROR
}

func newError(template string, params ...interface{}) object.Object {
	return object.NewError(fmt.Sprintf(template, params...))
}

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

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt, env)
		if result != nil {
			switch result := result.(type) {
			case *object.Error:
				return result
			case *object.Return:
				return result.Value
			}
		}
	}
	return result
}

func evalBlockStatement(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)

		if result != nil {
			switch result.Type() {
			case object.RETURN:
				return result
			case object.ERROR:
				return result
			}
		}
	}

	return result
}

func evalMinusOperatorExpression(obj object.Object) object.Object {
	if obj.Type() != object.INT {
		return newError("unknown operator: -%s", obj.Type())
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
	return newError("unknown operator: !%s", obj.Type())
}

func evalPrefixExpression(operator string, obj object.Object) object.Object {
	switch operator {
	case token.MINUS:
		return evalMinusOperatorExpression(obj)
	case token.BANG:
		return evalBangOperatorExpression(obj)
	}
	return newError("unknown operator: %s%s", operator, obj.Type())
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
	default:
		return newError("unknown operator: %s%s", operator, object.INT)
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
	return newError("unknown operator: %s %s %s", object.BOOL, operator, object.BOOL)
}

func evalInfixOperatorExpression(
	operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INT && right.Type() == object.INT:
		return evalInfixIntegerExpression(operator, left, right)
	case left.Type() == object.BOOL && right.Type() == object.BOOL:
		return evalInfixBooleanExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}
	return newError("unsupported types: %s %s %s", left.Type(), operator, right.Type())
}

func isTruthy(obj object.Object) bool {
	if obj == object.FALSE || obj == object.NULL {
		return false
	}
	if obj.Type() == object.INT {
		intObj, _ := obj.(*object.Integer)
		return intObj.Value != 0
	}
	return true
}

func evalIfConditionalExpression(ifExpression *ast.IfExpression, env *object.Environment) object.Object {
	evaluatedCondition := Eval(ifExpression.Condition, env)
	if isError(evaluatedCondition) {
		return evaluatedCondition
	}
	if isTruthy(evaluatedCondition) {
		return Eval(&ifExpression.Consequence, env)
	}
	if ifExpression.Alternative != nil {
		return Eval(ifExpression.Alternative, env)
	}
	return object.NULL
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements, env)
	case *ast.LetStatement:
		{
			value := Eval(node.Value, env)
			if isError(value) {
				return value
			}
			env.Set(node.Name.Value, value)
		}
	case *ast.Identifier:
		{
			value, ok := env.Get(node.Value)
			if !ok {
				return newError("reference error: %s is not defined", node.Value)
			}
			return value
		}
	case *ast.ReturnStatement:
		{
			value := Eval(node.ReturnValue, env)
			if isError(value) {
				return value
			}
			return &object.Return{Value: value}
		}
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.PrefixExpression:
		{
			operator := node.Operator
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			return evalPrefixExpression(operator, right)
		}
	case *ast.InfixExpression:
		{
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			return evalInfixOperatorExpression(node.Operator, left, right)
		}
	case *ast.IfExpression:
		return evalIfConditionalExpression(node, env)
	case *ast.IntegerLiteral:
		return object.NewInteger(node.Value)
	case *ast.BooleanLiteral:
		return booleanToObject(node.Value)
	}

	return object.NULL
}
