package lexer

import (
	"errors"
	"monkey/token"
)

type Lexer struct {
	input         string
	ch            byte
	position      int
	readPosition  int
	interpolation bool
	paused        bool
}

func New(input string) *Lexer {
	l := &Lexer{input: input, interpolation: false, paused: false}
	l.readChar()
	return l
}

func (l *Lexer) flipInterpolation() {
	l.interpolation = !l.interpolation
}

func (l *Lexer) pause() {
	l.interpolation = false
	l.paused = true
}

func (l *Lexer) unPause() {
	l.interpolation = true
	l.paused = true
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	if l.paused && l.ch == '}' {
		tok = newToken(token.RBRACE, l.ch)
		if !isQuote(l.peekChar()) {
			l.unPause()
			l.readChar()
		}
		return tok
	}
	if l.interpolation {
		if isSingleQuote(l.ch) {
			tok = newToken(token.ISTRING, l.ch)
			l.flipInterpolation()
			l.readChar()
			return tok
		}
		if l.ch == '{' && l.peekChar() != '}' {
			l.pause()
			l.readChar()
			exp := l.NextToken()
			return exp
		}
		tok = newToken(token.BYTES, l.ch)
		l.readChar()
		return tok
	}

	l.skipWhitespace()
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			tok = token.Token{token.EQ, string(l.ch) + string(l.ch)}
			l.readChar()
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '.':
		tok = newToken(token.DOT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '-':
		if l.peekChar() == '>' {
			tok = token.Token{token.ARROW, string(l.ch) + string(l.peekChar())}
			l.readChar()
		} else {
			tok = newToken(token.MINUS, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			tok = token.Token{token.NEQ, string(l.ch) + string(l.peekChar())}
			l.readChar()
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '%':
		tok = newToken(token.MOD, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else if isQuote(l.ch) {
			if s, err := l.readString(); err == nil {
				tok.Type = token.STRING
				tok.Literal = s
				return tok
			}
		} else if isSingleQuote(l.ch) {
			tok = newToken(token.ISTRING, l.ch)
			l.flipInterpolation()
			l.readChar()
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
	}
	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readString() (string, error) {
	start := l.position + 1
	for {
		l.readChar()
		if isQuote(l.ch) {
			l.readChar()
			break
		}
		if l.ch == 0 {
			err := errors.New("")
			return "", err
		}
	}
	return l.input[start : l.position-1], nil
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isQuote(ch byte) bool {
	return ch == 34
}

func isSingleQuote(ch byte) bool {
	return ch == 0x27
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
