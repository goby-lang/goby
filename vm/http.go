package vm

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/goby-lang/goby/vm/errors"
)

var (
	httpRequestClass  *RClass
	httpResponseClass *RClass
	httpClientClass   *RClass
)

// Class methods --------------------------------------------------------
func builtinHTTPClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Sends a GET request to the target and returns the HTTP response as a string.
			Name: "get",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					arg0, ok := args[0].(*StringObject)
					if !ok {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect argument 0 to be string, got: %s", args[0].Class().Name)
					}

					uri, err := url.Parse(arg0.value)

					if len(args) > 1 {
						var arr []string

						for i, v := range args[1:] {
							argn, ok := v.(*StringObject)
							if !ok {
								return t.vm.initErrorObject(errors.ArgumentError, "Splat arguments must be a string, got: %s for argument %d", v.Class().Name, i)
							}
							arr = append(arr, argn.value)
						}

						uri.Path = path.Join(arr...)
					}

					resp, err := http.Get(uri.String())
					if err != nil {
						return t.vm.initErrorObject(errors.HTTPError, "Could not complete request, %s", err)
					}
					if resp.StatusCode != http.StatusOK {
						return t.vm.initErrorObject(errors.HTTPError, "Non-200 response, %s (%d)", resp.Status, resp.StatusCode)
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
						return t.vm.initErrorObject(errors.ArgumentError, errors.WrongNumberOfArgumentFormat, 3, strconv.Itoa(len(args)))
					}

					arg0, ok := args[0].(*StringObject)
					if !ok {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect argument 0 to be string, got: %s", args[0].Class().Name)
					}
					host := arg0.value

					arg1, ok := args[1].(*StringObject)
					if !ok {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect argument 1 to be string, got: %s", args[0].Class().Name)
					}
					contentType := arg1.value

					arg2, ok := args[2].(*StringObject)
					if !ok {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect argument 0 to be string, got: %s", args[0].Class().Name)
					}
					body := arg2.value

					resp, err := http.Post(host, contentType, strings.NewReader(body))
					if err != nil {
						return t.vm.initErrorObject(errors.HTTPError, "Could not complete request, %s", err)
					}
					if resp.StatusCode != http.StatusOK {
						return t.vm.initErrorObject(errors.HTTPError, "Non-200 response, %s (%d)", resp.Status, resp.StatusCode)
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
						return t.vm.initErrorObject(errors.ArgumentError, "Expect 1 argument. got=%v", strconv.Itoa(len(args)))
					}

					host, ok := args[0].(*StringObject)
					if !ok {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect argument 0 to be string, got: %s", args[0].Class().Name)
					}

					resp, err := http.Head(host.value)
					if err != nil {
						return t.vm.initErrorObject(errors.HTTPError, "Could not complete request, %s", err)
					}
					if resp.StatusCode != http.StatusOK {
						return t.vm.initErrorObject(errors.HTTPError, "Non-200 response, %s (%d)", resp.Status, resp.StatusCode)
					}

					ret := t.vm.initHashObject(map[string]Object{})

					for k, v := range resp.Header {
						ret.Pairs[k] = t.vm.initStringObject(strings.Join(v, " "))
					}

					return ret
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func initHTTPClass(vm *VM) {
	net := vm.loadConstant("Net", true)
	http := vm.initializeClass("HTTP", false)
	http.setBuiltinMethods(builtinHTTPClassMethods(), true)
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
	builtinHTTPRequestInstanceMethods := []*BuiltinMethodObject{}

	requestClass.setBuiltinMethods(builtinHTTPRequestInstanceMethods, false)

	httpRequestClass = requestClass
	return requestClass
}

func initResponseClass(vm *VM, hc *RClass) *RClass {
	responseClass := vm.initializeClass("Response", false)
	hc.setClassConstant(responseClass)
	builtinHTTPResponseInstanceMethods := []*BuiltinMethodObject{}

	responseClass.setBuiltinMethods(builtinHTTPResponseInstanceMethods, false)

	httpResponseClass = responseClass
	return responseClass
}

func initClientClass(vm *VM, hc *RClass) *RClass {
	clientClass := vm.initializeClass("Client", false)
	hc.setClassConstant(clientClass)

	clientClass.setBuiltinMethods(builtinHTTPClientInstanceMethods(), false)

	httpClientClass = clientClass
	return clientClass
}
