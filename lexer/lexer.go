package lexer

import "github.com/rusilkoirala/lipi/token"

// Lexer goes through the source one byte at a time and makes tokens
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

// line starts at 1 and readChar loads the first character
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

// NextToken gives the next token every time you call it
func (l *Lexer) NextToken() token.Token {
	l.skipWhitespaceAndComments()

	// remember where token startss
	line, col := l.line, l.column
	var tok token.Token

	switch l.ch {
	// = and ! can be two characters, so peek first (== should not become = and =)
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
		// words and numbers already moved past their last character,
		// so they return early and skip the readChar at the bottom
		if isLetter(l.ch) {
			lit := l.readIdentifier()
			tok = newToken(token.LookupIdent(lit), lit)
			tok.Line, tok.Column = line, col
			return tok
		} else if isDigit(l.ch) {
			tok = newToken(token.INT, l.readNumber())
			tok.Line, tok.Column = line, col
			return tok
		}
		// no idea what this character is
		tok = newToken(token.ILLEGAL, string(l.ch))
	}
	l.readChar()
	tok.Line, tok.Column = line, col
	return tok
}

// lipi dont care about spaces or comments so they never become tokens
func (l *Lexer) skipWhitespaceAndComments() {
	for {
		switch {
		case l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r':
			l.readChar()
		case l.ch == '/' && l.peekChar() == '/':
			// comment goes till the end of the line
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
		default:
			return
		}
	}
}

// reads till the closing quote, false means the file ended first
// it only looks for the quote byte so devanagari inside strings just works
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

// small helper, Line and Column get set in NextToken
func newToken(t token.TokenType, literal string) token.Token {
	return token.Token{
		Type:    t,
		Literal: literal,
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// only english letters for names, so naam works but devanagari names dont (yet)
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// keep reading while its a digit
func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

// names can have digits after the first letter: jod2
func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}
