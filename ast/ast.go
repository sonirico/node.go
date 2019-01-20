package ast

import "node.go/token"

type Node interface {
	Literal() string
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
func (l *LetStatement) Literal() string {
	return l.Token.Literal
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) Literal() string {
	return i.Value
}

type Program struct {
	Statements []Statement
}

func (p *Program) Literal() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].Literal()
	} else {
		return ""
	}
}
