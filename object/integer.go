package object

import "fmt"

type IntegerClass struct {
	*BaseClass
}

type IntegerObject struct {
	Class *IntegerClass
	Value int64
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