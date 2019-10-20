package token

// Type is used to determine token type
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
	TernaryOperator    = "?"
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

var operators = map[string]Type{
	"=":   Assign,
	"+":   Plus,
	"+=":  PlusEq,
	"-":   Minus,
	"-=":  MinusEq,
	"!":   Bang,
	"*":   Asterisk,
	"**":  Pow,
	"/":   Slash,
	".":   Dot,
	"&&":  And,
	"||":  Or,
	"||=": OrEq,
	"%":   Modulo,

	"=~":  Match,
	"<":   LT,
	"<=":  LTE,
	">":   GT,
	">=":  GTE,
	"<=>": COMP,

	"==": Eq,
	"!=": NotEq,
	"..": Range,

	"::": ResolutionOperator,
	"?":  TernaryOperator,
}

var separators = map[string]Type{
	",": Comma,
	";": Semicolon,
	":": Colon,
	"|": Bar,

	"(": LParen,
	")": RParen,
	"{": LBrace,
	"}": RBrace,
	"[": LBracket,
	"]": RBracket,
}

// LookupIdent is used for keyword identification
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Ident
}

func getOperatorType(literal string) Type {
	if t, ok := operators[literal]; ok {
		return t
	}
	return Ident
}

func getSeparatorType(literal string) Type {
	if t, ok := separators[literal]; ok {
		return t
	}
	return Ident
}

// CreateOperator - Factory method for creating operator types token from literal string
func CreateOperator(literal string, line int) Token {
	return Token{Type: getOperatorType(literal), Literal: literal, Line: line}
}

// CreateSeparator - Factory method for creating separator types token from literal string
func CreateSeparator(literal string, line int) Token {
	return Token{Type: getSeparatorType(literal), Literal: literal, Line: line}
}
