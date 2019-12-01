package vm

import (
	"testing"

	"github.com/dlclark/regexp2"
)

// Verification helpers for tests

// VerifyExpected checks if the given Object is the expected class
func VerifyExpected(t *testing.T, i int, evaluated Object, expected interface{}) {
	t.Helper()
	if isError(evaluated) {
		t.Errorf("At test case %d: %s", i, evaluated.ToString())
		return
	}

	switch expected := expected.(type) {
	case int:
		verifyIntegerObject(t, i, evaluated, expected)
	case float64:
		verifyFloatObject(t, i, evaluated, expected)
	case string:
		verifyStringObject(t, i, evaluated, expected)
	case bool:
		verifyBooleanObject(t, i, evaluated, expected)
	case []interface{}:
		verifyArrayObject(t, i, evaluated, expected)
	case nil:
		verifyNullObject(t, i, evaluated)
	default:
		t.Errorf("Unknown type %T at case %d", expected, i)
	}
}

func verifyIntegerObject(t *testing.T, i int, obj Object, expected int) bool {
	t.Helper()
	switch result := obj.(type) {
	case *IntegerObject:
		if result.value != expected {
			t.Errorf("At test case %d: object has wrong value. expect=%d, got=%d", i, expected, result.value)
			return false
		}

		return true
	case *Error:
		t.Errorf("At test case %d: %s", i, result.Message())
		return false
	default:
		t.Errorf("At test case %d: object is not Integer. got=%s (%+v).", i, obj.Class().Name, obj)
		return false
	}
}

func verifyFloatObject(t *testing.T, i int, obj Object, expected float64) bool {
	t.Helper()
	switch result := obj.(type) {
	case *FloatObject:
		if result.value != expected {
			t.Errorf("At test case %d: object has wrong value. expect=%f, got=%f", i, expected, result.value)
			return false
		}

		return true
	case *Error:
		t.Errorf("At test case %d: %s", i, result.Message())
		return false
	default:
		t.Errorf("At test case %d: object is not Float. got=%s (%+v).", i, obj.Class().Name, obj)
		return false
	}
}

func verifyNullObject(t *testing.T, i int, obj Object) bool {
	t.Helper()
	switch result := obj.(type) {
	case *NullObject:
		return true
	case *Error:
		t.Errorf("At test case %d: %s", i, result.Message())
		return false
	default:
		t.Errorf("At test case %d: object is not NULL. got=%s (%+v)", i, obj.Class().Name, obj)
		return false
	}
}

func verifyStringObject(t *testing.T, i int, obj Object, expected string) bool {
	t.Helper()
	var fuzStr string
	switch result := obj.(type) {
	case *StringObject:
		re, _ := regexp2.Compile("(?<=#<[a-zA-Z0-9_]+:)[0-9]{12}(?=[ ]>?)", 0)
		fuzStr, _ = re.Replace(result.value, "##OBJECTID##", 0, -1)
		if fuzStr != expected {
			t.Errorf("At test case %d: object has wrong value. expect=%q, got=%q", i, expected, result.value)
			return false
		}
		return true
	case *Error:
		t.Errorf(result.Message())
		return false
	default:
		t.Errorf("At test case %d: object is not String. got=%s (%+v).", i, obj.Class().Name, obj)
		return false
	}
}

func verifyBooleanObject(t *testing.T, i int, obj Object, expected bool) bool {
	t.Helper()
	switch result := obj.(type) {
	case *BooleanObject:
		if result.value != expected {
			t.Errorf("At test case %d: object has wrong value. expect=%t, got=%t", i, expected, result.value)
			return false
		}

		return true
	case *Error:
		t.Errorf(result.Message())
		return false
	default:
		t.Errorf("At test case %d: object is not Boolean. got=%s (%+v).", i, obj.Class().Name, obj)
		return false
	}
}

