package vm

import (
	"fmt"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

type BlockObject struct {
	*baseObj
	instructionSet *instructionSet
	ep             *normalCallFrame
	self           Object
}

// Class methods --------------------------------------------------------
func builtinBlockClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if blockFrame == nil {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Can't initialize block object without block argument")
					}

					return t.vm.initBlockObject(blockFrame.instructionSet, blockFrame.ep, blockFrame.self)
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinBlockInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "call",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					block := receiver.(*BlockObject)
					c := newNormalCallFrame(block.instructionSet, block.instructionSet.filename, sourceLine)
					c.ep = block.ep
					c.self = block.self
					c.isBlock = true

					return t.builtinMethodYield(c, args...).Target
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initBlockClass() *RClass {
	class := vm.initializeClass(classes.BlockClass, false)
	class.setBuiltinMethods(builtinBlockClassMethods(), true)
	class.setBuiltinMethods(builtinBlockInstanceMethods(), false)
	return class
}

func (vm *VM) initBlockObject(is *instructionSet, ep *normalCallFrame, self Object) *BlockObject {
	return &BlockObject{
		baseObj:        &baseObj{class: vm.topLevelClass(classes.BlockClass)},
		instructionSet: is,
		ep:             ep,
		self:           self,
	}
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (bo *BlockObject) Value() interface{} {
	return bo.instructionSet
}

// toString returns the object's name as the string format
func (bo *BlockObject) toString() string {
	return fmt.Sprintf("<Block: %s>", bo.instructionSet.filename)
}

// toJSON just delegates to toString
func (bo *BlockObject) toJSON(t *Thread) string {
	return bo.toString()
}

// copy returns the duplicate of the Array object
func (bo *BlockObject) copy() Object {
	newC := &BlockObject{baseObj: &baseObj{class: bo.class}, instructionSet: bo.instructionSet}
	return newC
}
