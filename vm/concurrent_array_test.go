package vm

import (
	"testing"
)

func TestConcurrentArrayClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		require 'concurrent/array'
		Concurrent::Array.class.name`, "Class"},
		{`
		require 'concurrent/array'
		Concurrent::Array.superclass.name`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayIndex(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/array'
		Concurrent::Array.new([])[1]
		`, nil},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 3])[100]
		`, nil},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 10, 5])[2]
		`, 10},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "a", 10, 5])[1]
		`, "a"},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "a", 10, "b"])[-2]
		`, 10},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, "a", 10, 5])
		a[0]
		`, 1},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, "a", 10, 5])
		a[2] = a[1]
		a[2]
		`, "a"},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, "a", 10, 5])
		a[-2] = a[1]
		a[-2]
		`, "a"},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, "a", 10, 5])
		a[-4] = a[1]
		a[-4]
		`, "a"},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([])
		a[10] = 100
		a[10]
		`, 100},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([])
		a[10] = 100
		a[0]
		`, nil},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2 ,3 ,5 , 10])
		a[0] = a[1] + a[2] + a[3] * a[4]
		a[0]
		`, 55},
		{`
		require 'concurrent/array'
		code = []
		code[2] = 'Continue'
		code[3] = 'Switching Protocols'
		code[5] = 'OK'
		code.to_s
		`, `[nil, nil, "Continue", "Switching Protocols", nil, "OK"]`},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayIndexFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "a", 10, "b"])[-5]
		`, "ArgumentError: Index value -5 too small for array. minimum: -4", 1},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "a", 10, "b"], 1)[-5]
		`, "ArgumentError: Expect 1 or less argument(s). got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayComparisonOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) == 123`, false},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) == "123"`, false},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) == "124"`, false},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) == (1..3)`, false},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) == { a: 1, b: 2 }`, false},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) == Concurrent::Array.new([1, "String", true, 2..5])`, true},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) == Concurrent::Array.new([1, "String", false, 2..5])`, false},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) == Concurrent::Array.new(["String", 1, false, 2..5])`, false}, // Array has order issue
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, { a: 1, b: 2 }, "Goby" ]) == Concurrent::Array.new([1, { a: 1, b: 2 }, "Goby"])`, true},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, { a: 1, b: 2 }, "Goby" ]) == Concurrent::Array.new([1, { b: 2, a: 1 }, "Goby"])`, true},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, { a: 1, b: 2 }, "Goby" ]) == Concurrent::Array.new([1, { a: 1, b: 2, c: 3 }, "Goby"])`, false}, // Array of hash has no order issue
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, { a: 1, b: 2 }, "Goby" ]) == Concurrent::Array.new([1, { a: 2, b: 2, a: 1 }, "Goby"])`, true}, // Array of hash key will be overwritten if duplicated
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) == Integer`, false},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) != 123`, true},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) != "123"`, true},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) != "124"`, true},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) != (1..3)`, true},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) != { a: 1, b: 2 }`, true},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) != Concurrent::Array.new([1, "String", true, 2..5])`, false},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) != Concurrent::Array.new([1, "String", false, 2..5])`, true},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) != Concurrent::Array.new(["String", 1, false, 2..5])`, true}, // Array has order issue
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, { a: 1, b: 2 }, "Goby" ]) != Concurrent::Array.new([1, { a: 1, b: 2 }, "Goby"])`, false},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, { a: 1, b: 2 }, "Goby" ]) != Concurrent::Array.new([1, { b: 2, a: 1 }, "Goby"])`, false},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, { a: 1, b: 2 }, "Goby" ]) != Concurrent::Array.new([1, { a: 1, b: 2, c: 3 }, "Goby"])`, true}, // Array of hash has no order issue
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, { a: 1, b: 2 }, "Goby" ]) != Concurrent::Array.new([1, { a: 2, b: 2, a: 1 }, "Goby"])`, false}, // Array of hash key will be overwritten if duplicated
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "String", true, 2..5]) != Integer`, true},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
		vm.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayAnyMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 3]).any? do |e|
			e == 2
		end
		`, true},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 3]).any? do |e|
			e
		end
		`, true},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 3]).any? do |e|
			e == 5
		end
		`, false},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 3]).any? do |e|
			nil
		end
		`, false},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([]).any? do |e|
			true
		end
		`, false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayAnyMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		Concurrent::Array.new([]).any?`, "InternalError: Can't yield without a block", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayAtMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/array'
		Concurrent::Array.new([]).at(1)
		`, nil},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 10, 5]).at(2)
		`, 10},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "a", 10, 5]).at(1)
		`, "a"},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "a", 10, 5]).at(4)
		`, nil},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "a", 10, 5]).at(-2)
		`, 10},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, "a", 10, 5])
		a.at(0)
		`, 1},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, "a", 10, 5])
		a[2] = a.at(1)
		a[2]
		`, "a"},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 3, 5, 10])
		a[0] = a.at(1) + a.at(2) + a.at(3) * a.at(4)
		a.at(0)
		`, 55},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayAtMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 3]).at`, "ArgumentError: Expect 1 argument(s). got: 0", 1},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 3]).at(2, 3)`, "ArgumentError: Expect 1 argument(s). got: 2", 1},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 3]).at(true)`, "TypeError: Expect argument to be Integer. got: Boolean", 1},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 3]).at(1..3)`, "TypeError: Expect argument to be Integer. got: Range", 1},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "a", 10, 5]).at(-5)`, "ArgumentError: Index value -5 too small for array. minimum: -4", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayClearMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 3])
		a.clear
		`, []interface{}{}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([])
		a.clear
		`, []interface{}{}},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		verifyConcurrentArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
		vm.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayClearMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		Concurrent::Array.new(['Maxwell', 'Alexius']).clear(123)`, "ArgumentError: Expect 0 argument(s). got: 1", 1},
		{`
		require 'concurrent/array'
		Concurrent::Array.new(['Taipei', 101]).clear(1, 0, 1)`, "ArgumentError: Expect 0 argument(s). got: 3", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayConcatMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.concat([3], [4])
		`, []interface{}{1, 2, 3, 4}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([])
		a.concat([1], [2], ["a", "b"], [3], [4])
		`, []interface{}{1, 2, "a", "b", 3, 4}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.concat()
		`, []interface{}{1, 2}},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		verifyConcurrentArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
		vm.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayConcatMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.concat(3)
		`, "TypeError: Expect argument to be Array. got: Integer", 1},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([])
		a.concat("a")
		`, "TypeError: Expect argument to be Array. got: String", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayCountMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.count
		`, 2},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.count(1)
		`, 1},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new(["a", "bb", "c", "db", "bb", 2])
		a.count("bb")
		`, 2},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([true, true, true, false, true])
		a.count(true)
		`, 4},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([])
		a.count(true)
		`, 0},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 3, 4, 5, 6, 7, 8])
		a.count do |i|
			i > 3
		end
		`, 5},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new(["a", "bb", "c", "db", "bb"])
		a.count do |i|
			i.size > 1
		end
		`, 3},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([]).count do |i|
			i.size > 1
		end
		`, 0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayCountMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.count(3, 3)
		`, "ArgumentError: Expect 1 or less argument(s). got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayDeleteAtMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/array'
		Concurrent::Array.new([]).delete_at(1)
		`, nil},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 10, 5]).delete_at(2)
		`, 10},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "a", 10, 5]).delete_at(1)
		`, "a"},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "a", 10, 5]).delete_at(4)
		`, nil},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "a", 10, 5]).delete_at(-2)
		`, 10},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, "a", 10, 5]).delete_at(-5)
		`, nil},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}

	testsArray := []struct {
		input    string
		expected []interface{}
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 10, 5])
		a.delete_at(2)
		a

		`, []interface{}{1, 2, 5}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, "a", 10, 5])
		a.delete_at(4)
		a
		`, []interface{}{1, "a", 10, 5}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, "a", 10, 5])
		a.delete_at(-2)
		a
		`, []interface{}{1, "a", 5}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, "a", 10, 5])
		a.delete_at(-5)
		a
		`, []interface{}{1, "a", 10, 5}},
	}

	for i, tt := range testsArray {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		verifyConcurrentArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
		vm.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayDeleteAtMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 3]).delete_at`, "ArgumentError: Expect 1 argument(s). got: 0", 1},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 3]).delete_at(2, 3)`, "ArgumentError: Expect 1 argument(s). got: 2", 1},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 3]).delete_at(true)`, "TypeError: Expect argument to be Integer. got: Boolean", 1},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 3]).delete_at(1..3)`, "TypeError: Expect argument to be Integer. got: Range", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayEachMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
		require 'concurrent/array'
		sum = 0
		Concurrent::Array.new([1, 2, 3, 4, 5]).each do |i|
			sum = sum + i
		end
		sum
		`, 15},
		{`
		require 'concurrent/array'
		sum = 0
		Concurrent::Array.new([]).each do |i|
			sum += i
		end
		sum
		`, 0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayEachMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		Concurrent::Array.new(['M', 'A', 'X', 'W', 'E', 'L', 'L']).each`, "InternalError: Can't yield without a block", 1},
		{`
		require 'concurrent/array'
		Concurrent::Array.new(['T', 'A', 'I', 'P', 'E', 'I']).each(101) do |char|
			puts char
		end
		`, "ArgumentError: Expect 0 argument(s). got: 1", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayEachIndexMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
		require 'concurrent/array'
		sum = 0
		Concurrent::Array.new([2, 3, 40, 5, 22]).each_index do |i|
			sum = sum + i
		end
		sum
		`, 10},
		{`
		require 'concurrent/array'
		sum = 0
		Concurrent::Array.new([]).each_index do |i|
			sum += i
		end
		sum
		`, 0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayEachIndexMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		Concurrent::Array.new(['M', 'A', 'X', 'W', 'E', 'L', 'L']).each_index`, "InternalError: Can't yield without a block", 1},
		{`
		require 'concurrent/array'
		Concurrent::Array.new(['T', 'A', 'I', 'P', 'E', 'I']).each_index(101) do |char|
			puts char
		end
		`, "ArgumentError: Expect 0 argument(s). got: 1", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayEmptyMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 3]).empty?
		`, false},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([nil]).empty?
		`, false},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([]).empty?
		`, true},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([[]])
		a.empty?
		`, false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayEmptyMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 3]).empty?(123)`, "ArgumentError: Expect 0 argument(s). got: 1", 1},
		{`
		require 'concurrent/array'
		Concurrent::Array.new(['T', 'A', 'I', 'P', 'E', 'I']).empty?(1, 0, 1)`, "ArgumentError: Expect 0 argument(s). got: 3", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayFirstMethod(t *testing.T) {
	testsInt := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.first
		`, 1},
	}

	for i, tt := range testsInt {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
	}

	testsArray := []struct {
		input    string
		expected []interface{}
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([3, 4, 5, 1, 6])
		a.first(2)
		`, []interface{}{3, 4}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new(["a", "b", "d", "q"])
		a.first(2)
		`, []interface{}{"a", "b"}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new(["M", "A", "X", "W", "E", "L", "L"])
		a.first(7)`, []interface{}{"M", "A", "X", "W", "E", "L", "L"}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new(["M", "A", "X", "W", "E", "L", "L"])
		a.first(11)`, []interface{}{"M", "A", "X", "W", "E", "L", "L"}},
	}

	for i, tt := range testsArray {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyConcurrentArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayFirstMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.first("a")
		`, "TypeError: Expect argument to be Integer. got: String", 1},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.first(1, 2, 3)
		`, "ArgumentError: Expect 1 or less argument(s). got: 3", 1},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.first(-1)
		`, "ArgumentError: Expect argument to be positive value. got: -1", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayFlattenMethod(t *testing.T) {
	testsArray := []struct {
		input    string
		expected []interface{}
	}{
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2]).flatten
		`, []interface{}{1, 2}},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, [3, 4, 5]]).flatten
		`, []interface{}{1, 2, 3, 4, 5}},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([[1, 2], [3, 4], [5, 6]]).flatten
		`, []interface{}{1, 2, 3, 4, 5, 6}},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([[[1, 2], [[[3, 4]], [5, 6]]]]).flatten
		`, []interface{}{1, 2, 3, 4, 5, 6}},
	}

	for i, tt := range testsArray {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		verifyConcurrentArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
		vm.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayFlattenMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.flatten(1)
		`, "ArgumentError: Expect 0 argument(s). got: 1", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayJoinMethod(t *testing.T) {
	testsInt := []struct {
		input    string
		expected string
	}{
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2]).join
		`, "12"},
		{`
		require 'concurrent/array'
		Concurrent::Array.new(["1", 2]).join
		`, "12"},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2]).join(",")
		`, "1,2"},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, [3, 4]]).join(",")
		`, "1,2,3,4"},
	}

	for i, tt := range testsInt {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayJoinMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.join(",", "-")
		`, "ArgumentError: Expect 0 to 1 argument(s). got: 2", 1},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.join(1)
		`, "TypeError: Expect argument to be String. got: Integer", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayLastMethod(t *testing.T) {
	testsArray := []struct {
		input    string
		expected []interface{}
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([3, 4, 5, 1, 6])
		a.last(3)
		`, []interface{}{5, 1, 6}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new(["a", "b", "d", "q"])
		a.last(2)
		`, []interface{}{"d", "q"}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new(["M", "A", "X", "W", "E", "L", "L"])
		a.last(7)
		`, []interface{}{"M", "A", "X", "W", "E", "L", "L"}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new(["M", "A", "X", "W", "E", "L", "L"])
		a.last(10)
		`, []interface{}{"M", "A", "X", "W", "E", "L", "L"}},
	}

	for i, tt := range testsArray {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		verifyConcurrentArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
		vm.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayLastMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.last("l")
		`, "TypeError: Expect argument to be Integer. got: String", 1},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.last(1, 2, 3)
		`, "ArgumentError: Expect 1 or less argument(s). got: 3", 1},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.last(-1)
		`, "ArgumentError: Expect argument to be positive value. got: -1", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayLengthMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 3]).length
		`, 3},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([nil]).length
		`, 1},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([]).length
		`, 0},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([-10, "123", [1,2,3], 1, 2, 3])
		a.length
		`, 6},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayLengthMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2, 3]).length(10)`, "ArgumentError: Expect 0 argument(s). got: 1", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayMapMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 7])
		a.map do |i|
			i + 3
		end
		`, []interface{}{4, 5, 10}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([true, false, true, false, true ])
		a.map do |i|
			!i
		end
		`, []interface{}{false, true, false, true, false}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new(["1", "sss", "qwe"])
		a.map do |i|
			i + "1"
		end
		`, []interface{}{"11", "sss1", "qwe1"}},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([]).map do |i|
		end
		`, []interface{}{}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyConcurrentArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayPlusMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		// Make sure the result is an entirely new array.
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		b = [3, 4]
		c = a + b
		a[0] = -1
		b[0] = -1
		c
		`, []interface{}{1, 2, 3, 4}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([])
		b = []
		a + b
		`, []interface{}{}},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		verifyConcurrentArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
		vm.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayPlusMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2]) + true`, "TypeError: Expect argument to be Array. got: Boolean", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayPopMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 3]).pop
		a
		`, 3},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 3])
		a.pop
		a.length
		`, 2},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([]).pop
		`, nil},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayPushMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 3])
		a.push("test")
		a[3]
		`, "test"},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 3])
		a.push("test")
		a.length
		`, 4},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([])
		a.push(nil)
		a[0]
		`, nil},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([])
		a.push("foo")
		a.push(1)
		a.push(234)
		a[0]
		`, "foo"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayReduceMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 7])
		a.reduce do |sum, n|
			sum + n
		end
		`, 10},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 7])
		a.reduce(10) do |sum, n|
			sum + n
		end
		`, 20},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new(["This ", "is a ", "test!"])
		a.reduce do |prev, s|
			prev + s
		end
		`, "This is a test!"},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new(["this ", "is a ", "test!"])
		a.reduce("Yes, ") do |prev, s|
			prev + s
		end
		`, "Yes, this is a test!"},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([]).reduce("foo") do |i|
			true
		end
		`, "foo"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayReduceMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.reduce(1)
		`, "InternalError: Can't yield without a block", 1},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.reduce(1, 2) do |prev, n|
			prev + n
		end
		`, "ArgumentError: Expect 1 or less argument(s). got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayReverseMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 3])
		a.reverse
		`, []interface{}{3, 2, 1}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([])
		a.reverse
		`, []interface{}{}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyConcurrentArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayReverseEachMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		require 'concurrent/array'
		str = ""
		Concurrent::Array.new(["a", "b", "c"]).reverse_each do |char|
			str += char
		end
		str
		`, "cba"},
		{`
		require 'concurrent/array'
		str = ""
		Concurrent::Array.new([]).reverse_each do |i|
			str += char
		end
		str
		`, ""},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayReverseEachMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		Concurrent::Array.new(['M', 'A']).reverse_each`, "InternalError: Can't yield without a block", 1},
		{`
		require 'concurrent/array'
		Concurrent::Array.new(['T', 'A']).reverse_each(101) do |char|
			puts char
		end
		`, "ArgumentError: Expect 0 argument(s). got: 1", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayRotateMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.rotate
		`, []interface{}{2, 1}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 3, 4])
		a.rotate(2)
		`, []interface{}{3, 4, 1, 2}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyConcurrentArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayRotateMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.rotate("a")
		`, "TypeError: Expect argument to be Integer. got: String", 1},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.rotate(1, 2, 3)`, "ArgumentError: Expect 1 or less argument(s). got: 3", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArraySelectMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 3, 4, 5])
		a.select do |i|
			i > 3
		end
		`, []interface{}{4, 5}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([true, false, true, false, true ])
		a.select do |i|
			i
		end
		`, []interface{}{true, true, true}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new(["test", "not2", "3", "test", "5"])
		a.select do |i|
			i == "test"
		end
		`, []interface{}{"test", "test"}},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([]).select do |i|
			true
		end
		`, []interface{}{}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyConcurrentArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayShiftMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 3]).shift
		a
		`, 1},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 3])
		a.pop
		a.length
		`, 2},
		{`
		require 'concurrent/array'
		Concurrent::Array.new([]).shift
		`, nil},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayShiftMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2])
		a.shift(3, 3, 4, 5)
		`, "ArgumentError: Expect 0 argument(s). got: 4", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayStarMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 3])
		a * 2
		`, []interface{}{1, 2, 3, 1, 2, 3}},
		// Make sure the result is an entirely new array.
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 3])
		(a * 2)[0] = -1
		a
		`, []interface{}{1, 2, 3}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 3])
		a * 0
		`, []interface{}{}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([])
		a * 2
		`, []interface{}{}},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		verifyConcurrentArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
		vm.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayStarMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		Concurrent::Array.new([1, 2]) * nil`, "TypeError: Expect argument to be Integer. got: Null", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayUnshiftMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 3])
		a.unshift(0)
		a[0]
		`, 0},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([1, 2, 3])
		a.unshift(0)
		a.length
		`, 4},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([])
		a.unshift(nil)
		a[0]
		`, nil},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([])
		a.unshift("foo")
		a.unshift(1, 2)
		a[0]
		`, 1},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([])
		a.unshift("foo")
		a.unshift(1, 2)
		a[1]
		`, 2},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([])
		a.unshift("foo")
		a.unshift(1, 2)
		a[2]
		`, "foo"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayValuesAtMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new(["a", "b", "c"])
		a.values_at(1)
		`, []interface{}{"b"}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new(["a", "b", "c"])
		a.values_at(-1, 3)
		`, []interface{}{"c", nil}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new(["a", "b", "c"])
		a.values_at()
		`, []interface{}{}},
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new([])
		a.values_at(1, -1)
		`, []interface{}{nil, nil}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyConcurrentArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentArrayValuesAtMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/array'
		a = Concurrent::Array.new(["a", "b", "c"])
		a.values_at("-")
		`, "TypeError: Expect argument to be Integer. got: String", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}
