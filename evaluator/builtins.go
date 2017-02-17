package evaluator

import (
	"github.com/st0012/rooby/object"
)

func builtins(c *object.Class) map[string]*object.BuiltInMethod {
	var bis map[string]*object.BuiltInMethod

	bis = make(map[string]*object.BuiltInMethod)

	bis["new"] = &object.BuiltInMethod{
		Fn: func(args ...object.Object) object.Object {
			return &object.BaseObject{Class: c, InstanceVariables: object.NewEnvironment()}
		},
		Des: "Initialize class's instance",
	}

	return bis
}
