package vm

import (
	"io/ioutil"
	"os"
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

func TestSizeMethod(t *testing.T) {
	//creat in tmp file with size
	err := ioutil.WriteFile("/tmp/testSize", []byte("test\nadd some data\n"), 0644)
	if err != nil {
		panic(err)
	}
	fileStat, err := os.Stat("/tmp/testSize")
	if err != nil {
		panic(err)
	}

	defer os.Remove("/tmp/testSize")

	//check it
	tests := []struct {
		input    string
		expected int
	}{
		{`
		require("file")
		File.size("/tmp/testSize")
		`, int(fileStat.Size())},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestChmodMethod(t *testing.T) {
	//creat in tmp file with size
	err := ioutil.WriteFile("/tmp/testChmod", []byte("test\nadd some data\n"), 0644)
	if err != nil {
		panic(err)
	}

	defer os.Remove("/tmp/testSize")

	tests := []struct {
		input    string
		expected int
	}{
		{`
		require("file")
		File.chmod(0755, "/tmp/testChmod")
		`, 1},
	}

	for _, tt := range tests {
		testEval(t, tt.input)
	}

	fileStat, err := os.Stat("/tmp/testChmod")
	if err != nil {
		panic(err)
	}
	if fileStat.Mode() != 0755 {
		t.Errorf("Filemod incorect expected %o got=%o", 0755, fileStat.Mode())
	}

}
