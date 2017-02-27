package initializer

import (
	"fmt"
	"github.com/st0012/rooby/ast"
	"github.com/st0012/rooby/object"
)

func InitializeProgram() {
	// Initialize order matters
	initializeClassClass()
	initializeObjectClass()
	initializeNullClass()
	initializeStringClass()
	initializeIntegerClass()
	initializeBooleanClass()
}

func InitializeMainObject() *object.RObject {
	obj := &object.RObject{Class: ObjectClass, InstanceVariables: object.NewEnvironment()}
	scope := &object.Scope{Self: obj, Env: object.NewEnvironment()}
	obj.Scope = scope
	return obj
}

func InitializeClass(name *ast.Constant, scope *object.Scope) *object.RClass {
	class := &object.RClass{BaseClass: &object.BaseClass{Name: name, Methods: object.NewEnvironment(), Class: ClassClass, SuperClass: ObjectClass}}
	classScope := &object.Scope{Self: class, Env: object.NewClosedEnvironment(scope.Env)}
	class.Scope = classScope

	return class
}

func InitializeInstance(c *object.RClass) *object.RObject {
	instance := &object.RObject{Class: c, InstanceVariables: object.NewEnvironment()}

	return instance
}

func checkArgumentLen(args []object.Object, class object.Class, method_name string) *object.Error {
	if len(args) > 1 {
		return &object.Error{Message: fmt.Sprintf("Too many arguments for %s#%s", class.ReturnName().Value, method_name)}
	}

	return nil
}

func wrongTypeError(c object.Class) *object.Error {
	return &object.Error{Message: fmt.Sprintf("expect argument to be %s type", c.ReturnName().Value)}
}
