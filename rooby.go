package main

import (
	"github.com/st0012/Rooby/code_generator"
	"github.com/st0012/Rooby/lexer"
	"github.com/st0012/Rooby/parser"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"flag"
	"github.com/st0012/GVM"
	"github.com/st0012/Rooby/evaluator"
)

func main() {
	compileOptionPtr := flag.Bool("c", false, "Compile to bytecode")
	evalOptionPtr := flag.Bool("eval", true, "Eval program directly without using VM")

	flag.Parse()

	filepath := flag.Arg(0)

	file, err := ioutil.ReadFile(filepath)
	check(err)
	input := string(file)

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	p.CheckErrors()

	if *evalOptionPtr && !*compileOptionPtr {
		evaluator.Eval(program, evaluator.MainObj.Scope)
		return
	}

	cg := code_generator.New(program)

	bytecodes := cg.GenerateByteCode(program)

	if !*compileOptionPtr {
		gvm.Exec(bytecodes)
		return
	}

	writeByteCode(bytecodes, filepath)
}

func writeByteCode(bytecodes string, filepath string) {
	dir, filename := path.Split(filepath)
	filename = strings.Split(filename, ".")[0]
	f, err := os.Create(dir + filename + ".gbc")

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
