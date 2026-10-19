package ast

import (
	"strings"

	"github.com/rusilkoirala/lipi/token"
)

// Every node in the tree implements Node.
type Node interface {
	TokenLiteral() string
	String() string // prints the node back as code for debugging + testing
}

// Statements do something: rakha x = 5;  farkau x;
type Statement interface {
	Node
	statementNode()
}

// Expressions produce a value: 5, a + b, jod(2, 3)
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node: a list of statements.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out strings.Builder
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// ---------- Statements ----------

// rakha <Name> = <Value>;
type RakhaStatement struct {
	Token token.Token // the 'rakha' token
	Name  *Identifier
	Value Expression
}

func (rs *RakhaStatement) statementNode()       {}
func (rs *RakhaStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *RakhaStatement) String() string {
	var out strings.Builder
	out.WriteString(rs.TokenLiteral() + " ")
	out.WriteString(rs.Name.String())
	out.WriteString(" = ")
	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// farkau <ReturnValue>;
type FarkauStatement struct {
	Token       token.Token // the 'farkau' token
	ReturnValue Expression
}

func (fs *FarkauStatement) statementNode()       {}
func (fs *FarkauStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FarkauStatement) String() string {
	var out strings.Builder
	out.WriteString(fs.TokenLiteral() + " ")
	if fs.ReturnValue != nil {
		out.WriteString(fs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// A line that is just an expression: bhana("hi");  or  x + 5;
type ExpressionStatement struct {
	Token      token.Token // first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// { ...statements... }
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out strings.Builder
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// ---------- Simple expressions ----------

// naam, jod, a, b
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// 5, 42
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// "Ram"
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return `"` + sl.Value + `"` }

// sahi, galat
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// ---------- Operator expressions ----------

// !sahi, -5
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}

// a + b, x > 4
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

// ---------- Control flow + functions ----------

// yadi (<Condition>) { <Consequence> } natra { <Alternative> }
type IfExpression struct {
	Token       token.Token // the 'yadi' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // nil if there is no natra
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out strings.Builder
	out.WriteString("yadi ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString(" natra ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// kaam(a, b) { ... }   or named:  kaam jod(a, b) { ... }
type FunctionLiteral struct {
	Token      token.Token // the 'kaam' token
	Name       string      // "" for anonymous functions
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	var out strings.Builder
	out.WriteString(fl.TokenLiteral())
	if fl.Name != "" {
		out.WriteString(" " + fl.Name)
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

// jod(2, 3), bhana("hi")
type CallExpression struct {
	Token     token.Token // the '(' token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	return ce.Function.String() + "(" + strings.Join(args, ", ") + ")"
}
