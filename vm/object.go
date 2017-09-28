package vm

import (
	"fmt"
	"strconv"
)

// Object represents all objects in Goby, including Array, Integer or even Method and Error.
type Object interface {
	Class() *RClass
	Value() interface{}
	SingletonClass() *RClass
	SetSingletonClass(*RClass)
	findMethod(string) Object
	toString() string
	toJSON() string
	id() int
	instanceVariableGet(string) (Object, bool)
	instanceVariableSet(string, Object) Object
}

// baseObj ==============================================================

type baseObj struct {
	class             *RClass
	singletonClass    *RClass
	InstanceVariables *environment
}

// Polymorphic helper functions -----------------------------------------

// Class will return object's class
func (b *baseObj) Class() *RClass {
	if b.class == nil {
		panic(fmt.Sprint("Object doesn't have class."))
	}
	return b.class
}

// SingletonClass returns object's singleton class
func (b *baseObj) SingletonClass() *RClass {
	return b.singletonClass
}

// SetSingletonClass sets object's singleton class
func (b *baseObj) SetSingletonClass(c *RClass) {
	b.singletonClass = c
}

func (b *baseObj) instanceVariableGet(name string) (Object, bool) {
	v, ok := b.InstanceVariables.get(name)

	if !ok {
		return NULL, false
	}

	return v, true
}

func (b *baseObj) instanceVariableSet(name string, value Object) Object {
	b.InstanceVariables.set(name, value)

	return value
}

func (b *baseObj) findMethod(methodName string) (method Object) {
	if b.SingletonClass() != nil {
		method = b.SingletonClass().lookupMethod(methodName)
	}

	if method == nil {
		method = b.Class().lookupMethod(methodName)
	}

	return
}

func (b *baseObj) id() int {
	r, e := strconv.ParseInt(fmt.Sprintf("%p", b), 0, 64)
	if e != nil {
		panic(e.Error())
	}
	return int(r)
}

// Pointer ==============================================================

// Pointer is used to point to an object. Variables should hold pointer instead of holding a object directly.
type Pointer struct {
	Target      Object
	isNamespace bool
}

func (p *Pointer) returnClass() *RClass {
	return p.Target.(*RClass)
}

// RObject ==============================================================

// RObject represents any non built-in class's instance.
type RObject struct {
	*baseObj
	InitializeMethod *MethodObject
}

// Polymorphic helper functions -----------------------------------------

// toString returns the object's name as the string format
func (ro *RObject) toString() string {
	return "<Instance of: " + ro.class.Name + ">"
}

// toJSON just delegates to toString
func (ro *RObject) toJSON() string {
	return ro.toString()
}

// Value returns object's string format
func (ro *RObject) Value() interface{} {
	return ro.toString()
}
