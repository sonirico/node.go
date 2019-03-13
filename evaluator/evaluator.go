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

func evalInfixStringExpression(operator string, left object.Object, right object.Object) object.Object {
	if operator != token.PLUS {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	leftStr, ok := left.(*object.String)
	if !ok {
		return newError("type mismatch: %s %s %s. Bad left operand", left.Type(), operator, right.Type())
	}
	rightStr, ok := right.(*object.String)
	if !ok {
		return newError("type mismatch: %s %s %s. Bad left operand", left.Type(), operator, right.Type())
	}
	return object.NewString(leftStr.Value + rightStr.Value)
}

func evalInfixOperatorExpression(
	operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INT && right.Type() == object.INT:
		return evalInfixIntegerExpression(operator, left, right)
	case left.Type() == object.BOOL && right.Type() == object.BOOL:
		return evalInfixBooleanExpression(operator, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalInfixStringExpression(operator, left, right)
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

func evalIdentifierExpression(ident *ast.Identifier, env *object.Environment) object.Object {
	if value, ok := env.Get(ident.Value); ok {
		return value
	}
	if builtin, ok := object.LookUpBuiltin(ident.Value); ok {
		return builtin
	}
	return newError("reference error: %s is not defined", ident.Value)
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expression := range expressions {
		argument := Eval(expression, env)
		if isError(argument) {
			return []object.Object{argument}
		}
		result = append(result, argument)
	}

	return result
}

func applyFunction(function object.Object, arguments []object.Object) object.Object {
	if funcObj, ok := function.(*object.Function); ok {
		extendedEnv := extendFunctionEnvironment(funcObj, arguments)
		funcResult := evalBlockStatement(funcObj.Body.Statements, extendedEnv)
		return unwrapReturnValue(funcResult)
	}
	if builtin, ok := function.(*object.Builtin); ok {
		return builtin.Fn(arguments...)
	}
	return newError("not a function")
}

func extendFunctionEnvironment(function *object.Function, arguments []object.Object) *object.Environment {
	extendedEnv := object.NewEnclosedEnvironment(function.Env)
	// TODO: Check if arguments are less than actual Parameters
	for index, param := range function.Parameters {
		extendedEnv.Set(param.Value, arguments[index])
	}
	return extendedEnv
}

func unwrapReturnValue(result object.Object) object.Object {
	if result.Type() == object.RETURN {
		returnVal, _ := result.(*object.Return)
		return returnVal.Value
	}
	return result
}

func evalArrayIndexExpression(container *object.Array, indexExpression ast.Node, env *object.Environment) object.Object {
	indexObj := Eval(indexExpression, env)
	if indexObj.Type() != object.INT {
		return newError("type error: %s cannot be used as index of %s",
			indexObj.Type(), object.ARRAY)
	}
	index := indexObj.(*object.Integer)
	indexValue := index.Value
	if indexValue >= 0 && indexValue < int64(len(container.Items)) {
		return container.Items[indexValue]
	}
	return object.NULL
}

func evalHashLiteralExpression(astPairs map[ast.Expression]ast.Expression, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	hash := &object.Hash{Pairs: pairs}

	for keyNode, valueNode := range astPairs {
		evalKey := Eval(keyNode, env)
		if isError(evalKey) {
			return evalKey
		}
		key, ok := evalKey.(object.Hashable)
		if !ok {
			return newError("Got unhashable type as hash key: %s", evalKey.Type())
		}
		evalValue := Eval(valueNode, env)
		if isError(evalValue) {
			return evalValue
		}
		hashed := key.HashKey()
		pairs[hashed] = object.HashPair{Key: evalKey, Value: evalValue}
	}

	return hash
}

func evalIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object {
	container := Eval(node.Container, env)
	switch obj := container.(type) {
	case *object.Array:
		return evalArrayIndexExpression(obj, node.Index, env)
	default:
		return newError("type error: %s cannot be used as index expression", container.Type())
	}
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
		return evalIdentifierExpression(node, env)
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
	case *ast.CallExpression:
		{
			evalFunc := Eval(node.Function, env)
			if isError(evalFunc) {
				return evalFunc
			}
			evalArgs := evalExpressions(node.Arguments, env)
			if len(evalArgs) == 1 && isError(evalArgs[0]) {
				return evalArgs[0]
			}

			return applyFunction(evalFunc, evalArgs)
		}
	case *ast.IndexExpression:
		return evalIndexExpression(node, env)
	case *ast.IfExpression:
		return evalIfConditionalExpression(node, env)
	case *ast.IntegerLiteral:
		return object.NewInteger(node.Value)
	case *ast.BooleanLiteral:
		return booleanToObject(node.Value)
	case *ast.FunctionLiteral:
		return object.NewFunction(node.Parameters, node.Body, env)
	case *ast.StringLiteral:
		return object.NewString(node.Value)
	case *ast.ArrayLiteral:
		evalItems := evalExpressions(node.Items, env)
		if len(evalItems) == 1 && isError(evalItems[0]) {
			return evalItems[0]
		}
		return object.NewArray(evalItems)
	case *ast.HashLiteral:
		return evalHashLiteralExpression(node.Pairs, env)
	}

	return object.NULL
}
