package main

import (
	"io/ioutil"
	"testing"
	"github.com/rooby-lang/rooby/bytecode"
	"github.com/rooby-lang/rooby/vm"
)

func TestRequireFile(t *testing.T) {
	filename := "require_test/main.ro"
	result := execFile(filename)

	testIntegerObject(t, result, 10)
}

func execFile(filename string) vm.Object {
	filepath := "./test_fixtures/" + filename
	file, err := ioutil.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	program := buildAST(file)
	g := bytecode.NewGenerator(program)
	bytecodes := g.GenerateByteCode(program)
	v := vm.New()
	v.ExecBytecodes(bytecodes)
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