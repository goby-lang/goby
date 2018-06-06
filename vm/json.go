package vm

import (
	"encoding/json"
	"strconv"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

type jsonObj map[string]interface{}

const (
	// JSONError is for JSON-specific error
	JSONError = "JSONError"

	cantParseStringAsJSON = "Can't parse string %s as JSON: %s"
)

// Class methods --------------------------------------------------------
func builtinJSONClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "parse",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%v", strconv.Itoa(len(args)))
					}

					j, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
					}

					var obj jsonObj
					var objs []jsonObj

					jsonString := j.value

					err := json.Unmarshal([]byte(jsonString), &obj)

					if err != nil {
						err = json.Unmarshal([]byte(jsonString), &objs)

						if err != nil {
							return t.vm.InitErrorObject(JSONError, sourceLine, cantParseStringAsJSON, jsonString, err.Error())
						}

						var objects []Object

						for _, obj := range objs {
							objects = append(objects, t.vm.convertJSONToHashObj(obj))
						}

						return t.vm.InitArrayObject(objects)
					}

					return t.vm.convertJSONToHashObj(obj)
				}
			},
		},
		{
			Name: "validate",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%v", strconv.Itoa(len(args)))
					}

					j, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
					}

					var obj jsonObj
					var objs []jsonObj

					jsonString := j.value

					err := json.Unmarshal([]byte(jsonString), &obj)

					if err != nil {
						err = json.Unmarshal([]byte(jsonString), &objs)

						if err != nil {
							return FALSE
						}

						return TRUE
					}

					return TRUE
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinJSONInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func initJSONClass(vm *VM) {
	class := vm.initializeClass("JSON", false)
	class.setBuiltinMethods(builtinJSONClassMethods(), true)
	class.setBuiltinMethods(builtinJSONInstanceMethods(), false)
	vm.objectClass.setClassConstant(class)
}

// Polymorphic helper functions -----------------------------------------

func (v *VM) convertJSONToHashObj(j jsonObj) Object {
	objectMap := map[string]Object{}

	for key, jsonValue := range j {
		switch jsonValue := jsonValue.(type) {
		// Array of json objects
		case []map[string]interface{}:
			objs := []Object{}

			for _, value := range jsonValue {
				objs = append(objs, v.convertJSONToHashObj(value))
			}

			objectMap[key] = v.InitArrayObject(objs)
		case []interface{}:
			objs := []Object{}

			for _, elem := range jsonValue {
				switch e := elem.(type) {
				case map[string]interface{}:
					objs = append(objs, v.convertJSONToHashObj(e))
				default:
					objs = append(objs, v.InitObjectFromGoType(e))
				}
			}

			objectMap[key] = v.InitArrayObject(objs)
			// Single json object
		case map[string]interface{}:
			objectMap[key] = v.convertJSONToHashObj(jsonValue)
		case float64:
			// TODO: Find a better way to distinguish between Float & Integer because default GO JSON package
			// TODO: support only for parsing float out regardless of integer or float type data of JSON value
			if jsonValue == float64(int(jsonValue)) {
				objectMap[key] = v.InitIntegerObject(int(jsonValue))
			} else {
				objectMap[key] = v.initFloatObject(jsonValue)
			}
		default:
			objectMap[key] = v.InitObjectFromGoType(jsonValue)
		}
	}

	return v.InitHashObject(objectMap)
}
