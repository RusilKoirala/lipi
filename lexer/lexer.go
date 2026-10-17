package lexer

import "github.com/rusilkoirala/lipi/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
		line:  1,
	}
	l.readChar()
	return l
}

// readChar moves forward one character and keeps line/column up to date.
func (l *Lexer) readChar() {
	if l.ch == '\n' {
		l.line++
		l.column = 0
	}

	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
}

// peekChar looks at next character without moving

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespaceAndComments()

	// remember where token startss
	line, col := l.line, l.column
	var tok token.Token

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = newToken(token.EQ, "==")
		} else {
			tok = newToken(token.ASSIGN, "=")
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = newToken(token.NOT_EQ, "!=")
		} else {
			tok = newToken(token.BANG, "!")
		}
	case '+':
		tok = newToken(token.PLUS, "+")
	case '-':
		tok = newToken(token.MINUS, "-")
	case '*':
		tok = newToken(token.ASTERISK, "*")
	case '/':
		tok = newToken(token.SLASH, "/")
	case '<':
		tok = newToken(token.LT, "<")
	case '>':
		tok = newToken(token.GT, ">")
	case ',':
		tok = newToken(token.COMMA, ",")
	case ';':
		tok = newToken(token.SEMICOLON, ";")
	case '(':
		tok = newToken(token.LPAREN, "(")
	case ')':
		tok = newToken(token.RPAREN, ")")
	case '{':
		tok = newToken(token.LBRACE, "{")
	case '}':
		tok = newToken(token.RBRACE, "}")
	case '"':
		str, ok := l.readString()
		if ok {
			tok = newToken(token.STRING, str)
		} else {
			tok = newToken(token.ILLEGAL, str) // string never closed
		}
	case 0:
		tok = newToken(token.EOF, "")

	default:
		if isLetter(l.ch) {

		} else if isDigit(l.ch) {

		}
	}
	l.readChar()
	tok.Line, tok.Column = line, col
	return tok
}

func (l *Lexer) skipWhitespaceAndComments() {

}

func (l *Lexer) readString() (string, bool) {
	start := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' {
			return l.input[start:l.position], true
		}

		if l.ch == 0 {
			return l.input[start:l.position], false
		}
	}
}

func newToken(t token.TokenType, literal string) token.Token {
	return token.Token{
		Type:    t,
		Literal: literal,
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}
