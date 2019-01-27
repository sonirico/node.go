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

// PREFIX expression
type PrefixExpression struct {
	Token token.Token

	Operator string

	Right Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("(")
	buffer.WriteString(pe.Operator)
	buffer.WriteString(pe.Right.String())
	buffer.WriteString(")")

	return buffer.String()
}

// INFIX expression
type InfixExpression struct {
	Token token.Token

	Left Expression

	Operator string

	Right Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("(")
	buffer.WriteString(ie.Left.String())
	buffer.WriteString(" ")
	buffer.WriteString(ie.Operator)
	buffer.WriteString(" ")
	buffer.WriteString(ie.Right.String())
	buffer.WriteString(")")

	return buffer.String()
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

// LITERALS

// Boolean literal
type BooleanLiteral struct {
	Token token.Token

	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) TokenLiteral() string {
	return bl.Token.Literal
}
func (bl *BooleanLiteral) String() string {
	return bl.TokenLiteral()
}

// Integer literal
type IntegerLiteral struct {
	Token token.Token

	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// PROGRAM - The root node!
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
