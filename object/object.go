package object

import (
	"bytes"
	"github.com/st0012/rooby/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ         = "INTEGER"
	STRING_OBJ          = "STRING"
	BOOLEAN_OBJ         = "BOOLEAN"
	NULL_OBJ            = "NULL"
	RETURN_VALUE_OBJ    = "RETURN_VALUE"
	ERROR_OBJ           = "ERROR"
	METHOD_OBJ          = "METHOD"
	CLASS_OBJ           = "CLASS"
	BASE_OBJECT_OBJ     = "BASE_OBJECT"
	BUILD_IN_METHOD_OBJ = "BUILD_IN_METHOD"
	MAIN_OBJ            = "MAIN_OBJECT"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

type Method struct {
	Name       string
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Scope      *Scope
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

type Class struct {
	Name            *ast.Constant
	Scope           *Scope
	InstanceMethods *Environment
	ClassMethods    *Environment
	SuperClass      *Class
}

func (c *Class) Type() ObjectType {
	return CLASS_OBJ
}

func (c *Class) Inspect() string {
	return "<Class:" + c.Name.Value + ">"
}

type Main struct {
	Env *Environment
}

func (m *Main) Type() ObjectType {
	return MAIN_OBJ
}

func (m *Main) Inspect() string {
	return "Main Object"
}

type BuiltInMethod struct {
	Fn  func(args ...Object) Object
	Des string
}

func (bim *BuiltInMethod) Type() ObjectType {
	return BUILD_IN_METHOD_OBJ
}

func (bim *BuiltInMethod) Inspect() string {
	return bim.Des
}

type Scope struct {
	Env  *Environment
	Self Object
}
