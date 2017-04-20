package main

import (
	"flag"
	"fmt"
	"github.com/st0012/Rooby/ast"
	"github.com/st0012/Rooby/bytecode"
	"github.com/st0012/Rooby/lexer"
	"github.com/st0012/Rooby/parser"
	"github.com/st0012/Rooby/vm"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func main() {
	compileOptionPtr := flag.Bool("c", false, "Compile to bytecode")

	flag.Parse()

	filepath := flag.Arg(0)

	var fileExt string
	dir, filename := path.Split(filepath)
	splitedFN := strings.Split(filename, ".")

	if len(splitedFN) <= 1 {
		fmt.Printf("Only support eval/compile single file now.")
		return
	}

	filename = splitedFN[0]
	fileExt = splitedFN[1]

	file, err := ioutil.ReadFile(filepath)
	check(err)

	switch fileExt {
	case "ro":
		program := buildAST(file)

		g := bytecode.NewGenerator(program)
		bytecodes := g.GenerateByteCode(program)

		if !*compileOptionPtr {
			execBytecode(bytecodes)
			return
		}

		writeByteCode(bytecodes, dir, filename)

	case "robc":
		bytecodes := string(file)
		execBytecode(bytecodes)
	default:
		fmt.Printf("Unknown file extension: %s", fileExt)
	}
}


func writeByteCode(bytecodes, dir, filename string) {
	f, err := os.Create(dir + filename + ".robc")

	if err != nil {
		panic(err)
	}

	f.WriteString(bytecodes)
}

func execBytecode(bytecodes string) {
	p := vm.NewBytecodeParser()
	v := vm.New()
	p.VM = v
	p.Parse(bytecodes)
	cf := vm.NewCallFrame(v.LabelTable[vm.Program]["ProgramStart"][0])
	cf.Self = vm.MainObj
	v.CallFrameStack.Push(cf)
	v.Exec()
}

func buildAST(file []byte) *ast.Program {
	input := string(file)
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	p.CheckErrors()

	return program
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
