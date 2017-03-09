package object

import (
	"bytes"
	"fmt"
	"github.com/st0012/Rooby/ast"
	"strings"
)

var MainObj *RObject

type ObjectType string

const (
	INTEGER_OBJ         = "INTEGER"
	ARRAY_OBJ           = "ARRAY"
	HASH_OBJ            = "HASH"
	STRING_OBJ          = "STRING"
	BOOLEAN_OBJ         = "BOOLEAN"
	NULL_OBJ            = "NULL"
	RETURN_VALUE_OBJ    = "RETURN_VALUE"
	ERROR_OBJ           = "ERROR"
	METHOD_OBJ          = "METHOD"
	CLASS_OBJ           = "CLASS"
	BASE_OBJECT_OBJ     = "BASE_OBJECT"
	BUILD_IN_METHOD_OBJ = "BUILD_IN_METHOD"
)

func init() {
	initTopLevelClasses()
	initNull()
	initBool()
	initInteger()
	initString()
	initMainObj()
}

func initMainObj() {
	builtInClasses := []Class{StringClass, BooleanClass, IntegerClass}

	obj := &RObject{Class: ObjectClass, InstanceVariables: NewEnvironment()}
	scope := &Scope{Self: obj, Env: NewEnvironment()}

	for _, class := range builtInClasses {
		scope.Env.Set(class.ReturnName(), class)
	}

	obj.Scope = scope
	MainObj = obj
}

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

func (m *Method) ExtendEnv(args []Object) *Environment {
	e := NewClosedEnvironment(m.Scope.Env)

	for i, arg := range args {
		argName := m.Parameters[i].Value
		e.Set(argName, arg)
	}

	return e
}

type BuiltinMethodBody func(...Object) Object

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

type Scope struct {
	Env  *Environment
	Self Object
}

func newError(format string, args ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}

func checkArgumentLen(args []Object, class Class, method_name string) *Error {
	if len(args) > 1 {
		return &Error{Message: fmt.Sprintf("Too many arguments for %s#%s", class.ReturnName(), method_name)}
	}

	return nil
}

func wrongTypeError(c Class) *Error {
	return &Error{Message: fmt.Sprintf("expect argument to be %s type", c.ReturnName())}
}
