package token

// Type is used to determite token type
type Type string

// Token is structure for identifying input stream of characters
type Token struct {
	Type    Type
	Literal string
	Line    int
}

// Literals
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	CONSTANT         = "CONSTANT"
	IDENT            = "IDENT"
	InstanceVariable = "INSTANCE_VAR"
	INT              = "INT"
	STRING           = "STRING"
	COMMENT          = "COMMENT"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	DOT      = "."
	INCR     = "++"
	DECR     = "--"

	LT = "<"
	GT = ">"

	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	BAR       = "|"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	EQ    = "=="
	NotEq = "!="

	CLASS  = "CLASS"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
	DEF    = "DEF"
	SELF   = "SELF"
	END    = "END"
	WHILE  = "WHILE"
	DO     = "DO"
	YIELD  = "YIELD"
)

var keywords = map[string]Type{
	"def":    DEF,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"self":   SELF,
	"end":    END,
	"while":  WHILE,
	"do":     DO,
	"yield":  YIELD,
}

// LookupIdent is used for keyword identification
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
