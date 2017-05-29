package vm

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func initializeHTTPClass(vm *VM) {
	net := initializeClass("Net", true)
	http := initializeClass("HTTP", false)
	http.setBuiltInMethods(builtinHTTPClassMethods, true)
	net.constants[http.Name] = &Pointer{http}
	vm.constants["Net"] = &Pointer{Target: net}
}

var builtinHTTPClassMethods = []*BuiltInMethodObject{
	{
		// Sends a GET request to the target and returns the HTTP response as a string.
		Name: "get",
		Fn: func(receiver Object) builtinMethodBody {
			return func(v *VM, args []Object, blockFrame *callFrame) Object {
				var path string

				domain := args[0].(*StringObject).Value

				if len(args) > 1 {
					path = args[1].(*StringObject).Value
				}

				if !strings.HasPrefix(path, "/") {
					path = "/" + path
				}

				resp, err := http.Get(domain + path)

				if err != nil {
					v.returnError(err.Error())
				}

				content, err := ioutil.ReadAll(resp.Body)
				resp.Body.Close()

				if err != nil {
					v.returnError(err.Error())
				}

				return initializeString(string(content))
			}
		},
	},
}
