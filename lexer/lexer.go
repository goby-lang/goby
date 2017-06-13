package lexer

import (
	"github.com/goby-lang/goby/token"
	"github.com/looplab/fsm"
)

// Lexer is used for tokenizing programs
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	FSM          *fsm.FSM
}

// New initializes a new lexer with input string
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	l.FSM = fsm.NewFSM(
		"initial",
		fsm.Events{
			{Name: "method", Src: []string{"initial"}, Dst: "method"},
			{Name: "initialize", Src: []string{"method", "initial"}, Dst: "initial"},
		},
		fsm.Callbacks{},
	)
	return l
}

// NextToken makes lexer tokenize next character(s)
func (l *Lexer) NextToken() token.Token {

	var tok token.Token
	l.skipWhitespace()
	switch l.ch {
	case '"', byte('\''):
		tok.Literal = l.readString(l.ch)
		tok.Type = token.String
		tok.Line = l.line
		return tok
	case '=':
		if l.peekChar() == '=' {
			currentByte := l.ch
			l.readChar()
			tok = token.Token{Type: token.Eq, Literal: string(currentByte) + string(l.ch), Line: l.line}
		} else {
			tok = newToken(token.Assign, l.ch, l.line)
		}
	case '-':
		if l.peekChar() == '-' {
			tok.Literal = "--"
			tok.Line = l.line
			tok.Type = token.Decr
			l.readChar()
			l.readChar()
			return tok
		}
		tok = newToken(token.Minus, l.ch, l.line)
	case '!':
		if l.peekChar() == '=' {
			currentByte := l.ch
			l.readChar()
			tok = token.Token{Type: token.NotEq, Literal: string(currentByte) + string(l.ch), Line: l.line}
		} else {
			tok = newToken(token.Bang, l.ch, l.line)
		}
	case '/':
		tok = newToken(token.Slash, l.ch, l.line)
	case '*':
		if l.peekChar() == '*' {
			l.readChar()
			tok = token.Token{Type: token.Pow, Literal: "**", Line: l.line}
		} else {
			tok = newToken(token.Asterisk, l.ch, l.line)
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			if l.peekChar() == '>' {
				l.readChar()
				tok = token.Token{Type: token.COMP, Literal: "<=>", Line: l.line}
			} else {
				tok = token.Token{Type: token.LTE, Literal: "<=", Line: l.line}
			}
		} else {
			tok = newToken(token.LT, l.ch, l.line)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: ">=", Line: l.line}
		} else {
			tok = newToken(token.GT, l.ch, l.line)
		}
	case ';':
		tok = newToken(token.Semicolon, l.ch, l.line)
	case '(':
		tok = newToken(token.LParen, l.ch, l.line)
	case ')':
		tok = newToken(token.RParen, l.ch, l.line)
	case ',':
		tok = newToken(token.Comma, l.ch, l.line)
	case '+':
		if l.peekChar() == '+' {
			tok.Literal = "++"
			tok.Line = l.line
			tok.Type = token.Incr
			l.readChar()
			l.readChar()
			return tok
		}
		tok = newToken(token.Plus, l.ch, l.line)
	case '{':
		tok = newToken(token.LBrace, l.ch, l.line)
	case '}':
		tok = newToken(token.RBrace, l.ch, l.line)
	case '[':
		tok = newToken(token.LBracket, l.ch, l.line)
	case ']':
		tok = newToken(token.RBracket, l.ch, l.line)
	case '.':
		tok = newToken(token.Dot, l.ch, l.line)
		l.FSM.Event("method")
	case ':':
		if l.peekChar() == ':' {
			l.readChar()
			tok = token.Token{Type: token.ResolutionOperator, Literal: "::", Line: l.line}
		} else if isLetter(l.peekChar()) {
			tok.Literal = l.readSymbol()
			tok.Type = token.String
			tok.Line = l.line
			return tok
		} else {

			tok = newToken(token.Colon, l.ch, l.line)
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = token.Token{Type: token.Or, Literal: "||", Line: l.line}
		} else {
			tok = newToken(token.Bar, l.ch, l.line)
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = token.Token{Type: token.And, Literal: "&&", Line: l.line}
		}
	case '%':
		tok = newToken(token.Modulo, l.ch, l.line)
	case '#':
		tok.Literal = l.absorbComment()
		tok.Type = token.Comment
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
				tok.Type = token.Constant
				tok.Line = l.line
				l.FSM.Event("initialize")
			} else {
				tok.Literal = l.readIdentifier()
				if l.FSM.Is("method") {
					if tok.Literal == "self" {
						tok.Type = token.LookupIdent(tok.Literal)
					} else {
						tok.Type = token.Ident
					}
					l.FSM.Event("initialize")

				} else if l.FSM.Is("initial") {
					tok.Type = token.LookupIdent(tok.Literal)
					if tok.Literal == "def" {
						l.FSM.Event("method")
					} else {
						l.FSM.Event("initialize")
					}
				}
				tok.Line = l.line
			}
			return tok
		} else if isInstanceVariable(l.ch) {
			if isLetter(l.peekChar()) {
				tok.Literal = l.readInstanceVariable()
				tok.Type = token.InstanceVariable
				tok.Line = l.line
				return tok
			}

			return newToken(token.Illegal, l.ch, l.line)
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.Int
			tok.Line = l.line
			return tok
		}

		tok = newToken(token.Illegal, l.ch, l.line)
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

	// Strings like "" or ''
	if l.ch == ch {
		l.readChar()
		return ""
	}

	position := l.position // currently at string's first letter

	for l.peekChar() != ch {
		l.readChar()

		if l.peekChar() == 0 {
			panic("Unterminated string meets end of file")
		}
	}

	l.readChar()                           // currently at string's last letter
	result := l.input[position:l.position] // get full string
	l.readChar()                           // move to string's later quote
	return result
}


func (l *Lexer) readSymbol() string {
	l.readChar()

	position := l.position // currently at string's first letter

	for !(l.peekChar() == ' ' || l.peekChar() == '\n' || l.peekChar() == '\r' || l.peekChar() == '\t') {
		l.readChar()

		if l.peekChar() == 0 {
			panic("Unterminated string meets end of file")
		}
	}

	l.readChar()                           // currently at string's last letter
	result := l.input[position:l.position] // get full string
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

func newToken(tokenType token.Type, ch byte, line int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line}
}
