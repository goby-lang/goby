package vm

import (
	"bytes"
	"fmt"

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

func (m *MethodObject) argTypes() []int {
	return m.instructionSet.argTypes
}

func (m *MethodObject) lastArgType() int {
	if len(m.argTypes()) > 0 {
		return m.argTypes()[len(m.argTypes())-1]
	}

	return -1
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

// Returns the object's name as the string format
func (bim *BuiltinMethodObject) toString() string {
	return "<BuiltinMethod: " + bim.Name + ">"
}

// Alias of toString
func (bim *BuiltinMethodObject) toJSON() string {
	return bim.toString()
}
