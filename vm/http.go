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
)

// Class methods --------------------------------------------------------
func builtinHTTPClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
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

					uri, err := url.Parse(args[0].(*StringObject).value)
					if err != nil {
						return t.vm.initErrorObject(errors.ArgumentError, err.Error())
					}

					contentType := args[1].(*StringObject).value

					body := args[2].(*StringObject).value

					resp, err := http.Post(uri.String(), contentType, strings.NewReader(body))

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
