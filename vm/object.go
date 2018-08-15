package vm

import (
	"fmt"
	"strconv"

	"github.com/goby-lang/goby/compiler/bytecode"
)

// Object represents all objects in Goby, including Array, Integer or even Method and Error.
type Object interface {
	Class() *RClass
	Value() interface{}
	SingletonClass() *RClass
	SetSingletonClass(*RClass)
	findMethod(string) Object
	findMethodMissing(bool) Object
	ToString() string
	Inspect() string
	ToJSON(t *Thread) string
	id() int
	InstanceVariableGet(string) (Object, bool)
	InstanceVariableSet(string, Object) Object
	isTruthy() bool
}

// BaseObj ==============================================================

type BaseObj struct {
	class             *RClass
	singletonClass    *RClass
	InstanceVariables *environment
}

func NewBaseObject(v *VM, class string) *BaseObj {
	return &BaseObj{
		class:             v.TopLevelClass(class),
		InstanceVariables: newEnvironment(),
	}
}

// Polymorphic helper functions -----------------------------------------

// Class will return object's class
func (b *BaseObj) Class() *RClass {
	if b.class == nil {
		panic(fmt.Sprint("Object doesn't have class."))
	}
	return b.class
}

// SingletonClass returns object's singleton class
func (b *BaseObj) SingletonClass() *RClass {
	return b.singletonClass
}

// SetSingletonClass sets object's singleton class
func (b *BaseObj) SetSingletonClass(c *RClass) {
	b.singletonClass = c
}

func (b *BaseObj) InstanceVariableGet(name string) (Object, bool) {
	v, ok := b.InstanceVariables.get(name)

	if !ok {
		return NULL, false
	}

	return v, true
}

func (b *BaseObj) InstanceVariableSet(name string, value Object) Object {
	b.InstanceVariables.set(name, value)

	return value
}

func (b *BaseObj) findMethod(methodName string) (method Object) {
	if b.SingletonClass() != nil {
		method = b.SingletonClass().lookupMethod(methodName)
	}

	if method == nil {
		method = b.Class().lookupMethod(methodName)
	}

	return
}

func (b *BaseObj) findMethodMissing(searchAncestor bool) (method Object) {
	methodMissing := "method_missing"

	if b.SingletonClass() != nil {
		method, _ = b.SingletonClass().Methods.get(methodMissing)
	}

	if method == nil {
		method, _ = b.Class().Methods.get(methodMissing)
	}

	if method == nil && searchAncestor {
		method = b.findMethod(methodMissing)
	}

	return
}

func (b *BaseObj) id() int {
	r, e := strconv.ParseInt(fmt.Sprintf("%p", b), 0, 64)
	if e != nil {
		panic(e.Error())
	}
	return int(r)
}

func (b *BaseObj) isTruthy() bool {
	return true
}

// Pointer ==============================================================

// Pointer is used to point to an object. Variables should hold pointer instead of holding a object directly.
type Pointer struct {
	Target      Object
	isNamespace bool
}

func (p *Pointer) returnClass() *RClass {
	return p.Target.(*RClass)
}

// RObject ==============================================================

// RObject represents any non built-in class's instance.
type RObject struct {
	*BaseObj
	InitializeMethod *MethodObject
}

// Polymorphic helper functions -----------------------------------------

// ToString returns the object's name as the string format
func (ro *RObject) ToString() string {
	return "<Instance of: " + ro.class.Name + ">"
}

func (ro *RObject) Inspect() string {
	return ro.ToString()
}

// ToJSON just delegates to ToString
func (ro *RObject) ToJSON(t *Thread) string {
	customToJSONMethod := ro.findMethod("to_json").(*MethodObject)

	if customToJSONMethod != nil {
		t.Stack.Push(&Pointer{Target: ro})
		callObj := newCallObject(ro, customToJSONMethod, t.Stack.pointer, 0, &bytecode.ArgSet{}, nil, customToJSONMethod.instructionSet.instructions[0].sourceLine)
		t.evalMethodObject(callObj)
		result := t.Stack.Pop().Target
		return result.ToString()
	}
	return ro.ToString()
}

// Value returns object's string format
func (ro *RObject) Value() interface{} {
	return ro.ToString()
}
