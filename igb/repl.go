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
	prompt1   = "\033[32m»\033[0m "
	prompt2   = "\033[31m*\033[0m "
	pad       = "  "
	echo      = "#=>"
	interrupt = "^C"
	exit      = "exit"
	help      = "help"
	reset     = "reset"

	readyToExec = "readyToExec"
	Waiting     = "waiting"
	waitEnded   = "waitEnded"

	emojis = "😀😁😂🤣😃😄😅😆😉😊😋😎😍😘😗😙😚🙂🤗🤔😐😑😶🙄😏😮😪😴😌😛😜😝🤤🙃🤑😲😭😳🤧😇🤠🤡🤥🤓😈👿👹👺💀👻👽🤖💩😺😸😹😻😼😽"
)

// StartIgb starts goby's REPL.
func StartIgb(version string) {
reset:
	var sm = fsm.NewFSM(
		readyToExec,
		fsm.Events{
			{Name: Waiting, Src: []string{waitEnded, readyToExec}, Dst: Waiting},
			{Name: waitEnded, Src: []string{Waiting}, Dst: waitEnded},
			{Name: readyToExec, Src: []string{waitEnded, readyToExec}, Dst: readyToExec},
		},
		fsm.Callbacks{},
	)

	var completer = readline.NewPrefixCompleter(
		readline.PcItem(help),
		readline.PcItem(reset),
		readline.PcItem(exit),
	)

	var cmds []string
	var err error
	var stack = 0

	rl, err := readline.NewEx(&readline.Config{
		Prompt:              prompt1,
		HistoryFile:         "/tmp/readline_goby.tmp",
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

	println("Goby", version, fortune(), fortune(), fortune())

	// Initialize VM
	v := vm.New(os.Getenv("GOBY_ROOT"), []string{})
	v.SetClassISIndexTable("")
	v.SetMethodISIndexTable("")
	v.InitForREPL()

	// Initialize parser, lexer is not important here
	p := parser.New(lexer.New(""))

	program, _ := p.ParseProgram()

	// Initialize code generator, and it will behavior a little different in REPL mode.
	g := bytecode.NewGenerator()
	g.REPL = true
	g.InitTopLevelScope(program)

	for {
		rl.Config.UniqueEditLine = true
		line, err := rl.Readline()
		rl.Config.UniqueEditLine = false

		line = strings.TrimSpace(line)

		if err != nil {
			switch {
			case err == io.EOF:
				println(line + "")
				return
			case err == readline.ErrInterrupt: // Pressing Ctrl-C
				if len(line) == 0 && cmds == nil {
					println("")
					println("Bye!")
					return
				}
				// Erasing command buffer
				println("")
				stack = 0
				rl.SetPrompt(prompt1)
				sm.Event(Waiting)
				cmds = nil
				continue
			}
		}

		switch {
		case line == help:
			println(prompt(stack) + line)
			usage(rl.Stderr(), completer)
			continue
		case line == reset:
			rl = nil
			cmds = nil
			println(prompt(stack) + line)
			println("Restarting Igb...")
			goto reset
		case line == exit:
			println(prompt(stack) + line)
			println("Bye!")
			return
		case line == "":
			println(prompt(stack) + indent(stack) + line)
			continue
		}

		p.Lexer = lexer.New(line)
		program, perr := p.ParseProgram()

		if perr != nil {
			if perr.IsEOF() {
				if !sm.Is(Waiting) {
					sm.Event(Waiting)
				}

				println(prompt(stack) + indent(stack) + line)
				stack++
				rl.SetPrompt(prompt(stack) + indent(stack))
				cmds = append(cmds, line)
				continue
			}

			// If cmds is empty, it means that user just typed 'end' without corresponding statement/expression
			if perr.IsUnexpectedEnd() && len(cmds) != 0 {
				stack--
				rl.SetPrompt(prompt(stack) + indent(stack))
				sm.Event(waitEnded)
				cmds = append(cmds, line)
			} else {
				println(prompt(stack) + indent(stack) + line)
				stack = 0
				rl.SetPrompt(prompt1)
				fmt.Println(perr.Message)
				cmds = nil
				continue
			}

		}

		if sm.Is(Waiting) {
			println(prompt(stack) + indent(stack) + line)
			rl.SetPrompt(prompt(stack) + indent(stack))
			cmds = append(cmds, line)
			continue
		}

		if sm.Is(waitEnded) {
			p.Lexer = lexer.New(string(strings.Join(cmds, "\n")))

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
				println(prompt(stack) + indent(stack) + line)
				continue
			}

			// If everything goes well, reset state and statements buffer
			rl.SetPrompt(prompt(stack))
			sm.Event(readyToExec)
			cmds = nil
		}
		if sm.Is(readyToExec) {
			println(prompt(stack) + line)
			instructions := g.GenerateInstructions(program.Statements)
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
