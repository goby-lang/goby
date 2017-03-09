package object

import (
	"bytes"
	"strings"
)

var (
	ArrayClass *RArray
)

type RArray struct {
	*BaseClass
}

type ArrayObject struct {
	Class    *RArray
	Elements []Object
}

func (a *ArrayObject) Type() ObjectType {
	return ARRAY_OBJ
}

func (a *ArrayObject) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("Array:")
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

func (a *ArrayObject) ReturnClass() Class {
	return a.Class
}

func (a *ArrayObject) Length() int {
	return len(a.Elements)
}

func (a *ArrayObject) Pop() Object {
	value := a.Elements[len(a.Elements)-1]
	a.Elements = a.Elements[:len(a.Elements)-1]
	return value
}

func (a *ArrayObject) Push(objs []Object) *ArrayObject {
	a.Elements = append(a.Elements, objs...)
	return a
}

func InitializeArray(elements []Object) *ArrayObject {
	return &ArrayObject{Elements: elements, Class: ArrayClass}
}
