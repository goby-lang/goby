package object

import (
	"bytes"
	"strings"
)

type ArrayClass struct {
	*BaseClass
}

type ArrayObject struct {
	Class    *ArrayClass
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
