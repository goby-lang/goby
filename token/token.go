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
	Pow      = "**"
	Slash    = "/"
	Dot      = "."
	Incr     = "++"
	Decr     = "--"
	And      = "&&"
	Or       = "||"

	LT   = "<"
	LTE  = "<="
	GT   = ">"
	GTE  = ">="
	COMP = "<=>"

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

	True   = "TRUE"
	False  = "FALSE"
	Null   = "Null"
	If     = "IF"
	Else   = "ELSE"
	Return = "RETURN"
	Def    = "DEF"
	Self   = "SELF"
	End    = "END"
	While  = "WHILE"
	Next   = "NEXT"
	Do     = "DO"
	Yield  = "YIELD"

	ResolutionOperator = "::"
)

var keywords = map[string]Type{
	"def":    Def,
	"true":   True,
	"false":  False,
	"nil":    Null,
	"if":     If,
	"else":   Else,
	"return": Return,
	"self":   Self,
	"end":    End,
	"while":  While,
	"next":   Next,
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
