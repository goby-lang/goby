package vm

import (
	"testing"
)

func TestRangeClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Range.class.name`, "Class"},
		{`Range.superclass.name`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestRangeComparisonOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`(1..3) == (1..3)`, true},
		{`(1..3) == (1..4)`, false},
		{`(1..3) == 123`, false},
		{`(1..3) == "123"`, false},
		{`(1..3) == "124"`, false},
		{`(1..3) == { a: 1, b: 2 }`, false},
		{`(1..3) == [1, "String", true, 2..5]`, false},
		{`(1..3) == Integer`, false},
		{`(1..3) != (1..3)`, false},
		{`(1..3) != (1..4)`, true},
		{`(1..3) != 123`, true},
		{`(1..3) != "123"`, true},
		{`(1..3) != "124"`, true},
		{`(1..3) != { a: 1, b: 2 }`, true},
		{`(1..3) != [1, "String", true, 2..5]`, true},
		{`(1..3) != Integer`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestRangeBsearchMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		ary = [0, 4, 7, 10, 12]
		(0..4).bsearch do |i|
			ary[i] >= 0
		end
		`, 0},
		{`
		ary = [0, 4, 7, 10, 12]
		(0..4).bsearch do |i|
			ary[i] >= 4
		end
		`, 1},
		{`
		ary = [0, 4, 7, 10, 12]
		(2..4).bsearch do |i|
			ary[i] >= 4
		end
		`, 2},
		{`
		ary = [0, 4, 7, 10, 12]
		(4..4).bsearch do |i|
			ary[i] >= 4
		end
		`, 4},
		{`
		ary = [0, 4, 7, 10, 12]
		(0..4).bsearch do |i|
			ary[i] >= 6
		end
		`, 2},
		{`
		ary = [0, 4, 7, 10, 12]
		(0..4).bsearch do |i|
			ary[i] >= 8
		end
		`, 3},
		{`
		ary = [0, 4, 7, 10, 12]
		(0..2).bsearch do |i|
			ary[i] >= 8
		end
		`, nil},
		{`
		ary = [0, 4, 7, 10, 12]
		(0..4).bsearch do |i|
			ary[i] >= 100
		end
		`, nil},
		{`
		ary = [0, 4, 7, 10, 12]
		(4..0).bsearch do |i|
			ary[i] >= 4
		end
		`, nil},
		{`
		ary = [0, 4, 7, 10, 12]
		(-1..3).bsearch do |i|
			ary[i] >= 4
		end
		`, nil},
		{`
		ary = [0, 4, 7, 10, 12]
		(1..-2).bsearch do |i|
			ary[i] >= 4
		end
		`, nil},
		{`
		ary = [0, 4, 7, 10, 12]
		(-5..-2).bsearch do |i|
			ary[i] >= 4
		end
		`, nil},
		{`
		ary = [0, 100, 100, 100, 200]
		(0..4).bsearch do |i|
			100 - ary[i]
		end
		`, 2},
		{`
		ary = [0, 100, 100, 100, 200]
		(0..4).bsearch do |i|
			200 - ary[i]
		end
		`, 4},
		{`
		ary = [0, 100, 100, 100, 200]
		(0..4).bsearch do |i|
			0 - ary[i]
		end
		`, 0},
		{`
		ary = [0, 100, 100, 100, 200]
		(2..4).bsearch do |i|
			100 - ary[i]
		end
		`, 3},
		{`
		ary = [0, 100, 100, 100, 200]
		(2..4).bsearch do |i|
			0 - ary[i]
		end
		`, nil},
		{`
		ary = [0, 100, 100, 100, 200]
		(-1..4).bsearch do |i|
			0 - ary[i]
		end
		`, nil},
		{`
		ary = [0, 100, 100, 100, 200]
		(4..0).bsearch do |i|
			0 - ary[i]
		end
		`, nil},
		{`
		ary = [0, 100, 100, 100, 200]
		(2..-1).bsearch do |i|
			0 - ary[i]
		end
		`, nil},
		{`
		ary = [0, 100, 100, 100, 200]
		(-5..-1).bsearch do |i|
			0 - ary[i]
		end
		`, nil},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestRangeBsearchMethodFail(t *testing.T) {
	v := initTestVM()
	testsFail := []struct {
		input    string
		errorMsg string
	}{
		{`
		ary = [0, 4, 7, 10, 12]
		(0..4).bsearch do |i|
			"Binary Search"
		end
		`, "TypeError: Expect Integer or Boolean type. got=String"},
	}

	for i, tt := range testsFail {
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.errorMsg, getFilename())
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestRangeEachMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		r = 0
		(1..(1+4)).each do |i|
		  r = r + i
		end
		r
		`, 15},
		{`
		r = 0
		a = 1
		b = 5
		(a..b).each do |i|
		  r = r + i
		end
		r
		`, 15},
		{`
		r = 0
		a = 5
		b = 1
		(a..b).each do |i|
		  r = r + i
		end
		r
		`, 15},
		{`
		r = 0
		a = -1
		b = -5
		(a..b).each do |i|
		  r = r + i
		end
		r
		`, -15},
		{`
		r = 0
		a = -5
		b = -1
		(a..b).each do |i|
		  r = r + i
		end
		r
		`, -15},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestRangeFirstMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		(1..5).first
		`, 1},
		{`
		(5..1).first
		`, 5},
		{`
		(-2..3).first
		`, -2},
		{`
		(-5..-7).first
		`, -5},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestRangeIncludeMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		(5..10).include?(10)
		`, true},
		{`
		(5..10).include?(11)
		`, false},
		{`
		(5..10).include?(7)
		`, true},
		{`
		(5..10).include?(5)
		`, true},
		{`
		(5..10).include?(4)
		`, false},
		{`
		(-5..1).include?(-2)
		`, true},
		{`
		(-5..-2).include?(-2)
		`, true},
		{`
		(-5..-3).include?(-2)
		`, false},
		{`
		(1..-5).include?(-2)
		`, true},
		{`
		(-2..-5).include?(-2)
		`, true},
		{`
		(-3..-5).include?(-2)
		`, false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestRangeLastMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		(1..5).last
		`, 5},
		{`
		(5..1).last
		`, 1},
		{`
		(-2..3).last
		`, 3},
		{`
		(-5..-7).last
		`, -7},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestRangeSizeMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		(1..5).size
		`, 5},
		{`
		(3..9).size
		`, 7},
		{`
		(-1..-5).size
		`, 5},
		{`
		(-1..7).size
		`, 9},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestRangeStepMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		sum = 0
		(2..9).step(3) do |i|
		  sum = sum + i
		end
		sum
		`, 15},
		{`
		sum = 0
		(2..-9).step(3) do |i|
		  sum = sum + i
		end
		sum
		`, 0},
		{`
		sum = 0
		a = 2
		b = 9
		c = 3
		(a..b).step(c) do |i|
		  sum = sum + i
		 end
		 sum
		`, 15},
		{`
		sum = 0
		a = -1
		b = 5
		c = 2
		(a..b).step(c) do |i|
		  sum = sum + i
		 end
		 sum
		`, 8},
		{`
		sum = 0
		a = -1
		b = -5
		c = 2
		(a..b).step(c) do |i|
		  sum = sum + i
		 end
		 sum
		`, 0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestRangeToStringMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		(1..5).to_s
		`, "(1..5)"},
		{`
		(-1..-5).to_s
		`, "(-1..-5)"},
		{`
		(-1..5).to_s
		`, "(-1..5)"},
		{`
		(1..-5).to_s
		`, "(1..-5)"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestRangeToArrayMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		(1..5).to_a.length
		`, 5},
		{`
		(1..5).to_a[2]
		`, 3},
		{`
		(-1..-5).to_a.length
		`, 5},
		{`
		(-1..-5).to_a[2]
		`, -3},
		{`
		(-1..3).to_a.length
		`, 5},
		{`
		(-1..3).to_a[2]
		`, 1},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
