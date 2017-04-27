package main

import (
	"github.com/rooby-lang/rooby/bytecode"
	"github.com/rooby-lang/rooby/parser"
	"github.com/rooby-lang/rooby/vm"
	"io/ioutil"
	"path"
	"runtime"
	"testing"
)

//func TestRequireRelative(t *testing.T) {
//	filename := "main.ro"
//	fileDir := "require_test"
//	result := execFile(fileDir, filename)
//
//	testIntegerObject(t, result, 50)
//}

func execFile(fileDir, filename string) vm.Object {
	_, currentPath, _, _ := runtime.Caller(0)
	// dir is now project root: rooby/
	dir, _ := path.Split(currentPath)
	// dir is now rooby/test_fixtures/FILE_DIR
	dir = path.Join(dir, "./test_fixtures/", fileDir)

	filepath := path.Join(dir, filename)
	file, err := ioutil.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	program := parser.BuildAST(file)
	g := bytecode.NewGenerator(program)
	bytecodes := g.GenerateByteCode(program)

	v := vm.New(dir)
	v.ExecBytecodes(bytecodes, filepath)
	return v.GetExecResult()
}

func testIntegerObject(t *testing.T, obj vm.Object, expected int) bool {
	result, ok := obj.(*vm.IntegerObject)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v).", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. expect=%d, got=%d", expected, result.Value)
		return false
	}

	return true
}
