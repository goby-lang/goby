package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gooby-lang/gooby/compiler"
	"github.com/gooby-lang/gooby/compiler/parser"
	"github.com/gooby-lang/gooby/igb"
	_ "github.com/gooby-lang/gooby/native/db"
	_ "github.com/gooby-lang/gooby/native/plugin"
	_ "github.com/gooby-lang/gooby/native/result"
	_ "github.com/gooby-lang/gooby/native/ripper"
	"github.com/gooby-lang/gooby/vm"
	"github.com/pkg/profile"
)

const Version string = vm.Version

func main() {
	profileCPUOptionPtr := flag.Bool("profile-cpu", false, "Profile cpu usage")
	profileMemOptionPtr := flag.Bool("profile-mem", false, "Profile memory allocation")
	versionOptionPtr := flag.Bool("v", false, "Show current Gooby version")
	interactiveOptionPtr := flag.Bool("i", false, "Run interactive gooby")
	issueOptionPtr := flag.Bool("e", false, "Generate reporting format")

	flag.Parse()

	if *interactiveOptionPtr {
		igb.StartIgb(Version)
		os.Exit(0)
	}

	if *profileCPUOptionPtr {
		defer profile.Start().Stop()
	}

	if *profileMemOptionPtr {
		defer profile.Start(profile.MemProfile).Stop()
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
		filePath := flag.Arg(1)
		fileInfo, err := os.Stat(filePath)
		reportErrorAndExit(err)

		dir := extractDirFromFilePath(filePath, fileInfo)
		v, err := vm.New(dir, args)

		if fileInfo.Mode().IsDir() {
			fileInfos, err := ioutil.ReadDir(filePath)
			reportErrorAndExit(err)

			for _, fileInfo := range fileInfos {
				fp := filepath.Join(dir, fileInfo.Name())
				reportErrorAndExit(err)

				err := runSpecFile(v, fp)
				reportErrorAndExit(err)
			}
		} else {
			err := runSpecFile(v, filePath)
			reportErrorAndExit(err)
		}

		instructionSets, err := compiler.CompileToInstructions("Spec.run", parser.NormalMode)
		v.ExecInstructions(instructionSets, filePath)
		return
	default:
		fp = flag.Arg(0)

		if !strings.Contains(fp, ".") {
			flag.Usage()
			os.Exit(0)
		}
	}

	// Execute files normally
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

func extractDirFromFilePath(filePath string, fileInfo os.FileInfo) string {
	if fileInfo.Mode().IsDir() {
		dir, err := filepath.Abs(filePath)
		reportErrorAndExit(err)
		return dir
	}

	filePath, err := filepath.Abs(filePath)
	reportErrorAndExit(err)
	dir, _, _ := extractFileInfo(filePath)
	return dir
}

func readFile(filepath string) (file []byte) {
	file, err := ioutil.ReadFile(filepath)
	reportErrorAndExit(err)
	return
}

func runSpecFile(v *vm.VM, fp string) (err error) {
	file := readFile(fp)
	instructionSets, err := compiler.CompileToInstructions(string(file), parser.NormalMode)

	if err != nil {
		return
	}

	v.ExecInstructions(instructionSets, fp)
	return
}

func reportErrorAndExit(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
