package vm

import (
	"bytes"
	"github.com/st0012/Rooby/ast"
	"strings"
)

type Method struct {
	Name           string
	InstructionSet *InstructionSet
	Argc           int
	Parameters     []*ast.Identifier
	Body           *ast.BlockStatement
	Scope          *Scope
}

func (m *Method) Type() ObjectType {
	return METHOD_OBJ
}

func (m *Method) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range m.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(m.Name)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(m.Body.String())
	out.WriteString("\n}\n")

	return out.String()
}

func (m *Method) ExtendEnv(args []Object) *Environment {
	e := NewClosedEnvironment(m.Scope.Env)

	for i, arg := range args {
		argName := m.Parameters[i].Value
		e.Set(argName, arg)
	}

	return e
}

type BuiltinMethodBody func([]Object, *Method) Object

type BuiltInMethod struct {
	Fn   func(receiver Object) BuiltinMethodBody
	Name string
}

func (bim *BuiltInMethod) Type() ObjectType {
	return BUILD_IN_METHOD_OBJ
}

func (bim *BuiltInMethod) Inspect() string {
	return bim.Name
}
