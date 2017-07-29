package metago

import (
	"fmt"
	"reflect"
)

// CallFunc will call the target method on receiver with given args and returns the execution result.
// If no result is returned from target method, CallFunc will still returns a nil.
func CallFunc(receiver interface{}, methodName string, args ...interface{}) interface{} {
	var ptr reflect.Value
	var value reflect.Value

	value, ok := receiver.(reflect.Value)

	if !ok {
		value = reflect.ValueOf(receiver)
	}

	funcArgs := WrapArguments(args...)

	ptr, value = getReflectPtrAndValue(value, receiver)

	method := value.MethodByName(methodName)

	if method.IsValid() {
		return UnwrapReflectValues(method.Call(funcArgs))
	}

	method = ptr.MethodByName(methodName)

	if method.IsValid() {
		return UnwrapReflectValues(method.Call(funcArgs))
	}

	panic(fmt.Sprintf("%T type objects don't have %s method.", value.Interface(), methodName))
}

// WrapArguments receives a sequence of arguments and wrap each one of them into reflect.Value
func WrapArguments(args ...interface{}) []reflect.Value {
	funcArgs := []reflect.Value{}

	for _, arg := range args {
		value, wrapped := arg.(reflect.Value)

		if wrapped {
			funcArgs = append(funcArgs, value)
			continue
		}

		funcArgs = append(funcArgs, reflect.ValueOf(arg))
	}

	return funcArgs
}

// UnwrapReflectValues unwraps given result(s) from reflect.Value to interface
// If result is an empty slice of reflect.Value, it returns nil
func UnwrapReflectValues(result interface{}) interface{} {
	switch result := result.(type) {
	case []reflect.Value:
		if len(result) == 0 {
			return nil
		} else if len(result) == 1 {
			return result[0].Interface()
		} else {
			values := []interface{}{}

			for _, v := range result {
				values = append(values, v.Interface())
			}

			return values
		}
	default:
		return result
	}
}

func getReflectPtrAndValue(value reflect.Value, rawValue interface{}) (ptr, v reflect.Value) {
	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(rawValue)) // create new pointer
		temp := ptr.Elem()                          // create variable to value of pointer
		temp.Set(value)                             // set value of variable to our passed in value
	}

	return ptr, value
}
