package compiler

import (
	"testing"

	"github.com/gooby-lang/gooby/compiler/bytecode"
	"github.com/gooby-lang/gooby/compiler/parser"
)

type testInstruction struct {
	actionName string
	opCode     uint8
	sourceLine int
	paramsLen  int
}

func verifyInstructions(i *bytecode.Instruction, e testInstruction, t *testing.T) {
	if i.ActionName() != e.actionName || i.Opcode != e.opCode || i.SourceLine() != e.sourceLine || len(i.Params) != e.paramsLen {
		t.Fatalf("Line %d: expect ActionName: `%s`, opCode: %d, SourceLine: %d, paramsLen: %d. got: ActionName: `%s`, opCode: %d, SourceLine: %d, paramsLen: %d", i.Line(),
			e.actionName, e.opCode, e.sourceLine, e.paramsLen,
			i.ActionName(), i.Opcode, i.SourceLine(), len(i.Params))
	}
}

func TestCompileToInstructionsNormalMode(t *testing.T) {

	is, err := CompileToInstructions(`
def bar(a)
  99 + a
end
while true do
end`, parser.NormalMode)

	if err != nil {
		t.Fatal(err.Error())
	}

	tests := []struct {
		line     int
		expected testInstruction
	}{
		{
			0,
			testInstruction{actionName: "putobject", opCode: 10, sourceLine: 3, paramsLen: 1},
		},
		{
			1,
			testInstruction{actionName: "getlocal", opCode: 0, sourceLine: 3, paramsLen: 2},
		},
		{
			2,
			testInstruction{actionName: "send", opCode: 24, sourceLine: 3, paramsLen: 4},
		},
		{
			3,
			testInstruction{actionName: "leave", opCode: 29, sourceLine: 2, paramsLen: 0},
		},
	}
	for _, tt := range tests {
		i := is[0].Instructions[tt.line]
		verifyInstructions(i, tt.expected, t)
	}
}

func TestCompileToInstructionsTESTMode(t *testing.T) {

	is, err := CompileToInstructions(`
module Foo
end
`, parser.TestMode)

	if err != nil {
		t.Fatal(err.Error())
	}

	tests := []struct {
		line     int
		expected testInstruction
	}{
		{
			0,
			testInstruction{actionName: "leave", opCode: 29, sourceLine: 2, paramsLen: 0},
		},
	}
	for _, tt := range tests {
		i := is[0].Instructions[tt.line]
		verifyInstructions(i, tt.expected, t)
	}
}

func TestCompileToInstructionsREPLMode(t *testing.T) {

	is, err := CompileToInstructions(`
def bar(a)
  99 + a
end
while true do
end
`, parser.REPLMode)

	if err != nil {
		t.Fatal(err.Error())
	}

	tests := []struct {
		line     int
		expected testInstruction
	}{
		{
			0,
			testInstruction{actionName: "putobject", opCode: 10, sourceLine: 3, paramsLen: 1},
		},
		{
			1,
			testInstruction{actionName: "getlocal", opCode: 0, sourceLine: 3, paramsLen: 2},
		},
		{
			2,
			testInstruction{actionName: "send", opCode: 24, sourceLine: 3, paramsLen: 4},
		},
		{
			3,
			testInstruction{actionName: "leave", opCode: 29, sourceLine: 2, paramsLen: 0},
		},
	}
	for _, tt := range tests {
		i := is[0].Instructions[tt.line]
		verifyInstructions(i, tt.expected, t)
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
	}

	for _, tt := range tests {
		_, err := CompileToInstructions(tt.input, parser.REPLMode)

		if err.Error() != tt.expected {
			t.Fatalf("Expect `%s` error. got: %s", tt.expected, err.Error())
		}
	}
}
