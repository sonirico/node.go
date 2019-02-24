package object

import (
	"bytes"
	"node.go/ast"
	"strings"
)

type Function struct {
	Parameters []*ast.Identifier
	Body       ast.BlockStatement
	Env        *Environment
}

func NewFunction(params []*ast.Identifier, body ast.BlockStatement, env *Environment) *Function {
	return &Function{Parameters: params, Body: body, Env: env}
}

func (f *Function) Type() Type {
	return FUNCTION
}

func (f *Function) Inspect() string {
	var buffer bytes.Buffer
	var params []string

	for _, param := range f.Parameters {
		params = append(params, param.String())
	}

	buffer.WriteString("fn")
	buffer.WriteString("(")
	buffer.WriteString(strings.Join(params, ", "))
	buffer.WriteString(")")
	buffer.WriteString(" ")
	buffer.WriteString(f.Body.String())

	return buffer.String()
}
