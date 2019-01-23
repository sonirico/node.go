package parser

import (
	"fmt"
	"node.go/ast"
	"node.go/lexer"
	"node.go/token"
	"strconv"
)

// Precedence priorities
const (
	_ int = iota
	LOWEST
	EQUALS  // ==
	SUM     // + and -
	PRODUCT // * and /
	POWER   // ^
	PREFIX  // ! and - E.g: !true, -2
	CALL    // add(1, 2)
)

type (
	PrefixParserFunc func() ast.Expression
	InfixParserFunc  func() ast.Expression
)

type Parser struct {
	lexer *lexer.Lexer

	errors []string

	currentToken token.Token
	peekToken    token.Token

	prefixParserFunctions map[token.TokenType]PrefixParserFunc
	infixParserFunctions  map[token.TokenType]InfixParserFunc
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:  lexer,
		errors: []string{},
	}
	// Read to token so as to have initialised both currentToken and peekToken
	parser.nextToken()
	parser.nextToken()

	parser.prefixParserFunctions = make(map[token.TokenType]PrefixParserFunc)
	parser.infixParserFunctions = make(map[token.TokenType]InfixParserFunc)

	// Prefix parsers
	parser.registerPrefixFunction(token.IDENTIFIER, parser.parseIdentifierExpression)
	parser.registerPrefixFunction(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefixFunction(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefixFunction(token.MINUS, parser.parsePrefixExpression)

	// Infix parsers
	parser.registerInfixFunction(token.INT, parser.parseIntegerLiteral)

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
	msg := fmt.Sprintf("Expected next token to be of type '%s'. Got '%s' -> %s",
		tokenType, p.peekToken.Type, p.peekToken.Literal)
	p.addError(msg)
}

func (p *Parser) addError(msg string) {
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

func (p *Parser) registerPrefixFunction(tokenType token.TokenType, fn PrefixParserFunc) {
	p.prefixParserFunctions[tokenType] = fn
}

func (p *Parser) registerInfixFunction(tokenType token.TokenType, fn InfixParserFunc) {
	p.infixParserFunctions[tokenType] = fn
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

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	// Well, there is no expression parsing yet xD
	// TODO: Implement expression parsing.
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	// 1 + 1; <-- This is and expression statement. Produces a value!
	// 1 + 1 <-- This is also a valid expression. Rather handy to skip the semicolon in the repl!

	// There is no other way but descending recursively into the ast tree...
	stmt.Expression = p.parseExpression(LOWEST)

	// Make the semicolon token optional
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	expr.Right = p.parseExpression(PREFIX)

	return expr
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseIdentifierExpression() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	il := &ast.IntegerLiteral{Token: p.currentToken}

	base := 10
	value, err := strconv.ParseInt(p.currentToken.Literal, base, 64)

	if err != nil {
		msg := fmt.Sprintf("unable to parse '%s' as integer", p.currentToken.Literal)
		p.addError(msg)
		return nil
	}

	il.Value = value

	return il
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	tokenType := p.currentToken.Type
	prefixParserFunction := p.prefixParserFunctions[tokenType]

	if prefixParserFunction == nil {
		msg := fmt.Sprintf("there is not registered prefix parser function for token type %q", tokenType)
		p.addError(msg)
		return nil
	}

	leftExpression := prefixParserFunction()

	return leftExpression
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
