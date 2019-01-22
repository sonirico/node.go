package ast

import (
	"bytes"
	"fmt"
	"node.go/token"
)

type Node interface {
	TokenLiteral() string // Token literal value

	String() string
}

type Statement interface {
	Node

	statementNode()
}

type Expression interface {
	Node

	expressionNode()
}

// LET statement
type LetStatement struct {
	Token token.Token

	Name *Identifier

	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(ls.TokenLiteral())
	buffer.WriteString(" ")
	buffer.WriteString(ls.Name.String())
	buffer.WriteString(fmt.Sprintf(" %s ", token.ASSIGNMENT))

	if ls.Value != nil {
		buffer.WriteString(ls.Value.String())
	}

	buffer.WriteString(";")

	return buffer.String()
}

// RETURN statement
type ReturnStatement struct {
	Token token.Token

	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReturnStatement) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(rs.TokenLiteral())
	buffer.WriteString(" ")

	if rs.ReturnValue != nil {
		buffer.WriteString(rs.ReturnValue.String())
	}

	buffer.WriteString(";")

	return buffer.String()
}

// STATEMENT expression
type ExpressionStatement struct {
	Token token.Token

	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IDENTIFIER expression
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Value
}
func (i *Identifier) String() string {
	return i.Value
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string {
	var buffer bytes.Buffer

	for _, stmt := range p.Statements {
		buffer.WriteString(stmt.String())
	}

	return buffer.String()
}
