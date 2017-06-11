package vm

import (
	"bytes"
	"fmt"
)

var methodClass *RMethod

func init() {
	methods := newEnvironment()

	bc := &BaseClass{Name: "Method", Methods: methods, ClassMethods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass, superClass: objectClass}
	mc := &RMethod{BaseClass: bc}
	methodClass = mc
}

// RMethod represents all method's class. Currently has no methods.
type RMethod struct {
	*BaseClass
}

// MethodObject represents methods defined using goby.
type MethodObject struct {
	class          *RMethod
	Name           string
	instructionSet *instructionSet
	argc           int
}

// Inspect returns method's name, params count and instruction set.
func (m *MethodObject) Inspect() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("<Method: %s (%d params)\n>", m.Name, m.argc))
	out.WriteString(m.instructionSet.inspect())

	return out.String()
}

func (m *MethodObject) toJSON() string {
	return m.Inspect()
}

func (m *MethodObject) returnClass() Class {
	return m.class
}

type builtinMethodBody func(*thread, []Object, *callFrame) Object

// BuiltInMethodObject represents methods defined in go.
type BuiltInMethodObject struct {
	class *RMethod
	Name  string
	Fn    func(receiver Object) builtinMethodBody
}

// Inspect just returns built in method's name.
func (bim *BuiltInMethodObject) Inspect() string {
	return "<BuiltInMethod: " + bim.Name + ">"
}

func (bim *BuiltInMethodObject) toJSON() string {
	return bim.Inspect()
}

func (bim *BuiltInMethodObject) returnClass() Class {
	return bim.class
}
