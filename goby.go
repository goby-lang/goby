package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/goby-lang/goby/compiler"
	"github.com/goby-lang/goby/compiler/parser"
	"github.com/goby-lang/goby/igb"
	"github.com/goby-lang/goby/vm"
	"github.com/pkg/profile"
)

const Version string = vm.Version

func main() {
	profileOptionPtr := flag.Bool("p", false, "Profile program execution")
	versionOptionPtr := flag.Bool("v", false, "Show current Goby version")
	interactiveOptionPtr := flag.Bool("i", false, "Run interactive goby")
	issueOptionPtr := flag.Bool("e", false, "Run interactive goby")

	flag.Parse()

	if *interactiveOptionPtr {
		igb.StartIgb(Version)
		os.Exit(0)
	}

	if *profileOptionPtr {
		defer profile.Start().Stop()
	}

	if *versionOptionPtr {
		fmt.Println(Version)
		os.Exit(0)
	}

	fp := flag.Arg(0)

	if fp == "" || !strings.Contains(fp, ".") {
		flag.Usage()
		os.Exit(0)
	}

	args := flag.Args()[1:]

	dir, _, fileExt := extractFileInfo(fp)
	file, ok := readFile(fp)

	if !ok {
		return
	}

	switch fileExt {
	case "gb", "rb":
		instructionSets, err := compiler.CompileToInstructions(string(file), parser.NormalMode)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var v *vm.VM

		if *issueOptionPtr {
			fmt.Println("Will report first issue...\n")
			v, err = vm.InitIssueReportVM(dir, args)
			defer vm.PrintError(v)
		} else {
			v, err = vm.New(dir, args)
		}

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fp, err := filepath.Abs(fp)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		v.ExecInstructions(instructionSets, fp)
	default:
		fmt.Printf("Unknown file extension: %s", fileExt)
	}
}

func extractFileInfo(fp string) (dir, filename, fileExt string) {
	dir, filename = filepath.Split(fp)
	dir, _ = filepath.Abs(dir)
	fileExt = filepath.Ext(fp)
	splited := strings.Split(filename, ".")
	filename, fileExt = splited[0], splited[1]
	return
}

func readFile(filepath string) (file []byte, ok bool) {
	file, err := ioutil.ReadFile(filepath)

	if err != nil {
		fmt.Println(err.Error())
		return []byte{}, false
	}

	return file, true
}
