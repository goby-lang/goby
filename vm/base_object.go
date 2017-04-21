package vm

import "fmt"

// BaseObject is an interface that implements basic functions any object requires.
type BaseObject interface {
	ReturnClass() Class
	Object
}

// RObject represents any non built-in class's instance.
type RObject struct {
	Class             *RClass
	InstanceVariables *Environment
	Scope             *Scope
	InitializeMethod  *Method
}

func (ro *RObject) Type() ObjectType {
	return BASE_OBJECT_OBJ
}

func (ro *RObject) Inspect() string {
	return "<Instance of: " + ro.Class.Name + ">"
}

// ReturnClass will return object's class
func (ro *RObject) ReturnClass() Class {
	if ro.Class == nil {
		panic(fmt.Sprintf("Object %s doesn't have class.", ro.Inspect()))
	}
	return ro.Class
}

