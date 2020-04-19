package vm

import (
	"testing"
)

func TestRangeEnumeratorEnumerationWithElements(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		iterated_values = []
	
		enumerator = RangeEnumerator.new((1..3))
	
		while enumerator.has_next? do
			iterated_values.push(enumerator.next)
		end
	
		iterated_values
		`,
			[]interface{}{1, 2, 3},
		},
		{`
		iterated_values = []
	
		enumerator = RangeEnumerator.new((3..1))
	
		while enumerator.has_next? do
			iterated_values.push(enumerator.next)
		end
	
		iterated_values
		`,
			[]interface{}{3, 2, 1},
		},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
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

func TestRangeEnumeratorRaiseErrorWhenNoElementsOnNext(t *testing.T) {
	tests := []errorTestCase{
		{`
			enumerator = RangeEnumerator.new((1..1))
			enumerator.next
			enumerator.next
			`,
			"StopIteration: \"No more elements!\"",
			1, 1,
		},
		{`
			enumerator = RangeEnumerator.new((1..0))
			enumerator.next
			enumerator.next
			enumerator.next
			`,
			"StopIteration: \"No more elements!\"",
			1, 1,
		},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, tt.expectedSP)

	}

}
