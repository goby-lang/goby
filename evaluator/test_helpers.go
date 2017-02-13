package evaluator

import (
	"testing"
	"github.com/st0012/rooby/object"
	"github.com/st0012/rooby/parser"
	"github.com/st0012/rooby/lexer"
)

func testEval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	return Eval(program, object.NewEnvironment())
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
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
	result, ok := obj.(*object.Integer)
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
	result, ok := obj.(*object.String)
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
	result, ok := obj.(*object.Class)
	if !ok {
		t.Errorf("object is not a Class. got=%T (%+v", obj, obj)
		return false
	}

	if result.Name.Value != expected {
		t.Errorf("expect Class's name to be %s. got=%s", expected, result.Name.Value)
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
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
