package vm

import (
	"os"
	"os/exec"
	"testing"
)

func TestSpecSuccessWithExitCode0(t *testing.T) {
	input := `
require "spec"

Spec.describe Spec do
  it "passes" do
	expect(1).to eq(1)
  end
end

Spec.run
`
	if os.Getenv("TEST_SPEC_NOT_EXIT") == "1" {
		v := initTestVM()
		result := v.testEval(t, input, getFilename())
		verifyExpected(t, 0, result, 10)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestSpecSuccessWithExitCode0")
	cmd.Env = append(os.Environ(), "TEST_SPEC_NOT_EXIT=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		t.Fatalf("Spec should exit with status 0.")
	}
}

func TestSpecFailWithExitCode1(t *testing.T) {
	input := `
require "spec"

Spec.describe Spec do
  it "fails and exit with code 1" do
	expect(1).to eq(2)
  end
end

Spec.run
`
	if os.Getenv("TEST_SPEC_EXIT") == "1" {
		v := initTestVM()
		v.testEval(t, input, getFilename())
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestSpecFailWithExitCode1")
	cmd.Env = append(os.Environ(), "TEST_SPEC_EXIT=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Spec fail should exit with status 1.")
}
