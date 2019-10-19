package token

import "testing"

func TestLookupIdentFalse(t *testing.T) {
	token := LookupIdent("nonexist")
	if token != Ident {
		t.Fatalf("Expect %s got %s", Ident, token)
	}
}

func TestLookupIdentTrue(t *testing.T) {
	var keywords = map[string]Type{
		"def":       Def,
		"true":      True,
		"false":     False,
		"if":        If,
		"elsif":     ElsIf,
		"else":      Else,
		"when":      When,
		"case":      Case,
		"return":    Return,
		"next":      Next,
		"self":      Self,
		"end":       End,
		"while":     While,
		"do":        Do,
		"yield":     Yield,
		"nil":       Null,
		"get_block": GetBlock,
	}

	for name, token := range keywords {
		test := LookupIdent(name)
		if test != token {
			t.Fatalf("Expect %s got %s", token, test)
		}
	}
}

func TestCreateOperatorIdentFalse(t *testing.T) {
	line := 123
	token := CreateOperator("nonexist", line)
	if token.Type != Ident {
		t.Fatalf("Expect token type %s, got: %s", Ident, token.Type)
	}
	if token.Line != line {
		t.Fatalf("Expect token line %v, got: %v", line, token.Line)
	}
}

func TestCreateOperatorIdentTrue(t *testing.T) {
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
	}

	line := 123

	for name, tokenType := range operators {
		tok := CreateOperator(name, line)
		if tok.Type != tokenType {
			t.Fatalf("Expect token type %s, got: %s", tokenType, tok.Type)
		}
		if tok.Line != line {
			t.Fatalf("Expect token line %v, got: %v", line, tok.Line)
		}
	}
}

func TestCreateSeparatorIdentFalse(t *testing.T) {
	line := 123
	token := CreateSeparator("nonexist", line)
	if token.Type != Ident {
		t.Fatalf("Expect token type %s, got: %s", Ident, token.Type)
	}
	if token.Line != line {
		t.Fatalf("Expect token line %v, got: %v", line, token.Line)
	}
}

func TestCreateSeparatorIdentTrue(t *testing.T) {
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

	line := 123

	for name, tokenType := range separators {
		tok := CreateSeparator(name, line)
		if tok.Type != tokenType {
			t.Fatalf("Expect token type %s, got: %s", tokenType, tok.Type)
		}
		if tok.Line != line {
			t.Fatalf("Expect token line %v, got: %v", line, tok.Line)
		}
	}
}
