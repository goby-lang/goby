package vm

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var (
	stringClass *RString
)

// RString is the built in string class
type RString struct {
	*BaseClass
}

// StringObject represents string instances
type StringObject struct {
	Class *RString
	Value string
}

func (s *StringObject) objectType() objectType {
	return stringObj
}

func (s *StringObject) Inspect() string {
	return s.Value
}

func (s *StringObject) returnClass() Class {
	if s.Class == nil {
		panic(fmt.Sprintf("String %s doesn't have class.", s.Inspect()))
	}

	return s.Class
}

func (s *StringObject) equal(e *StringObject) bool {
	return s.Value == e.Value
}

var (
	stringTable = make(map[string]*StringObject)
)

func initializeString(value string) *StringObject {
	addr, ok := stringTable[value]

	if !ok {
		s := &StringObject{Value: value, Class: stringClass}
		stringTable[value] = s
		return s
	}

	return addr
}

var builtinStringMethods = []*BuiltInMethod{
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				rightValue := right.Value
				return &StringObject{Value: leftValue + rightValue, Class: stringClass}
			}
		},
		Name: "+",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*IntegerObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				if right.Value < 0 {
					return newError("Second argument must be greater than or equal to 0 String#*")
				}

				var result string

				for i := 0; i < right.Value; i++ {
					result += leftValue
				}

				return &StringObject{Value: result, Class: stringClass}
			}
		},
		Name: "*",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				rightValue := right.Value

				if leftValue > rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: ">",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				rightValue := right.Value

				if leftValue < rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: "<",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				rightValue := right.Value

				if leftValue == rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: "==",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				rightValue := right.Value

				if leftValue < rightValue {
					return initilaizeInteger(-1)
				}
				if leftValue > rightValue {
					return initilaizeInteger(1)
				}

				return initilaizeInteger(0)
			}
		},
		Name: "<=>",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(stringClass)
				}

				rightValue := right.Value

				if leftValue != rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: "!=",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				str := []byte(receiver.(*StringObject).Value)
				start := string(str[0])
				rest := string(str[1:])
				result := strings.ToUpper(start) + strings.ToLower(rest)

				return initializeString(result)
			}
		},
		Name: "capitalize",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				return initializeString(strings.ToUpper(str))
			}
		},
		Name: "upcase",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				return initializeString(strings.ToLower(str))
			}
		},
		Name: "downcase",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				return initilaizeInteger(len(str))
			}
		},
		Name: "size",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				return initilaizeInteger(len(str))
			}
		},
		Name: "length",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value
				var revert string
				for i := len(str) - 1; i >= 0; i-- {
					revert += string(str[i])
				}

				return initializeString(revert)
			}
		},
		Name: "reverse",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				return initializeString(str)
			}
		},
		Name: "to_s",
	},
	{
		Fn: func(receiver Object) builtinMethodBody {
			return func(vm *VM, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value
				parsedStr, err := strconv.ParseInt(str, 10, 0)

				if err == nil {
					return initilaizeInteger(int(parsedStr))
				}

				var digits string
				for _, char := range str {
					if unicode.IsDigit(char) {
						digits += string(char)
					} else {
						break
					}
				}

				if len(digits) > 0 {
					parsedStr, _ = strconv.ParseInt(digits, 10, 0)
					return initilaizeInteger(int(parsedStr))
				}

				return initilaizeInteger(0)
			}
		},
		Name: "to_i",
	},
}

func initString() {
	methods := newEnvironment()

	for _, m := range builtinStringMethods {
		methods.set(m.Name, m)
	}

	bc := &BaseClass{Name: "String", Methods: methods, ClassMethods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass, superClass: objectClass}
	sc := &RString{BaseClass: bc}
	stringClass = sc
}
