package vm

import (
	"net/url"
	"strconv"
)

func initializeURIClass(vm *VM) {
	uri := vm.initializeClass("URI", true)
	http := vm.initializeClass("HTTP", false)
	https := vm.initializeClass("HTTPS", false)
	https.superClass = http
	https.pseudoSuperClass = http
	uri.constants[http.Name] = &Pointer{http}
	uri.constants[https.Name] = &Pointer{https}
	uri.setBuiltInMethods(builtinURIClassMethods, true)

	attrs := []Object{
		vm.initStringObject("host"),
		vm.initStringObject("path"),
		vm.initStringObject("port"),
		vm.initStringObject("query"),
		vm.initStringObject("scheme"),
		vm.initStringObject("user"),
		vm.initStringObject("password"),
	}

	http.setAttrReader(attrs)
	http.setAttrWriter(attrs)

	vm.constants["URI"] = &Pointer{Target: uri}
}

var builtinURIClassMethods = []*BuiltInMethodObject{
	{
		// Returns a Net::HTTP or Net::HTTPS's instance (depends on the url scheme).
		//
		// ```ruby
		// u = URI.parse("https://example.com")
		// u.scheme # => "https"
		// u.host # => "example.com"
		// u.port # => 80
		// u.path # => "/"
		// ```
		Name: "parse",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				uri := args[0].(*StringObject).Value
				uriModule := t.vm.constants["URI"].Target.(*RClass)
				u, err := url.Parse(uri)

				if err != nil {
					t.returnError(err.Error())
				}

				uriAttrs := map[string]Object{
					"@user":     NULL,
					"@password": NULL,
					"@query":    NULL,
					"@path":     t.vm.initStringObject("/"),
				}

				// Scheme
				uriAttrs["@scheme"] = t.vm.initStringObject(u.Scheme)

				// Host
				uriAttrs["@host"] = t.vm.initStringObject(u.Host)

				// Port
				if len(u.Port()) == 0 {
					switch u.Scheme {
					case "http":
						uriAttrs["@port"] = t.vm.initIntegerObject(80)
					case "https":
						uriAttrs["@port"] = t.vm.initIntegerObject(443)
					}
				} else {
					p, err := strconv.ParseInt(u.Port(), 0, 64)

					if err != nil {
						t.returnError(err.Error())
					}

					uriAttrs["@port"] = t.vm.initIntegerObject(int(p))
				}

				// Path
				if len(u.Path) != 0 {
					uriAttrs["@path"] = t.vm.initStringObject(u.Path)
				}

				// Query
				if len(u.RawQuery) != 0 {
					uriAttrs["@query"] = t.vm.initStringObject(u.RawQuery)
				}

				// User
				if u.User != nil {
					if len(u.User.Username()) != 0 {
						uriAttrs["@user"] = t.vm.initStringObject(u.User.Username())
					}

					if p, ok := u.User.Password(); ok {
						uriAttrs["@password"] = t.vm.initStringObject(p)
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
