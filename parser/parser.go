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
	EQUALS      // ==
	LESSGREATER // >, <, <=, >=
	SUM         // + and -
	MODULE      // %
	PRODUCT     // * and /
	POWER       // ^
	PREFIX      // ! or -
	CALL        // add(1, 2)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LTE:      LESSGREATER,
	token.GTE:      LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.PERCENT:  MODULE,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.POWER:    POWER,
	token.LPAREN:   CALL,
}

func getPrecedence(tokenType token.TokenType) int {
	if precedence, ok := precedences[tokenType]; ok {
		return precedence
	}

	// Token is not bound to a precedence
	return LOWEST
}

type (
	PrefixParserFunc func() ast.Expression
	InfixParserFunc  func(ast.Expression) ast.Expression
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
	parser.registerPrefixFunction(token.TRUE, parser.parseBooleanLiteral)
	parser.registerPrefixFunction(token.FALSE, parser.parseBooleanLiteral)
	parser.registerPrefixFunction(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefixFunction(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefixFunction(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefixFunction(token.LPAREN, parser.parseGroupedExpression)
	parser.registerPrefixFunction(token.IF, parser.parseIfExpression)
	parser.registerPrefixFunction(token.FUNC, parser.parseFunctionExpression)

	// Infix parsers
	parser.registerInfixFunction(token.PLUS, parser.parseInfixExpression)
	parser.registerInfixFunction(token.MINUS, parser.parseInfixExpression)
	parser.registerInfixFunction(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfixFunction(token.SLASH, parser.parseInfixExpression)
	parser.registerInfixFunction(token.PERCENT, parser.parseInfixExpression)
	parser.registerInfixFunction(token.POWER, parser.parseInfixExpression)
	parser.registerInfixFunction(token.EQ, parser.parseInfixExpression)
	parser.registerInfixFunction(token.NOT_EQ, parser.parseInfixExpression)
	parser.registerInfixFunction(token.LT, parser.parseInfixExpression)
	parser.registerInfixFunction(token.GT, parser.parseInfixExpression)
	parser.registerInfixFunction(token.LTE, parser.parseInfixExpression)
	parser.registerInfixFunction(token.GTE, parser.parseInfixExpression)
	parser.registerInfixFunction(token.PERCENT, parser.parseInfixExpression)
	parser.registerInfixFunction(token.LPAREN, parser.parseCallExpression)

	return parser
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) currentPrecedence() int {
	return getPrecedence(p.currentToken.Type)
}

func (p *Parser) peekPrecedence() int {
	return getPrecedence(p.peekToken.Type)
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

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
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

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token: p.currentToken,
		Left:  left,
	}

	expr.Operator = p.currentToken.Literal

	precedence := p.currentPrecedence()

	p.nextToken()

	expr.Right = p.parseExpression(precedence)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		return expr
	}

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

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.currentToken, Value: p.currentTokenIs(token.TRUE)}
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

func (p *Parser) parseGroupedExpression() ast.Expression {
	exp := &ast.PrefixExpression{Token: p.currentToken, Operator: ""}

	p.nextToken()

	exp.Right = p.parseExpression(LOWEST)

	if !p.expectPeekToken(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseBlockStatement() ast.BlockStatement {
	bs := ast.BlockStatement{Token: p.currentToken}
	bs.Statements = []ast.Statement{}

	p.nextToken()

	for !p.currentTokenIs(token.RBRACE) {
		stmt := p.parseExpressionStatement()
		if stmt != nil {
			bs.Statements = append(bs.Statements, stmt)
		}
		p.nextToken()
	}

	return bs
}

func (p *Parser) parseIfExpression() ast.Expression {
	ifExp := &ast.IfExpression{Token: p.currentToken}

	if !p.expectPeekToken(token.LPAREN) {
		return nil
	}

	p.nextToken()

	ifExp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeekToken(token.RPAREN) {
		return nil
	}

	if !p.expectPeekToken(token.LBRACE) {
		return nil
	}

	ifExp.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeekToken(token.LBRACE) {
			return nil
		}

		stmt := p.parseBlockStatement()
		ifExp.Alternative = &stmt

		p.nextToken()
	}

	return ifExp
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var params []*ast.Identifier

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()
	ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	params = append(params, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		params = append(
			params,
			&ast.Identifier{
				Token: p.currentToken,
				Value: p.currentToken.Literal,
			})
	}

	if !p.expectPeekToken(token.RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	funcExp := &ast.FunctionLiteral{Token: p.currentToken}

	if !p.expectPeekToken(token.LPAREN) {
		return nil
	}

	funcExp.Parameters = p.parseFunctionParameters()

	if !p.expectPeekToken(token.LBRACE) {
		return nil
	}

	funcExp.Body = p.parseBlockStatement()

	return funcExp
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	callExp := &ast.CallExpression{Token: p.currentToken, Function: function}
	callExp.Arguments = []ast.Expression{}

	if !p.currentTokenIs(token.LPAREN) {
		return nil
	}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return callExp
	}

	p.nextToken()

	callExp.Arguments = append(callExp.Arguments, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		callExp.Arguments = append(callExp.Arguments, p.parseExpression(LOWEST))
	}

	if !p.expectPeekToken(token.RPAREN) {
		return nil
	}

	return callExp
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

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infixParserFunc := p.infixParserFunctions[p.peekToken.Type]

		if infixParserFunc == nil {
			return leftExpression
		}

		p.nextToken()
		leftExpression = infixParserFunc(leftExpression)
	}

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
