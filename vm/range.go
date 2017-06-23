package vm

import (
	"fmt"
)

var (
	rangeClass *RClass
)

type RangeObject struct {
	Class *RClass
	Start int
	End   int
}

func (ro *RangeObject) toString() string {
	return fmt.Sprintf("(%d..%d)", ro.Start, ro.End)
}

func (ro *RangeObject) toJSON() string {
	return ro.toString()
}

func (ro *RangeObject) returnClass() Class {
	return ro.Class
}

func (ro *RangeObject) toArray() *ArrayObject {
	elems := []Object{}

	for i := ro.Start; i <= ro.End; i++ {
		elems = append(elems, initIntegerObject(i))
	}

	return initArrayObject(elems)
}

func initRangeObject(start, end int) *RangeObject {
	return &RangeObject{Class: rangeClass, Start: start, End: end}
}

func initRangeClass() {
	bc := &BaseClass{Name: "Range", ClassMethods: newEnvironment(), Methods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass, superClass: objectClass}
	rc := &RClass{BaseClass: bc}
	rc.setBuiltInMethods(builtInRangeInstanceMethods, false)
	rangeClass = rc
}

var builtInRangeInstanceMethods = []*BuiltInMethodObject{}
