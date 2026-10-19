package ast

import (
	"strconv"
	"testing"

	"github.com/rusilkoirala/lipi/token"
)

// small helpers so the test trees are readable

func ident(name string) *Identifier {
	return &Identifier{Token: token.Token{Type: token.IDENT, Literal: name}, Value: name}
}

func num(v int64) *IntegerLiteral {
	return &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: strconv.FormatInt(v, 10)}, Value: v}
}

func str(s string) *StringLiteral {
	return &StringLiteral{Token: token.Token{Type: token.STRING, Literal: s}, Value: s}
}

func infix(left Expression, op string, right Expression) *InfixExpression {
	return &InfixExpression{Token: token.Token{Type: token.TokenType(op), Literal: op}, Left: left, Operator: op, Right: right}
}

func call(fn Expression, args ...Expression) *CallExpression {
	return &CallExpression{Token: token.Token{Type: token.LPAREN, Literal: "("}, Function: fn, Arguments: args}
}

func block(stmts ...Statement) *BlockStatement {
	return &BlockStatement{Token: token.Token{Type: token.LBRACE, Literal: "{"}, Statements: stmts}
}

func exprStmt(e Expression) *ExpressionStatement {
	return &ExpressionStatement{Token: token.Token{Literal: e.TokenLiteral()}, Expression: e}
}

func TestString(t *testing.T) {
	tests := []struct {
		node     Node
		expected string
	}{
		// rakha naam = "Ram";
		{
			&RakhaStatement{
				Token: token.Token{Type: token.LET, Literal: "rakha"},
				Name:  ident("naam"),
				Value: str("Ram"),
			},
			`rakha naam = "Ram";`,
		},
		// kaam jod(a, b) { farkau a + b; }
		{
			&FunctionLiteral{
				Token:      token.Token{Type: token.FUNCTION, Literal: "kaam"},
				Name:       "jod",
				Parameters: []*Identifier{ident("a"), ident("b")},
				Body: block(&FarkauStatement{
					Token:       token.Token{Type: token.RETURN, Literal: "farkau"},
					ReturnValue: infix(ident("a"), "+", ident("b")),
				}),
			},
			`kaam jod(a, b) farkau (a + b);`,
		},
		// yadi (jod(2, 3) > 4) { bhana("thulo!"); } natra { bhana("sano"); }
		{
			&IfExpression{
				Token:       token.Token{Type: token.IF, Literal: "yadi"},
				Condition:   infix(call(ident("jod"), num(2), num(3)), ">", num(4)),
				Consequence: block(exprStmt(call(ident("bhana"), str("thulo!")))),
				Alternative: block(exprStmt(call(ident("bhana"), str("sano")))),
			},
			`yadi (jod(2, 3) > 4) bhana("thulo!") natra bhana("sano")`,
		},
	}

	for i, tt := range tests {
		if got := tt.node.String(); got != tt.expected {
			t.Errorf("tests[%d] wrong String().\nexpected: %s\n     got: %s", i, tt.expected, got)
		}
	}
}
