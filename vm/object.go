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
	obj := &RObject{Class: objectClass, InstanceVariables: newEnvironment()}

	mainObj = obj
}

type Object interface {
	returnClass() Class
	objectType() objectType
	Inspect() string
}

type Pointer struct {
	Target Object
}

// RObject represents any non built-in class's instance.
type RObject struct {
	Class             *RClass
	InstanceVariables *environment
	InitializeMethod  *MethodObject
}

func (ro *RObject) objectType() objectType {
	return baseObject
}

func (ro *RObject) Inspect() string {
	return "<Instance of: " + ro.Class.Name + ">"
}

// returnClass will return object's class
func (ro *RObject) returnClass() Class {
	if ro.Class == nil {
		panic(fmt.Sprintf("Object %s doesn't have class.", ro.Inspect()))
	}
	return ro.Class
}

func checkArgumentLen(args []Object, class Class, methodName string) *Error {
	if len(args) > 1 {
		return &Error{Message: fmt.Sprintf("Too many arguments for %s#%s", class.ReturnName(), methodName)}
	}

	return nil
}

func wrongTypeError(c Class) *Error {
	return &Error{Message: fmt.Sprintf("expect argument to be %s type", c.ReturnName())}
}
