package main

import (
	"github.com/st0012/rooby/evaluator"
	"github.com/st0012/rooby/lexer"
	"github.com/st0012/rooby/object"
	"github.com/st0012/rooby/parser"
	"io/ioutil"
	"os"
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

	mainObj := object.InitializeMainObject()
	evaluator.Eval(program, mainObj.Scope)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
