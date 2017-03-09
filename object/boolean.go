package object

import (
	"fmt"
)

var (
	BooleanClass *RBool
	TRUE         *BooleanObject
	FALSE        *BooleanObject
)

type RBool struct {
	*BaseClass
}

type BooleanObject struct {
	Class *RBool
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
