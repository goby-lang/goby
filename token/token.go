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
	Illegal = "ILLEGAL"
	EOF     = "EOF"

	Constant         = "CONSTANT"
	Ident            = "IDENT"
	InstanceVariable = "INSTANCE_VAR"
	Int              = "INT"
	String           = "STRING"
	Comment          = "COMMENT"

	Assign   = "="
	Plus     = "+"
	Minus    = "-"
	Bang     = "!"
	Asterisk = "*"
	Slash    = "/"
	Dot      = "."
	Incr     = "++"
	Decr     = "--"

	LT = "<"
	GT = ">"

	Comma     = ","
	Semicolon = ";"
	Colon     = ":"
	Bar       = "|"

	LParen   = "("
	RParen   = ")"
	LBrace   = "{"
	RBrace   = "}"
	LBracket = "["
	RBracket = "]"

	Eq    = "=="
	NotEq = "!="

	Class  = "CLASS"
	True   = "TRUE"
	False  = "FALSE"
	If     = "IF"
	Else   = "ELSE"
	Return = "RETURN"
	Def    = "DEF"
	Self   = "SELF"
	End    = "END"
	While  = "WHILE"
	Do     = "DO"
	Yield  = "YIELD"
)

var keywords = map[string]Type{
	"def":    Def,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
	"self":   Self,
	"end":    End,
	"while":  While,
	"do":     Do,
	"yield":  Yield,
}

// LookupIdent is used for keyword identification
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Ident
}
