package vm

import (
	"bytes"
	"fmt"
	"github.com/goby-lang/goby/vm/classes"
)

func (vm *VM) initMethodClass() *RClass {
	return vm.initializeClass(classes.MethodClass, false)
}

type builtinMethodBody func(*thread, []Object, *callFrame) Object

// MethodObject represents methods defined using goby.
type MethodObject struct {
	*baseObj
	Name           string
	instructionSet *instructionSet
	argc           int
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

// BuiltInMethodObject represents methods defined in go.
type BuiltInMethodObject struct {
	*baseObj
	Name string
	Fn   func(receiver Object) builtinMethodBody
}

// Polymorphic helper functions -----------------------------------------
func (bim *BuiltInMethodObject) toString() string {
	return "<BuiltInMethod: " + bim.Name + ">"
}

func (bim *BuiltInMethodObject) toJSON() string {
	return bim.toString()
}
