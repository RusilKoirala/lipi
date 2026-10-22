//go:build js && wasm

// Exposes the Lipi interpreter to JavaScript for the online playground.
// Build: GOOS=js GOARCH=wasm go build -o website/public/lipi.wasm ./wasm
package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"github.com/rusilkoirala/lipi/evaluator"
	"github.com/rusilkoirala/lipi/lexer"
	"github.com/rusilkoirala/lipi/object"
	"github.com/rusilkoirala/lipi/parser"
)

// run executes one program and returns what bhana printed plus any Galti.
func run(code string) (output string, errs []string) {
	var out bytes.Buffer
	evaluator.Out = &out

	defer func() {
		if r := recover(); r != nil {
			output = out.String()
			errs = []string{fmt.Sprintf("Galti: interpreter bigriyo (%v)", r)}
		}
	}()

	p := parser.New(lexer.New(code))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		return out.String(), p.Errors()
	}

	result := evaluator.Eval(program, object.NewEnvironment())
	if e, ok := result.(*object.Error); ok {
		return out.String(), []string{e.Message}
	}
	return out.String(), nil
}

func main() {
	evaluator.MaxCallDepth = 1500 // the browser's stack is smaller than native Go's

	// runLipi(code) -> { output: string, errors: string[] }
	js.Global().Set("runLipi", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return map[string]any{"output": "", "errors": []any{}}
		}
		output, errs := run(args[0].String())

		list := make([]any, len(errs))
		for i, e := range errs {
			list[i] = e
		}
		return map[string]any{"output": output, "errors": list}
	}))

	select {} // keep Go running so JavaScript can keep calling runLipi
}
