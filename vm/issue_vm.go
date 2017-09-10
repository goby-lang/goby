package vm

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
)

func InitIssueReportVM(dir string, args []string) (*VM, error) {
	v, err := New(dir, args)
	v.mode = TestMode

	return v, err
}

func PrintError(v *VM) {
	eval := v.mainThread.stack.top().Target
	err, ok := eval.(*Error)
	if !ok {
		fmt.Println("No error detected")
	}
	fmt.Printf("# Help! %s\n", err.Type)
	fmt.Println(err.Message)

	fmt.Printf("### Goby version\n%s\n", Version)
	fmt.Printf("### Go version\n%s\n", runtime.Version())
	fmt.Printf("### Operating system\n%s\n", runtime.GOOS)
	fmt.Printf("### GOBY_ROOT\n%s\n", os.Getenv("GOBY_ROOT"))

	t := v.mainThread
	cf := t.callFrameStack.top()

	// If program counter is 0 means we need to trace back to previous call frame
	if cf.pc == 0 {
		t.callFrameStack.pop()
		cf = t.callFrameStack.top()
	}

	file := cf.instructionSet.filename
	line := cf.instructionSet.instructions[cf.pc-1].Line

	f, osErr := os.Open(string(file))
	if osErr != nil {
		fmt.Println("Could not open problem file")
	}

	scanner := bufio.NewScanner(f)

	scanner.Split(bufio.ScanLines)

	currLine := 0
	for ; currLine < line-20; currLine++ {
		scanner.Scan()
	}

	fmt.Println("``` ruby")
	for ; currLine < line+20 && scanner.Scan(); currLine++ {
		fmt.Printf("%s\n", scanner.Text())
	}
	fmt.Println("```\n")
}
