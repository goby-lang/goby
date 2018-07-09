package vm

import (
	"bytes"
	"fmt"

	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/vm/classes"
)

// MethodObject represents methods defined using goby.
type MethodObject struct {
	*BaseObj
	Name           string
	instructionSet *instructionSet
	argc           int
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initMethodClass() *RClass {
	return vm.initializeClass(classes.MethodClass)
}

// Polymorphic helper functions -----------------------------------------
func (m *MethodObject) ToString() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("<Method: %s (%d params)\n>", m.Name, m.argc))
	out.WriteString(m.instructionSet.inspect())

	return out.String()
}

func (m *MethodObject) ToJSON(t *Thread) string {
	return m.ToString()
}

// Value returns method object's string format
func (m *MethodObject) Value() interface{} {
	return m.ToString()
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
	*BaseObj
	Name string
	Fn   builtinMethodBody
}

// Method is a callable function
type Method = func(receiver Object, line int, t *Thread, args []Object) Object

// ExternalBuiltinMethod is a function that builds a BuiltinMethodObject from an external function
func ExternalBuiltinMethod(name string, m Method) *BuiltinMethodObject {
	return &BuiltinMethodObject{
		Name: name,
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, c *normalCallFrame) Object {
			return m(receiver, sourceLine, t, args)
		},
	}
}

type builtinMethodBody func(Object, int, *Thread, []Object, *normalCallFrame) Object

// Polymorphic helper functions -----------------------------------------

// ToString returns the object's name as the string format
func (bim *BuiltinMethodObject) ToString() string {
	return "<BuiltinMethod: " + bim.Name + ">"
}

// ToJSON just delegates to `ToString`
func (bim *BuiltinMethodObject) ToJSON(t *Thread) string {
	return bim.ToString()
}

// Value returns builtin method object's function
func (bim *BuiltinMethodObject) Value() interface{} {
	return bim.Fn
}
