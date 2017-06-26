package igb

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/goby-lang/goby/bytecode"
	"github.com/goby-lang/goby/lexer"
	"github.com/goby-lang/goby/parser"
	"github.com/goby-lang/goby/vm"
	"io"
	"os"

	"github.com/goby-lang/goby/Godeps/_workspace/src/github.com/looplab/fsm"
)

const (
	prompt = ">> "

	initial   = "initial"
	wait      = "wait"
	waitEnded = "waitEnded"
)

var sm = fsm.NewFSM(
	"initial",
	fsm.Events{
		{Name: wait, Src: []string{initial}, Dst: wait},
		{Name: waitEnded, Src: []string{wait}, Dst: waitEnded},
		{Name: initial, Src: []string{waitEnded, initial}, Dst: initial},
	},
	fsm.Callbacks{},
)
var stmts = bytes.Buffer{}

// Start starts goby's REPL.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	// Initialize VM
	v := vm.New(os.Getenv("GOBY_ROOT"), []string{})
	v.SetClassISIndexTable("")
	v.SetMethodISIndexTable("")
	v.InitForREPL()

	// Initialize parser, lexer is not important here
	l := lexer.New("")
	p := parser.New(l)
	program, _ := p.ParseProgram()

	// Initialize code generator, and it will behavior a little different in REPL mode.
	g := bytecode.NewGenerator()
	g.REPL = true
	g.InitTopLevelScope(program)

	for {
		out.Write([]byte(prompt))
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p.Lexer = l
		program, err := p.ParseProgram()

		if err != nil {
			if err.IsEOF() {
				if sm.Is(initial) {
					sm.Event(wait)
				}

				appendStmt(line)
				continue
			}

			if err.IsUnexpectedEnd() {
				sm.Event(waitEnded)
				appendStmt(line)
			} else {
				fmt.Println(err.Message)
				continue
			}

		}

		if sm.Is(wait) {
			appendStmt(line)
			continue
		}

		if sm.Is(waitEnded) {
			l := lexer.New(stmts.String())
			p.Lexer = l

			// Test if current input can be properly parsed.
			program, err = p.ParseProgram()

			/*
			 This could mean there still are statements not ended, for example:

			 ```ruby
			 class Foo
			   def bar
			   end # This make state changes to WaitEnded
			 # But here still needs an "end"
			 ```
			*/

			if err != nil {
				if !err.IsEOF() {
					fmt.Println(err.Message)
				}
				continue
			}

			// If everything goes well, reset state and statements buffer
			sm.Event(initial)
			stmts.Reset()
		}

		if sm.Is(initial) {
			bytecodes := g.GenerateByteCode(program.Statements)
			g.ResetInstructionSets()
			v.REPLExec(bytecodes)
			out.Write([]byte(v.GetREPLResult() + "\n"))
		}
	}
}

func appendStmt(line string) {
	stmts.WriteString(line + "\n")
}
