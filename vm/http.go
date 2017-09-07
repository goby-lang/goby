package vm

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

var (
	httpRequestClass  *RClass
	httpResponseClass *RClass
	httpClientClass *RClass
)

func initHTTPClass(vm *VM) {
	net := vm.loadConstant("Net", true)
	http := vm.initializeClass("HTTP", false)
	http.setBuiltInMethods(builtinHTTPClassMethods(), true)
	initRequestClass(vm, http)
	initResponseClass(vm, http)
	initClientClass(vm, http)

	net.setClassConstant(http)

	// Use Goby code to extend request and response classes.
	vm.execGobyLib("net/http/response.gb")
	vm.execGobyLib("net/http/request.gb")
	vm.execGobyLib("net/http/client.gb")
}

func initRequestClass(vm *VM, hc *RClass) *RClass {
	requestClass := vm.initializeClass("Request", false)
	hc.setClassConstant(requestClass)
	builtinHTTPRequestInstanceMethods := []*BuiltInMethodObject{}

	requestClass.setBuiltInMethods(builtinHTTPRequestInstanceMethods, false)

	httpRequestClass = requestClass
	return requestClass
}

func initResponseClass(vm *VM, hc *RClass) *RClass {
	responseClass := vm.initializeClass("Response", false)
	hc.setClassConstant(responseClass)
	builtinHTTPResponseInstanceMethods := []*BuiltInMethodObject{}

	responseClass.setBuiltInMethods(builtinHTTPResponseInstanceMethods, true)

	httpResponseClass = responseClass
	return responseClass
}

func initClientClass(vm *VM, hc *RClass) *RClass {
	clientClass := vm.initializeClass("Client", false)
	hc.setClassConstant(clientClass)

	clientClass.setBuiltInMethods(builtinHTTPClientClassMethods(), false)

	httpClientClass = clientClass
	return clientClass
}

func builtinHTTPClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			// Sends a GET request to the target and returns the HTTP response as a string.
			Name: "get",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					uri, err := url.Parse(args[0].(*StringObject).value)

					if len(args) > 1 {
						var arr []string

						for _, v := range args[1:] {
							arr = append(arr, v.(*StringObject).value)
						}

						uri.Path = path.Join(arr...)
					}

					resp, err := http.Get(uri.String())
					if err != nil {
						return t.vm.initErrorObject(HTTPError, err.Error())
					}

					if resp.StatusCode != http.StatusOK {
						return t.vm.initErrorObject(HTTPResponseError, resp.Status)
					}

					content, err := ioutil.ReadAll(resp.Body)
					resp.Body.Close()

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					return t.vm.initStringObject(string(content))
				}
			},
		}, {
			// Sends a POST request to the target with type header and body. Returns the HTTP response as a string.
			Name: "post",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 3 {
						return t.vm.initErrorObject(ArgumentError, "Expect 3 arguments. got=%v", strconv.Itoa(len(args)))
					}

					host := args[0].(*StringObject).value

					contentType := args[1].(*StringObject).value

					body := args[2].(*StringObject).value

					resp, err := http.Post(host, contentType, strings.NewReader(body))
					if err != nil {
						return t.vm.initErrorObject(HTTPError, err.Error())
					}
					if resp.StatusCode != http.StatusOK {
						return t.vm.initErrorObject(HTTPResponseError, resp.Status)
					}

					content, err := ioutil.ReadAll(resp.Body)
					resp.Body.Close()

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					return t.vm.initStringObject(string(content))
				}
			},
		}, {
			// Sends a POST request to the target with type header and body. Returns the HTTP response as a string.
			Name: "head",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(ArgumentError, "Expect 1 arguments. got=%v", strconv.Itoa(len(args)))
					}

					host := args[0].(*StringObject).value

					resp, err := http.Head(host)
					if err != nil {
						return t.vm.initErrorObject(HTTPError, err.Error())
					}
					if resp.StatusCode != http.StatusOK {
						return t.vm.initErrorObject(HTTPResponseError, resp.Status)
					}

					content, err := ioutil.ReadAll(resp.Body)
					resp.Body.Close()

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					return t.vm.initStringObject(string(content))
				}
			},
		},
	}
}
