package vm

import (
	"testing"
)

func TestEaxtnameMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		require("file")
		File.extname("loop.gb")
		`, ".gb"},
		{`
		require("file")
		File.extname("text.txt")
		`, ".txt"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestBasenameMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		require("file")
		File.basename("/home/goby/plugin/test.gb")
		`, "test.gb"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestSplitMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected *ArrayObject
	}{
		{`
		require("file")
		File.split("/home/goby/plugin/test.gb")
		`, initializeArray([]Object{initializeString("/home/goby/plugin/"), initializeString("test.gb")})},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testArrayObject(t, evaluated, tt.expected)
	}
}

func TestJoinMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		require("file")
		File.join("test1", "test2", "test3")
		`, "test1/test2/test3"},
		{`
		require("file")
		File.join("goby", "plugin")
		`, "goby/plugin"},
		{`
		require("file")
		File.join("plugin")
		`, "plugin"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

//@TODO add test for size and chmod form a847c8b41f29657b380c1731ec36a660dbf49bc4
