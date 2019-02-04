package evaluator

import (
	"fmt"
	"node.go/ast"
	"node.go/object"
)

func Eval(node ast.Node) object.Object {
	var result object.Object
	switch node := node.(type) {
	case *ast.Program:
		{
			for _, stmt := range node.Statements {
				result = Eval(stmt)
			}
		}
	case *ast.ExpressionStatement:
		result = Eval(node.Expression)
		break
	case *ast.IntegerLiteral:
		result = object.NewInteger(node.Value)
		break
	}
	if result == nil {
		fmt.Println(fmt.Sprintf("There is no an evaluator function for %q", node))
	}
	return result
}
