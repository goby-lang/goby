package vm

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/compiler/lexer"
	"github.com/goby-lang/goby/compiler/parser"
	"github.com/goby-lang/goby/vm/errors"
)

func TestVM_REPLExec(t *testing.T) {
	tests := []struct {
		inputs   []string
		expected interface{}
	}{
		{
			[]string{
				`
				a, b = [3, 6]
				a + b
				`,
			},
			9,
		}, {
			[]string{
				`
				a, _, c = [1, 2, 3]
				c
				`,
			},
			3,
		}, {
			[]string{`
def foo(x)
  yield(x + 10)
end
def bar(y)
  foo(y) do |f|
	yield(f)
  end
end
def baz(z)
  bar(z + 100) do |b|
	yield(b)
  end
end
a = 0
baz(100) do |b|
  a = b
end
a
			`,
				`
class Foo
  def bar
	100
  end
end
module Baz
  class Bar
	def bar
	  Foo.new.bar
	end
  end
end
				`,
				`
Baz::Bar.new.bar + a
`,
			},
			310,
		},
		{
			[]string{`

def foo
  123
end
`, `
foo
`},
			123},
		{
			[]string{
				`
class Foo
  def bar(x)
    x + 10
  end
end
`,
				`
Foo.new.bar(90)
`},
			100},
		{
			[]string{
				`
def foo
  123
end
`,
				`
foo
`,
				`
def foo
  345
end
`,

				`
foo
`,
			}, 345},
	}

	for i, test := range tests {
		v := initTestVM()
		v.InitForREPL()

		// Initialize parser, lexer is not important here
		p := parser.New(lexer.New(""))
		p.Mode = parser.REPLMode

		program, _ := p.ParseProgram()

		// Initialize code generator, and it will behavior a little different in REPL mode.
		g := bytecode.NewGenerator()
		g.REPL = true
		g.InitTopLevelScope(program)

		for _, input := range test.inputs {
			p := parser.New(lexer.New(input))
			p.Mode = parser.REPLMode

			program, _ := p.ParseProgram()
			sets := g.GenerateInstructions(program.Statements)

			v.REPLExec(sets)
		}

		evaluated := v.GetExecResult()
		VerifyExpected(t, i, evaluated, test.expected)
		// Because REPL should maintain a base call frame so that the whole program won't exit
		v.checkCFP(t, i, 1)
	}
}

func TestVM_REPLExecFail(t *testing.T) {

	tests := []struct {
		inputs   []string
		expected string
	}{
		{
			[]string{
				`raise ArgumentError`,
			},
			fmt.Sprintf("InternalError: '%s'", errors.ArgumentError),
		},
		{
			[]string{
				"NonExistentBuiltinMethod",
			},
			"NameError: uninitialized constant NonExistentBuiltinMethod",
		},
		{
			[]string{
				"Hash.notExist",
			},
			"UndefinedMethodError: Undefined Method 'notExist' for Hash",
		},
	}

	for i, test := range tests {
		v := initTestVM()
		v.InitForREPL()

		// Initialize parser, lexer is not important here
		p := parser.New(lexer.New(""))
		p.Mode = parser.REPLMode

		program, _ := p.ParseProgram()

		// Initialize code generator, and it will behavior a little different in REPL mode.
		g := bytecode.NewGenerator()
		g.REPL = true
		g.InitTopLevelScope(program)

		for _, input := range test.inputs {
			p := parser.New(lexer.New(input))
			p.Mode = parser.REPLMode

			// prevent parse errors from panicking tests
			program, Err := p.ParseProgram()
			if Err != nil {
				t.Fatalf("At case %d unexpected parse error %q", i, Err.Message)
			}
			sets := g.GenerateInstructions(program.Statements)

			v.REPLExec(sets)
		}

		evaluated := v.GetExecResult()

		if evaluated.toString() != test.expected {
			t.Fatalf("At case %d expected %s got %s", i,
				test.expected, evaluated.toString())
		}

		// Because REPL should maintain a base call frame so that the whole program won't exit
		v.checkCFP(t, i, 1)
	}
}

func TestLoadingGobyLibraryFail(t *testing.T) {
	vm := initTestVM()

	libFileFullPath := filepath.Join(vm.projectRoot, "lib/_library_not_existing.gb")
	expectedErrorMessage := fmt.Sprintf("open %s: no such file or directory", libFileFullPath)

	err := vm.mainThread.execGobyLib("_library_not_existing.gb")

	if err.Error() != expectedErrorMessage {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}

func (v *VM) checkCFP(t *testing.T, index, expectedCFP int) {
	t.Helper()
	if v.mainThread.callFrameStack.pointer != expectedCFP {
		t.Errorf("At case %d expect main thread's cfp to be %d. got: %d", index, expectedCFP, v.mainThread.callFrameStack.pointer)
	}
}

func (v *VM) checkSP(t *testing.T, index, expectedSp int) {
	t.Helper()
	if v.mainThread.Stack.pointer != expectedSp {
		fmt.Println(v.mainThread.Stack.inspect())
		t.Errorf("At case %d expect main thread's sp to be %d. got: %d", index, expectedSp, v.mainThread.Stack.pointer)
	}

}
