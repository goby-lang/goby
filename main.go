package main

import (
	"fmt"
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

	env := object.NewEnvironment()
	mainObj := &object.Main{Env: env}
	scope := &object.Scope{Self: mainObj, Env: env}
	result := evaluator.Eval(program, scope).Inspect()

	fmt.Print(result)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
