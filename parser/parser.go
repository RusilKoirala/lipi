package parser

import (
	"github.com/rusilkoirala/lipi/ast"
	"github.com/rusilkoirala/lipi/lexer"
	"github.com/rusilkoirala/lipi/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x or !x
	CALL        // jod(x)
)

var precedences = map[token.TokenType]int{}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}
	p.NextToken()
}
