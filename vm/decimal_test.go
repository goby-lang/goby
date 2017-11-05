package vm

import (
	"testing"
)

func TestDecimalClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Decimal.class.name`, "Class"},
		{`Decimal.superclass.name`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalConversionString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`'13.5'.to_d.to_s`, "13.5"},
		{`'13.5'.to_d.fraction`, "27/2"},
		{`'13.5'.to_d.inverse.fraction`, "2/27"},
		{`'20/40'.to_d.reduction`, "1/2"},
		{`'40/20'.to_d.fraction`, "2/1"},
		{`'40/20'.to_d.reduction`, "2"},
		{`'-13.5'.to_d.numerator.to_s`, "-27"},
		{`'-13.5'.to_d.denominator.to_s`, "2"},
		{`'129.30928304982039482039842'.to_d.numerator.to_s`, "6465464152491019741019921"},
		{`'129.30928304982039482039842'.to_d.denominator.to_s`, "50000000000000000000000"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalConversionNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'-13.5'.to_d.to_i`, -13},
		{`'-13.5'.to_d.to_f`, -13.5},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalArithmeticOperationWithDecimal(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`('13.5'.to_d  +  '3.5'.to_d).to_s`, "17"},
		{`('13.5'.to_d  +  '3.2'.to_d).to_s`, "16.7"},
		{`('13.5'.to_d  -  '3.2'.to_d).to_s`, "10.3"},
		{`('13.5'.to_d  *  '3.2'.to_d).to_s`, "43.2"},
		{`('13.5'.to_d  /  '3.75'.to_d).to_s`, "3.6"},
		{`('16.0'.to_d  ** '3.5'.to_d).to_s`, "16384"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalArithmeticOperationWithInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`('13.5'.to_d  +  3).to_s`, "16.5"},
		{`('13.5'.to_d  -  3).to_s`, "10.5"},
		{`('13.5'.to_d  *  3).to_s`, "40.5"},
		{`('13.5'.to_d  /  3).to_s`, "4.5"},
		{`('13.5'.to_d  ** 3).to_s`, "2460.375"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalArithmeticOperationWithFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`('16.1'.to_d  + "1.1".to_d).to_s`, "17.2"},
		{`('16.1'.to_d  + "1.1".to_f).to_s`, "17.200000000000000088817841970012523233890533447265625"},
		{`('16.1'.to_d  - "1.1".to_d).to_s`, "15"},
		{`('16.1'.to_d  - "1.1".to_f).to_s`, "14.999999999999999911182158029987476766109466552734375"},
		{`('16.1'.to_d  * "1.1".to_d).to_s`, "17.71"},
		{`('16.1'.to_d  * "1.1".to_f).to_s`, "17.7100000000000014299672557172016240656375885009765625"},
		{`('16.1'.to_d  / "1.1".to_d).to_s`, "14.636363636363636363636363636363636363636363636363636363636364"},
		{`('16.1'.to_d  / "1.1".to_f).to_s`, "14.636363636363635181845243208924373053647177484459953103948567"},
		{`('16.1'.to_d  ** "1.1".to_d).to_s`, "21.257317715840930105741790612228214740753173828125"},
		{`('16.1'.to_d  ** "1.1".to_f).to_s`, "21.257317715840930105741790612228214740753173828125"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalArithmeticOperationFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`'1'.to_d + "p"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_d - "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_d / "t"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalComparisonWithFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'1.5'.to_d >   '2.5'.to_d`, false},
		{`'2.5'.to_d >   '1.5'.to_d`, true},
		{`'3.5'.to_d >   '3.5'.to_d`, false},
		{`'1.5'.to_d <   '2.5'.to_d`, true},
		{`'2.5'.to_d <   '1.5'.to_d`, false},
		{`'3.5'.to_d <   '3.5'.to_d`, false},
		{`'1.5'.to_d >=  '2.5'.to_d`, false},
		{`'2.5'.to_d >=  '1.5'.to_d`, true},
		{`'3.5'.to_d >=  '3.5'.to_d`, true},
		{`'1.5'.to_d <=  '2.5'.to_d`, true},
		{`'2.5'.to_d <=  '1.5'.to_d`, false},
		{`'3.5'.to_d <=  '3.5'.to_d`, true},
		{`'1.5'.to_d <=> '2.5'.to_d`, -1},
		{`'2.5'.to_d <=> '1.5'.to_d`, 1},
		{`'3.5'.to_d <=> '3.5'.to_d`, 0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalComparisonWithInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'1'.to_d >   2`, false},
		{`'2'.to_d >   1`, true},
		{`'3'.to_d >   3`, false},
		{`'1'.to_d <   2`, true},
		{`'2'.to_d <   1`, false},
		{`'3'.to_d <   3`, false},
		{`'1'.to_d >=  2`, false},
		{`'2'.to_d >=  1`, true},
		{`'3'.to_d >=  3`, true},
		{`'1'.to_d <=  2`, true},
		{`'2'.to_d <=  1`, false},
		{`'3'.to_d <=  3`, true},
		{`'1'.to_d <=> 2`, -1},
		{`'2'.to_d <=> 1`, 1},
		{`'3'.to_d <=> 3`, 0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalComparisonFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`'1'.to_d > "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_d >= "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_d < "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_d <= "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_d <=> "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalEquality(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'123.5'.to_d  ==  '123.5'.to_d`, true},
		//{`'123'.to_d    ==  123`, true},
		{`'123.5'.to_d  ==  '124'.to_d`, false},
		{`'123.5'.to_d  ==  "123.5"`, false},
		{`'123.5'.to_d  ==  (1..3)`, false},
		{`'123.5'.to_d  ==  { a: 1 }`, false},
		{`'123.5'.to_d  ==  [1]`, false},
		{`'123.5'.to_d  ==  Float`, false},
		{`'123.5'.to_d  !=  '123.5'.to_d`, false},
		{`'123.5'.to_d  !=  123`, true},
		{`'123.5'.to_d  !=  '124'.to_d`, true},
		{`'123.5'.to_d  !=  "123.5"`, true},
		{`'123.5'.to_d  !=  (1..3)`, true},
		{`'123.5'.to_d  !=  { a: 1 }`, true},
		{`'123.5'.to_d  !=  [1]`, true},
		{`'123.5'.to_d  !=  Float`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalConversions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		//{`'100.3'.to_d.to_i`, 100},
		{`'100.3'.to_d.to_s`, "100.3"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestArrayConversionWithDecimal(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		"129.30928304982039482039842".to_d.to_a[0].to_s
		`, "6465464152491019741019921"},
		{`
		"129.30928304982039482039842".to_d.to_a[1].to_s
		`, "50000000000000000000000"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestArrayConversionWithInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		"355/133".to_d.to_ai
		`, []interface{}{355, 133}},
		{`
		"-355/133".to_d.to_ai
		`, []interface{}{-355, 133}},
		{`
		"129.3".to_d.to_ai
		`, []interface{}{1293, 10}},
		{`
		"129.30928304982039482039842".to_d.to_ai
		`, []interface{}{-8964879735843078383, -9123183826594430976}},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		testArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
		vm.checkSP(t, i, 1)
	}
}
