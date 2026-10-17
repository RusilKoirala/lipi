package lexer

import (
	"testing"

	"github.com/rusilkoirala/lipi/token"
)

func TestNextToken(t *testing.T) {
	input := `rakha naam = "Ram";
bhana("Namaste, " + naam);

kaam jod(a, b) {
  farkau a + b;
}

yadi (jod(2, 3) > 4) {
  bhana("thulo!");
} natra {
  bhana("sano");
}
!-/*5 < 10 == sahi != galat; // yo comment ho
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "rakha"},
		{token.IDENT, "naam"},
		{token.ASSIGN, "="},
		{token.STRING, "Ram"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "bhana"},
		{token.LPAREN, "("},
		{token.STRING, "Namaste, "},
		{token.PLUS, "+"},
		{token.IDENT, "naam"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.FUNCTION, "kaam"},
		{token.IDENT, "jod"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "farkau"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.IF, "yadi"},
		{token.LPAREN, "("},
		{token.IDENT, "jod"},
		{token.LPAREN, "("},
		{token.INT, "2"},
		{token.COMMA, ","},
		{token.INT, "3"},
		{token.RPAREN, ")"},
		{token.GT, ">"},
		{token.INT, "4"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "bhana"},
		{token.LPAREN, "("},
		{token.STRING, "thulo!"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "natra"},
		{token.LBRACE, "{"},
		{token.IDENT, "bhana"},
		{token.LPAREN, "("},
		{token.STRING, "sano"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.TRUE, "sahi"},
		{token.NOT_EQ, "!="},
		{token.FALSE, "galat"},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong type. expected=%q, got=%q (literal %q)",
				i, tt.expectedType, tok.Type, tok.Literal)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - wrong literal. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLineAndColumn(t *testing.T) {
	input := "rakha x = 5;\n  rakha y = 10;"

	tests := []struct {
		literal string
		line    int
		column  int
	}{
		{"rakha", 1, 1},
		{"x", 1, 7},
		{"=", 1, 9},
		{"5", 1, 11},
		{";", 1, 12},
		{"rakha", 2, 3},
		{"y", 2, 9},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Literal != tt.literal || tok.Line != tt.line || tok.Column != tt.column {
			t.Fatalf("tests[%d] - expected %q at %d:%d, got %q at %d:%d",
				i, tt.literal, tt.line, tt.column, tok.Literal, tok.Line, tok.Column)
		}
	}
}
