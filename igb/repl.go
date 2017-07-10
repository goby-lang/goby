package igb

import (
	"io"
	"strings"

	"github.com/chzyer/readline"
)

const (
	//readyToExec = "readyToExec"
	//Waiting     = "waiting"
	//waitEnded   = "waitEnded"

	prompt    = "\033[31m»\033[0m "
	echo      = "#=>"
	interrupt = "^C"
	exit      = "exit"
	help      = "help"
)

var completer = readline.NewPrefixCompleter(
	readline.PcItem(help),
	readline.PcItem(exit),
)

//var sm = fsm.NewFSM(
//	readyToExec,
//	fsm.Events{
//		{Name: Waiting, Src: []string{waitEnded, readyToExec}, Dst: Waiting},
//		{Name: waitEnded, Src: []string{Waiting}, Dst: waitEnded},
//		{Name: readyToExec, Src: []string{waitEnded, readyToExec}, Dst: readyToExec},
//	},
//	fsm.Callbacks{},
//)

//var stmts = bytes.Buffer{}

// Start starts goby's REPL.
func Start(version string) {
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

	//log.SetOutput(rl.Stderr())

	println("Goby", version)

	//// Initialize VM
	//v := vm.New(os.Getenv("GOBY_ROOT"), []string{})
	//v.SetClassISIndexTable("")
	//v.SetMethodISIndexTable("")
	//v.InitForREPL()
	//
	//// Initialize parser, lexer is not important here
	//l := lexer.New("")
	//p := parser.New(l)
	//program, _ := p.ParseProgram()
	//
	//// Initialize code generator, and it will behavior a little different in REPL mode.
	//g := bytecode.NewGenerator()
	//g.REPL = true
	//g.InitTopLevelScope(program)

	for {
		line, err := rl.Readline()

		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case line == help:
			usage(rl.Stderr())
		case line == exit:
			println("Bye!")
			return
		case line == "":
		default:
			println(echo, line)
			//l := lexer.New(line)
			//p.Lexer = l
			//program, _ := p.ParseProgram()
			//instructions := g.GenerateInstructions(program.Statements)
			//g.ResetInstructionSets()
			//v.REPLExec(instructions)
			//
			//r := v.GetREPLResult()
			//println(echo, r)
		}
	}
}

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}
