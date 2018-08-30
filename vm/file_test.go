package vm

import (
	"os/exec"
	"testing"
)

func TestFileObject(t *testing.T) {
	setup()
	defer teardown()

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
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFileObjectFail(t *testing.T) {
	setup()
	defer teardown()

	testsFail := []errorTestCase{
		{`f = File.new("fictitious.gb")`,
			`IOError: open fictitious.gb: no such file or directory`, 1},
		{`f = File.new("fictitious/")`,
			`IOError: open fictitious/: no such file or directory`, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

// Tests for class methods
func TestFileBasenameMethod(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		input    string
		expected string
	}{
		{`
				File.basename("/home/goby/plugin/test.gb")
		`, "test.gb"},
		// tests for non-existent file/dir
		{`
				File.basename("/home/goby/plugin/fictitious.gb")
		`, "fictitious.gb"},
		{`
				File.basename("/home/goby/plugin/fictitious/")
		`, "fictitious"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFileBasenameMethodFail(t *testing.T) {
	setup()
	defer teardown()

	testsFail := []errorTestCase{
		{`File.basename`,
			`ArgumentError: Expect 1 argument(s). got: 0`, 1},
		{`File.basename("test1.txt", "test2.txt")`,
			`ArgumentError: Expect 1 argument(s). got: 2`, 1},
		{`File.basename(1)`,
			`TypeError: Expect argument to be String. got: Integer`, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestFileChmodMethod(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		input    string
		expected int
	}{
		{`
		path = "/tmp/goby/chmod_out.txt"
		File.open(path, "r+", 0755)
		File.chmod(0777, path)
		`, 1},
		{`
		File.open("/tmp/goby/out1.txt", "w", 0755)
		File.open("/tmp/goby/out2.txt", "w", 0744)
		File.open("/tmp/goby/out3.txt", "w", 0644)
		File.chmod(0777, "/tmp/goby/out1.txt", "/tmp/goby/out2.txt", "/tmp/goby/out3.txt")
		`, 3},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFileChmodMethodFail(t *testing.T) {
	setup()
	defer teardown()

	testsFail := []errorTestCase{
		{`File.chmod`,
			`ArgumentError: Expect 2 or more argument(s). got: 0`, 1},
		{`File.chmod(0755)`,
			`ArgumentError: Expect 2 or more argument(s). got: 1`, 1},
		{`File.chmod(0755, "/tmp/goby/fictitious.gb")`,
			`IOError: chmod /tmp/goby/fictitious.gb: no such file or directory`, 1},
		{`
		File.open("/tmp/goby/out_chmod.txt", "w", 0755)
		File.chmod(0777, "/tmp/goby/out_chmod.txt", "/tmp/goby/fictitious.gb")
		`, `IOError: chmod /tmp/goby/fictitious.gb: no such file or directory`, 1},
		{`File.chmod("string", "filePath")`,
			`TypeError: Expect argument #1 to be Integer. got: String`, 1},
		{`
		File.open("/tmp/goby/out_chmod.txt", "w", 0755)
		File.chmod(-999, "/tmp/goby/out_chmod.txt")
		`, `ArgumentError: Invalid chmod number. got: -999`, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestFileDeleteMethod(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		File.open("/tmp/goby/out1.txt", "w", 0755)
		File.open("/tmp/goby/out2.txt", "w", 0755)
		File.open("/tmp/goby/out3.txt", "w", 0755)

		File.delete("/tmp/goby/out1.txt", "/tmp/goby/out2.txt", "/tmp/goby/out3.txt")
		`, 3},
		{`
		File.open("/tmp/goby/out.txt", "w", 0755)
		File.delete("/tmp/goby/out.txt")
		File.exist?("/tmp/goby/out.txt")
		`, false},
		{`
		File.delete
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

func TestFileDeleteMethodFail(t *testing.T) {
	setup()
	defer teardown()

	testsFail := []errorTestCase{
		{`File.delete("/tmp/goby/non-existent.txt")`,
			`IOError: remove /tmp/goby/non-existent.txt: no such file or directory`, 1},
		{`File.delete 1`,
			`TypeError: Expect argument #1 to be String. got: Integer`, 1},
		{`f = "/tmp/goby/out.txt"; File.open(f, "w", 0755);File.delete(f, 1)`,
			`TypeError: Expect argument #2 to be String. got: Integer`, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestFileExistMethod(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		input    string
		expected bool
	}{
		{`
		File.exist?("/tmp/goby/non-existent.txt")
		`, false},
		{`
		File.open("/tmp/goby/out1.txt", "w", 0755)
		File.exist?("/tmp/goby/out1.txt")
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

func TestFileExistMethodFail(t *testing.T) {
	setup()
	defer teardown()

	testsFail := []errorTestCase{
		{`File.exist?`,
			`ArgumentError: Expect 1 argument(s). got: 0`, 1},
		{`File.exist?("test1.txt", "test2.txt")`,
			`ArgumentError: Expect 1 argument(s). got: 2`, 1},
		{`File.exist? 1`,
			`TypeError: Expect argument to be String. got: Integer`, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
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
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFileExtnameMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`File.extname`,
			`ArgumentError: Expect 1 argument(s). got: 0`, 1},
		{`File.extname("test1.txt", "test2.txt")`,
			`ArgumentError: Expect 1 argument(s). got: 2`, 1},
		{`File.extname 1`,
			`TypeError: Expect argument to be String. got: Integer`, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
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
		{`
		File.join
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

func TestFileJoinMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`File.join(1)`,
			`TypeError: Expect argument to be String. got: Integer`, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestFileNewMethod(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		input    string
		expected string
	}{
		{`
		File.open("/tmp/goby/test.gb", "r+");a = File.new("/tmp/goby/test.gb")
		a.class.name
		`, "File"},
		{`
		File.open("/tmp/goby/test.gb", "r+");a = File.new("/tmp/goby/test.gb", "w")
		a.class.name
		`, "File"},
		{`
		File.open("/tmp/goby/test.gb", "r+");a = File.new("/tmp/goby/test.gb", "w", 0777)
		a.class.name
		`, "File"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFileNewMethodFail(t *testing.T) {
	setup()
	defer teardown()

	testsFail := []errorTestCase{
		{`
		File.new()
		`, `ArgumentError: Expect 1 to 3 argument(s). got: 0`, 1},
		{`
		File.new("/tmp/goby/test.gb", "w", 0777, "a")
		`, `ArgumentError: Expect 1 to 3 argument(s). got: 4`, 1},
		{`
		File.new(1)
		`, `TypeError: Expect argument #1 to be String. got: Integer`, 1},
		{`
		File.new("/tmp/goby/test.gb", 1, 0777)
		`, `TypeError: Expect argument #2 to be String. got: Integer`, 1},
		{`
		File.new("/tmp/goby/test.gb", "p", 0777)
		`, `ArgumentError: Unknown file mode: p`, 1},
		{`
		File.new("/tmp/goby/test.gb", "w", "e")
		`, `TypeError: Expect argument #3 to be Integer. got: String`, 1},
		{`
		File.new("/tmp/goby/test.gb", "w", -99999)
		`, `ArgumentError: Invalid chmod number. got: -99999`, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestFileSizeMethod(t *testing.T) {
	input := `
	File.size("../test_fixtures/file_test/size.gb")
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	VerifyExpected(t, 0, evaluated, 22)
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestFileSizeMethodFail(t *testing.T) {
	setup()
	defer teardown()

	testsFail := []errorTestCase{
		{`
		File.size()
		`, `ArgumentError: Expect 1 argument(s). got: 0`, 1},
		{`
		File.size("../test_fixtures/file_test/size.gb","/tmp/goby/test.gb")
		`, `ArgumentError: Expect 1 argument(s). got: 2`, 1},
		{`
		File.size(1)
		`, `TypeError: Expect argument to be String. got: Integer`, 1},
		{`
		File.size("/tmp/goby/fictitious.gb")
		`, `IOError: stat /tmp/goby/fictitious.gb: no such file or directory`, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
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
		verifyArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFileSplitMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		File.split()
		`, `ArgumentError: Expect 1 argument(s). got: 0`, 1},
		{`
		File.split("/home/goby/plugin/test.gb", "/home/goby/plugin/test.gb")
		`, `ArgumentError: Expect 1 argument(s). got: 2`, 1},
		{`
		File.split(1)
		`, `TypeError: Expect argument to be String. got: Integer`, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

// Tests for instance methods

func TestFileCloseMethod(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		f = File.new("/tmp/goby/out.txt", "w", 0755)
		f.close
		f.close
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

func TestFileNameMethod(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		input    string
		expected string
	}{
		{`
		l = ""
		File.open("/tmp/goby/out.txt", "w", 0755) do |f|
		  l = f.name
		end
		l
		`, "/tmp/goby/out.txt"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFileReadMethod(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		input    string
		expected string
	}{
		{`
		l = ""
		File.open("/tmp/goby/out.txt", "w", 0755) do |f|
		  f.write("Hello, Goby!")
			l = f.read
		end
		l
		`, "Hello, Goby!"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFileInstanceSizeMethod(t *testing.T) {
	input := `
		l = 0
		File.open("../test_fixtures/file_test/size.gb", "r", 0755) do |f|
			l = f.size
		end
		l
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	VerifyExpected(t, 0, evaluated, 22)
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestFileWriteMethod(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		l = 0
		File.open("/tmp/goby/out.txt", "w", 0755) do |f|
		  l = f.write("12345")
		end

		l
		`, 5},
		{`
		File.open("/tmp/goby/out.txt", "w", 0755) do |f|
		  f.write("Goby is awesome!!!")
		end

		File.new("/tmp/goby/out.txt").read
		`, "Goby is awesome!!!"},
		{`
		File.open("/tmp/goby/out.txt", "w", 0755)
		File.new("/tmp/goby/out.txt").size
		`, 0},
		{`
		File.open("/tmp/goby/out.txt", "w", 0755)
		File.exist?("/tmp/goby/out.txt")
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

// Helper functions -----------------------------------------------------
func setup() {
	// initialize test directory
	exec.Command("rm", "-rf", "/tmp/goby/*").Run()
	exec.Command("mkdir", "/tmp/goby").Run()
}

func teardown() {
	// initialize test directory
	exec.Command("rm", "-rf", "/tmp/goby/*").Run()
}
