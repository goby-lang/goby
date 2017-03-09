package evaluator_test

import (
	"github.com/st0012/Rooby/evaluator"
	"github.com/st0012/Rooby/lexer"
	"github.com/st0012/Rooby/object"
	"github.com/st0012/Rooby/parser"
	"testing"
)

func testEval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	mainObj := object.MainObj
	return evaluator.Eval(program, mainObj.Scope)
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
	if obj != object.NULL {
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
