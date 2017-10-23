package vm

import "testing"

func TestJSONValidateMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		JSON.validate('{"Name": "Stan"}')
	`, true},
		{`
		JSON.validate('{"Name": "Stan}')
	`, false},
		{`
		JSON.validate('{"Name": Stan}')
	`, false},
		{`
		JSON.validate('{Name: "Stan"')
	`, false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEvalWithRequire(t, tt.input, getFilename(), "json")
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestJSONObjectParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		h = JSON.parse('{"Name": "Stan"}')
		h["Name"]`, "Stan"},
		{`
		h = JSON.parse('{"Age": 23}')
		h["Age"]`, 23},
		{`
		h = JSON.parse('
		  {
			"Project": {
			  "Name": "Goby"
			}
		  }
		')
		h["Project"]["Name"]`, "Goby"},
		{`
		h = JSON.parse('
		  {
			"Project": {
			  "Name": "Goby",
			  "Months": 7
			}
		  }
		')
		h["Project"]["Months"]`, 7},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEvalWithRequire(t, tt.input, getFilename(), "json")
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestJSONObjectArrayParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		a = JSON.parse('[{"Name": "Stan"}]')
		h = a.first
		h["Name"]`, "Stan"},
		{`
		a = JSON.parse('[{"Age": 23}]')
		h = a.first
		h["Age"]`, 23},
		{`
		a = JSON.parse('
		  [{
			"Projects": [{
			  "Name": "Goby"
			}]
		  }]
		')
		h = a.first
		h["Projects"][0]["Name"]`, "Goby"},
		{`
		a = JSON.parse('
		  [{
			"Projects": [{
			  "Name": "Goby",
			  "Months": 7
			}]
		  }]
		')
		h = a.first
		h["Projects"][0]["Months"]`, 7},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEvalWithRequire(t, tt.input, getFilename(), "json")
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
