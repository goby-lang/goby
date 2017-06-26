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

const(
	Initial = "initial"
	Wait = "wait"
	WaitEnded = "waitEnded"
)

var sm = fsm.NewFSM(
	"initial",
	/*
		Initial state is default state
		Nosymbol state helps us identify tok ':' is for symbol or hash value
		Method state helps us identify 'class' literal is a keyword or an identifier
		Reference: https://github.com/looplab/fsm
	*/
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
	v := vm.New(os.Getenv("GOBY_ROOT"), []string{})
	v.InitForREPL()
	l := lexer.New("")
	p := parser.New(l)
	g := bytecode.NewGenerator()
	g.REPL = true
	program, _ := p.ParseProgram()
	bytecodes := g.GenerateByteCode(program, true) // Set generator's new scope
	v.ExecBytecodes(bytecodes, "")
	v.SetClassISIndexTable("")
	v.SetMethodISIndexTable("")

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
