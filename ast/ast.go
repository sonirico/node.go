package ast

import (
	"bytes"
	"fmt"
	"node.go/token"
	"strings"
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

// BLOCK statement
type BlockStatement struct {
	Token token.Token

	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStatement) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	if len(bs.Statements) > 0 {
		for _, stmt := range bs.Statements {
			buffer.WriteString(stmt.String())
		}
	}
	buffer.WriteString("}")
	return buffer.String()
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

	if pe.Operator != "" {
		buffer.WriteString("(")
		buffer.WriteString(pe.Operator)
		buffer.WriteString(pe.Right.String())
		buffer.WriteString(")")
	} else {
		buffer.WriteString(pe.Right.String())
	}

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

// IF expression
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IfExpression) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("if (")
	buffer.WriteString(ie.Condition.String())
	buffer.WriteString(")")

	if len(ie.Consequence.Statements) > 0 {
		buffer.WriteString(" ")
		buffer.WriteString(ie.Consequence.String())
	}

	if ie.Alternative != nil && len(ie.Alternative.Statements) > 0 {
		buffer.WriteString(" else ")
		buffer.WriteString(ie.Alternative.String())
	}

	return buffer.String()
}

type IndexExpression struct {
	Token     token.Token
	Container Expression
	Index     Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Container.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

// LITERALS

// Function literal

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}
func (fl *FunctionLiteral) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("fn")
	buffer.WriteString("(")
	if fl.Parameters != nil && len(fl.Parameters) > 0 {
		var params []string
		for _, param := range fl.Parameters {
			params = append(params, param.String())
		}
		buffer.WriteString(strings.Join(params, ", "))
	}
	buffer.WriteString(")")
	buffer.WriteString(" ")
	buffer.WriteString(fl.Body.String())

	return buffer.String()
}

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

// String literal
type StringLiteral struct {
	Token token.Token

	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}
func (sl *StringLiteral) String() string {
	return fmt.Sprintf(`\"%s\"`, sl.Value)
}

// Call expression
type CallExpression struct {
	Token token.Token

	Function Expression

	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
func (ce *CallExpression) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(ce.Function.String())
	buffer.WriteString("(")
	var arguments []string
	if ce.Arguments != nil {
		for _, ParamExpression := range ce.Arguments {
			arguments = append(arguments, ParamExpression.String())
		}
	}
	buffer.WriteString(strings.Join(arguments, ", "))
	buffer.WriteString(")")

	return buffer.String()
}

// ARRAY LITERAL
type ArrayLiteral struct {
	Token token.Token

	Items []Expression
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	var items []string

	for _, item := range al.Items {
		items = append(items, item.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(items, ", "))
	out.WriteString("]")

	return out.String()
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
