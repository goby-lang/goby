package vm

import (
	"bytes"
	"fmt"
)

// Method represents methods defined using rooby.
type Method struct {
	Name           string
	instructionSet *instructionSet
	argc           int
	scope          *scope
}

func (m *Method) objectType() objectType {
	return methodObj
}

func (m *Method) Inspect() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("<Method: %s (%d params)\n>", m.Name, m.argc))
	out.WriteString(m.instructionSet.inspect())

	return out.String()
}

type builtinMethodBody func(*VM, []Object, *callFrame) Object

// BuiltInMethod represents methods defined in go.
type BuiltInMethod struct {
	Fn   func(receiver Object) builtinMethodBody
	Name string
}

func (bim *BuiltInMethod) objectType() objectType {
	return buildInMethodObj
}

func (bim *BuiltInMethod) Inspect() string {
	return bim.Name
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
