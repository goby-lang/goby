package vm

import (
	"net/url"
	"strconv"

	"github.com/goby-lang/goby/vm/errors"
)

// Class methods --------------------------------------------------------
func builtinURIClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
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
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				uri := args[0].(*StringObject).value
				uriModule := t.vm.topLevelClass("URI")
				u, err := url.Parse(uri)

				if err != nil {
					return t.vm.InitErrorObject(errors.InternalError, sourceLine, err.Error())
				}

				uriAttrs := map[string]Object{
					"@user":     NULL,
					"@password": NULL,
					"@query":    NULL,
					"@path":     t.vm.InitStringObject("/"),
				}

				// Scheme
				uriAttrs["@scheme"] = t.vm.InitStringObject(u.Scheme)

				// Host
				uriAttrs["@host"] = t.vm.InitStringObject(u.Host)

				// Port
				if len(u.Port()) == 0 {
					switch u.Scheme {
					case "http":
						uriAttrs["@port"] = t.vm.InitIntegerObject(80)
					case "https":
						uriAttrs["@port"] = t.vm.InitIntegerObject(443)
					}
				} else {
					p, err := strconv.ParseInt(u.Port(), 0, 64)

					if err != nil {
						return t.vm.InitErrorObject(errors.InternalError, sourceLine, err.Error())
					}

					uriAttrs["@port"] = t.vm.InitIntegerObject(int(p))
				}

				// Path
				if len(u.Path) != 0 {
					uriAttrs["@path"] = t.vm.InitStringObject(u.Path)
				}

				// Query
				if len(u.RawQuery) != 0 {
					uriAttrs["@query"] = t.vm.InitStringObject(u.RawQuery)
				}

				// User
				if u.User != nil {
					if len(u.User.Username()) != 0 {
						uriAttrs["@user"] = t.vm.InitStringObject(u.User.Username())
					}

					if p, ok := u.User.Password(); ok {
						uriAttrs["@password"] = t.vm.InitStringObject(p)
					}
				}

				var c *RClass

				if u.Scheme == "https" {
					c = uriModule.getClassConstant("HTTPS")
				} else {
					c = uriModule.getClassConstant("HTTP")
				}

				i := c.initializeInstance()

				for varName, value := range uriAttrs {
					i.InstanceVariables.set(varName, value)
				}

				return i

			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func initURIClass(vm *VM) {
	uri := vm.initializeModule("URI")
	http := vm.initializeClass("HTTP")
	https := vm.initializeClass("HTTPS")
	https.superClass = http
	https.pseudoSuperClass = http
	uri.setClassConstant(http)
	uri.setClassConstant(https)
	uri.setBuiltinMethods(builtinURIClassMethods(), true)

	attrs := []Object{
		vm.InitStringObject("host"),
		vm.InitStringObject("path"),
		vm.InitStringObject("port"),
		vm.InitStringObject("query"),
		vm.InitStringObject("scheme"),
		vm.InitStringObject("user"),
		vm.InitStringObject("password"),
	}

	http.setAttrReader(attrs)
	http.setAttrWriter(attrs)

	vm.objectClass.setClassConstant(uri)
}
