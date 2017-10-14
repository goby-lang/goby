package vm

import (
	"bytes"
	"fmt"
	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/vm/classes"
)

// MethodObject represents methods defined using goby.
type MethodObject struct {
	*baseObj
	Name           string
	instructionSet *instructionSet
	argc           int
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initMethodClass() *RClass {
	return vm.initializeClass(classes.MethodClass, false)
}

// Polymorphic helper functions -----------------------------------------
func (m *MethodObject) toString() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("<Method: %s (%d params)\n>", m.Name, m.argc))
	out.WriteString(m.instructionSet.inspect())

	return out.String()
}

func (m *MethodObject) toJSON() string {
	return m.toString()
}

// Value returns method object's string format
func (m *MethodObject) Value() interface{} {
	return m.toString()
}

func (m *MethodObject) paramTypes() []int {
	return m.instructionSet.paramTypes.Types()
}

func (m *MethodObject) isSplatArgIncluded() bool {
	for _, argType := range m.paramTypes() {
		if argType == bytecode.SplatArg {
			return true
		}
	}

	return false
}

func (m *MethodObject) isKeywordArgIncluded() bool {
	for _, argType := range m.paramTypes() {
		if argType == bytecode.OptionalKeywordArg || argType == bytecode.RequiredKeywordArg {
			return true
		}
	}

	return false
}

//  BuiltinMethodObject =================================================

// BuiltinMethodObject represents methods defined in go.
type BuiltinMethodObject struct {
	*baseObj
	Name string
	Fn   func(receiver Object) builtinMethodBody
}

type builtinMethodBody func(*thread, []Object, *callFrame) Object

// Polymorphic helper functions -----------------------------------------

// toString returns the object's name as the string format
func (bim *BuiltinMethodObject) toString() string {
	return "<BuiltinMethod: " + bim.Name + ">"
}

// toJSON just delegates to `toString`
func (bim *BuiltinMethodObject) toJSON() string {
	return bim.toString()
}

// Value returns builtin method object's function
func (bim *BuiltinMethodObject) Value() interface{} {
	return bim.Fn
}

type callObject struct {
	receiver    Object
	method      *MethodObject
	receiverPtr int
	argPtr      int
	argCount    int
	argSet      *bytecode.ArgSet
	blockFrame  *callFrame
	// argIndex + argPr == current argument's position
	argIndex     int
	lastArgIndex int
	callFrame    *callFrame
}

func newCallObject(receiver Object, method *MethodObject, receiverPtr, argCount int, argSet *bytecode.ArgSet, blockFrame *callFrame) *callObject {
	cf := newCallFrame(method.instructionSet)
	cf.self = receiver
	cf.blockFrame = blockFrame

	return &callObject{
		receiver:    receiver,
		method:      method,
		receiverPtr: receiverPtr,
		argPtr:      receiverPtr + 1,
		argCount:    argCount,
		argSet:      argSet,
		// This is only for normal/optioned arguments
		lastArgIndex: -1,
		callFrame:    cf,
	}
}

func (co *callObject) InstructionSet() *instructionSet {
	return co.method.instructionSet
}

func (co *callObject) ParamTypes() []int {
	return co.InstructionSet().paramTypes.Types()
}

func (co *callObject) ParamNames() []string {
	return co.InstructionSet().paramTypes.Names()
}

func (co *callObject) ArgTypes() []int {
	return co.argSet.Types()
}

func (co *callObject) MethodName() string {
	return co.method.Name
}

func (co *callObject) assignNormalArguments(stack []*Pointer) {
	for i, paramType := range co.ParamTypes() {
		if paramType == bytecode.NormalArg {
			co.callFrame.insertLCL(i, 0, stack[co.argPosition()].Target)
			co.argIndex++
		}
	}
}

func (co *callObject) assignNormalAndOptionedArguments(paramIndex int, stack []*Pointer) {
	/*
		Find first usable value as normal argument, for example:

		```ruby
		  def foo(x, y:); end

		  foo(y: 100, 10)
		```

		In the example we can see that 'x' is the first parameter,
		but in the method call it's the second argument.

		This loop is for skipping other types of arguments and get the correct argument index.
	*/
	for argIndex, at := range co.ArgTypes() {
		if co.lastArgIndex < argIndex && (at == bytecode.NormalArg || at == bytecode.OptionedArg) {
			co.callFrame.insertLCL(paramIndex, 0, stack[co.argPtr+argIndex].Target)

			// Store latest index value (and compare them to current argument index)
			// This is to make sure we won't get same argument's index twice.
			co.lastArgIndex = argIndex
			break
		}
	}
}

func (co *callObject) assignKeywordArguments(paramIndex int, stack []*Pointer) (paramName string, success bool) {
	paramName = co.ParamNames()[paramIndex]
	argIndex := co.argSet.FindIndex(paramName)

	if argIndex != -1 {
		co.callFrame.insertLCL(paramIndex, 0, stack[co.argPtr+argIndex].Target)
		success = true
	}

	return
}

func (co *callObject) assignSplatArgument(stack []*Pointer, arr *ArrayObject) {
	index := len(co.ParamTypes()) - 1

	for co.argIndex < co.argCount {
		arr.Elements = append(arr.Elements, stack[co.argPosition()].Target)
		co.argIndex++
	}

	co.callFrame.insertLCL(index, 0, arr)
}

func (co *callObject) minimumArgNumber() (n int) {
	for _, at := range co.ParamTypes() {
		if at == bytecode.NormalArg {
			n++
		}
	}

	return
}

func (co *callObject) argPosition() int {
	return co.argPtr + co.argIndex
}
