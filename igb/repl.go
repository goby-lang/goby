package igb

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/compiler/lexer"
	"github.com/goby-lang/goby/compiler/parser"
	"github.com/goby-lang/goby/vm"
	"github.com/looplab/fsm"
)

const (
	prmpt1    = "»"
	prmpt2    = "¤"
	prompt1   = "\033[32m" + prmpt1 + "\033[0m "
	prompt2   = "\033[31m" + prmpt2 + "\033[0m "
	pad       = "  "
	echo      = "\033[33m#»\033[0m"
	interrupt = "^C"
	semicolon = ";"
	exit      = "exit"
	help      = "help"
	reset     = "reset"

	readyToExec = "readyToExec"
	Waiting     = "waiting"
	waitEnded   = "waitEnded"
	waitExited  = "waitExited"

	emojis = "😀😁😂🤣😃😄😅😆😉😊😋😎😍😘😗😙😚🙂🤗🤔😐😑😶🙄😏😮😪😴😌😛😜😝🤤🙃🤑😲😭😳🤧😇🤠🤡🤥🤓😈👿👹👺💀👻👽🤖💩😺😸😹😻😼😽"
)

type Igb struct {
	sm        *fsm.FSM
	rl        *readline.Instance
	completer *readline.PrefixCompleter
	line      string
	cmds      []string
	stack     int
}

type Ivm struct {
	v *vm.VM
	p *parser.Parser
	g *bytecode.Generator
}

