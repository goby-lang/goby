package vm

import "testing"

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestEvalInfixBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
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
		{"(false || true) && (\"string\" == \"string\")", true},
		{"((10 > 3) && (3 < 4)) || ((10 == 10) || false)", true},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestInitializeBoolean(t *testing.T) {
	if !TRUE.Value {
		t.Errorf("expected 'true'. got=%t", TRUE.Value)
	}

	if FALSE.Value {
		t.Errorf("expected 'false'. got=%t", FALSE.Value)
	}
}
