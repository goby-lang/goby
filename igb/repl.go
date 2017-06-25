package igb

import (
	"bufio"
	"fmt"
	"github.com/goby-lang/goby/bytecode"
	"github.com/goby-lang/goby/lexer"
	"github.com/goby-lang/goby/parser"
	"github.com/goby-lang/goby/vm"
	"io"
	"os"
)

const PROMT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	v := vm.New(os.Getenv("GOBY_ROOT"), []string{})
	l := lexer.New("")
	p := parser.New(l)
	g := bytecode.NewGenerator()

	for {
		fmt.Printf(PROMT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p.Lexer = l
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		bytecodes := g.GenerateByteCode(program)
		fmt.Println(bytecodes)
		v.ReplExec(bytecodes, os.Getenv("GOBY_ROOT"))
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
