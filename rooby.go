package main

import (
	"github.com/st0012/Rooby/evaluator"
	"github.com/st0012/Rooby/lexer"
	"github.com/st0012/Rooby/parser"
	"io/ioutil"
	"os"
	"github.com/st0012/Rooby/object"
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

	mainObj := object.MainObj
	evaluator.Eval(program, mainObj.Scope)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
