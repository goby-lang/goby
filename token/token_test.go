package token

import "testing"

func TestLookupIdentFalse(t *testing.T) {
	token := LookupIdent("nonexist")
	if token != IDENT {
		t.Fatalf("Expect %s got %s", IDENT, token)
	}
}

func TestLookupIdentTrue(t *testing.T) {
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

	for name, token := range keywords {
		test := LookupIdent(name)
		if test != token {
			t.Fatalf("Expect %s got %s", token, test)
		}
	}
}
