package main

import (
	"flag"
	"github.com/st0012/GVM"
	"github.com/st0012/Rooby/code_generator"
	"github.com/st0012/Rooby/lexer"
	"github.com/st0012/Rooby/parser"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	execOptionPtr := flag.Bool("c", false, "Compile to bytecode")

	flag.Parse()

	filepath := flag.Arg(0)

	file, err := ioutil.ReadFile(filepath)
	check(err)
	input := string(file)

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	p.CheckErrors()
	cg := code_generator.New(program)

	bytecodes := cg.GenerateByteCode(program)

	if !*execOptionPtr {
		gvm.Exec(bytecodes)
		return
	}
	writeByteCode(bytecodes, filepath)
	//evaluator.Eval(program, evaluator.MainObj.Scope)
}

func writeByteCode(bytecodes string, filepath string) {
	filepath = strings.Split(filepath, ".")[0] + ".gbc"
	f, err := os.Create(filepath)

	if err != nil {
		panic(err)
	}

	f.WriteString(bytecodes)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
