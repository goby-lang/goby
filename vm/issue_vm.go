package vm

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
)

// InitIssueReportVM initializes a vm in test mode for issue reporting
func InitIssueReportVM(dir string, args []string) (*VM, error) {
	v, err := New(dir, args)
	v.mode = TestMode

	return v, err
}

// PrintError prints an error report string given a vm which evaluated to and Error object
func PrintError(v *VM) {
	eval := v.mainThread.stack.top().Target
	err, ok := eval.(*Error)
	if !ok {
		fmt.Println("No error detected")
	}
	fmt.Printf("# %s\n", err.Type)
	fmt.Println(err.Message)

	fmt.Printf("### Goby version\n%s\n", Version)
	fmt.Printf("### GOBY_ROOT\n%s\n", os.Getenv("GOBY_ROOT"))
	fmt.Printf("### Go version\n%s\n", runtime.Version())
	fmt.Printf("### GOROOT\n%s\n", os.Getenv("GOROOT"))
	fmt.Printf("### GOPATH\n%s\n", os.Getenv("GOPATH"))
	fmt.Printf("### Operating system\n%s\n", runtime.GOOS)

	t := v.mainThread
	cf := t.callFrameStack.top()

	// If program counter is 0 means we need to trace back to previous call frame
	if cf.pc == 0 {
		t.callFrameStack.pop()
		cf = t.callFrameStack.top()
	}

	file := cf.instructionSet.filename
	line := cf.instructionSet.instructions[cf.pc-1].Line

	// Print lines in file surrounding error in markdown code block
	f, osErr := os.Open(string(file))
	if osErr != nil {
		fmt.Println("Could not open problem file")
	}

	scanner := bufio.NewScanner(f)

	scanner.Split(bufio.ScanLines)

	currLine := 0
	// Skip lines until at least 20 lines from error
	for ; currLine < line-20; currLine++ {
		scanner.Scan()
	}
	fmt.Println("``` ruby")
	// Print until 20 lines past error
	for ; currLine < line+20 && scanner.Scan(); currLine++ {
		fmt.Printf("%s\n", scanner.Text())
	}
	fmt.Println("```\n")
}
