package lexer

import (
	"github.com/rooby-lang/Rooby/token"
)

// Lexer is used for tokenizing programs
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
}

// New initializes a new lexer with input string
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// NextToken makes lexer tokenize next character(s)
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '"', byte('\''):
		tok.Literal = l.readString(l.ch)
		tok.Type = token.STRING
		tok.Line = l.line
		return tok
	case '=':
		if l.peekChar() == '=' {
			currentByte := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(currentByte) + string(l.ch), Line: l.line}
		} else {
			tok = newToken(token.ASSIGN, l.ch, l.line)
		}
	case '-':
		if l.peekChar() == '-' {
			tok.Literal = "--"
			tok.Line = l.line
			tok.Type = token.DECR
			l.readChar()
			l.readChar()
			return tok
		}
		tok = newToken(token.MINUS, l.ch, l.line)
	case '!':
		if l.peekChar() == '=' {
			currentByte := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(currentByte) + string(l.ch), Line: l.line}
		} else {
			tok = newToken(token.BANG, l.ch, l.line)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch, l.line)
	case '*':
		tok = newToken(token.ASTERISK, l.ch, l.line)
	case '<':
		tok = newToken(token.LT, l.ch, l.line)
	case '>':
		tok = newToken(token.GT, l.ch, l.line)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, l.line)
	case '(':
		tok = newToken(token.LPAREN, l.ch, l.line)
	case ')':
		tok = newToken(token.RPAREN, l.ch, l.line)
	case ',':
		tok = newToken(token.COMMA, l.ch, l.line)
	case '+':
		if l.peekChar() == '+' {
			tok.Literal = "++"
			tok.Line = l.line
			tok.Type = token.INCR
			l.readChar()
			l.readChar()
			return tok
		}
		tok = newToken(token.PLUS, l.ch, l.line)
	case '{':
		tok = newToken(token.LBRACE, l.ch, l.line)
	case '}':
		tok = newToken(token.RBRACE, l.ch, l.line)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, l.line)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, l.line)
	case '.':
		tok = newToken(token.DOT, l.ch, l.line)
	case ':':
		tok = newToken(token.COLON, l.ch, l.line)
	case '|':
		tok = newToken(token.BAR, l.ch, l.line)
	case '#':
		tok.Literal = l.absorbComment()
		tok.Type = token.COMMENT
		tok.Line = l.line
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.line
	default:
		if isLetter(l.ch) {
			if 'A' <= l.ch && l.ch <= 'Z' {
				tok.Literal = l.readConstant()
				tok.Type = token.CONSTANT
				tok.Line = l.line
			} else {
				tok.Literal = l.readIdentifier()
				tok.Type = token.LookupIdent(tok.Literal)
				tok.Line = l.line
			}

			return tok
		} else if isInstanceVariable(l.ch) {
			if isLetter(l.peekChar()) {
				tok.Literal = l.readInstanceVariable()
				tok.Type = token.INSTANCE_VARIABLE
				tok.Line = l.line
				return tok
			}

			return newToken(token.ILLEGAL, l.ch, l.line)
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			tok.Line = l.line
			return tok
		}

		tok = newToken(token.ILLEGAL, l.ch, l.line)
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		if l.ch == '\n' {
			l.line++
		}
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readConstant() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readInstanceVariable() string {
	position := l.position
	for isLetter(l.ch) || isInstanceVariable(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString(ch byte) string {
	l.readChar()
	position := l.position // currently at string's first letter

	for l.peekChar() != ch {
		l.readChar()
	}

	l.readChar()                           // currently at string's last letter
	result := l.input[position:l.position] // get full string
	l.readChar()                           // move to string's later quote
	return result
}

func (l *Lexer) absorbComment() string {
	p := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	result := l.input[p:l.position]
	return result
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// ascii code's null
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
	// Peek shouldn't increment positions.
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isInstanceVariable(ch byte) bool {
	return ch == '@'
}

func newToken(tokenType token.TokenType, ch byte, line int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line}
}
