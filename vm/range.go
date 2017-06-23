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

var builtInRangeInstanceMethods = []*BuiltInMethodObject{
	{
		Name: "first",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				ran := receiver.(*RangeObject)
				return initIntegerObject(ran.Start)
			}
		},
	},
	{
		Name: "last",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				ran := receiver.(*RangeObject)
				return initIntegerObject(ran.End)
			}
		},
	},
	{
		Name: "each",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				ran := receiver.(*RangeObject)

				if blockFrame == nil {
					t.returnError("Can't yield without a block")
				}

				for i := ran.Start; i <= ran.End; i++ {
					obj := initIntegerObject(i)
					t.builtInMethodYield(blockFrame, obj)
				}
				return ran
			}
		},
	},
	{
		Name: "to_a",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				ran := receiver.(*RangeObject)

				return ran.toArray()
			}
		},
	},
}
