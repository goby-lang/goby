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

func TestFileExtname(t *testing.T) {
	filename := "extname.gb"
	fileDir := "file_test"

	result := execFile(fileDir, filename)

	testStringObject(t, result, ".gb")
}

func TestFileJoin(t *testing.T) {
	filename := "join.gb"
	fileDir := "file_test"

	result := execFile(fileDir, filename)

	testStringObject(t, result, "home/goby/test")
}

func TestFileSplit(t *testing.T) {
	filename := "split.gb"
	fileDir := "file_test"

	result := execFile(fileDir, filename)

	expected := &vm.ArrayObject{Elements: []vm.Object{
		&vm.StringObject{Value: "/home/goby/"},
		&vm.StringObject{Value: ".settings"},
	}}
	testArrayObject(t, result, expected)
}

func TestFileBasename(t *testing.T) {
	filename := "basename.gb"
	fileDir := "file_test"

	result := execFile(fileDir, filename)

	testStringObject(t, result, "test.gb")
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

func testStringObject(t *testing.T, obj vm.Object, expected string) bool {
	switch result := obj.(type) {
	case *vm.StringObject:
		if result.Value != expected {
			t.Errorf("object has wrong value. expect=%d, got=%d", expected, result.Value)
			return false
		}

		return true
	case *vm.Error:
		t.Error(result.Message)
		return false
	default:
		t.Errorf("object is not String. got=%T (%+v).", obj, obj)
		return false
	}
}

func testArrayObject(t *testing.T, obj vm.Object, expected *vm.ArrayObject) bool {
	result, ok := obj.(*vm.ArrayObject)
	if !ok {
		t.Errorf("object is not Array. got=%T (%+v)", obj, obj)
		return false
	}

	if len(result.Elements) != len(expected.Elements) {
		t.Fatalf("Don't equals length of array. Expect %d, got=%d", len(expected.Elements), len(result.Elements))
	}

	for i := 0; i < len(result.Elements); i++ {
		intObj, ok := expected.Elements[i].(*vm.IntegerObject)
		if ok {
			testIntegerObject(t, result.Elements[i], intObj.Value)
			continue
		}
		str, ok := expected.Elements[i].(*vm.StringObject)
		if ok {
			testStringObject(t, result.Elements[i], str.Value)
			continue
		}

		b, ok := expected.Elements[i].(*vm.BooleanObject)
		if ok {
			testBooleanObject(t, result.Elements[i], b.Value)
			continue
		}

		t.Fatalf("object is wrong type %T", expected.Elements[i])
	}

	return true
}

func testBooleanObject(t *testing.T, obj vm.Object, expected bool) bool {
	switch result := obj.(type) {
	case *vm.BooleanObject:
		if result.Value != expected {
			t.Errorf("object has wrong value. expect=%d, got=%d", expected, result.Value)
			return false
		}

		return true
	case *vm.Error:
		t.Error(result.Message)
		return false
	default:
		t.Errorf("object is not Boolean. got=%T (%+v).", obj, obj)
		return false
	}
}
