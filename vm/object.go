package vm

import (
	"fmt"
)

var MainObj *RObject

type ObjectType string

const (
	INTEGER_OBJ         = "INTEGER"
	ARRAY_OBJ           = "ARRAY"
	HASH_OBJ            = "HASH"
	STRING_OBJ          = "STRING"
	BOOLEAN_OBJ         = "BOOLEAN"
	NULL_OBJ            = "NULL"
	RETURN_VALUE_OBJ    = "RETURN_VALUE"
	ERROR_OBJ           = "ERROR"
	METHOD_OBJ          = "METHOD"
	CLASS_OBJ           = "CLASS"
	BASE_OBJECT_OBJ     = "BASE_OBJECT"
	BUILD_IN_METHOD_OBJ = "BUILD_IN_METHOD"
)

func init() {
	initTopLevelClasses()
	initNull()
	initBool()
	initInteger()
	initString()
	initMainObj()
}

func initMainObj() {
	builtInClasses := []Class{StringClass, booleanClass, IntegerClass, arrayClass, HashClass, NullClass, classClass}

	obj := &RObject{Class: objectClass, InstanceVariables: NewEnvironment()}
	scope := &Scope{Self: obj, Env: NewEnvironment()}

	for _, class := range builtInClasses {
		scope.Env.Set(class.ReturnName(), class)
	}

	obj.Scope = scope
	MainObj = obj
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Pointer struct {
	Target Object
}

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

func checkArgumentLen(args []Object, class Class, method_name string) *Error {
	if len(args) > 1 {
		return &Error{Message: fmt.Sprintf("Too many arguments for %s#%s", class.ReturnName(), method_name)}
	}

	return nil
}

func wrongTypeError(c Class) *Error {
	return &Error{Message: fmt.Sprintf("expect argument to be %s type", c.ReturnName())}
}
