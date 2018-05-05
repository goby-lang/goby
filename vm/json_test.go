package vm

import "testing"

func TestJSONValidateMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require "json"
		JSON.validate('{"Name": "Stan"}')
	`, true},
		{`
		require "json"
		JSON.validate('{"Name": "Stan}')
	`, false},
		{`
		require "json"
		JSON.validate('{"Name": Stan}')
	`, false},
		{`
		require "json"
		JSON.validate('{Name: "Stan"')
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

func TestJSONObjectParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require "json"
		h = JSON.parse('{"Name": "Stan"}')
		h["Name"]`, "Stan"},
		{`
		require "json"
		h = JSON.parse('{"Age": 23}')
		h["Age"]`, 23},
		{`
		require "json"
		h = JSON.parse('
		  {
			"Project": {
			  "Name": "Goby"
			}
		  }
		')
		h["Project"]["Name"]`, "Goby"},
		{`
		require "json"
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
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
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
		require "json"
		a = JSON.parse('[{"Name": "Stan"}]')
		h = a.first
		h["Name"]`, "Stan"},
		{`
		require "json"
		a = JSON.parse('[{"Age": 23}]')
		h = a.first
		h["Age"]`, 23},
		{`
		require "json"
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
		require "json"
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
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
