package evaluator_test

import (
	"github.com/st0012/rooby/evaluator"
	"github.com/st0012/rooby/lexer"
	"github.com/st0012/rooby/object"
	"github.com/st0012/rooby/parser"
	"testing"

	"github.com/st0012/rooby/initializer"
)

func testEval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	mainObj := initializer.MainObj
	return evaluator.Eval(program, mainObj.Scope)
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.BooleanObject)
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

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.IntegerObject)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. expect=%d, got=%d", expected, result.Value)
		return false
	}

	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.StringObject)
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

func testClassObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.RClass)
	if !ok {
		t.Errorf("object is not a Class. got=%T (%+v", obj, obj)
		return false
	}

	if result.Name != expected {
		t.Errorf("expect Class's name to be %s. got=%s", expected, result.Name)
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != initializer.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}

	return true
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

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
