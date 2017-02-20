package object

import (
	"fmt"
)


type Main struct {
	Env *Environment
}

func (m *Main) Type() ObjectType {
	return MAIN_OBJ
}

func (m *Main) Inspect() string {
	return "Main Object"
}

func InitializeMainObject() *Main {
	env := NewEnvironment()

	for key, value := range builtinGlobalMethods {
		env.Set(key, value)
	}

	return &Main{Env: env}
}

var builtinGlobalMethods = map[string]*BuiltInMethod {
	"puts": &BuiltInMethod{
		Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
}

