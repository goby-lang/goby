package igb

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/compiler/lexer"
	"github.com/goby-lang/goby/compiler/parser"
	parserErr "github.com/goby-lang/goby/compiler/parser/errors"
	"github.com/goby-lang/goby/vm"
	"github.com/looplab/fsm"
	"github.com/mattn/go-colorable"
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
	help      = "help"
	reset     = "reset"

	readyToExec = "readyToExec"
	Waiting     = "waiting"
	waitEnded   = "waitEnded"
	waitExited  = "waitExited"

	emojis = "😀😁😂🤣😃😄😅😆😉😊😋😎😍😘😗😙😚🙂🤗🤔😐😑😶🙄😏😮😪😴😌😛😜😝🤤🙃🤑😲😭😳🤧😇🤠🤡🤥🤓😈👿👹👺💀👻👽🤖💩😺😸😹😻😼😽"
)

// iGb holds internal states of iGb.
type iGb struct {
	sm        *fsm.FSM
	rl        *readline.Instance
	completer *readline.PrefixCompleter
	lines     string
	cmds      []string
	indents   int
	tokenErr  bool
	firstWhen bool
}

// iVM holds VM only for iGb.
type iVM struct {
	v *vm.VM
	p *parser.Parser
	g *bytecode.Generator
}

var out io.Writer

func init() {
	if runtime.GOOS == "windows" {
		out = colorable.NewColorableStderr()
	} else {
		out = os.Stderr
	}
}

func println(s ...string) {
	out.Write([]byte(strings.Join(s, " ") + "\n"))
}

