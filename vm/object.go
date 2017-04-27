package vm

import (
	"fmt"
)

var mainObj *RObject

type objectType string

const (
	integerObj       = "INTEGER"
	arrayObj         = "ARRAY"
	hashObj          = "HASH"
	stringObj        = "STRING"
	booleanObj       = "BOOLEAN"
	nullObj          = "NULL"
	returnValueObj   = "RETURN_VALUE"
	errorObj         = "ERROR"
	methodObj        = "METHOD"
	classObj         = "CLASS"
	baseObject       = "BASE_OBJECT"
	buildInMethodObj = "BUILD_IN_METHOD"
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
	builtInClasses := []Class{stringClass, booleanClass, integerClass, arrayClass, hashClass, nullClass, classClass}

	obj := &RObject{Class: objectClass, InstanceVariables: NewEnvironment()}
	scope := &Scope{Self: obj, Env: NewEnvironment()}

	for _, class := range builtInClasses {
		scope.Env.Set(class.ReturnName(), class)
	}

	obj.Scope = scope
	mainObj = obj
}

type Object interface {
	objectType() objectType
	Inspect() string
}

type Pointer struct {
	Target Object
}

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() objectType {
	return returnValueObj
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
