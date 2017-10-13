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
	argCount    int
	argSet      *bytecode.ArgSet
	blockFrame  *callFrame
}

func newCallObject(receiver Object, method *MethodObject, receiverPtr, argCount int, argSet *bytecode.ArgSet, blockFrame *callFrame) *callObject {
	return &callObject{
		receiver:    receiver,
		method:      method,
		receiverPtr: receiverPtr,
		argCount:    argCount,
		argSet:      argSet,
		blockFrame:  blockFrame,
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
