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

func init() {
	_, err := os.Stat("./goby")
	if err == nil {
		err := exec.Command("rm", "./goby").Run()
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Building Goby binary")
	err = exec.Command("make", "build").Run()
	if err != nil {
		panic(err)
	}
}

func execGoby(t *testing.T, args ...string) (in io.WriteCloser, out io.ReadCloser, e io.ReadCloser) {
	t.Helper()

	cmd := exec.Command("./goby", args...)

	in, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Error getting stdin\n%s", err.Error())
	}

	out, err = cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Error getting stdout\n%s", err.Error())
	}

	e, err = cmd.StderrPipe()
	if err != nil {
		t.Fatalf("Error getting stderr\n%s", err.Error())
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

	_, out, _ := execGoby(t, "-e", "samples/error-report.gb")

	byt, err := ioutil.ReadAll(out)
	if err != nil {
		t.Fatalf("Couldn't read from pipe: %s", err.Error())
	}

	if !strings.Contains(string(byt), partialReport()) {
		t.Fatalf("Interpreter -e output incorrect:\nExpected \n`%s` in string `\n%s`", partialReport(), string(byt))
	}
}

func TestArgI(t *testing.T) {

	in, out, _ := execGoby(t, "-i")

	fmt.Fprintln(in, `puts "hello world"`)
	fmt.Fprintln(in, `exit`)

	expectedOutput := "hello world\n"

	byt, err := ioutil.ReadAll(out)
	if err != nil {
		t.Fatalf("Couldn't read from pipe: %s", err.Error())
	}

	if !strings.HasSuffix(string(byt), expectedOutput) {
		t.Fatalf("Interpreter output incorrect. Expect '%s' to contain '%s'", string(byt), expectedOutput)
	}
}

func TestArgV(t *testing.T) {

	_, out, _ := execGoby(t, "-v")

	byt, err := ioutil.ReadAll(out)
	if err != nil {
		t.Fatalf("Couldn't read from pipe: %s", err.Error())
	}

	if !strings.Contains(string(byt), vm.Version) {
		t.Fatalf("Interpreter -v output incorrect:\nExpected '%s' in string '%s'.", vm.Version, string(byt))
	}
}

func TestArgProfileCPU(t *testing.T) {

	_, out, _ := execGoby(t, "-profile-cpu", "samples/one_thousand_threads.gb")

	byt, err := ioutil.ReadAll(out)
	if err != nil {
		t.Fatalf("Couldn't read from pipe: %s", err.Error())
	}

	if string(byt) != "500500\n" {
		t.Fatalf("Test failed, excpected 500500, got %s", string(byt))
	}
}

func TestArgProfileMem(t *testing.T) {

	_, out, _ := execGoby(t, "-profile-mem", "samples/one_thousand_threads.gb")

	byt, err := ioutil.ReadAll(out)
	if err != nil {
		t.Fatalf("Couldn't read from pipe: %s", err.Error())
	}

	if string(byt) != "500500\n" {
		t.Fatalf("Test failed, excpected 500500, got %s", string(byt))
	}
}

func TestExecFileWithError(t *testing.T) {
	expectedError := "NoMethodError: Undefined Method 'foo' for "

	_, _, stderr := execGoby(t, "test_fixtures/file_with_error.gb")

	output, _ := ioutil.ReadAll(stderr)

	if !strings.Contains(string(output), expectedError) {
		t.Fatalf("Expect to see error: '%s'. But got: '%s' instead", expectedError, string(output))
	}
}

func TestTestCommand(t *testing.T) {
	// Folder name with slash
	_, out, _ := execGoby(t, "test", "test_fixtures/test_command_test/")

	byt, err := ioutil.ReadAll(out)
	if err != nil {
		t.Fatalf("Couldn't read from pipe: %s", err.Error())
	}

	if !strings.Contains(string(byt), "Spec test 2") {
		t.Fatalf("Test files by giving folder name with slash failed, got: %s", string(byt))
	}

	// Folder name
	_, out, _ = execGoby(t, "test", "test_fixtures/test_command_test")

	byt, err = ioutil.ReadAll(out)
	if err != nil {
		t.Fatalf("Couldn't read from pipe: %s", err.Error())
	}

	if !strings.Contains(string(byt), "Spec test 2") {
		t.Fatalf("Test files by giving folder name failed, got: %s", string(byt))
	}

	// File name
	_, out, _ = execGoby(t, "test", "test_fixtures/test_command_test/test_spec.gb")

	byt, err = ioutil.ReadAll(out)
	if err != nil {
		t.Fatalf("Couldn't read from pipe: %s", err.Error())
	}

	if !strings.Contains(string(byt), "Spec") {
		t.Fatalf("Test files by giving file name failed, got: %s", string(byt))
	}
}
