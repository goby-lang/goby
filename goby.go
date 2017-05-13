package main

import (
	"flag"
	"fmt"
	"github.com/pkg/profile"
	"github.com/goby-lang/goby/bytecode"
	"github.com/goby-lang/goby/parser"
	"github.com/goby-lang/goby/vm"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"log"
	"path/filepath"
)

func main() {
	compileOptionPtr := flag.Bool("c", false, "Compile to bytecode")
	profileOptionPtr := flag.Bool("p", false, "Profile program execution")

	flag.Parse()

	if *profileOptionPtr {
		defer profile.Start().Stop()
	}

	filepath := flag.Arg(0)
	args := flag.Args()[1:]

	dir, filename, fileExt := extractFileInfo(filepath)
	file := readFile(filepath)

	switch fileExt {
	case "ro":
		program := parser.BuildAST(file)

		g := bytecode.NewGenerator(program)
		bytecodes := g.GenerateByteCode(program)

		if !*compileOptionPtr {
			v := vm.New(dir, args)
			v.ExecBytecodes(bytecodes, filepath)
			return
		}

		writeByteCode(bytecodes, dir, filename)
	case "robc":
		bytecodes := string(file)
		v := vm.New(dir, args)
		v.ExecBytecodes(bytecodes, filepath)
	default:
		fmt.Printf("Unknown file extension: %s", fileExt)
	}
}

func sourcePath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		log.Fatal(err)
	}

	return dir
}

func extractFileInfo(filepath string) (dir, filename, fileExt string) {
	dir, filename = path.Split(filepath)
	splitedFN := strings.Split(filename, ".")

	if len(splitedFN) <= 1 {
		fmt.Printf("Only support eval/compile single file now.")
		return
	}

	filename = splitedFN[0]
	fileExt = splitedFN[1]
	return
}

func writeByteCode(bytecodes, dir, filename string) {
	f, err := os.Create(dir + filename + ".robc")

	if err != nil {
		panic(err)
	}

	f.WriteString(bytecodes)
}

func readFile(filepath string) []byte {
	file, err := ioutil.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	return file
}
