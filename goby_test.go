package main

import (
	"fmt"
	"github.com/goby-lang/goby/vm"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"io"
)

func init() {
	_, err := os.Stat("goby.go")
	if err != nil {
		panic(fmt.Errorf("Goby not bulit"))
	}

	fmt.Println("Testing on prebuilt ./goby")

	//cmd := exec.Command("go", "build", ".")
	//err := cmd.Run()
	//if err != nil {
	//	fmt.Println("could not build binary\n", err.Error())
	//	panic(err)
	//}
}

func execGoby(t *testing.T, args ...string) (in io.WriteCloser, out io.ReadCloser) {
	cmd := exec.Command("./goby", args...)

	var err error
	in, err = cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Error getting stdin\n%s", err.Error())
	}

	out, err = cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Error getting stdout\n%s", err.Error())
	}

	err = cmd.Start()
	if err != nil {
		t.Fatalf("Error running goby\n%s", err.Error())
	}

	return
}

func partialReport() (md string) {

	md += fmt.Sprintf("### Goby version\n%s\n", vm.Version)
	md += fmt.Sprintf("### GOBY_ROOT\n%s\n", os.Getenv("GOBY_ROOT"))
	md += fmt.Sprintf("### Go version\n%s\n", runtime.Version())
	md += fmt.Sprintf("### GOROOT\n%s\n", os.Getenv("GOROOT"))
	md += fmt.Sprintf("### GOPATH\n%s\n", os.Getenv("GOPATH"))
	md += fmt.Sprintf("### Operating system\n%s\n", runtime.GOOS)

	return
}

func TestArgE(t *testing.T) {

	_, out := execGoby(t, "-e", "samples/error-report.gb")

	byt, err := ioutil.ReadAll(out)
	if err != nil {
		t.Fatalf("Couldn't read from pipe: %s", err.Error())
	}

	if !strings.Contains(string(byt), partialReport()) {
		t.Fatalf("Interpreter output incorect")
	}
}

func TestArgI(t *testing.T) {

	in, out := execGoby(t, "-i")

	fmt.Fprintln(in, `puts "hello world"`)
	fmt.Fprintln(in, `exit`)

	byt, err := ioutil.ReadAll(out)
	if err != nil {
		t.Fatalf("Couldn't read from pipe: %s", err.Error())
	}

	if strings.HasSuffix(string(byt), "hello world\nBye") {
		t.Fatalf("Interpreter output incorect")
	}
}

func TestArgV(t *testing.T) {

	_, out := execGoby(t, "-v")

	byt, err := ioutil.ReadAll(out)
	if err != nil {
		t.Fatalf("Couldn't read from pipe: %s", err.Error())
	}

	if string(byt) != vm.Version {
		t.Fatalf("Interpreter output incorect")
	}
}

func TestArgP(t *testing.T) {

	_, out := execGoby(t, "-p", "samples/one_thousand_threads.gb")

	byt, err := ioutil.ReadAll(out)
	if err != nil {
		t.Fatalf("Couldn't read from pipe: %s", err.Error())
	}

	if string(byt) != "500500\n" {
		t.Fatalf("Test failed, excpected 500500, got %s", string(byt))
	}
}
