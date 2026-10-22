package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rusilkoirala/lipi/evaluator"
	"github.com/rusilkoirala/lipi/lexer"
	"github.com/rusilkoirala/lipi/object"
	"github.com/rusilkoirala/lipi/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Lipi 0.1")
		fmt.Println("Prayog: lipi <file.lip>")
		os.Exit(1)
	}

	path := os.Args[1]
	if filepath.Ext(path) != ".lip" {
		fmt.Fprintf(os.Stderr, "Galti: '%s' .lip file hoina\n", path)
		os.Exit(1)
	}

	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Galti: '%s' file khulena\n", path)
		os.Exit(1)
	}

	// 1. lexer + parser
	p := parser.New(lexer.New(string(src)))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		for _, e := range p.Errors() {
			fmt.Fprintln(os.Stderr, e)
		}
		os.Exit(1)
	}

	// 2. run it
	result := evaluator.Eval(program, object.NewEnvironment())
	if errObj, ok := result.(*object.Error); ok {
		fmt.Fprintln(os.Stderr, errObj.Message)
		os.Exit(1)
	}
}
