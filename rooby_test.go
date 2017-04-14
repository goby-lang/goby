package main

import (
	"github.com/st0012/Rooby/code_generator"
	"github.com/st0012/Rooby/lexer"
	"github.com/st0012/Rooby/parser"
	"testing"
	"github.com/st0012/Rooby/vm"
)

func testEval(t *testing.T, input string) vm.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	cg := code_generator.New(program)
	bytecodes := cg.GenerateByteCode(program)
	return testExec(bytecodes)
}

func testExec(bytecodes string) vm.Object {
	p := vm.NewBytecodeParser()
	v := vm.New()
	p.VM = v
	p.Parse(bytecodes)
	cf := vm.NewCallFrame(v.LabelTable[vm.PROGRAM]["ProgramStart"][0])
	cf.Self = vm.MainObj
	v.CallFrameStack.Push(cf)
	v.Exec()

	return v.Stack.Top()
}

func testClassObject(t *testing.T, obj vm.Object, expected string) bool {
	result, ok := obj.(*vm.RClass)
	if !ok {
		t.Errorf("object is not a Class. got=%T (%+v", obj, obj)
		return false
	}

	if result.Name != expected {
		t.Errorf("expect Class's name to be %s. got=%s", expected, result.Name)
	}

	return true
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

func testNullObject(t *testing.T, obj vm.Object) bool {
	if obj != vm.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}

func testStringObject(t *testing.T, obj vm.Object, expected string) bool {
	result, ok := obj.(*vm.StringObject)
	if !ok {
		t.Errorf("object is not a String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. expect=%s, got=%s", expected, result.Value)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj vm.Object, expected bool) bool {
	result, ok := obj.(*vm.BooleanObject)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. expect=%d, got=%d", expected, result.Value)
		return false
	}

	return true
}

func isError(obj vm.Object) bool {
	if obj != nil {
		return obj.Type() == vm.ERROR_OBJ
	}
	return false
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}