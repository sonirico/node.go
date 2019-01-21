package ast

import "node.go/token"

type Node interface {
	TokenLiteral() string // Token literal value
}

type Statement interface {
	Node

	statementNode()
}

type Expression interface {
	Node

	expressionNode()
}

type LetStatement struct {
	Token token.Token

	Name *Identifier

	Value string
}

func (l *LetStatement) statementNode() {}
func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
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

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
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
