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
	issueOptionPtr := flag.Bool("e", false, "Generate reporting format")

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

	var fp string

	switch flag.Arg(0) {
	case "":
		flag.Usage()
		os.Exit(0)
	case "test":
		args := flag.Args()[1:]
		fp := flag.Arg(1)
		fi, err := os.Stat(fp)
		reportErrorAndExit(err)

		var dir string

		// To get the directory path
		if fi.Mode().IsDir() {
			dir, err = filepath.Abs(fi.Name())
			reportErrorAndExit(err)
		} else {
			fp, err = filepath.Abs(fp)
			reportErrorAndExit(err)
			fmt.Println(fp)
			dir, _, _ = extractFileInfo(fp)
		}

		v, err := vm.New(dir, args)

		execTestFile := func(dir, fp string) (err error) {
			file := readFile(fp)
			instructionSets, err := compiler.CompileToInstructions(string(file), parser.NormalMode)

			if err != nil {
				return
			}

			v.ExecInstructions(instructionSets, fp)
			return
		}

		if fi.Mode().IsDir() {
			fileInfos, err := ioutil.ReadDir(fp)
			reportErrorAndExit(err)

			for _, fileInfo := range fileInfos {
				fp := filepath.Join(dir, fileInfo.Name())
				reportErrorAndExit(err)

				err := execTestFile(dir, fp)
				reportErrorAndExit(err)
			}
		} else {
			err := execTestFile(dir, fp)
			reportErrorAndExit(err)
		}

		instructionSets, err := compiler.CompileToInstructions("Spec.test", parser.NormalMode)
		v.ExecInstructions(instructionSets, fp)
		return
	default:
		fp = flag.Arg(0)

		if !strings.Contains(fp, ".") {
			flag.Usage()
			os.Exit(0)
		}
	}

	dir, _, fileExt := extractFileInfo(fp)
	file := readFile(fp)

	switch fileExt {
	case "gb", "rb":
		args := flag.Args()[1:]
		instructionSets, err := compiler.CompileToInstructions(string(file), parser.NormalMode)
		reportErrorAndExit(err)

		var v *vm.VM

		if *issueOptionPtr {
			fmt.Println("Will generate issue report on error...")
			v, err = vm.InitIssueReportVM(dir, args)
			defer vm.PrintError(v)
		} else {
			v, err = vm.New(dir, args)
		}
		reportErrorAndExit(err)

		fp, err := filepath.Abs(fp)
		reportErrorAndExit(err)

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

func readFile(filepath string) (file []byte) {
	file, err := ioutil.ReadFile(filepath)
	reportErrorAndExit(err)
	return
}

func reportErrorAndExit(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}