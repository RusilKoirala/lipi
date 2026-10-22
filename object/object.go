package object

import (
	"fmt"
	"strings"

	"github.com/rusilkoirala/lipi/ast"
)

type ObjectType string

// these names show up in Galti messages
const (
	INTEGER  = "sankhya"
	STRING   = "string"
	BOOLEAN  = "sahi/galat"
	NULL     = "khali"
	RETURN   = "farkau"
	ERROR    = "galti"
	FUNCTION = "kaam"
	BUILTIN  = "builtin"
)

// Every value in Lipi is an Object.
type Object interface {
	Type() ObjectType
	Inspect() string // how bhana() prints it
}

// 5
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// "Ram"
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING }
func (s *String) Inspect() string  { return s.Value }

// sahi / galat
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN }
func (b *Boolean) Inspect() string {
	if b.Value {
		return "sahi"
	}
	return "galat"
}

// "nothing" (e.g. a yadi with no natra that didn't run)
type Null struct{}

func (n *Null) Type() ObjectType { return NULL }
func (n *Null) Inspect() string  { return "khali" }

// wraps the value from farkau so it can travel up and out of blocks
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Galti: line 3 ma ...
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR }
func (e *Error) Inspect() string  { return e.Message }

// kaam jod(a, b) { ... }
type Function struct {
	Name       string
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment // the scope where the function was made (closures!)
}

func (f *Function) Type() ObjectType { return FUNCTION }
func (f *Function) Inspect() string {
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.Value)
	}
	return "kaam " + f.Name + "(" + strings.Join(params, ", ") + ")"
}

// functions written in Go: bhana, lambai
type BuiltinFunction func(line int, args ...Object) Object

type Builtin struct {
	Name string
	Fn   BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN }
func (b *Builtin) Inspect() string  { return "builtin " + b.Name }
