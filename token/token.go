package token

// TokenType is just a string so its easy to print while debugging
type TokenType string

// i am just defining the basic struct
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

const (
	// ILLEGAL = a character lipi dont know, EOF = file ended
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// names and values: naam, 5, "Ram"
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

	// comparing things
	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="

	// punctuation
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

// THE NEPALI PART! change a word here and the interpreter changes :D
// (website highlighting has its own list in website/src so change it there too)
var keywords = map[string]TokenType{
	"kaam":   FUNCTION,
	"rakha":  LET,
	"sahi":   TRUE,
	"galat":  FALSE,
	"yadi":   IF,
	"natra":  ELSE,
	"farkau": RETURN,
}

// keyword gets its own token type, anything else is a normal name
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
