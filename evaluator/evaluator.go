package evaluator

import (
	"fmt"

	"github.com/rusilkoirala/lipi/ast"
	"github.com/rusilkoirala/lipi/object"
)

// There is only ever one sahi, one galat and one khali.
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

// A function calling itself forever would crash Go, so we stop it here.
// The browser playground sets this lower because its stack is smaller.
var MaxCallDepth = 10000

var callDepth = 0

// Eval walks the tree and runs it.
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// ----- statements -----
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlock(node, env)

	case *ast.RakhaStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return NULL

	case *ast.FarkauStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	// ----- simple values -----
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return toBool(node.Value)

	case *ast.Identifier:
		if val, ok := env.Get(node.Value); ok {
			return val
		}
		if b, ok := builtins[node.Value]; ok {
			return b
		}
		return newError(node.Token.Line, "'%s' bhanne naam bhetiyena", node.Value)

	// ----- operators + yadi -----
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefix(node, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfix(node, left, right)

	case *ast.IfExpression:
		return evalIf(node, env)

	// ----- functions -----
	case *ast.FunctionLiteral:
		return &object.Function{Name: node.Name, Parameters: node.Parameters, Body: node.Body, Env: env}

	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		args := []object.Object{}
		for _, a := range node.Arguments {
			val := Eval(a, env)
			if isError(val) {
				return val
			}
			args = append(args, val)
		}
		return applyFunction(node, fn, args)
	}

	return NULL
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = NULL

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		switch r := result.(type) {
		case *object.ReturnValue:
			return r.Value // farkau at top level: stop the program
		case *object.Error:
			return r // Galti: stop the program
		}
	}
	return result
}

func evalBlock(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object = NULL

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		// keep farkau wrapped so it bubbles up through nested blocks
		if result.Type() == object.RETURN || result.Type() == object.ERROR {
			return result
		}
	}
	return result
}

// !x  -x
func evalPrefix(node *ast.PrefixExpression, right object.Object) object.Object {
	switch node.Operator {
	case "!":
		return toBool(!isTruthy(right))
	case "-":
		n, ok := right.(*object.Integer)
		if !ok {
			return newError(node.Token.Line, "%s lai '-' garna mildaina", right.Type())
		}
		return &object.Integer{Value: -n.Value}
	}
	return newError(node.Token.Line, "'%s' bujhina", node.Operator)
}

// a + b, a > b, "Namaste, " + naam ...
func evalInfix(node *ast.InfixExpression, left, right object.Object) object.Object {
	line := node.Token.Line
	op := node.Operator

	// sankhya op sankhya
	l, lok := left.(*object.Integer)
	r, rok := right.(*object.Integer)
	if lok && rok {
		return evalIntegerInfix(line, op, l.Value, r.Value)
	}

	// string op string
	ls, lok := left.(*object.String)
	rs, rok := right.(*object.String)
	if lok && rok {
		switch op {
		case "+":
			return &object.String{Value: ls.Value + rs.Value}
		case "==":
			return toBool(ls.Value == rs.Value)
		case "!=":
			return toBool(ls.Value != rs.Value)
		}
		return newError(line, "string ma '%s' chaldaina", op)
	}

	// anything else (sahi/galat, mixed types)
	switch op {
	case "==":
		return toBool(left == right)
	case "!=":
		return toBool(left != right)
	case "+":
		return newError(line, "%s ra %s jodna mildaina", left.Type(), right.Type())
	}
	return newError(line, "%s %s %s garna mildaina", left.Type(), op, right.Type())
}

func evalIntegerInfix(line int, op string, a, b int64) object.Object {
	switch op {
	case "+":
		return &object.Integer{Value: a + b}
	case "-":
		return &object.Integer{Value: a - b}
	case "*":
		return &object.Integer{Value: a * b}
	case "/":
		if b == 0 {
			return newError(line, "0 le bhag garna mildaina")
		}
		return &object.Integer{Value: a / b}
	case "<":
		return toBool(a < b)
	case ">":
		return toBool(a > b)
	case "==":
		return toBool(a == b)
	case "!=":
		return toBool(a != b)
	}
	return newError(line, "sankhya ma '%s' chaldaina", op)
}

// yadi (cond) { ... } natra { ... }
func evalIf(node *ast.IfExpression, env *object.Environment) object.Object {
	cond := Eval(node.Condition, env)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(node.Consequence, env)
	}
	if node.Alternative != nil {
		return Eval(node.Alternative, env)
	}
	return NULL
}

// jod(2, 3)
func applyFunction(node *ast.CallExpression, fn object.Object, args []object.Object) object.Object {
	line := node.Token.Line

	switch fn := fn.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return newError(line, "'%s' lai %d ota input chahinchha, %d ota aayo",
				node.Function.String(), len(fn.Parameters), len(args))
		}

		if callDepth >= MaxCallDepth {
			return newError(line, "'%s' le aafailai dherai patak bolayo (recursion rokiyena)", node.Function.String())
		}
		callDepth++
		defer func() { callDepth-- }()

		// new scope for this call: a = 2, b = 3
		callEnv := object.NewEnclosedEnvironment(fn.Env)
		for i, param := range fn.Parameters {
			callEnv.Set(param.Value, args[i])
		}

		result := Eval(fn.Body, callEnv)
		if rv, ok := result.(*object.ReturnValue); ok {
			return rv.Value // unwrap farkau here so it stops at the function
		}
		return result

	case *object.Builtin:
		return fn.Fn(line, args...)
	}

	return newError(line, "'%s' kaam hoina, bolauna mildaina", node.Function.String())
}

// ---------- helpers ----------

func toBool(b bool) *object.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

// only galat and khali count as "false"
func isTruthy(obj object.Object) bool {
	return obj != FALSE && obj != NULL
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR
}

func newError(line int, format string, a ...any) *object.Error {
	msg := fmt.Sprintf(format, a...)
	return &object.Error{Message: fmt.Sprintf("Galti: line %d ma %s", line, msg)}
}
