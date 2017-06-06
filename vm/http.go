package vm

import (
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	httpRequestClass  *RClass
	httpResponseClass *RClass
)

func initializeHTTPClass(vm *VM) {
	net := vm.loadConstant("Net", true)
	http := initializeClass("HTTP", false)
	http.setBuiltInMethods(builtinHTTPClassMethods, true)
	initializeRequestClass(http)
	initializeResponseClass(http)

	net.constants[http.Name] = &Pointer{http}

	// Use Goby code to extend request and response classes.
	vm.execGobyLib("net/http/response.gb")
	vm.execGobyLib("net/http/request.gb")
}

func initializeRequestClass(hc *RClass) *RClass {
	requestClass := initializeClass("Request", false)
	hc.constants["Request"] = &Pointer{requestClass}
	builtinHTTPRequestInstanceMethods := []*BuiltInMethodObject{}

	requestClass.setBuiltInMethods(builtinHTTPRequestInstanceMethods, false)

	httpRequestClass = requestClass
	return requestClass
}

func initializeResponseClass(hc *RClass) *RClass {
	responseClass := initializeClass("Response", false)
	hc.constants["Response"] = &Pointer{responseClass}
	builtinHTTPResponseInstanceMethods := []*BuiltInMethodObject{}

	responseClass.setBuiltInMethods(builtinHTTPResponseInstanceMethods, false)

	httpResponseClass = responseClass
	return responseClass
}

var builtinHTTPClassMethods = []*BuiltInMethodObject{
	{
		// Sends a GET request to the target and returns the HTTP response as a string.
		Name: "get",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
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
					t.returnError(err.Error())
				}

				content, err := ioutil.ReadAll(resp.Body)
				resp.Body.Close()

				if err != nil {
					t.returnError(err.Error())
				}

				return initializeString(string(content))
			}
		},
	},
}
