package main

import (
	"github.com/st0012/Rooby/lexer"
	"github.com/st0012/Rooby/parser"
	"io/ioutil"
	"os"
	"github.com/st0012/Rooby/code_generator"
	"path"
	"strings"
)

func main() {
	filepath := os.Args[1]

	file, err := ioutil.ReadFile(filepath)
	check(err)
	input := string(file)

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	p.CheckErrors()

	bytecodes := code_generator.GenerateByteCode(program)
	writeByteCode(bytecodes, filepath)
	//evaluator.Eval(program, evaluator.MainObj.Scope)
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
