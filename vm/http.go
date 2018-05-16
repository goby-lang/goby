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
			// Sends a GET request to the target and returns the HTTP response as a string. Will error on non-200 responses, for more control over http requests look at the `start` method.
			Name: "get",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					arg0, ok := args[0].(*StringObject)
					if !ok {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect argument 0 to be string, got: %s", args[0].Class().Name)
					}

					uri, err := url.Parse(arg0.value)

					if len(args) > 1 {
						var arr []string

						for i, v := range args[1:] {
							argn, ok := v.(*StringObject)
							if !ok {
								return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Splat arguments must be a string, got: %s for argument %d", v.Class().Name, i)
							}
							arr = append(arr, argn.value)
						}

						uri.Path = path.Join(arr...)
					}

					resp, err := http.Get(uri.String())
					if err != nil {
						return t.vm.InitErrorObject(errors.HTTPError, sourceLine, "Could not complete request, %s", err)
					}
					if resp.StatusCode != http.StatusOK {
						return t.vm.InitErrorObject(errors.HTTPError, sourceLine, "Non-200 response, %s (%d)", resp.Status, resp.StatusCode)
					}

					content, err := ioutil.ReadAll(resp.Body)
					resp.Body.Close()

					if err != nil {
						return t.vm.InitErrorObject(errors.InternalError, sourceLine, err.Error())
					}

					return t.vm.initStringObject(string(content))
				}
			},
		}, {
			// Sends a POST request to the target with type header and body. Returns the HTTP response as a string. Will error on non-200 responses, for more control over http requests look at the `start` method.
			Name: "post",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 3 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentFormat, 3, len(args))
					}

					arg0, ok := args[0].(*StringObject)
					if !ok {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect argument 0 to be string, got: %s", args[0].Class().Name)
					}
					host := arg0.value

					arg1, ok := args[1].(*StringObject)
					if !ok {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect argument 1 to be string, got: %s", args[0].Class().Name)
					}
					contentType := arg1.value

					arg2, ok := args[2].(*StringObject)
					if !ok {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect argument 2 to be string, got: %s", args[0].Class().Name)
					}
					body := arg2.value

					resp, err := http.Post(host, contentType, strings.NewReader(body))
					if err != nil {
						return t.vm.InitErrorObject(errors.HTTPError, sourceLine, "Could not complete request, %s", err)
					}
					if resp.StatusCode != http.StatusOK {
						return t.vm.InitErrorObject(errors.HTTPError, sourceLine, "Non-200 response, %s (%d)", resp.Status, resp.StatusCode)
					}

					content, err := ioutil.ReadAll(resp.Body)
					resp.Body.Close()

					if err != nil {
						return t.vm.InitErrorObject(errors.InternalError, sourceLine, err.Error())
					}

					return t.vm.initStringObject(string(content))
				}
			},
		}, {
			// Sends a HEAD request to the target with type header and body. Returns the HTTP headers as a map[string]string. Will error on non-200 responses, for more control over http requests look at the `start` method.
			Name: "head",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {
					arg0, ok := args[0].(*StringObject)
					if !ok {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect argument 0 to be string, got: %s", args[0].Class().Name)
					}

					uri, err := url.Parse(arg0.value)

					if len(args) > 1 {
						var arr []string

						for i, v := range args[1:] {
							argn, ok := v.(*StringObject)
							if !ok {
								return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Splat arguments must be a string, got: %s for argument %d", v.Class().Name, i)
							}
							arr = append(arr, argn.value)
						}

						uri.Path = path.Join(arr...)
					}

					resp, err := http.Head(uri.String())
					if err != nil {
						return t.vm.InitErrorObject(errors.HTTPError, sourceLine, "Could not complete request, %s", err)
					}
					if resp.StatusCode != http.StatusOK {
						return t.vm.InitErrorObject(errors.HTTPError, sourceLine, "Non-200 response, %s (%d)", resp.Status, resp.StatusCode)
					}

					ret := t.vm.InitHashObject(map[string]Object{})

					for k, v := range resp.Header {
						ret.Pairs[k] = t.vm.initStringObject(strings.Join(v, " "))
					}

					return ret
				}
			},
		}, {
			// Starts an HTTP client. This method requires a block which takes a Net::HTTP::Client object. The return value of this method is the last evaluated value of the provided block.
			Name: "start",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *Thread, args []Object, blockFrame *normalCallFrame) Object {

					if len(args) != 0 {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Expect 0 arguments. got=%v", strconv.Itoa(len(args)))
					}

					gobyClient := httpClientClass.initializeInstance()

					result := t.builtinMethodYield(blockFrame, gobyClient)

					if err, ok := result.Target.(*Error); ok {
						return err //an Error object
					}

					return result.Target
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
	vm.mainThread.execGobyLib("net/http/response.gb")
	vm.mainThread.execGobyLib("net/http/request.gb")
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
