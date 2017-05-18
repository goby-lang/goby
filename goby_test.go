package main

import (
	"github.com/goby-lang/goby/bytecode"
	"github.com/goby-lang/goby/parser"
	"github.com/goby-lang/goby/vm"
	"io/ioutil"
	"path"
	"runtime"
	"testing"
)

func TestRequireRelative(t *testing.T) {
	filename := "main.gb"
	fileDir := "require_test"
	result := execFile(fileDir, filename)

	testIntegerObject(t, result, 160)
}

func execFile(fileDir, filename string) vm.Object {
	_, currentPath, _, _ := runtime.Caller(0)
	// dir is now project root: goby/
	dir, _ := path.Split(currentPath)
	// dir is now goby/test_fixtures/FILE_DIR
	dir = path.Join(dir, "./test_fixtures/", fileDir)

	filepath := path.Join(dir, filename)
	file, err := ioutil.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	program := parser.BuildAST(file)
	g := bytecode.NewGenerator(program)
	bytecodes := g.GenerateByteCode(program)

	v := vm.New(dir, []string{})
	v.ExecBytecodes(bytecodes, filepath)
	return v.GetExecResult()
}

func testIntegerObject(t *testing.T, obj vm.Object, expected int) bool {
	switch result := obj.(type) {
	case *vm.IntegerObject:
		if result.Value != expected {
			t.Errorf("object has wrong value. expect=%d, got=%d", expected, result.Value)
			return false
		}

		return true
	case *vm.Error:
		t.Error(result.Message)
		return false
	default:
		t.Errorf("object is not Integer. got=%T (%+v).", obj, obj)
		return false
	}
}