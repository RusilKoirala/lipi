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

// version is set when a release is built:
// go build -ldflags "-X main.version=v0.1.0"
var version = "dev"

// the help text, in Nepali of course
func usage() {
	fmt.Printf("Lipi %s\n\n", version)
	fmt.Println("Prayog:")
	fmt.Println("  lipi <file.lip>     program chalau")
	fmt.Println("  lipi --version      version hera")
	fmt.Println("  lipi --help         yo sahayata")
}

func main() {
	// no file given (no REPL yet)
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	// flags first
	switch os.Args[1] {
	case "-v", "--version", "version":
		fmt.Printf("lipi %s\n", version)
		return
	case "-h", "--help", "help":
		usage()
		return
	}

	// anything else should be a .lip file
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
	// only a Galti comes back that we care about, bhana already printed the rest
	result := evaluator.Eval(program, object.NewEnvironment())
	if errObj, ok := result.(*object.Error); ok {
		fmt.Fprintln(os.Stderr, errObj.Message)
		os.Exit(1)
	}
}
