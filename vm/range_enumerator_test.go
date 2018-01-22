package vm

import (
	"testing"
)

func TestRangeEnumeratorEnumerationWithElements(t *testing.T) {
	input := `
	iterated_values = []

	enumerator = RangeEnumerator.new((1..3))

	while enumerator.has_next? do
		iterated_values.push(enumerator.next)
	end

	iterated_values
	`

	expected := []interface{}{1, 2, 3}

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	verifyArrayObject(t, i, evaluated, expected)
	v.checkCFP(t, i, 0)
	v.checkSP(t, i, 1)
}

func TestRangeEnumeratorEnumerationWithOneElement(t *testing.T) {
	input := `
	iterated_values = []

	enumerator = RangeEnumerator.new((1..1))

	while enumerator.has_next? do
		iterated_values.push(enumerator.next)
	end

	iterated_values
	`

	expected := []interface{}{1}

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	verifyArrayObject(t, i, evaluated, expected)
	v.checkCFP(t, i, 0)
	v.checkSP(t, i, 1)
}

func TestRangeEnumeratorHasNextForReverseRange(t *testing.T) {
	input := `
	iterated_values = []

	enumerator = RangeEnumerator.new((3..1))

	enumerator.has_next?
	`

	expected := false

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	verifyExpected(t, i, evaluated, expected)
	v.checkCFP(t, i, 0)
	v.checkSP(t, i, 1)
}

func TestRangeEnumeratorRaiseErrorWhenNoElementsOnNext(t *testing.T) {
	testCase := errorTestCase{`
	enumerator = RangeEnumerator.new((1..1))
	enumerator.next
	enumerator.next
	`, "StopIteration: 'No more elements!'", 2}

	v := initTestVM()
	evaluated := v.testEval(t, testCase.input, getFilename())
	checkErrorMsg(t, i, evaluated, testCase.expected)
	v.checkCFP(t, i, testCase.expectedCFP)
	v.checkSP(t, i, 2)
}
