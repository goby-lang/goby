package vm

import (
	"testing"
)

// TODO: Add failed tests
func TestFileObject(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		f = File.new("../test_fixtures/file_test/size.gb")
		f.name
		`, "../test_fixtures/file_test/size.gb"},
		{`
		f = File.new("../test_fixtures/file_test/size.gb")
		f.size
		`, 22},
		{`
		f = File.new("../test_fixtures/file_test/size.gb")
		f.close
		`, nil},
		{`
		f = File.new("../test_fixtures/file_test/size.gb")
		f.read
		`, "this file's size is\n22"},
		{`
		file = ""
		File.open("../test_fixtures/file_test/size.gb", "r", 0755) do |f|
	 	  file = f.read
	 	end
	 	file
		`, "this file's size is\n22"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

// Method test

func TestFileBasenameMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
				File.basename("/home/goby/plugin/test.gb")
		`, "test.gb"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFileDeleteMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		File.open("/tmp/out1.txt", "w", 0755)
		File.open("/tmp/out2.txt", "w", 0755)
		File.open("/tmp/out3.txt", "w", 0755)

		File.delete("/tmp/out1.txt", "/tmp/out2.txt", "/tmp/out3.txt")
		`, 3},
		{`
		File.open("/tmp/out.txt", "w", 0755)
		File.delete("/tmp/out.txt")
		File.exist?("/tmp/out.txt")
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

func TestFileExtnameMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		File.extname("loop.gb")
		`, ".gb"},
		{`
		File.extname("text.txt")
		`, ".txt"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFileJoinMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		File.join("test1", "test2", "test3")
		`, "test1/test2/test3"},
		{`
		File.join("goby", "plugin")
		`, "goby/plugin"},
		{`
		File.join("plugin")
		`, "plugin"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFileWriteMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		l = 0
		File.open("/tmp/out.txt", "w", 0755) do |f|
		  l = f.write("12345")
		end

		l
		`, 5},
		{`
		File.open("/tmp/out.txt", "w", 0755) do |f|
		  f.write("Goby is awesome!!!")
		end

		File.new("/tmp/out.txt").read
		`, "Goby is awesome!!!"},
		{`
		File.open("/tmp/out.txt", "w", 0755)
		File.new("/tmp/out.txt").size
		`, 0},
		{`
		File.open("/tmp/out.txt", "w", 0755)
		File.exist?("/tmp/out.txt")
		`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFileSizeMethod(t *testing.T) {
	input := `
	File.size("../test_fixtures/file_test/size.gb")
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	checkExpected(t, 0, evaluated, 22)
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestFileSplitMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		File.split("/home/goby/plugin/test.gb")
		`, []interface{}{"/home/goby/plugin/", "test.gb"}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

//@TODO add test for chmod form a847c8b41f29657b380c1731ec36a660dbf49bc4
