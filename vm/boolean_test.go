package vm

import "testing"

func TestBooleanClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Boolean.class.name`, "Class"},
		{`Boolean.superclass.name`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"true", true},
		{"false", false},
		{"'true'", "true"},
		{"'false'", "false"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestBooleanComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"'true' == true", false},
		{"'false' == false", false},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"true == 10", false},
		{"true != 10", true},
		{"false == 10", false},
		{"false != 10", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestBooleanLogicalExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"true && true", true},
		{"false && true", false},
		{"true && false", false},
		{"false && false", false},
		{"true || true", true},
		{"false || true", true},
		{"true || false", true},
		{"false || false", false},
		{"100 > 10 && true == false", false},
		{"true && true == true", true},
		{"true && 1", 1},
		{"false && 1", false},
		{"'true' && 1", 1},
		{"'false' && 1", 1},
		{`
		h = {}
		h && true
		`, true},
		{"(false || true) && (\"string\" == \"string\")", true},
		{"((10 > 3) && (3 < 4)) || ((10 == 10) || false)", true},
		{`
		a = 0

		# a = 10 shouldn't be executed
		false && a = 10
		a
		`, 0},
		{`
		a = 0

		# a = 10 shouldn't be executed
		nil && a = 10
		a
		`, 0},
		{`
		a = 0

		# a = 10 should be executed
		true && a = 10
		a
		`, 10},
		{`
		a = 0

		# a = 10 should be executed
		false || a = 10
		a
		`, 10},
		{`
		a = 0

		# a = 10 shouldn't be executed
		true || a = 10
		a
		`, 0},
		{`
		a = 0
		false || true && a = 10
		a
		`, 10},
		{`
		a = false || 10
		a
		`, 10},
		{`
		a = nil || 10
		a
		`, 10},
		{`
		a = true || 10
		a
		`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestBooleanAssignmentByOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`a = true;  a ||= 123;       a;`, true},
		{`a = true;  a ||= "string";  a;`, true},
		{`a = true;  a ||= false;     a;`, true},
		{`a = true;  a ||= (1..4);    a;`, true},
		{`a = true;  a ||= { b: 1 };  a;`, true},
		{`a = true;  a ||= Object;    a;`, true},
		{`a = true;  a ||= [1, 2, 3]; a;`, true},
		{`a = true;  a ||= nil;       a;`, true},
		{`a = true;  a ||= nil || 1;  a;`, true},
		{`a = true;  a ||= 1 || nil;  a;`, true},
		{`a = false; a ||= 123;       a;`, 123},
		{`a = false; a ||= "string";  a;`, "string"},
		{`a = false; a ||= false;     a;`, false},
		{`a = false; a ||= (1..4);    a.to_s;`, "(1..4)"},
		{`a = false; a ||= { b: 1 };  a["b"];`, 1},
		{`a = false; a ||= Object;    a.name;`, "Object"},
		{`a = false; a ||= [1, 2, 3]; a[0];`, 1},
		{`a = false; a ||= [1, 2, 3]; a[1];`, 2},
		{`a = false; a ||= [1, 2, 3]; a[2];`, 3},
		{`a = false; a ||= nil;       a;`, nil},
		{`a = false; a ||= nil || 1;  a;`, 1},
		{`a = false; a ||= 1 || nil;  a;`, 1},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestInitializeBoolean(t *testing.T) {
	if !TRUE.value {
		t.Errorf("expected 'true'. got=%t", TRUE.value)
	}

	if FALSE.value {
		t.Errorf("expected 'false'. got=%t", FALSE.value)
	}
}
