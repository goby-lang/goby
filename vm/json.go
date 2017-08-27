package vm

import (
	"encoding/json"
	"strconv"
)

type jsonObj map[string]interface{}

func initJSONClass(vm *VM) {
	class := vm.initializeClass("JSON", false)
	class.setBuiltInMethods(builtInJSONClassMethods(), true)
	class.setBuiltInMethods(builtInJSONInstanceMethods(), false)
	vm.objectClass.setClassConstant(class)
}

func builtInJSONClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "parse",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(ArgumentError, "Expect 1 argument. got=%v", strconv.Itoa(len(args)))
					}

					j, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, stringClass, args[0].Class().Name)
					}

					var obj jsonObj
					var objs []jsonObj

					jsonString := j.value

					err := json.Unmarshal([]byte(jsonString), &obj)

					if err != nil {
						err = json.Unmarshal([]byte(jsonString), &objs)

						if err != nil {
							return t.vm.initErrorObject(InternalError, "Can't parse string %s as json: %s", jsonString, err.Error())
						}

						var objects []Object

						for _, obj := range objs {
							objects = append(objects, t.vm.convertJSONToHashObj(obj))
						}

						return t.vm.initArrayObject(objects)
					}

					return t.vm.convertJSONToHashObj(obj)
				}
			},
		},
		{
			Name: "validate",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(ArgumentError, "Expect 1 argument. got=%v", strconv.Itoa(len(args)))
					}

					j, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, stringClass, args[0].Class().Name)
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

func builtInJSONInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{}
}

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

			objectMap[key] = v.initArrayObject(objs)
		case []interface{}:
			objs := []Object{}

			for _, elem := range jsonValue {
				switch e := elem.(type) {
				case map[string]interface{}:
					objs = append(objs, v.convertJSONToHashObj(e))
				default:
					objs = append(objs, v.initObjectFromGoType(e))
				}
			}

			objectMap[key] = v.initArrayObject(objs)
		// Single json object
		case map[string]interface{}:
			objectMap[key] = v.convertJSONToHashObj(jsonValue)
		default:
			objectMap[key] = v.initObjectFromGoType(jsonValue)
		}
	}

	return v.initHashObject(objectMap)
}
