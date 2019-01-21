package parser

import (
	"fmt"
	"node.go/ast"
	"node.go/lexer"
	"node.go/token"
)

type Parser struct {
	lexer *lexer.Lexer

	errors []string

	currentToken token.Token
	peekToken    token.Token
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:  lexer,
		errors: []string{},
	}
	// Read to token so as to have initialised both currentToken and peekToken
	parser.nextToken()
	parser.nextToken()

	return parser
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) currentTokenIs(tokenType token.TokenType) bool {
	return p.currentToken.Type == tokenType
}

func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) peekError(tokenType token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be '%s'. Got '%s'",
		tokenType, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeekToken(tokenType token.TokenType) bool {
	if p.peekTokenIs(tokenType) {
		p.nextToken()
		return true
	} else {
		p.peekError(tokenType)
		return false
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeekToken(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	// Empty Let definitions
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		return stmt
	}

	if !p.expectPeekToken(token.ASSIGNMENT) {
		return nil
	}

	// TODO: Implement expression parsing
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}

	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}

	return program
}
