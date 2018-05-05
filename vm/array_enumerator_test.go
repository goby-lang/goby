package vm

import (
	"testing"
)

func TestArrayEnumeratorEnumerationWithoutElements(t *testing.T) {
	input := `
	enumerator = ArrayEnumerator.new([])
	enumerator.has_next?
	`

	expected := false

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	VerifyExpected(t, i, evaluated, expected)
	v.checkCFP(t, i, 0)
	v.checkSP(t, i, 1)
}

func TestArrayEnumeratorEnumerationWithElements(t *testing.T) {
	input := `
	iterated_values = []

	enumerator = ArrayEnumerator.new([1, 2, 4])

	while enumerator.has_next? do
		iterated_values.push(enumerator.next)
	end

	iterated_values
	`

	expected := []interface{}{1, 2, 4}

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	verifyArrayObject(t, i, evaluated, expected)
	v.checkCFP(t, i, 0)
	v.checkSP(t, i, 1)
}

func TestArrayEnumeratorRaiseErrorWhenNoElementsOnNext(t *testing.T) {
	testCase := errorTestCase{`
	ArrayEnumerator.new([]).next
	`, "StopIteration: 'No more elements!'", 2}

	v := initTestVM()
	evaluated := v.testEval(t, testCase.input, getFilename())
	checkErrorMsg(t, i, evaluated, testCase.expected)
	v.checkCFP(t, i, testCase.expectedCFP)
	v.checkSP(t, i, 2)
}
