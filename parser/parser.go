package parser

import (
	"fmt"
	"strconv"

	"github.com/rusilkoirala/lipi/ast"
	"github.com/rusilkoirala/lipi/lexer"
	"github.com/rusilkoirala/lipi/token"
)

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
	}
	return nil
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
	stmt.Value = p.parseValue()

	p.expectSemicolon()
	return stmt
}

// farkau x;
func (p *Parser) parseFarkau() ast.Statement {
	stmt := &ast.FarkauStatement{Token: p.curToken}

	p.nextToken()
	stmt.ReturnValue = p.parseValue()

	p.expectSemicolon()
	return stmt
}

// For now only simple values. Part 5 replaces this with the Pratt parser.
func (p *Parser) parseValue() ast.Expression {
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
	}
	p.addError(p.curToken.Line, fmt.Sprintf("'%s' bujhina", p.curToken.Literal))
	return nil
}
