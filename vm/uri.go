package vm

import (
	"net/url"
	"strconv"
)

func initializeURIClass(vm *VM) {
	uri := initializeClass("URI", true)
	http := initializeClass("HTTP", false)
	https := initializeClass("HTTPS", false)
	https.superClass = http
	https.pseudoSuperClass = http
	uri.constants[http.Name] = &Pointer{http}
	uri.constants[https.Name] = &Pointer{https}

	for _, m := range builtinURIClassMethods {
		uri.ClassMethods.set(m.Name, m)
	}

	attrs := []Object{
		initializeString("host"),
		initializeString("path"),
		initializeString("port"),
		initializeString("query"),
		initializeString("scheme"),
		initializeString("user"),
		initializeString("password"),
	}

	http.setAttrReader(attrs)
	http.setAttrWriter(attrs)

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
					"@user":     NULL,
					"@password": NULL,
					"@query":    NULL,
					"@path":     initializeString("/"),
				}

				// Scheme
				uriAttrs["@scheme"] = initializeString(u.Scheme)

				// Host
				uriAttrs["@host"] = initializeString(u.Host)

				// Port
				if len(u.Port()) == 0 {
					switch u.Scheme {
					case "http":
						uriAttrs["@port"] = initilaizeInteger(80)
					case "https":
						uriAttrs["@port"] = initilaizeInteger(443)
					}
				} else {
					p, err := strconv.ParseInt(u.Port(), 0, 64)

					if err != nil {
						v.returnError(err.Error())
					}

					uriAttrs["@port"] = initilaizeInteger(int(p))
				}

				// Path
				if len(u.Path) != 0 {
					uriAttrs["@path"] = initializeString(u.Path)
				}

				// Query
				if len(u.RawQuery) != 0 {
					uriAttrs["@query"] = initializeString(u.RawQuery)
				}

				// User
				if u.User != nil {
					if len(u.User.Username()) != 0 {
						uriAttrs["@user"] = initializeString(u.User.Username())
					}

					if p, ok := u.User.Password(); ok {
						uriAttrs["@password"] = initializeString(p)
					}
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
