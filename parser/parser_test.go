package parser

import (
	"testing"

	"github.com/rusilkoirala/lipi/lexer"
)

func TestRakha(t *testing.T) {
	p := New(lexer.New(`rakha naam = "Ram"; farkau naam;`))
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("errors: %v", p.Errors())
	}

	expected := `rakha naam = "Ram";farkau naam;`
	if program.String() != expected {
		t.Fatalf("expected %s, got %s", expected, program.String())
	}
}

func TestSemicolonError(t *testing.T) {
	p := New(lexer.New("rakha a = 1;\nrakha b = 2;\nfarkau b\n"))
	p.ParseProgram()

	expected := "Galti: line 3 ma semicolon birsiyo"
	if len(p.Errors()) == 0 || p.Errors()[0] != expected {
		t.Fatalf("expected %q, got %v", expected, p.Errors())
	}
}

func TestMath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1 + 2 * 3;", "(1 + (2 * 3))"},
		{"(1 + 2) * 3;", "((1 + 2) * 3)"},
		{"-5 + 10;", "((-5) + 10)"},
		{"!sahi == galat;", "((!sahi) == galat)"},
		{"a + b > 4;", "((a + b) > 4)"},
		{"rakha x = 1 + 2;", "rakha x = (1 + 2);"},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			t.Fatalf("input %q errors: %v", tt.input, p.Errors())
		}
		if program.String() != tt.expected {
			t.Errorf("input %q: expected %s, got %s", tt.input, tt.expected, program.String())
		}
	}
}

func TestCallsIfAndFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"jod(1, 2 * 3);", "jod(1, (2 * 3))"},
		{"bhana();", "bhana()"},
		{"a + jod(b, c) * 2;", "(a + (jod(b, c) * 2))"},
		{`yadi (x > 4) { bhana("thulo!"); }`, `yadi (x > 4) bhana("thulo!")`},
		{"yadi (x) { 1; } natra { 2; }", "yadi x 1 natra 2"},
		{"rakha f = kaam(a, b) { farkau a + b; };", "rakha f = kaam(a, b) farkau (a + b);;"},
		{"kaam jod(a, b) { farkau a + b; }", "rakha jod = kaam jod(a, b) farkau (a + b);;"},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			t.Fatalf("input %q errors: %v", tt.input, p.Errors())
		}
		if program.String() != tt.expected {
			t.Errorf("input %q:\nexpected %s\n     got %s", tt.input, tt.expected, program.String())
		}
	}
}

func TestFullExample(t *testing.T) {
	input := `rakha naam = "Ram";
bhana("Namaste, " + naam);

kaam jod(a, b) {
  farkau a + b;
}

yadi (jod(2, 3) > 4) {
  bhana("thulo!");
} natra {
  bhana("sano");
}`

	p := New(lexer.New(input))
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("errors: %v", p.Errors())
	}
	if len(program.Statements) != 4 {
		t.Fatalf("expected 4 statements, got %d", len(program.Statements))
	}
}
