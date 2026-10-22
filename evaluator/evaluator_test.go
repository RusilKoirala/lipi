package evaluator

import (
	"testing"

	"github.com/rusilkoirala/lipi/lexer"
	"github.com/rusilkoirala/lipi/object"
	"github.com/rusilkoirala/lipi/parser"
)

func run(t *testing.T, input string) object.Object {
	t.Helper()
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors for %q: %v", input, p.Errors())
	}
	return Eval(program, object.NewEnvironment())
}

func TestEval(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5;", "5"},
		{"1 + 2 * 3;", "7"},
		{"(10 - 4) / 2;", "3"},
		{"-5 + 10;", "5"},
		{`"Namaste, " + "Ram";`, "Namaste, Ram"},
		{"3 > 2;", "sahi"},
		{"!sahi;", "galat"},
		{`"a" == "a";`, "sahi"},
		{"yadi (5 > 4) { 1; } natra { 2; }", "1"},
		{"yadi (galat) { 1; }", "khali"},
		{"rakha x = 5; rakha y = x * 2; y;", "10"},
		{`rakha naam = "Ram"; "Namaste, " + naam;`, "Namaste, Ram"},
		{"yadi (sahi) { farkau 10; 99; }", "10"},
	}

	for _, tt := range tests {
		got := run(t, tt.input).Inspect()
		if got != tt.expected {
			t.Errorf("input %q: expected %s, got %s", tt.input, tt.expected, got)
		}
	}
}

func TestErrors(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5 + sahi;", "Galti: line 1 ma sankhya ra sahi/galat jodna mildaina"},
		{"rakha naam = \"Ram\";\nnaam + 5;", "Galti: line 2 ma string ra sankhya jodna mildaina"},
		{"rakha x = 5;\nrakha y = z;", "Galti: line 2 ma 'z' bhanne naam bhetiyena"},
		{"10 / 0;", "Galti: line 1 ma 0 le bhag garna mildaina"},
		{"-sahi;", "Galti: line 1 ma sahi/galat lai '-' garna mildaina"},
	}

	for _, tt := range tests {
		result := run(t, tt.input)
		err, ok := result.(*object.Error)
		if !ok {
			t.Errorf("input %q: expected a Galti, got %s", tt.input, result.Inspect())
			continue
		}
		if err.Message != tt.expected {
			t.Errorf("input %q:\nexpected %s\n     got %s", tt.input, tt.expected, err.Message)
		}
	}
}

func TestFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"kaam jod(a, b) { farkau a + b; } jod(2, 3);", "5"},
		{"rakha dobar = kaam(x) { farkau x * 2; }; dobar(21);", "42"},
		{"kaam fib(n) { yadi (n < 2) { farkau n; } farkau fib(n - 1) + fib(n - 2); } fib(15);", "610"},
		{"kaam gunak(x) { farkau kaam(y) { farkau x * y; }; } rakha tebar = gunak(3); tebar(7);", "21"},
		{`lambai("Kathmandu");`, "9"},
	}

	for _, tt := range tests {
		got := run(t, tt.input).Inspect()
		if got != tt.expected {
			t.Errorf("input %q: expected %s, got %s", tt.input, tt.expected, got)
		}
	}
}
