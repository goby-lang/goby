package vm

import (
	"testing"
)

func TestLazyEnumeratorEachMethodWithoutEnumeratorBlock(t *testing.T) {
	input := `
	enumerator = [1, 2, 3].lazy
	result = []

	enumerator.each do |value|
		result.push(value)
	end

	result
	`

	expected := []interface{}{1, 2, 3}

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	verifyArrayObject(t, i, evaluated, expected)
	v.checkCFP(t, i, 0)
	v.checkSP(t, i, 1)
}

func TestLazyEnumeratorEachMethodWithEnumeratorBlock(t *testing.T) {
	input := `
	enumerator = LazyEnumerator.new(ArrayEnumerator.new([1, 2, 3])) do |value|
		2 * value
	end
	result = []

	enumerator.each do |value|
		result.push(value)
	end

	result
	`

	expected := []interface{}{2, 4, 6}

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	verifyArrayObject(t, i, evaluated, expected)
	v.checkCFP(t, i, 0)
	v.checkSP(t, i, 1)
}

func TestLazyEnumeratorMapMethod(t *testing.T) {
	input := `
	enumerator = [1, 2, 3].lazy.map do |value|
		2 * value
	end
	result = []

	enumerator.each do |value|
		result.push(value)
	end

	result
	`

	expected := []interface{}{2, 4, 6}

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	verifyArrayObject(t, i, evaluated, expected)
	v.checkCFP(t, i, 0)
	v.checkSP(t, i, 1)
}

func TestLazyEnumeratorNextAndHasNextMethods(t *testing.T) {
	input := `
	enumerator = [1, 2, 3].lazy
	result = []

  while enumerator.has_next? do
		result.push(enumerator.next)
  end

	result
	`

	expected := []interface{}{1, 2, 3}

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	verifyArrayObject(t, i, evaluated, expected)
	v.checkCFP(t, i, 0)
	v.checkSP(t, i, 1)
}

func TestLazyEnumeratorFirstMethod(t *testing.T) {
	input := `
	[1, 2, 3].lazy.first(2)
	`

	expected := []interface{}{1, 2}

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	verifyArrayObject(t, i, evaluated, expected)
	v.checkCFP(t, i, 0)
	v.checkSP(t, i, 1)
}

func TestLazyEnumeratorFirstMethodWithZeroValues(t *testing.T) {
	input := `
	iterated_values = []

	result = [1, 2, 3].lazy.map do |n|
		iterated_values.push(n)
	end.first(0)

	[iterated_values, result]
	`

	expected := [][]interface{}{{}, {}}

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	verifyBidimensionalArrayObject(t, i, evaluated, expected)
	v.checkCFP(t, i, 0)
	v.checkSP(t, i, 1)
}

func TestLazyEnumeratorFunctional(t *testing.T) {
	input := `
	accumulated_values = []

	result = [1,2,3].lazy.map do |i|
		accumulated_values.push(i)
		2 * i
	end.map do |i|
		accumulated_values.push(i)
		3 * i
	end.map do |i|
		accumulated_values.push(i)
		4 * i
	end.first(2)

	[accumulated_values, result]
	`

	expected := [][]interface{}{{1, 2, 6, 2, 4, 12}, {24, 48}}

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	verifyBidimensionalArrayObject(t, i, evaluated, expected)
	v.checkCFP(t, i, 0)
	v.checkSP(t, i, 1)
}
