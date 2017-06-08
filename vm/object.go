package vm

import (
	"fmt"
)

var mainObj *RObject

func init() {
	initTopLevelClasses()
	initNull()
	initBool()
	initInteger()
	initString()
	initArray()
	initHash()
	initMainObj()
}

func initMainObj() {
	obj := &RObject{Class: objectClass, InstanceVariables: newEnvironment()}

	mainObj = obj
}

// Object represents all objects in Goby, including Array, Integer or even Method and Error.
type Object interface {
	returnClass() Class
	Inspect() string
}

// Pointer is used to point to an object. Variables should hold pointer instead of holding a object directly.
type Pointer struct {
	Target Object
}

func (p *Pointer) returnClass() *RClass {
	return p.Target.(*RClass)
}

// RObject represents any non built-in class's instance.
type RObject struct {
	Class             *RClass
	InstanceVariables *environment
	InitializeMethod  *MethodObject
}

// Inspect tells which class it belongs to.
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
