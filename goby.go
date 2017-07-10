package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/goby-lang/goby/compiler"
	"github.com/goby-lang/goby/igb"
	"github.com/goby-lang/goby/vm"
	"github.com/pkg/profile"
)

// Version stores current Goby version
const Version string = "0.0.8"

func main() {
	profileOptionPtr := flag.Bool("p", false, "Profile program execution")
	versionOptionPtr := flag.Bool("v", false, "Show current Goby version")
	interactiveOptionPtr := flag.Bool("i", false, "Run interactive goby")

	flag.Parse()

	if *interactiveOptionPtr {
		igb.Start(Version)
		os.Exit(0)
	}

	if *profileOptionPtr {
		defer profile.Start().Stop()
	}

	if *versionOptionPtr {
		fmt.Println(Version)
		os.Exit(0)
	}

	filepath := flag.Arg(0)

	if filepath == "" {
		flag.Usage()
		os.Exit(0)
	}

	args := flag.Args()[1:]

	dir, _, fileExt := extractFileInfo(filepath)
	file, ok := readFile(filepath)

	if !ok {
		return
	}

	switch fileExt {
	case "gb", "rb":
		instructionSets, err := compiler.CompileToInstructions(string(file))

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		v := vm.New(dir, args)
		v.ExecInstructions(instructionSets, filepath)
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
