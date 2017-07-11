package igb

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/compiler/lexer"
	"github.com/goby-lang/goby/compiler/parser"
	"github.com/goby-lang/goby/vm"
	"github.com/looplab/fsm"
)

const (
	prompt    = "\033[32m»\033[0m "
	prompt2   = "\033[31m*\033[0m "
	echo      = "#=>"
	interrupt = "^C"
	exit      = "exit"
	help      = "help"

	readyToExec = "readyToExec"
	Waiting     = "waiting"
	waitEnded   = "waitEnded"
)

var sm = fsm.NewFSM(
	readyToExec,
	fsm.Events{
		{Name: Waiting, Src: []string{waitEnded, readyToExec}, Dst: Waiting},
		{Name: waitEnded, Src: []string{Waiting}, Dst: waitEnded},
		{Name: readyToExec, Src: []string{waitEnded, readyToExec}, Dst: readyToExec},
	},
	fsm.Callbacks{},
)

var cmds []string

var completer = readline.NewPrefixCompleter(
	readline.PcItem(help),
	readline.PcItem(exit),
)

// StartIgb starts goby's REPL.
func StartIgb(version string) {
	var err error
	rl, err := readline.NewEx(&readline.Config{
		Prompt:              prompt,
		HistoryFile:         "/tmp/readline.tmp",
		AutoComplete:        completer,
		InterruptPrompt:     interrupt,
		EOFPrompt:           exit,
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	log.SetOutput(rl.Stderr())

	println("Goby", version)

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
		line, err := rl.Readline()

		if err == io.EOF {
			break
		}

		// Pressing ctrl-C
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				rl.SetPrompt(prompt2)
				continue
			}
		}

		line = strings.TrimSpace(line)

		switch {
		case line == help:
			usage(rl.Stderr())
			continue
		case line == exit:
			println("Bye!")
			return
		case line == "":
			continue
		}

		l := lexer.New(line)
		p.Lexer = l
		program, perr := p.ParseProgram()

		if perr != nil {
			if perr.IsEOF() {
				if !sm.Is(Waiting) {
					rl.SetPrompt(prompt2)
					sm.Event(Waiting)
				}

				rl.SetPrompt(prompt2)
				cmds = append(cmds, line)
				continue
			}

			// If cmds is empty, it means that user just typed 'end' without corresponding statement/expression
			if perr.IsUnexpectedEnd() && len(cmds) != 0 {
				rl.SetPrompt(prompt2)
				sm.Event(waitEnded)
				cmds = append(cmds, line)
			} else {
				rl.SetPrompt(prompt)
				fmt.Println(perr.Message)
				continue
			}

		}

		if sm.Is(Waiting) {
			rl.SetPrompt(prompt2)
			cmds = append(cmds, line)
			continue
		}

		if sm.Is(waitEnded) {
			l := lexer.New(string(strings.Join(cmds, "\n")))
			p.Lexer = l

			// Test if current input can be properly parsed.
			program, perr = p.ParseProgram()

			/*
				   This could mean there still are statements not ended, for example:

				   ```ruby
				   class Foo
					 def bar
					 end # This make state changes to WaitEnded
				   # But here still needs an "end"
				   ```
			*/

			if perr != nil {
				if !perr.IsEOF() {
					fmt.Println(perr.Message)
				}
				continue
			}

			// If everything goes well, reset state and statements buffer
			rl.SetPrompt(prompt)
			sm.Event(readyToExec)
			cmds = []string{}
		}
		if sm.Is(readyToExec) {
			instructions := g.GenerateInstructions(program.Statements)
			g.ResetInstructionSets()
			v.REPLExec(instructions)

			r := v.GetREPLResult()
			println(echo, r)
		}
	}
}

// Polymorphic helper functions --------------------------------------------

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

// Other helper functions --------------------------------------------------

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("   "))
}