// StartIgb starts goby's REPL.
func StartIgb(version string) {
reset:
	var err error
	igb := initIgb()

	igb.rl, err = readline.NewEx(&readline.Config{
		Prompt:              prompt1,
		HistoryFile:         "/tmp/readline_goby.tmp",
		AutoComplete:        igb.completer,
		InterruptPrompt:     interrupt,
		EOFPrompt:           exit,
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	defer igb.rl.Close()

	if err != nil {
		fmt.Printf("Igb error: %s", err)
		return
	}

	log.SetOutput(igb.rl.Stderr())

	println("Goby", version, fortune(), fortune(), fortune())

	ivm := createVM()

	for {
		igb.rl.Config.UniqueEditLine = true
		igb.line, err = igb.rl.Readline()
		igb.rl.Config.UniqueEditLine = false

		igb.line = strings.TrimPrefix(igb.line, prmpt1)
		igb.line = strings.TrimPrefix(igb.line, prmpt2)
		igb.line = strings.TrimSpace(igb.line)

		if err != nil {
			switch {
			case err == io.EOF:
				println(igb.line + "")
				return
			case err == readline.ErrInterrupt: // Pressing Ctrl-C
				if len(igb.line) == 0 {
					if igb.cmds == nil {
						println("")
						println("Bye!")
						return
					}
				}
				// Erasing command buffer
				igb.stack = 0
				igb.rl.SetPrompt(prompt1)
				igb.sm.Event(waitExited)
				igb.cmds = nil
				println(" -- block cleared")
				continue
			}
		}

		switch {
		case strings.HasPrefix(igb.line, "#"):
			println(prompt(igb.stack) + igb.line)
			continue
		case igb.line == help:
			println(prompt(igb.stack) + igb.line)
			usage(igb.rl.Stderr(), igb.completer)
			continue
		case igb.line == reset:
			igb.rl = nil
			igb.cmds = nil
			println(prompt(igb.stack) + igb.line)
			println("Restarting Igb...")
			goto reset
		case igb.line == exit:
			println(prompt(igb.stack) + igb.line)
			println("Bye!")
			return
		case igb.line == "":
			println(prompt(igb.stack) + indent(igb.stack) + igb.line)
			continue
		}

		ivm.p.Lexer = lexer.New(igb.line)
		program, perr := ivm.p.ParseProgram()

		if perr != nil {
			if perr.IsEOF() {
				if !igb.sm.Is(Waiting) {
					igb.sm.Event(Waiting)
				}
				println(prompt(igb.stack) + indent(igb.stack) + igb.line)
				igb.stack++
				igb.rl.SetPrompt(prompt(igb.stack) + indent(igb.stack))
				igb.cmds = append(igb.cmds, igb.line)
				continue
			}

			// If igb.cmds is empty, it means that user just typed 'end' without corresponding statement/expression
			if perr.IsUnexpectedEnd() && len(igb.cmds) == 0 {
				println(prompt(igb.stack) + indent(igb.stack) + igb.line)
				igb.stack = 0
				igb.rl.SetPrompt(prompt1)
				fmt.Println(perr.Message)
				igb.cmds = nil
				continue
			}

			if perr.IsUnexpectedEnd() {
				if igb.stack > 1 {
					igb.stack--
					println(prompt(igb.stack) + indent(igb.stack) + igb.line)
					igb.sm.Event(Waiting)
					igb.rl.SetPrompt(prompt(igb.stack) + indent(igb.stack))
					igb.cmds = append(igb.cmds, igb.line)
					continue
				}
				igb.stack = 0
				igb.sm.Event(waitEnded)
				igb.rl.SetPrompt(prompt(igb.stack) + indent(igb.stack))
				igb.cmds = append(igb.cmds, igb.line)
			}
		}

		if igb.sm.Is(Waiting) && igb.stack > 0 {
			println(prompt(igb.stack) + indent(igb.stack) + igb.line)
			igb.rl.SetPrompt(prompt(igb.stack) + indent(igb.stack))
			igb.cmds = append(igb.cmds, igb.line)
			continue
		}

		if igb.sm.Is(waitEnded) {
			ivm.p.Lexer = lexer.New(string(strings.Join(igb.cmds, "\n")))

			// Test if current input can be properly parsed.
			program, perr = ivm.p.ParseProgram()

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
				println(prompt(igb.stack) + indent(igb.stack) + igb.line)
				continue
			}

			// If everything goes well, reset state and statements buffer
			igb.rl.SetPrompt(prompt(igb.stack))
			igb.sm.Event(readyToExec)
		}
		if igb.sm.Is(readyToExec) {
			println(prompt(igb.stack) + igb.line)
			instructions := ivm.g.GenerateInstructions(program.Statements)
			ivm.v.REPLExec(instructions)

			r := ivm.v.GetREPLResult()

			// Suppress echo back on trailing ';'
			if igb.cmds != nil {
				if t := igb.cmds[len(igb.cmds)-1]; string(t[len(t)-1]) != semicolon {
					println(echo, r)
				}
			} else {
				if string(igb.line[len(igb.line)-1]) != semicolon {
					println(echo, r)
				}
			}
			//}
			igb.cmds = nil
		}
	}
}

// Polymorphic helper functions --------------------------------------------

// Other helper functions --------------------------------------------------

func initIgb() *Igb {
	return &Igb{
		cmds:  nil,
		stack: 0,
		sm: fsm.NewFSM(
			readyToExec,
			fsm.Events{
				{Name: Waiting, Src: []string{waitEnded, readyToExec}, Dst: Waiting},
				{Name: waitEnded, Src: []string{Waiting}, Dst: waitEnded},
				{Name: waitExited, Src: []string{Waiting, waitEnded}, Dst: readyToExec},
				{Name: readyToExec, Src: []string{waitEnded, readyToExec}, Dst: readyToExec},
			},
			fsm.Callbacks{},
		),
		completer: readline.NewPrefixCompleter(
			readline.PcItem(help),
			readline.PcItem(reset),
			readline.PcItem(exit),
		),
	}
}

func createVM() Ivm {
	// Initialize VM
	ivm := Ivm{}
	ivm.v = vm.New(os.Getenv("GOBY_ROOT"), []string{})
	ivm.v.SetClassISIndexTable("")
	ivm.v.SetMethodISIndexTable("")
	ivm.v.InitForREPL()
	// Initialize parser, lexer is not important here
	ivm.p = parser.New(lexer.New(""))
	program, _ := ivm.p.ParseProgram()
	// Initialize code generator, and it will behavior a little different in REPL mode.
	ivm.g = bytecode.NewGenerator()
	ivm.g.REPL = true
	ivm.g.InitTopLevelScope(program)
	return ivm
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func usage(w io.Writer, c *readline.PrefixCompleter) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, c.Tree("   "))
}

func indent(c int) string {
	var s string
	for i := 0; i < c; i++ {
		s = s + pad
	}
	return s
}

func prompt(s int) string {
	if s > 0 {
		return prompt2
	}
	return prompt1
}

func fortune() string {
	var randSrc = rand.NewSource(time.Now().UnixNano())
	s := strings.Split(emojis, "")
	l := len(s)
	r := randSrc.Int63() % int64(l)
	return s[r]
}
