package igb

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	parserErr "github.com/gooby-lang/gooby/compiler/parser/errors"

	"github.com/chzyer/readline"
	"github.com/gooby-lang/gooby/compiler/bytecode"
	"github.com/gooby-lang/gooby/compiler/lexer"
	"github.com/gooby-lang/gooby/compiler/parser"
	"github.com/gooby-lang/gooby/vm"
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
	waiting     = "waiting"
	waitEnded   = "waitEnded"
	waitExited  = "waitExited"

	noMultiLineQuote     = "noMultiLineQuote"
	multiLineDoubleQuote = "multiLineDoubleQuote"
	multiLineSingleQuote = "multiLineSingleQuote"

	noNestedQuote = "noNestedQuote"
	nestedQuote   = "nestedQuote"

	emojis = "😀😁😂🤣😃😄😅😆😉😊😋😎😍😘😗😙😚🙂🤗🤔😐😑😶🙄😏😮😪😴😌😛😜😝🤤🙃🤑😲😭😳🤧😇🤠🤡🤥🤓😈👿👹👺💀👻👽🤖💩😺😸😹😻😼😽"
)

// iGb holds internal states of iGb.
type iGb struct {
	sm        *fsm.FSM
	qsm       *fsm.FSM
	nqsm      *fsm.FSM
	rl        *readline.Instance
	completer *readline.PrefixCompleter
	lines     string
	cmds      []string
	indents   int
	caseBlock bool
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

// StartIgb starts gooby's REPL.
func StartIgb(version string) {

reset:
	var err error
	igb := newIgb()

	igb.rl, err = readline.NewEx(&readline.Config{
		Prompt:              prompt1,
		HistoryFile:         filepath.Join(os.TempDir(), "readline_gooby.tmp"),
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

	println("Gooby", version, fortune(), fortune(), fortune())

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
				igb.eraseBuffer()
				println(" -- block cleared")
				continue
			}
		}

		// Multi-line quotation handling
		dq := checkOpenQuotes(igb.lines, `"`, `'`)
		sq := checkOpenQuotes(igb.lines, `'`, `"`)

		switch {
		case igb.qsm.Is(noMultiLineQuote):
			switch {
			case dq: // start multi-line double quote
				igb.qsm.Event(multiLineDoubleQuote)
				igb.startMultiLineQuote()
				continue
			case sq: // start multi-line single quote
				igb.qsm.Event(multiLineSingleQuote)
				igb.startMultiLineQuote()
				continue
			}

		case igb.qsm.Is(multiLineDoubleQuote):
			if igb.continueMultiLineQuote(dq) {
				continue
			}

		case igb.qsm.Is(multiLineSingleQuote):
			if igb.continueMultiLineQuote(sq) {
				continue
			}
		}

		// Command handling
		switch {
		case strings.HasPrefix(igb.lines, "#"):
			println(switchPrompt(igb.indents) + indent(igb.indents) + igb.lines)
			continue
		case igb.lines == help:
			println(switchPrompt(igb.indents) + igb.lines)
			usage(igb.rl.Stderr(), igb.completer)
			continue
		case igb.lines == reset:
			igb.rl = nil
			igb.cmds = nil
			println(switchPrompt(igb.indents) + igb.lines)
			println("Restarting iGb...")
			goto reset
		case igb.lines == "":
			println(switchPrompt(igb.indents) + indent(igb.indents) + igb.lines)
			continue
		}

		ivm.p.Lexer = lexer.New(igb.lines)
		program, pErr := ivm.p.ParseProgram()

		// Parse error handling
		if pErr != nil {
			switch {

			// To handle beginning of a block
			case pErr.IsEOF():
				if !igb.sm.Is(waiting) {
					igb.sm.Event(waiting)
				}
				igb.printLineAndIndentRight()
				continue

			// To handle 'case'
			case pErr.IsUnexpectedCase():
				println(switchPrompt(igb.indents) + indent(igb.indents) + igb.lines)
				igb.rl.SetPrompt(prompt2 + indent(igb.indents))
				igb.cmds = append(igb.cmds, igb.lines)
				igb.sm.Event(waiting)
				igb.caseBlock = true
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
				igb.cmds = append(igb.cmds, igb.lines)
				igb.rl.SetPrompt(prompt2 + indent(igb.indents))
				igb.sm.Event(waiting)
				igb.caseBlock = true
				continue

			// To handle such as 'else' or 'elsif'
			// The prompt should be `¤` even on the top level indentation when the line is `else` or `elif` or like that
			case pErr.IsUnexpectedToken():
				println(prompt2 + indent(igb.indents-1) + igb.lines)
				igb.rl.SetPrompt(prompt2 + indent(igb.indents))
				igb.cmds = append(igb.cmds, igb.lines)
				igb.sm.Event(waiting)
				continue

			// To handle empty line
			case pErr.IsUnexpectedEmptyLine(len(igb.cmds)):
				// If igb.cmds is empty, it means that user just typed 'end' without corresponding statement/expression
				println("exceptEmptyLine")
				println(switchPrompt(igb.indents) + indent(igb.indents) + igb.lines)
				igb.rl.SetPrompt(prompt1)
				fmt.Println(pErr.Message)
				igb.eraseBuffer()
				continue

			// To handle 'end'
			case pErr.IsUnexpectedEnd():
				if igb.indents > 1 {
					igb.indents--
					println(switchPrompt(igb.indents) + indent(igb.indents) + igb.lines)
					if igb.caseBlock {
						igb.indents++
					}
					igb.rl.SetPrompt(switchPrompt(igb.indents) + indent(igb.indents))
					igb.cmds = append(igb.cmds, igb.lines)
					igb.sm.Event(waiting)
					igb.caseBlock = false
					igb.firstWhen = false
					continue
				}

				// Exiting error handling
				igb.indents = 0
				igb.rl.SetPrompt(switchPrompt(igb.indents) + indent(igb.indents))
				igb.cmds = append(igb.cmds, igb.lines)
				igb.sm.Event(waitEnded)
				igb.caseBlock = false
				igb.firstWhen = false
			}
		}

		if igb.sm.Is(waiting) {
			// Indent = 0 but not ended
			if igb.caseBlock {
				igb.caseBlock = false
				println(prompt2 + indent(igb.indents) + igb.lines)
				igb.rl.SetPrompt(prompt2 + indent(igb.indents))
				igb.cmds = append(igb.cmds, igb.lines)
				continue
			}

			// Still indented
			if igb.indents > 0 {
				println(switchPrompt(igb.indents) + indent(igb.indents) + igb.lines)
				igb.rl.SetPrompt(switchPrompt(igb.indents) + indent(igb.indents))
				igb.cmds = append(igb.cmds, igb.lines)
				continue
			}
		}

		// Ending the block and prepare execution
		if igb.sm.Is(waitEnded) {
			ivm.p.Lexer = lexer.New(string(strings.Join(igb.cmds, "\n")))

			// Test if current input can be properly parsed.
			program, pErr = ivm.p.ParseProgram()

			if pErr != nil {
				handleParserError(pErr, igb)
				igb.eraseBuffer()
				continue
			}

			// If everything goes well, reset state and statements buffer
			igb.rl.SetPrompt(switchPrompt(igb.indents))
			igb.sm.Event(readyToExec)
		}

		// Execute the lines
		if igb.sm.Is(readyToExec) {
			println(switchPrompt(igb.indents) + igb.lines)

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
		println(switchPrompt(igb.indents) + indent(igb.indents) + igb.lines)
	}
}

// Polymorphic helper functions --------------------------------------------

// Returns true if indentation continues.
func (igb *iGb) continueMultiLineQuote(cdq bool) bool {
	if cdq { // end multi-line double quote
		igb.qsm.Event(noMultiLineQuote)
		if igb.nqsm.Is(nestedQuote) { // end nested-quote
			igb.nqsm.Event(noNestedQuote)
			igb.cmds = append(igb.cmds, igb.lines)
			println(switchPrompt(igb.indents) + igb.lines)
			igb.rl.SetPrompt(prompt2 + indent(igb.indents))
			return true
		}
	} else { // continue multi-line double quote
		println(prompt2 + igb.lines)
		igb.rl.SetPrompt(prompt2)
		igb.cmds = append(igb.cmds, igb.lines)
		return true
	}
	// exit multi-line double quote
	igb.cmds = append(igb.cmds, igb.lines)
	igb.rl.SetPrompt(prompt1)
	igb.sm.Event(waitEnded)
	return false
}

func (igb *iGb) eraseBuffer() {
	igb.indents = 0
	igb.rl.SetPrompt(prompt1)
	igb.sm.Event(waitExited)
	igb.qsm.Event(noMultiLineQuote)
	igb.cmds = nil
	igb.caseBlock = false
	igb.firstWhen = false
}

// Prints and add an indent.
func (igb *iGb) printLineAndIndentRight() {
	println(switchPrompt(igb.indents) + indent(igb.indents) + igb.lines)
	igb.indents++
	igb.rl.SetPrompt(switchPrompt(igb.indents) + indent(igb.indents))
	igb.cmds = append(igb.cmds, igb.lines)
}

// Starts multiple line quotation.
func (igb *iGb) startMultiLineQuote() {
	if igb.sm.Is(waiting) {
		igb.nqsm.Event(nestedQuote)
	} else {
		igb.sm.Event(waiting)
	}
	igb.cmds = append(igb.cmds, igb.lines)
	println(switchPrompt(igb.indents) + indent(igb.indents) + igb.lines)
	igb.rl.SetPrompt(prompt2)
}

// Other helper functions --------------------------------------------------

// filterInput just ignores Ctrl-z.
func filterInput(r rune) (rune, bool) {
	if r == readline.CharCtrlZ {
		return r, false
	}
	return r, true
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

// indent performs indentation with space padding.
func indent(c int) string {
	var s string
	for i := 0; i < c; i++ {
		s += pad
	}
	return s
}

// newIgb initializes iGb.
func newIgb() *iGb {
	return &iGb{
		cmds:    nil,
		indents: 0,
		// sm is for controlling the status of REPL.
		sm: fsm.NewFSM(
			readyToExec,
			fsm.Events{
				{Name: waiting, Src: []string{waitEnded, readyToExec}, Dst: waiting},
				{Name: waitEnded, Src: []string{waiting}, Dst: waitEnded},
				{Name: waitExited, Src: []string{waiting, waitEnded}, Dst: readyToExec},
				{Name: readyToExec, Src: []string{waitEnded, readyToExec}, Dst: readyToExec},
			},
			fsm.Callbacks{},
		),
		// qsm is for controlling the status of multi-line quotation.
		// Note that double and single multi-line quotations are exclusive and do not coexist.
		qsm: fsm.NewFSM(
			noMultiLineQuote,
			fsm.Events{
				{Name: noMultiLineQuote, Src: []string{multiLineDoubleQuote, multiLineSingleQuote}, Dst: noMultiLineQuote},
				{Name: multiLineDoubleQuote, Src: []string{noMultiLineQuote}, Dst: multiLineDoubleQuote},
				{Name: multiLineSingleQuote, Src: []string{noMultiLineQuote}, Dst: multiLineSingleQuote},
			},
			fsm.Callbacks{},
		),
		// nqsm is for controlling the status if the multi-line quotation is nested within other "waiting" statement.
		nqsm: fsm.NewFSM(
			noNestedQuote,
			fsm.Events{
				{Name: noNestedQuote, Src: []string{nestedQuote}, Dst: noNestedQuote},
				{Name: nestedQuote, Src: []string{noNestedQuote}, Dst: nestedQuote},
			},
			fsm.Callbacks{},
		),
		completer: readline.NewPrefixCompleter(
			readline.PcItem(help),
			readline.PcItem(reset),
		),
		caseBlock: false,
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

// Returns true if the specified quotation in the string is open.
func checkOpenQuotes(s, open, ignore string) bool {
	s = strings.Replace(s, `\\"`, "", -1)
	s = strings.Replace(s, `\\'`, "", -1)

	rq, _ := regexp2.Compile(`^[^`+ignore+`]*`+open, 0)
	isOpen, _ := rq.MatchString(s)
	if strings.Count(s, open)%2 == 1 && isOpen {
		return true
	}
	return false
}

// switchPrompt switches the prompt sign.
func switchPrompt(s int) string {
	if s > 0 {
		return prompt2
	}
	return prompt1
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

// usage shows help lines.
func usage(w io.Writer, c *readline.PrefixCompleter) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, c.Tree("   "))
}