// StartIgb starts goby's REPL.
func StartIgb(version string) {
reset:
	var err error
	igb := newIgb()

	igb.rl, err = readline.NewEx(&readline.Config{
		Prompt:              prompt1,
		HistoryFile:         filepath.Join(os.TempDir(), "readline_goby.tmp"),
		AutoComplete:        igb.completer,
		InterruptPrompt:     interrupt,
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	defer igb.rl.Close()

	if err != nil {
		fmt.Printf("iGb error: %s", err)
		return
	}

	log.SetOutput(igb.rl.Stderr())

	println("Goby", version, fortune(), fortune(), fortune())

	ivm, err := newIVM()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for {
		igb, err = readIgb(igb, err)

		// Interruption handling
		if err != nil {
			switch {
			case err == io.EOF:
				println(igb.lines + "")
				return
			case err == readline.ErrInterrupt: // Pressing Ctrl-C
				if len(igb.lines) == 0 {
					if igb.cmds == nil {
						println("")
						println("Bye!")
						return
					}
				}
				// Erasing command buffer
				igb.indents = 0
				igb.rl.SetPrompt(prompt1)
				igb.sm.Event(waitExited)
				igb.cmds = nil
				igb.tokenErr = false
				igb.firstWhen = false
				println(" -- block cleared")
				continue
			}
		}

		// Command handling
		switch {
		case strings.HasPrefix(igb.lines, "#"):
			println(prompt(igb.indents) + indent(igb.indents) + igb.lines)
			continue
		case igb.lines == help:
			println(prompt(igb.indents) + igb.lines)
			usage(igb.rl.Stderr(), igb.completer)
			continue
		case igb.lines == reset:
			igb.rl = nil
			igb.cmds = nil
			println(prompt(igb.indents) + igb.lines)
			println("Restarting iGb...")
			goto reset
		case igb.lines == "":
			println(prompt(igb.indents) + indent(igb.indents) + igb.lines)
			continue
		}

		ivm.p.Lexer = lexer.New(igb.lines)
		program, pErr := ivm.p.ParseProgram()

		// Parse error handling
		if pErr != nil {
			switch {
			// To handle beginning of a block
			case pErr.IsEOF():
				if !igb.sm.Is(Waiting) {
					igb.sm.Event(Waiting)
				}
				println(prompt(igb.indents) + indent(igb.indents) + igb.lines)
				igb.indents++
				igb.rl.SetPrompt(prompt(igb.indents) + indent(igb.indents))
				igb.cmds = append(igb.cmds, igb.lines)
				continue
			// To handle 'case'
			case pErr.IsUnexpectedCase():
				println(prompt(igb.indents) + indent(igb.indents) + igb.lines)
				igb.sm.Event(Waiting)
				igb.rl.SetPrompt(prompt2 + indent(igb.indents))
				igb.cmds = append(igb.cmds, igb.lines)
				igb.tokenErr = true
				igb.firstWhen = true
				continue
			// To handle 'when'
			case pErr.IsUnexpectedWhen():
				if igb.firstWhen {
					igb.firstWhen = false
				} else {
					igb.indents--
				}
				println(prompt2 + indent(igb.indents) + igb.lines)
				igb.indents++
				igb.sm.Event(Waiting)
				igb.rl.SetPrompt(prompt2 + indent(igb.indents))
				igb.cmds = append(igb.cmds, igb.lines)
				igb.tokenErr = true
				continue
			// To handle such as 'else' or 'elsif'
			case pErr.IsUnexpectedToken():
				igb.indents--
				println(prompt2 + indent(igb.indents) + igb.lines)
				igb.indents++
				igb.sm.Event(Waiting)
				igb.rl.SetPrompt(prompt2 + indent(igb.indents))
				igb.cmds = append(igb.cmds, igb.lines)
				igb.tokenErr = true
				continue
			// To handle empty line
			case pErr.IsUnexpectedEmptyLine(len(igb.cmds)):
				// If igb.cmds is empty, it means that user just typed 'end' without corresponding statement/expression
				println(prompt(igb.indents) + indent(igb.indents) + igb.lines)
				igb.indents = 0
				igb.rl.SetPrompt(prompt1)
				fmt.Println(pErr.Message)
				igb.cmds = nil
				continue
			// To handle 'end'
			case pErr.IsUnexpectedEnd():
				if igb.indents > 1 {
					igb.indents--
					println(prompt(igb.indents) + indent(igb.indents) + igb.lines)
					if igb.tokenErr {
						igb.indents++
					}
					igb.sm.Event(Waiting)
					igb.rl.SetPrompt(prompt(igb.indents) + indent(igb.indents))
					igb.cmds = append(igb.cmds, igb.lines)
					igb.tokenErr = false
					igb.firstWhen = false
					continue
				}
				igb.indents = 0
				igb.sm.Event(waitEnded)
				igb.rl.SetPrompt(prompt(igb.indents) + indent(igb.indents))
				igb.cmds = append(igb.cmds, igb.lines)
				igb.tokenErr = false
				igb.firstWhen = false
			}
		}

		if igb.sm.Is(Waiting) && igb.tokenErr {
			igb.tokenErr = false
			println(prompt2 + indent(igb.indents) + igb.lines)
			igb.rl.SetPrompt(prompt2 + indent(igb.indents))
			igb.cmds = append(igb.cmds, igb.lines)
			continue
		}

		if igb.sm.Is(Waiting) && igb.indents > 0 {
			println(prompt(igb.indents) + indent(igb.indents) + igb.lines)
			igb.rl.SetPrompt(prompt(igb.indents) + indent(igb.indents))
			igb.cmds = append(igb.cmds, igb.lines)
			continue
		}

		if igb.sm.Is(waitEnded) {
			ivm.p.Lexer = lexer.New(string(strings.Join(igb.cmds, "\n")))

			// Test if current input can be properly parsed.
			program, pErr = ivm.p.ParseProgram()

			if pErr != nil {
				handleParserError(pErr, igb)
				igb.cmds = nil
				igb.sm.Event(readyToExec)
				continue
			}

			// If everything goes well, reset state and statements buffer
			igb.rl.SetPrompt(prompt(igb.indents))
			igb.sm.Event(readyToExec)
		}

		if igb.sm.Is(readyToExec) {
			println(prompt(igb.indents) + igb.lines)

			if pErr != nil {
				handleParserError(pErr, igb)
				continue
			}

			instructions := ivm.g.GenerateInstructions(program.Statements)
			ivm.v.REPLExec(instructions)

			r := ivm.v.GetREPLResult()

			// Suppress echo back on trailing ';'
			if igb.cmds != nil {
				if t := igb.cmds[len(igb.cmds)-1]; string(t[len(t)-1]) != semicolon {
					println(echo, r)
				}
			} else {
				if string(igb.lines[len(igb.lines)-1]) != semicolon {
					println(echo, r)
				}
			}
			//}
			igb.cmds = nil
		}
	}
}

func handleParserError(e *parserErr.Error, igb *iGb) {
	if e != nil {
		if !e.IsEOF() {
			fmt.Println(e.Message)
		}
		println(prompt(igb.indents) + indent(igb.indents) + igb.lines)
	}
}

// Polymorphic helper functions --------------------------------------------

// Other helper functions --------------------------------------------------

// newIgb initializes iGb.
func newIgb() *iGb {
	return &iGb{
		cmds:    nil,
		indents: 0,
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
		),
		tokenErr:  false,
		firstWhen: false,
	}
}

//
// newIVM initializes iVM.
func newIVM() (ivm iVM, err error) {
	ivm = iVM{}
	v, err := vm.New(os.Getenv("GOBY_ROOT"), []string{})

	if err != nil {
		return ivm, err
	}
	ivm.v = v
	ivm.v.SetClassISIndexTable("")
	ivm.v.SetMethodISIndexTable("")
	ivm.v.InitForREPL()
	// Initialize parser, lexer is not important here
	ivm.p = parser.New(lexer.New(""))
	ivm.p.Mode = parser.REPLMode
	program, _ := ivm.p.ParseProgram()
	// Initialize code generator, and it will behavior a little different in REPL mode.
	ivm.g = bytecode.NewGenerator()
	ivm.g.REPL = true
	ivm.g.InitTopLevelScope(program)
	return ivm, nil
}

// readIgb fetches one line from input, with helps of Readline lib.
func readIgb(igb *iGb, err error) (*iGb, error) {
	igb.rl.Config.UniqueEditLine = true // required to update the previous prompt
	igb.lines, err = igb.rl.Readline()
	igb.rl.Config.UniqueEditLine = false

	igb.lines = strings.TrimSpace(igb.lines)
	igb.lines = strings.TrimPrefix(igb.lines, prmpt1)
	igb.lines = strings.TrimPrefix(igb.lines, prmpt2)
	igb.lines = strings.TrimSpace(igb.lines)
	return igb, err
}

// filterInput just ignores Ctrl-z.
func filterInput(r rune) (rune, bool) {
	switch r {
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

// usage shows help lines.
func usage(w io.Writer, c *readline.PrefixCompleter) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, c.Tree("   "))
}

// indent performs indentation with space padding.
func indent(c int) string {
	var s string
	for i := 0; i < c; i++ {
		s = s + pad
	}
	return s
}

// prompt switches prompt sign.
func prompt(s int) string {
	if s > 0 {
		return prompt2
	}
	return prompt1
}

// fortune is just a fun item to show slot machine: receiving rep-digit would imply your fortune ;-)
func fortune() string {
	if runtime.GOOS == "windows" {
		return ""
	}
	var randSrc = rand.NewSource(time.Now().UnixNano())
	s := strings.Split(emojis, "")
	l := len(s)
	r := randSrc.Int63() % int64(l)
	return s[r]
}
