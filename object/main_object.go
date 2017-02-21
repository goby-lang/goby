package object

import (
	"fmt"
	"github.com/st0012/rooby/ast"
)

var (
	ObjectClass = InitializeObjectClass()
)

func InitializeMainObject() *BaseObject {
	obj := &BaseObject{Class: ObjectClass, InstanceVariables: NewEnvironment()}
	scope := &Scope{Self: obj, Env: NewEnvironment()}
	obj.Scope = scope
	return obj
}

func InitializeObjectClass() *Class {
	instanceMethods := NewEnvironment()
	classMethods := NewEnvironment()

	for key, value := range builtinGlobalMethods {
		instanceMethods.Set(key, value)
		classMethods.Set(key, value)
	}

	name := &ast.Constant{Value: "Object"}
	class := &Class{Name: name, InstanceMethods: instanceMethods, ClassMethods: classMethods}

	return class
}

var builtinGlobalMethods = map[string]*BuiltInMethod{
	"puts": &BuiltInMethod{
		Fn: func(args ...Object) Object {

			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
		Des: "Print arguments",
	},
}
