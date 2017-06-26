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

const PROMT = ">> "

const (
	Initial   = "initial"
	Wait      = "wait"
	WaitEnded = "waitEnded"
)

var sm = fsm.NewFSM(
	"initial",
	fsm.Events{
		{Name: Wait, Src: []string{Initial}, Dst: Wait},
		{Name: WaitEnded, Src: []string{Wait}, Dst: WaitEnded},
		{Name: Initial, Src: []string{WaitEnded, Initial}, Dst: Initial},
	},
	fsm.Callbacks{},
)
var stmts = bytes.Buffer{}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	// Initialize VM
	v := vm.New(os.Getenv("GOBY_ROOT"), []string{})
	v.SetClassISIndexTable("")
	v.SetMethodISIndexTable("")
	v.InitForREPL()

	// Initialize code generator, and it will behavior a little different in REPL mode.
	g := bytecode.NewGenerator()
	g.REPL = true

	// Initialize parser, lexer is not important here
	l := lexer.New("")
	p := parser.New(l)
	program, _ := p.ParseProgram()

	// Set generator's new scope, program here is not important
	// The result bytecode should be just "<ProgramStart>", also can be ignored.
	_ = g.GenerateByteCode(program, true)

	for {
		fmt.Printf(PROMT)
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
				if sm.Is(Initial) {
					sm.Event(Wait)
				}
				stmts.WriteString(line + "\n")

				continue
			}

			if err.IsUnexpectedEnd() {
				stmts.WriteString(line + "\n")

				sm.Event(WaitEnded)
			} else {
				fmt.Println(err.Message)
				continue
			}

		}

		if sm.Is(Wait) {
			stmts.WriteString(line + "\n")
			continue
		}

		if sm.Is(WaitEnded) {
			l := lexer.New(stmts.String())
			p.Lexer = l
			program, err = p.ParseProgram()

			if err != nil {
				if err.IsEOF() {
					continue
				} else {
					fmt.Println(err.Message)
					continue
				}
			}

			sm.Event(Initial)
			stmts.Reset()
		}

		if sm.Is(Initial) {
			bytecodes := g.GenerateByteCode(program, false)
			g.ResetInstructionSets()
			v.REPLExec(bytecodes)
			fmt.Println(v.GetExecResultToString())
		}
	}
}
