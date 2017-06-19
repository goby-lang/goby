package vm

import (
	"bytes"
	"fmt"
)

var methodClass *RMethod

func initMethodClass() {
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

// toString returns method's name, params count and instruction set.
func (m *MethodObject) toString() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("<Method: %s (%d params)\n>", m.Name, m.argc))
	out.WriteString(m.instructionSet.inspect())

	return out.String()
}

func (m *MethodObject) toJSON() string {
	return m.toString()
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

// toString just returns built in method's name.
func (bim *BuiltInMethodObject) toString() string {
	return "<BuiltInMethod: " + bim.Name + ">"
}

func (bim *BuiltInMethodObject) toJSON() string {
	return bim.toString()
}

func (bim *BuiltInMethodObject) returnClass() Class {
	return bim.class
}
