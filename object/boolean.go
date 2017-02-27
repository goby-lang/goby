package object

import "fmt"

type BooleanClass struct {
	*BaseClass
}

type BooleanObject struct {
	Class *BooleanClass
	Value bool
}

func (b *BooleanObject) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *BooleanObject) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *BooleanObject) ReturnClass() Class {
	return b.Class
}
