package vm

import (
	"testing"
)

func TestEachThroughRange(t *testing.T) {
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestRangeToArray(t *testing.T) {
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestRangeToString(t *testing.T) {
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestFirstAndLast(t *testing.T) {
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestSize(t *testing.T) {
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestStepThroughRange(t *testing.T) {
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestInclude(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		(5..10).include(10)
		`, true},
		{`
		(5..10).include(11)
		`, false},
		{`
		(5..10).include(7)
		`, true},
		{`
		(5..10).include(5)
		`, true},
		{`
		(5..10).include(4)
		`, false},
		{`
		(-5..1).include(-2)
		`, true},
		{`
		(-5..-2).include(-2)
		`, true},
		{`
		(-5..-3).include(-2)
		`, false},
		{`
		(1..-5).include(-2)
		`, true},
		{`
		(-2..-5).include(-2)
		`, true},
		{`
		(-3..-5).include(-2)
		`, false},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}
