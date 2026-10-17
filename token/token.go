package token

type TokenType string

// i am just defining the basic struct
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	INT    = "INT"
	STRING = ""

	// some operators for some reason duh
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="

	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// keywords for basic things

	FUNCTION = "KAAM"
	LET      = "RAKHA"
	TRUE     = "SAHI"
	FALSE    = "GALAT"
	IF       = "YADI"
	ELSE     = "NATRA"
	RETURN   = "FARKAU"
)

var keywords = map[string]TokenType{
	"kaam":   FUNCTION,
	"rakha":  LET,
	"sahi":   TRUE,
	"galat":  FALSE,
	"yadi":   IF,
	"natra":  ELSE,
	"farkau": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
