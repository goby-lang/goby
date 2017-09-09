package vm

import (
	"github.com/goby-lang/goby/vm/errors"
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
	httpClientClass   *RClass
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

	responseClass.setBuiltInMethods(builtinHTTPResponseInstanceMethods, false)

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
					arg0, ok := args[0].(*StringObject)
					if !ok {
						return t.vm.initErrorObject(errors.ArgumentError, "Argument 0 must be a string")
					}

					uri, err := url.Parse(arg0.value)

					if len(args) > 1 {
						var arr []string

						for _, v := range args[1:] {
							argn, ok := v.(*StringObject)
							if !ok {
								return t.vm.initErrorObject(errors.ArgumentError, "Splat arguments must be a string")
							}
							arr = append(arr, argn.value)
						}

						uri.Path = path.Join(arr...)
					}

					resp, err := http.Get(uri.String())
					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
					}

					if resp.StatusCode != http.StatusOK {
						return t.vm.initErrorObject(errors.InternalError, resp.Status)
					}

					content, err := ioutil.ReadAll(resp.Body)
					resp.Body.Close()

					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
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
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 3 arguments. got=%v", strconv.Itoa(len(args)))
					}

					arg0, ok := args[0].(*StringObject)
					if !ok {
						return t.vm.initErrorObject(errors.ArgumentError, "Argument 0 must be a string")
					}
					host := arg0.value

					arg1, ok := args[1].(*StringObject)
					if !ok {
						return t.vm.initErrorObject(errors.ArgumentError, "Argument 1 must be a string")
					}
					contentType := arg1.value

					arg2, ok := args[2].(*StringObject)
					if !ok {
						return t.vm.initErrorObject(errors.ArgumentError, "Argument 2 must be a string")
					}
					body := arg2.value

					resp, err := http.Post(host, contentType, strings.NewReader(body))
					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
					}
					if resp.StatusCode != http.StatusOK {
						return t.vm.initErrorObject(errors.InternalError, resp.Status)
					}

					content, err := ioutil.ReadAll(resp.Body)
					resp.Body.Close()

					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
					}

					return t.vm.initStringObject(string(content))
				}
			},
		}, {
			// Sends a HEAD request to the target with type header and body. Returns the HTTP response as a string.
			Name: "head",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 1 arguments. got=%v", strconv.Itoa(len(args)))
					}

					host, ok := args[0].(*StringObject)
					if !ok {
						return t.vm.initErrorObject(errors.ArgumentError, "Argument 0 must be a string")
					}

					_, err := http.Head(host.value)
					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
					}

					//TODO: make return value a map of headers
					return t.vm.initStringObject("")
				}
			},
		},
	}
}
