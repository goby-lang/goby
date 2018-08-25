package compiler

import (
	"testing"

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

	if e, a := 99, ci[0].Instructions[0].Params[1]; e != a {
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

//func TestCompileToInstructionsNormalModeFail(t *testing.T) {
//
//	ci, err := CompileToInstructions(`
//def bar
//  99
//end `, parser.NormalMode)
//
//	if err != nil {
//		t.Fatal(err.Error())
//	}
//
//	if e, a := 1, ci[0].Instructions[0].AnchorLine(); e != a {
//		t.Fatalf("Expect `%d` for first instruction inspect. got: %d", e, a)
//	}
//}
