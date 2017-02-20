package evaluator

import (
	"github.com/st0012/rooby/object"
)

func builtinClassMethods(c *object.Class) map[string]*object.BuiltInMethod {
	var bim map[string]*object.BuiltInMethod

	bim = make(map[string]*object.BuiltInMethod)

	bim["new"] = &object.BuiltInMethod{
		Fn: func(args ...object.Object) object.Object {
			instance := &object.BaseObject{Class: c, InstanceVariables: object.NewEnvironment()}

			if instance.RespondTo("initialize") {
				evalInstanceMethod(instance, "initialize", args)
			}

			return instance
		},
		Des: "Initialize class's instance",
	}

	return bim
}
