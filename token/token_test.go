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

	for name, token := range keywords {
		test := LookupIdent(name)
		if test != token {
			t.Fatalf("Expect %s got %s", token, test)
		}
	}
}
