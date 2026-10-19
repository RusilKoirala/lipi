package parser

import (
	"fmt"
	"strconv"

	"github.com/rusilkoirala/lipi/ast"
	"github.com/rusilkoirala/lipi/lexer"
	"github.com/rusilkoirala/lipi/token"
)

// operator strength: bigger number = stronger
const (
	LOWEST  = iota
	COMPARE // == != < >
	SUM     // + -
	PRODUCT // * /
	PREFIX  // -x !x
	CALL    // jod(x)
)

func precedenceOf(t token.TokenType) int {
	switch t {
	case token.EQ, token.NOT_EQ, token.LT, token.GT:
		return COMPARE
	case token.PLUS, token.MINUS:
		return SUM
	case token.ASTERISK, token.SLASH:
		return PRODUCT
	case token.LPAREN:
		return CALL
	}
	return LOWEST
}

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token // the token we are looking at
	peekToken token.Token // the token after it
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Galti: line 3 ma ...
func (p *Parser) addError(line int, msg string) {
	p.errors = append(p.errors, fmt.Sprintf("Galti: line %d ma %s", line, msg))
}

// move forward only if the next token is the one we expect
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.addError(p.peekToken.Line, fmt.Sprintf("'%s' chahiyo", t))
	return false
}

func (p *Parser) expectSemicolon() {
	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
		return
	}
	p.addError(p.curToken.Line, "semicolon birsiyo")
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseRakha()
	case token.RETURN:
		return p.parseFarkau()
	case token.FUNCTION:
		if p.peekToken.Type == token.IDENT {
			return p.parseNamedFunction()
		}
	}
	return p.parseExpressionStatement()
}

// 5 + 5;
func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if stmt.Expression == nil {
		return nil
	}
	if p.curToken.Type == token.RBRACE {
		// yadi {...} natra {...} ended with '}', so ';' is optional
		if p.peekToken.Type == token.SEMICOLON {
			p.nextToken()
		}
		return stmt
	}
	p.expectSemicolon()
	return stmt
}

// rakha x = 5;
func (p *Parser) parseRakha() ast.Statement {
	stmt := &ast.RakhaStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	p.expectSemicolon()
	return stmt
}

// farkau x;
func (p *Parser) parseFarkau() ast.Statement {
	stmt := &ast.FarkauStatement{Token: p.curToken}

	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)

	p.expectSemicolon()
	return stmt
}

// The Pratt parser: read one value, then keep grabbing operators
// as long as they are stronger than the current one.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	left := p.parsePrefix()
	if left == nil {
		return nil
	}

	for precedence < precedenceOf(p.peekToken.Type) {
		p.nextToken()
		if p.curToken.Type == token.LPAREN {
			left = p.parseCall(left)
		} else {
			left = p.parseInfix(left)
		}
		if left == nil {
			return nil
		}
	}
	return left
}

// things that START an expression: 5, "Ram", sahi, x, -x, !x, ( ... )
func (p *Parser) parsePrefix() ast.Expression {
	switch p.curToken.Type {
	case token.INT:
		n, _ := strconv.ParseInt(p.curToken.Literal, 10, 64)
		return &ast.IntegerLiteral{Token: p.curToken, Value: n}
	case token.STRING:
		return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	case token.TRUE, token.FALSE:
		return &ast.Boolean{Token: p.curToken, Value: p.curToken.Type == token.TRUE}
	case token.IDENT:
		return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	case token.MINUS, token.BANG:
		expr := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
		p.nextToken()
		expr.Right = p.parseExpression(PREFIX)
		if expr.Right == nil {
			return nil
		}
		return expr

	case token.LPAREN:
		p.nextToken()
		expr := p.parseExpression(LOWEST)
		if !p.expectPeek(token.RPAREN) {
			return nil
		}
		return expr

	case token.IF:
		return p.parseIf()
	case token.FUNCTION:
		return p.parseFunction()
	}

	p.addError(p.curToken.Line, fmt.Sprintf("'%s' bujhina", p.curToken.Literal))
	return nil
}

// operator in the MIDDLE: left + right
func (p *Parser) parseInfix(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	prec := precedenceOf(p.curToken.Type)

	p.nextToken()
	expr.Right = p.parseExpression(prec)
	if expr.Right == nil {
		return nil
	}
	return expr
}

// { statements }
func (p *Parser) parseBlock() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	p.nextToken()

	for p.curToken.Type != token.RBRACE {
		if p.curToken.Type == token.EOF {
			p.addError(block.Token.Line, "'{' kholiyo tara '}' le band gariyena")
			return nil
		}
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

// yadi (x > 4) { ... } natra { ... }
func (p *Parser) parseIf() ast.Expression {
	expr := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expr.Condition = p.parseExpression(LOWEST)
	if expr.Condition == nil {
		return nil
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expr.Consequence = p.parseBlock()
	if expr.Consequence == nil {
		return nil
	}

	if p.peekToken.Type == token.ELSE {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expr.Alternative = p.parseBlock()
		if expr.Alternative == nil {
			return nil
		}
	}
	return expr
}

// kaam(a, b) { ... }   or   kaam jod(a, b) { ... }
func (p *Parser) parseFunction() ast.Expression {
	fn := &ast.FunctionLiteral{Token: p.curToken}

	if p.peekToken.Type == token.IDENT { // named function
		p.nextToken()
		fn.Name = p.curToken.Literal
	}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	params, ok := p.parseParams()
	if !ok {
		return nil
	}
	fn.Parameters = params

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	fn.Body = p.parseBlock()
	if fn.Body == nil {
		return nil
	}
	return fn
}

// (a, b, c)  -- curToken is '('
func (p *Parser) parseParams() ([]*ast.Identifier, bool) {
	params := []*ast.Identifier{}

	if p.peekToken.Type == token.RPAREN { // kaam() with no params
		p.nextToken()
		return params, true
	}

	if !p.expectPeek(token.IDENT) {
		return nil, false
	}
	params = append(params, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})

	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		if !p.expectPeek(token.IDENT) {
			return nil, false
		}
		params = append(params, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	if !p.expectPeek(token.RPAREN) {
		return nil, false
	}
	return params, true
}

// jod(1, 2)  -- curToken is '('
func (p *Parser) parseCall(fn ast.Expression) ast.Expression {
	call := &ast.CallExpression{Token: p.curToken, Function: fn}

	if p.peekToken.Type == token.RPAREN { // jod() with no args
		p.nextToken()
		return call
	}

	p.nextToken()
	arg := p.parseExpression(LOWEST)
	if arg == nil {
		return nil
	}
	call.Arguments = append(call.Arguments, arg)

	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		arg := p.parseExpression(LOWEST)
		if arg == nil {
			return nil
		}
		call.Arguments = append(call.Arguments, arg)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return call
}

// kaam jod(a, b) { ... }   becomes   rakha jod = kaam jod(a, b) { ... };
func (p *Parser) parseNamedFunction() ast.Statement {
	tok := p.curToken
	fn := p.parseFunction()
	if fn == nil {
		return nil
	}
	name := fn.(*ast.FunctionLiteral).Name

	if p.peekToken.Type == token.SEMICOLON { // optional ';' after '}'
		p.nextToken()
	}

	return &ast.RakhaStatement{
		Token: token.Token{Type: token.LET, Literal: "rakha", Line: tok.Line},
		Name:  &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: name, Line: tok.Line}, Value: name},
		Value: fn,
	}
}
