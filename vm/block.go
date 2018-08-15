package vm

import (
	"fmt"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// BlockObject represents an instance of `Block` class.
// In Goby, block literals can be used to define an "anonymous function" by using the `Block` class.
//
// A block literal consists of `do`-`end` and code snippets between them,
// containing optional "block parameters" surrounded by `| |`
// that can be referred to within the block as "block variables".
//
// `Block.new` can take a block literal, returning a "block" object.
//
// You can call `#call` method on the block object to execute the block whenever and wherever you want.
// You can even pass around the block objects across your codebase.
//
// ```ruby
// bl = Block.new do |array|
//   array.reduce do |sum, i|
//     sum + i
//   end
// end
//                       #=> <Block: REPL>
// bl.call([1, 2, 3, 4]) #=> 10
// ```
//
// You can even form a `closure` (note that you can do that without using `Block.new`):
//
// ```ruby
// n = 1
// bl = Block.new do
//   n = n + 1
// end
// #=> <Block: REPL>
// bl.call
// #=> 2
// bl.call
// #=> 3
// bl.call
// #=> 4
// ```
//
type BlockObject struct {
	*BaseObj
	instructionSet *instructionSet
	ep             *normalCallFrame
	self           Object
}

// Class methods --------------------------------------------------------
func builtinBlockClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// @param block literal
			// @return [Block]
			Name: "new",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if blockFrame == nil {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Can't initialize block object without block argument")
				}

				return t.vm.initBlockObject(blockFrame.instructionSet, blockFrame.ep, blockFrame.self)
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinBlockInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Executes the block and returns the result.
			// It can take arbitrary number of arguments and passes them to the block arguments of the block object,
			// keeping the order of the arguments.
			//
			// ```ruby
			// bl = Block.new do |array|
			//   array.reduce do |sum, i|
			//     sum + i
			//   end
			// end
			// #=> <Block: REPL>
			// bl.call([1, 2, 3, 4])     #=> 10
			// ```
			//
			// TODO: should check if the following behavior is OK or not
			// Note that the method does NOT check the number of the arguments and the number of block parameters.
			// * if the number of the arguments exceed, the rest will just be truncated:
			//
			// ```ruby
			// p = Block.new do |i, j, k|
			//   [i, j, k]
			// end
			// p.call(1, 2, 3, 4, 5)     #=> [1, 2, 3]
			// ```
			//
			// * if the number of the block parameters exceeds, the rest will just be filled with `nil`:
			//
			// ```ruby
			// p = Block.new do |i, j, k|
			//   [i, j, k]
			// end
			// p.call                    #=> [nil, nil, nil]
			// ```
			//
			// @param object [Object]...
			// @return [Object]
			Name: "call",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				block := receiver.(*BlockObject)
				c := newNormalCallFrame(block.instructionSet, block.instructionSet.filename, sourceLine)
				c.ep = block.ep
				c.self = block.self
				c.isBlock = true

				return t.builtinMethodYield(c, args...).Target
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initBlockClass() *RClass {
	class := vm.initializeClass(classes.BlockClass)
	class.setBuiltinMethods(builtinBlockClassMethods(), true)
	class.setBuiltinMethods(builtinBlockInstanceMethods(), false)
	return class
}

func (vm *VM) initBlockObject(is *instructionSet, ep *normalCallFrame, self Object) *BlockObject {
	return &BlockObject{
		BaseObj:        &BaseObj{class: vm.TopLevelClass(classes.BlockClass)},
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

// ToString returns the object's name as the string format
func (bo *BlockObject) ToString() string {
	return fmt.Sprintf("<Block: %s>", bo.instructionSet.filename)
}

// Inspect delegates to ToString
func (bo *BlockObject) Inspect() string {
	return bo.ToString()
}

// ToJSON just delegates to ToString
func (bo *BlockObject) ToJSON(t *Thread) string {
	return bo.ToString()
}

// copy returns the duplicate of the Array object
func (bo *BlockObject) copy() Object {
	newC := &BlockObject{BaseObj: &BaseObj{class: bo.class}, instructionSet: bo.instructionSet}
	return newC
}
