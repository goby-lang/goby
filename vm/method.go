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

func (m *MethodObject) objectType() objectType {
	return methodObj
}

// Inspect returns method's name, params count and instruction set.
func (m *MethodObject) Inspect() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("<Method: %s (%d params)\n>", m.Name, m.argc))
	out.WriteString(m.instructionSet.inspect())

	return out.String()
}

func (m *MethodObject) returnClass() Class {
	return m.class
}

type builtinMethodBody func(*VM, []Object, *callFrame) Object

// BuiltInMethodObject represents methods defined in go.
type BuiltInMethodObject struct {
	class *RMethod
	Name  string
	Fn    func(receiver Object) builtinMethodBody
}

func (bim *BuiltInMethodObject) objectType() objectType {
	return buildInMethodObj
}

// Inspect just returns built in method's name.
func (bim *BuiltInMethodObject) Inspect() string {
	return bim.Name
}

func (bim *BuiltInMethodObject) returnClass() Class {
	return bim.class
}

// builtInMethodYield is like invokeblock instruction for built in methods
func builtInMethodYield(vm *VM, blockFrame *callFrame, args ...Object) *Pointer {
	c := newCallFrame(blockFrame.instructionSet)
	c.blockFrame = blockFrame
	c.ep = blockFrame.ep
	c.self = blockFrame.self

	for i := 0; i < len(args); i++ {
		c.locals[0] = &Pointer{args[i]}
	}

	vm.callFrameStack.push(c)
	vm.startFromTopFrame()

	return vm.stack.top()
}
