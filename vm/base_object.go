package vm

import "fmt"

type BaseObject interface {
	ReturnClass() Class
	Object
}

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

func (ro *RObject) ReturnClass() Class {
	if ro.Class == nil {
		panic(fmt.Sprintf("Object %s doesn't have class.", ro.Inspect()))
	}
	return ro.Class
}

func InitializeInstance(c *RClass) *RObject {
	instance := &RObject{Class: c, InstanceVariables: NewEnvironment()}

	return instance
}