func verifyArrayObject(t *testing.T, index int, obj Object, expected []interface{}) bool {
	t.Helper()
	result, ok := obj.(*ArrayObject)
	if !ok {
		t.Errorf("At test case %d: object is not Array. got=%s (%+v)", index, obj.Class().Name, obj)
		return false
	}

	if len(result.Elements) != len(expected) {
		t.Errorf("Don't equals length of array. Expect %d, got=%d", len(expected), len(result.Elements))
	}

	for i := 0; i < len(result.Elements); i++ {
		VerifyExpected(t, index, result.Elements[i], expected[i])
	}

	return true
}

// Same as testHashObject(), but expects a ConcurrentArray.
func verifyConcurrentArrayObject(t *testing.T, index int, obj Object, expected []interface{}) bool {
	t.Helper()
	result, ok := obj.(*concurrentArrayObject)
	if !ok {
		t.Errorf("At test case %d: object is not ConcurrentArray. got=%s (%+v)", index, obj.Class().Name, obj)
		return false
	}

	if len(result.InternalArray.Elements) != len(expected) {
		t.Errorf("Don't equals length of array. Expect %d, got=%d", len(expected), len(result.InternalArray.Elements))
	}

	for i := 0; i < len(result.InternalArray.Elements); i++ {
		VerifyExpected(t, i, result.InternalArray.Elements[i], expected[i])
	}

	return true
}

// Same as testHashObject(), but expects a ConcurrentHash.
//
func verifyConcurrentHashObject(t *testing.T, index int, objectResult Object, expected map[string]interface{}) bool {
	t.Helper()
	result, ok := objectResult.(*ConcurrentHashObject)

	if !ok {
		t.Errorf("At test case %d: result is not ConcurrentHash. got=%s", index, objectResult.Class().Name)
		return false
	}

	pairs := make(map[string]Object)

	iterator := func(key, value interface{}) bool {
		pairs[key.(string)] = value.(Object)
		return true
	}

	result.internalMap.Range(iterator)

	return _checkHashPairs(t, pairs, expected)
}

// Tests a Hash Object, with a few limitations:
//
// - the tested hash must be shallow (no nested objects as values);
// - the test hash must have strings as keys;
// - the error message won't mention the key - only the value.
//
// The second limitation is currently the only Hash format in Goby, anyway.
//
func verifyHashObject(t *testing.T, index int, objectResult Object, expected map[string]interface{}) bool {
	t.Helper()
	result, ok := objectResult.(*HashObject)

	if !ok {
		t.Errorf("At test case %d: result is not Hash. got=%s", index, objectResult.Class().Name)
		return false
	}

	return _checkHashPairs(t, result.Pairs, expected)
}

// Testing API like testArrayObject(), but performed on bidimensional arrays.
//
// Input example:
//
//		evaluated = '[["a", 1], ["b", "2"]]'
//		expected = [][]interface{}{{"a", 1}, {"b", "2"}}
//		testBidimensionalArrayObject(t, i, evaluated, expected)
//
func verifyBidimensionalArrayObject(t *testing.T, index int, obj Object, expected [][]interface{}) bool {
	t.Helper()
	result, ok := obj.(*ArrayObject)
	if !ok {
		t.Errorf("At test case %d: object is not Array. got=%T (%+v)", index, obj, obj)
		return false
	}

	if len(result.Elements) != len(expected) {
		t.Errorf("Unexpected result size. Expect %d, got=%d", len(expected), len(result.Elements))
	}

	for i := 0; i < len(result.Elements); i++ {
		resultRow := result.Elements[i]
		expectedRow := expected[i]

		verifyArrayObject(t, index, resultRow, expectedRow)
	}

	return true
}

func isError(obj Object) bool {
	if obj != nil {
		_, ok := obj.(*Error)
		return ok
	}
	return false
}

// Internal helpers -----------------------------------------------------

func _checkHashPairs(t *testing.T, actual map[string]Object, expected map[string]interface{}) bool {
	if len(actual) != len(expected) {
		t.Errorf("Unexpected result size. Expected %d, got=%d", len(expected), len(actual))
	}

	for expectedKey, expectedValue := range expected {
		resultValue := actual[expectedKey]

		VerifyExpected(t, i, resultValue, expectedValue)
	}

	return true
}
