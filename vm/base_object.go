package vm

import "fmt"

// BaseObject is an interface that implements basic functions any object requires.
type BaseObject interface {
	returnClass() Class
	Object
}

// RObject represents any non built-in class's instance.
type RObject struct {
	Class             *RClass
	InstanceVariables *environment
	Scope             *scope
	InitializeMethod  *Method
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
