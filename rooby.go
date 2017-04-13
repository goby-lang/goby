package main

import (
	"flag"
	"github.com/st0012/Rooby/vm"
	"github.com/st0012/Rooby/code_generator"
	"github.com/st0012/Rooby/lexer"
	"github.com/st0012/Rooby/parser"
	"io/ioutil"
	"os"
	"strings"
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
		execBytecode(bytecodes)
		return
	}
	writeByteCode(bytecodes, filepath)
}

func execBytecode(bytecodes string) {
	p := vm.NewBytecodeParser()
	p.Parse(bytecodes)
	v := vm.New()
	p.VM = v
	cf := vm.NewCallFrame(v.LabelTable[vm.PROGRAM]["ProgramStart"][0])
	cf.Self = vm.MainObj
	v.CallFrameStack.Push(cf)
	v.Exec()
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
