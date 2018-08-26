package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goby-lang/goby/compiler/parser"
)

func TestCompileToInstructionsNormalMode(t *testing.T) {

	ci, err := CompileToInstructions(`
def bar(a)
  99 + a
end
while true do
end`, parser.NormalMode)

	if err != nil {
		t.Fatal(err.Error())
	}

	if e, a := uint8(10), ci[0].Instructions[0].Opcode; e != a {
		t.Fatalf("Expect `%d` for first instruction opcode. got: %d", e, a)
	}

	if e, a := "putobject: 99. source line: 3", ci[0].Instructions[0].Inspect(); e != a {
		t.Fatalf("Expect `%s` for first instruction inspect. got: %s", e, a)
	}

	if e, a := 4, ci[1].Instructions[3].AnchorLine(); e != a {
		t.Fatalf("Expect `%d` for first instruction inspect. got: %d", e, a)
	}

	if e, a := 99, ci[0].Instructions[0].Params[0]; e != a {
		t.Fatalf("Expect `%d` for first instruction first param. got: %d", e, a)
	}

	// TODO: change the following simple public functions to public variables
	if e, a := "bar", ci[0].Name(); e != a {
		t.Fatalf("Expect `%s` for instruction set name. got: %s", e, a)
	}

	if e, a := "Def", ci[0].Type(); e != a {
		t.Fatalf("Expect `%s` for instruction set type. got: %s", e, a)
	}

	if e, a := "putobject", ci[0].Instructions[0].ActionName(); e != a {
		t.Fatalf("Expect `%s` for first instruction action name. got: %s", e, a)
	}

	if e, a := 0, ci[0].Instructions[0].Line(); e != a {
		t.Fatalf("Expect `%d` for first instruction line. got: %d", e, a)
	}

	if e, a := 3, ci[0].Instructions[0].SourceLine(); e != a {
		t.Fatalf("Expect `%d` for first instruction source line. got: %d", e, a)
	}
}

func TestCompileToInstructionsNormalModePanic(t *testing.T) {
	assert.Panics(t, func() {
		ci, err := CompileToInstructions(`99`, parser.NormalMode)
		if err != nil {
			t.Fatal(err.Error())
		}

		ci[0].Instructions[0].AnchorLine()
	}, "The code did not panic")
}

func TestCompileToInstructionsTESTMode(t *testing.T) {

	ci, err := CompileToInstructions(`
module Foo
end
`, parser.TestMode)

	if err != nil {
		t.Fatal(err.Error())
	}

	if e, a := uint8(29), ci[0].Instructions[0].Opcode; e != a {
		t.Fatalf("Expect `%d` for first instruction opcode. got: %d", e, a)
	}

	if e, a := "leave: . source line: 2", ci[0].Instructions[0].Inspect(); e != a {
		t.Fatalf("Expect `%s` for first instruction inspect. got: %s", e, a)
	}

	// TODO: change the following simple public functions to public variables
	if e, a := "Foo", ci[0].Name(); e != a {
		t.Fatalf("Expect `%s` for instruction set name. got: %s", e, a)
	}

	if e, a := "DefClass", ci[0].Type(); e != a {
		t.Fatalf("Expect `%s` for instruction set type. got: %s", e, a)
	}

	if e, a := "leave", ci[0].Instructions[0].ActionName(); e != a {
		t.Fatalf("Expect `%s` for first instruction action name. got: %s", e, a)
	}

	if e, a := 0, ci[0].Instructions[0].Line(); e != a {
		t.Fatalf("Expect `%d` for first instruction line. got: %d", e, a)
	}

	if e, a := 2, ci[0].Instructions[0].SourceLine(); e != a {
		t.Fatalf("Expect `%d` for first instruction source line. got: %d", e, a)
	}
}

func TestCompileToInstructionsREPLMode(t *testing.T) {

	ci, err := CompileToInstructions(`
def bar(a)
  99 + a
end
while true do
end
`, parser.REPLMode)

	if err != nil {
		t.Fatal(err.Error())
	}

	if e, a := uint8(10), ci[0].Instructions[0].Opcode; e != a {
		t.Fatalf("Expect `%d` for first instruction opcode. got: %d", e, a)
	}

	if e, a := "putobject: 99. source line: 3", ci[0].Instructions[0].Inspect(); e != a {
		t.Fatalf("Expect `%s` for first instruction inspect. got: %s", e, a)
	}

	// TODO: change the following simple public functions to public variables
	if e, a := "bar", ci[0].Name(); e != a {
		t.Fatalf("Expect `%s` for instruction set name. got: %s", e, a)
	}

	if e, a := "Def", ci[0].Type(); e != a {
		t.Fatalf("Expect `%s` for instruction set type. got: %s", e, a)
	}

	if e, a := "putobject", ci[0].Instructions[0].ActionName(); e != a {
		t.Fatalf("Expect `%s` for first instruction action name. got: %s", e, a)
	}

	if e, a := 0, ci[0].Instructions[0].Line(); e != a {
		t.Fatalf("Expect `%d` for first instruction line. got: %d", e, a)
	}

	if e, a := 3, ci[0].Instructions[0].SourceLine(); e != a {
		t.Fatalf("Expect `%d` for first instruction source line. got: %d", e, a)
	}
}

// TODO: The tests needs to be updated with igb/repl.go
func TestCompileToInstructionsREPLModeFail(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
iff
end
`, "unexpected end Line: 2"},
		{`
case
when 1
  11
when 2
  22
else
  99
end
`, "expected next token to be WHEN, got INT(1) instead. Line: 2",
		},
	}

	for _, tt := range tests {
		_, err := CompileToInstructions(tt.input, parser.REPLMode)

		if err.Error() != tt.expected {
			t.Fatalf("Expect `%s` error. got: %s", tt.expected, err.Error())
		}
	}
}
