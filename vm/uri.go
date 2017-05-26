package vm

import (
	"net/url"
)

func initializeURIClass(vm *VM) {
	uri := initializeClass("URI", true)
	http := initializeClass("HTTP", false)
	https := initializeClass("HTTPS", false)
	uri.constants[http.Name] = &Pointer{http}
	uri.constants[https.Name] = &Pointer{https}

	for _, m := range builtinURIClassMethods {
		uri.ClassMethods.set(m.Name, m)
	}

	vm.constants["URI"] = &Pointer{Target: uri}
}

var builtinURIClassMethods = []*BuiltInMethod{
	{
		Name: "parse",
		Fn: func(receiver Object) builtinMethodBody {
			return func(v *VM, args []Object, blockFrame *callFrame) Object {
				uri := args[0].(*StringObject).Value
				uriModule := v.constants["URI"].Target.(*RClass)
				u, err := url.Parse(uri)

				if err != nil {
					v.returnError(err.Error())
				}

				uriAttrs := map[string]Object{
					"@host":     initializeString(u.Host),
					"@path":     initializeString(u.Path),
					"@port":     initializeString(u.Port()),
					"@query":    initializeString(u.RawQuery),
					"@scheme":   initializeString(u.Scheme),
					"@user":     NULL,
					"@password": NULL,
				}

				if len(u.User.Username()) != 0 {
					uriAttrs["@user"] = initializeString(u.User.Username())
				}

				if p, ok := u.User.Password(); ok {
					uriAttrs["@password"] = initializeString(p)
				}

				var c *RClass

				if u.Scheme == "https" {
					c = uriModule.constants["HTTPS"].Target.(*RClass)
				} else {
					c = uriModule.constants["HTTP"].Target.(*RClass)
				}

				i := c.initializeInstance()

				for varName, value := range uriAttrs {
					i.InstanceVariables.set(varName, value)
				}

				return i
			}
		},
	},
}
