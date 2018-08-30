package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/goby-lang/goby/vm"
)

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	//if ret == 0 {
	//	teardown()
	//}
	os.Exit(ret)
}

func setup() {
	_, err := os.Stat("./goby")
	if err != nil {
		fmt.Println("Goby binary not found, building")

		cmd := exec.Command("go", "build", ".")
		err = cmd.Run()
		if err != nil {
			fmt.Println("Could not build binary\n", err.Error())
			panic(err)
		}
		fmt.Println("Built. Testing ./goby")
	} else {
		fmt.Println("Using existing Goby binary. Testing ./goby")
	}
}

//func teardown() {
//}

func execGoby(t *testing.T, args ...string) (in io.WriteCloser, out io.ReadCloser) {
	exec.Command("rm ./goby")
	exec.Command("make build")
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
		t.Fatalf("Interpreter -e output incorect:\nExpected \n`%s` in string `\n%s`", partialReport(), string(byt))
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

	if !strings.Contains(string(byt), vm.Version) {
		t.Fatalf("Interpreter -v output incorect:\nExpected '%s' in string '%s'.", vm.Version, string(byt))
	}
}

func TestArgProfileCPU(t *testing.T) {

	_, out := execGoby(t, "-profile-cpu", "samples/one_thousand_threads.gb")

	byt, err := ioutil.ReadAll(out)
	if err != nil {
		t.Fatalf("Couldn't read from pipe: %s", err.Error())
	}

	if string(byt) != "500500\n" {
		t.Fatalf("Test failed, excpected 500500, got %s", string(byt))
	}
}

func TestArgProfileMem(t *testing.T) {

	_, out := execGoby(t, "-profile-mem", "samples/one_thousand_threads.gb")

	byt, err := ioutil.ReadAll(out)
	if err != nil {
		t.Fatalf("Couldn't read from pipe: %s", err.Error())
	}

	if string(byt) != "500500\n" {
		t.Fatalf("Test failed, excpected 500500, got %s", string(byt))
	}
}

func TestTestCommand(t *testing.T) {
	// Folder name with slash
	_, out := execGoby(t, "test", "test_fixtures/test_command_test/")

	byt, err := ioutil.ReadAll(out)
	if err != nil {
		t.Fatalf("Couldn't read from pipe: %s", err.Error())
	}

	if !strings.Contains(string(byt), "Spec test 2") {
		t.Fatalf("Test files by giving folder name with slash failed, got: %s", string(byt))
	}

	// Folder name
	_, out = execGoby(t, "test", "test_fixtures/test_command_test")

	byt, err = ioutil.ReadAll(out)
	if err != nil {
		t.Fatalf("Couldn't read from pipe: %s", err.Error())
	}

	if !strings.Contains(string(byt), "Spec test 2") {
		t.Fatalf("Test files by giving folder name failed, got: %s", string(byt))
	}

	// File name
	_, out = execGoby(t, "test", "test_fixtures/test_command_test/test_spec.gb")

	byt, err = ioutil.ReadAll(out)
	if err != nil {
		t.Fatalf("Couldn't read from pipe: %s", err.Error())
	}

	if !strings.Contains(string(byt), "Spec") {
		t.Fatalf("Test files by giving file name failed, got: %s", string(byt))
	}
}
