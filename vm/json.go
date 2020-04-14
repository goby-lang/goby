package vm

import (
	"encoding/json"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

type jsonObj map[string]interface{}

// Class methods --------------------------------------------------------
var builtinJSONClassMethods = []*BuiltinMethodObject{
	{
		Name: "parse",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
			}

			typeErr := t.vm.checkArgTypes(args, sourceLine, classes.StringClass)

			if typeErr != nil {
				return typeErr
			}

			jsonString := args[0].Value().(string)

			var obj jsonObj
			var objs []jsonObj

			err := json.Unmarshal([]byte(jsonString), &obj)

			if err != nil {
				err = json.Unmarshal([]byte(jsonString), &objs)

				if err != nil {
					return t.vm.InitErrorObject(errors.InternalError, sourceLine, "Can't parse string `%s` as json: %s", jsonString, err.Error())
				}

				var objects []Object

				for _, obj := range objs {
					objects = append(objects, t.vm.convertJSONToHashObj(obj))
				}

				return t.vm.InitArrayObject(objects)
			}

			return t.vm.convertJSONToHashObj(obj)

		},
	},
	{
		Name: "validate",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
			}

			typeErr := t.vm.checkArgTypes(args, sourceLine, classes.StringClass)

			if typeErr != nil {
				return typeErr
			}

			jsonString := args[0].Value().(string)

			var obj jsonObj
			var objs []jsonObj

			err := json.Unmarshal([]byte(jsonString), &obj)

			if err != nil {
				err = json.Unmarshal([]byte(jsonString), &objs)

				if err != nil {
					return FALSE
				}

				return TRUE
			}

			return TRUE

		},
	},
}

// Instance methods -----------------------------------------------------
var builtinJSONInstanceMethods = []*BuiltinMethodObject{}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func initJSONClass(vm *VM) {
	class := vm.initializeClass("JSON")
	class.setBuiltinMethods(builtinJSONClassMethods, true)
	class.setBuiltinMethods(builtinJSONInstanceMethods, false)
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
