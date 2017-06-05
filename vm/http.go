package vm

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type httpClass struct {
	*RClass
}

var (
	httpRequestClass  *RClass
	httpResponseClass *RClass
)

func initializeHTTPClass(vm *VM) {
	net := vm.loadConstant("Net", true)
	http := httpClass{initializeClass("HTTP", false)}
	http.setBuiltInMethods(builtinHTTPClassMethods, true)
	http.initializeRequestClass()
	http.initializeResponseClass()

	net.constants[http.Name] = &Pointer{http}
}

func (hc httpClass) initializeRequestClass() *RClass {
	requestClass := initializeClass("Request", false)
	hc.constants["Request"] = &Pointer{requestClass}

	attrs := []string{
		"body",
		"method",
		"path",
		"url",
	}

	requestClass.setAttrAccessor(attrs)

	builtinHTTPRequestInstanceMethods := []*BuiltInMethodObject{}

	requestClass.setBuiltInMethods(builtinHTTPRequestInstanceMethods, false)

	httpRequestClass = requestClass
	return requestClass
}

func (hc httpClass) initializeResponseClass() *RClass {
	responseClass := initializeClass("Response", false)
	hc.constants["Response"] = &Pointer{responseClass}

	attrs := []string{
		"body",
		"status",
		"header",
		"http_version",
		"request_http_version",
		"request",
	}

	responseClass.setAttrAccessor(attrs)

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
