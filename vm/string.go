package vm

import "fmt"

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

func (s *StringObject) Type() objectType {
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

				err := checkArgumentLen(args, stringClass, "+")

				if err != nil {
					return err
				}

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

				err := checkArgumentLen(args, stringClass, ">")
				if err != nil {
					return err
				}

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

				err := checkArgumentLen(args, stringClass, "<")
				if err != nil {
					return err
				}

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

				err := checkArgumentLen(args, stringClass, "==")

				if err != nil {
					return err
				}

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

				err := checkArgumentLen(args, stringClass, "!=")

				if err != nil {
					return err
				}

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
}

func initString() {
	methods := NewEnvironment()

	for _, m := range builtinStringMethods {
		methods.Set(m.Name, m)
	}

	bc := &BaseClass{Name: "String", Methods: methods, ClassMethods: NewEnvironment(), Class: classClass, SuperClass: objectClass}
	sc := &RString{BaseClass: bc}
	stringClass = sc
}
