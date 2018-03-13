package vm

import (
	"testing"
)

func TestSpecSuccessWithoutExit(t *testing.T) {
	input := `
require "spec"

Spec.describe Spec do
  it "fails and exit with code 1" do
	expect(1).to eq(1)
  end
end

Spec.test

10
`
	v := initTestVM()
	result := v.testEval(t, input, getFilename())
	verifyExpected(t, 0, result, 10)
}

func TestSpecFailAndExit(t *testing.T) {
	input := `
require "spec"

Spec.describe Spec do
  it "fails and exit with code 1" do
	expect(1).to eq(2)
  end
end

Spec.test

10
`
	v := initTestVM()
	result := v.testEval(t, input, getFilename())
	_, ok := result.(*IntegerObject)

	if ok {
		t.Fatal("Program should exit early because the spec failed")
	}
}
