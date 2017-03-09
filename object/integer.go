package object

import (
	"fmt"
)

var (
	IntegerClass *RInteger
)

type RInteger struct {
	*BaseClass
}

type IntegerObject struct {
	Class *RInteger
	Value int
}

func (i *IntegerObject) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *IntegerObject) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *IntegerObject) ReturnClass() Class {
	return i.Class
}

func InitilaizeInteger(value int) *IntegerObject {
	return &IntegerObject{Value: value, Class: IntegerClass}
}
