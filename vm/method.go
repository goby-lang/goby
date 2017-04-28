package vm

import (
	"bytes"
	"github.com/rooby-lang/rooby/ast"
	"strings"
)

// Method represents methods defined using rooby.
type Method struct {
	Name           string
	instructionSet *instructionSet
	argc           int
	parameters     []*ast.Identifier
	body           *ast.BlockStatement
	scope          *scope
}

func (m *Method) objectType() objectType {
	return methodObj
}

func (m *Method) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range m.parameters {
		params = append(params, p.String())
	}

	out.WriteString(m.Name)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(m.body.String())
	out.WriteString("\n}\n")

	return out.String()
}

func (m *Method) extendEnv(args []Object) *environment {
	e := closedEnvironment(m.scope.Env)

	for i, arg := range args {
		argName := m.parameters[i].Value
		e.set(argName, arg)
	}

	return e
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
	vm.start()

	return vm.stack.top()
}
