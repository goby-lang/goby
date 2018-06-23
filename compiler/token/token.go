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
	Float            = "FLOAT"
	String           = "STRING"
	Comment          = "COMMENT"

	UnderScore   = "_"
	Assign   = "="
	Plus     = "+"
	PlusEq   = "+="
	Minus    = "-"
	MinusEq  = "-="
	Bang     = "!"
	Asterisk = "*"
	Pow      = "**"
	Slash    = "/"
	Dot      = "."
	And      = "&&"
	Or       = "||"
	OrEq     = "||="
	Modulo   = "%"

	Match = "=~"
	LT    = "<"
	LTE   = "<="
	GT    = ">"
	GTE   = ">="
	COMP  = "<=>"

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
	Range = ".."

	True     = "TRUE"
	False    = "FALSE"
	Null     = "Null"
	If       = "IF"
	ElsIf    = "ELSIF"
	Else     = "ELSE"
	Case     = "CASE"
	When     = "WHEN"
	Return   = "RETURN"
	Next     = "NEXT"
	Break    = "BREAK"
	Def      = "DEF"
	Self     = "SELF"
	End      = "END"
	While    = "WHILE"
	Do       = "DO"
	Yield    = "YIELD"
	GetBlock = "GET_BLOCK"
	Class    = "CLASS"
	Module   = "MODULE"

	ResolutionOperator = "::"
)

var keywords = map[string]Type{
	"def":       Def,
	"true":      True,
	"false":     False,
	"nil":       Null,
	"if":        If,
	"elsif":     ElsIf,
	"else":      Else,
	"case":      Case,
	"when":      When,
	"return":    Return,
	"self":      Self,
	"end":       End,
	"while":     While,
	"do":        Do,
	"yield":     Yield,
	"next":      Next,
	"class":     Class,
	"module":    Module,
	"break":     Break,
	"get_block": GetBlock,
}

// LookupIdent is used for keyword identification
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Ident
}
