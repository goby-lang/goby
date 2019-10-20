package lexer

import (
	"github.com/goby-lang/goby/compiler/token"
	"github.com/looplab/fsm"
)

// Lexer is used for tokenizing programs
type Lexer struct {
	input        []rune
	position     int
	readPosition int
	ch           rune
	line         int
	FSM          *fsm.FSM
}

// New initializes a new lexer with input string
func New(input string) *Lexer {
	l := &Lexer{input: []rune(input)}
	l.readChar()
	l.FSM = fsm.NewFSM(
		"initial",
		/*
			Initial state is default state
			Nosymbol state helps us identify tok ':' is for symbol or hash value
			Method state helps us identify 'class' literal is a keyword or an identifier
			Reference: https://github.com/looplab/fsm
		*/
		fsm.Events{
			{Name: "nosymbol", Src: []string{"initial"}, Dst: "nosymbol"},
			{Name: "method", Src: []string{"initial"}, Dst: "method"},
			{Name: "initial", Src: []string{"method", "initial", "nosymbol"}, Dst: "initial"},
		},
		fsm.Callbacks{},
	)
	return l
}

// NextToken makes lexer tokenize next character(s)
func (l *Lexer) NextToken() token.Token {

	var tok token.Token
	l.resetNosymbol()

	l.skipWhitespace()
	switch l.ch {
	case '"', '\'':
		tok.Literal = l.readString(l.ch)
		tok.Type = token.String
		tok.Line = l.line
		return tok
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.CreateOperator("==", l.line)
		} else if l.peekChar() == '~' {
			l.readChar()
			tok = token.CreateOperator("=~", l.line)
		} else {
			tok = token.CreateOperator("=", l.line)
		}
	case '-':
		if l.peekChar() == '=' {
			tok = token.CreateOperator("-=", l.line)
			l.readChar()
			l.readChar()
			return tok
		}
		tok = token.CreateOperator("-", l.line)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.CreateOperator("!=", l.line)
		} else {
			tok = token.CreateOperator("!", l.line)
		}
	case '/':
		tok = token.CreateOperator("/", l.line)
	case '*':
		if l.peekChar() == '*' {
			l.readChar()
			tok = token.CreateOperator("**", l.line)
		} else {
			tok = token.CreateOperator("*", l.line)
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			if l.peekChar() == '>' {
				l.readChar()
				tok = token.CreateOperator("<=>", l.line)
			} else {
				tok = token.CreateOperator("<=", l.line)
			}
		} else {
			tok = token.CreateOperator("<", l.line)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.CreateOperator(">=", l.line)
		} else {
			tok = token.CreateOperator(">", l.line)
		}
	case ';', ',', '(', ')', '{', '}', '[', ']':
		tok = token.CreateSeparator(string(l.ch), l.line)
	case '+':
		if l.peekChar() == '=' {
			tok = token.CreateOperator("+=", l.line)
			l.readChar()
			l.readChar()
			return tok
		}
		tok = token.CreateOperator("+", l.line)
	case '.':
		if l.peekChar() == '.' {
			tok = token.CreateOperator("..", l.line)
			l.readChar()
			l.readChar()
			return tok
		}
		tok = token.CreateOperator(".", l.line)
		l.FSM.Event("method")
	case ':':
		if l.FSM.Is("nosymbol") {
			//e.g. {test: abc} || {test: :abc} || {test: 50}

			tok = token.CreateSeparator(":", l.line)

		} else {
			if l.peekChar() == ':' {
				l.readChar()
				tok = token.CreateOperator("::", l.line)

			} else if isLetter(l.peekChar()) {
				tok.Literal = string(l.readSymbol())
				tok.Type = token.String
				tok.Line = l.line
				return tok

			} else {
				tok = token.CreateSeparator(":", l.line)
			}
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = token.CreateOperator("||=", l.line)
			} else {
				tok = token.CreateOperator("||", l.line)
			}
		} else {
			tok = token.CreateSeparator("|", l.line)
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = token.CreateOperator("&&", l.line)
		}
	case '%', '?':
		tok = token.CreateOperator(string(l.ch), l.line)
	case '#':
		tok.Literal = string(l.absorbComment())
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
				tok.Literal = string(l.readConstant())
				tok.Type = token.Constant
				tok.Line = l.line
				l.FSM.Event("initial")
			} else {
				tok.Literal = string(l.readIdentifier())
				if l.FSM.Is("method") {
					if tok.Literal == "self" {
						tok.Type = token.LookupIdent(tok.Literal)
					} else {
						tok.Type = token.Ident
					}
					l.FSM.Event("initial")

				} else if l.FSM.Is("initial") {
					tok.Type = token.LookupIdent(tok.Literal)
					if tok.Literal == "def" {
						l.FSM.Event("method")
					} else {
						l.FSM.Event("initial")
					}
				}
				tok.Line = l.line
			}
			if tok.Type == token.Ident {
				l.FSM.Event("nosymbol")
			}
			return tok
		} else if isInstanceVariable(l.ch) {
			if isLetter(l.peekChar()) {
				tok.Literal = string(l.readInstanceVariable())
				tok.Type = token.InstanceVariable
				tok.Line = l.line
				return tok
			}

			return token.Token{Type: token.Illegal, Literal: string(l.ch), Line: l.line}
		} else if isDigit(l.ch) {
			tok.Literal = string(l.readNumber())
			tok.Type = token.Int
			tok.Line = l.line
			return tok
		}

		tok = token.Token{Type: token.Illegal, Literal: string(l.ch), Line: l.line}
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

func (l *Lexer) resetNosymbol() {

	if !l.FSM.Is("method") && l.ch != ':' {
		l.FSM.Event("initial")

	}

}

func (l *Lexer) readNumber() []rune {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() []rune {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '?' {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readConstant() []rune {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readInstanceVariable() []rune {
	position := l.position
	for isLetter(l.ch) || isInstanceVariable(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString(ch rune) string {
	l.readChar()

	// Empty strings case such as "" or ''
	if l.ch == ch {
		l.readChar()
		return ""
	}

	result := ""

	for {
		if isEscapedChar(l.ch) {
			result += escapedCharResult(ch, l.peekChar())
			l.readChar()
		} else {
			result += string(l.ch)
		}
		l.readChar()

		if l.ch == ch || l.peekChar() == 0 {
			break
		}
	}

	// fmt.Println(l.ch) <- Currently at string's last character
	l.readChar() // move to string's latter quote

	return result
}

func (l *Lexer) readSymbol() []rune {
	l.readChar()

	position := l.position // currently at string's first letter

	for isLetter(l.peekChar()) || isDigit(l.peekChar()) {
		l.readChar()
	}

	l.readChar()                           // currently at string's last letter
	result := l.input[position:l.position] // get full string
	return result
}

func (l *Lexer) absorbComment() []rune {
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

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
	// Peek shouldn't increment positions.
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isInstanceVariable(ch rune) bool {
	return ch == '@'
}

func isEscapedChar(ch rune) bool {
	return ch == '\\'
}

func escapedCharResult(quotedChar rune, peeked rune) string {
	if quotedChar == '"' {
		switch peeked {
		case 'n':
			return "\n"
		case 't':
			return "\t"
		case 'v':
			return "\v"
		case 'f':
			return "\f"
		case 'r':
			return "\r"
		case '\\':
			return "\\"
		case '"':
			return "\""
		case '\'':
			return "'"
		default:
			return "\\" + string(peeked)
		}
	}
	switch peeked {
	case '"':
		return "\\\""
	case '\'':
		return "'"
	default:
		return "\\" + string(peeked)
	}
}
