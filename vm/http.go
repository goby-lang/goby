package vm

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/gooby-lang/gooby/vm/errors"
)

const (
	invalidSplatArgument    = "Splat arguments must be a string, got: %s on argument #%d"
	couldNotCompleteRequest = "Could not complete request, %s"
	non200Response          = "Non-200 response, %s (%d)"
)

var (
	httpRequestClass  *RClass
	httpResponseClass *RClass
	httpClientClass   *RClass
)

// Class methods --------------------------------------------------------
var builtinHTTPClassMethods = []*BuiltinMethodObject{
	{
		// Sends a GET request to the target and returns the HTTP response as a string. Will error on non-200 responses, for more control over http requests look at the `start` method.
		Name: "get",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			arg0, ok := args[0].(*StringObject)
			if !ok {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongArgumentTypeFormatNum, 0, "String", args[0].Class().Name)
			}

			uri, err := url.Parse(arg0.value)

			if len(args) > 1 {
				var arr []string

				for i, v := range args[1:] {
					argn, ok := v.(*StringObject)
					if !ok {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, invalidSplatArgument, v.Class().Name, i)
					}
					arr = append(arr, argn.value)
				}

				uri.Path = path.Join(arr...)
			}

			resp, err := http.Get(uri.String())
			if err != nil {
				return t.vm.InitErrorObject(errors.HTTPError, sourceLine, couldNotCompleteRequest, err)
			}
			if resp.StatusCode != http.StatusOK {
				return t.vm.InitErrorObject(errors.HTTPError, sourceLine, non200Response, resp.Status, resp.StatusCode)
			}

			content, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()

			if err != nil {
				return t.vm.InitErrorObject(errors.InternalError, sourceLine, err.Error())
			}

			return t.vm.InitStringObject(string(content))

		},
	}, {
		// Sends a POST request to the target with type header and body. Returns the HTTP response as a string. Will error on non-200 responses, for more control over http requests look at the `start` method.
		Name: "post",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 3 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 3, len(args))
			}

			arg0, ok := args[0].(*StringObject)
			if !ok {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongArgumentTypeFormatNum, 0, "String", args[0].Class().Name)
			}
			host := arg0.value

			arg1, ok := args[1].(*StringObject)
			if !ok {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongArgumentTypeFormatNum, 1, "String", args[0].Class().Name)
			}
			contentType := arg1.value

			arg2, ok := args[2].(*StringObject)
			if !ok {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongArgumentTypeFormatNum, 2, "String", args[0].Class().Name)
			}
			body := arg2.value

			resp, err := http.Post(host, contentType, strings.NewReader(body))
			if err != nil {
				return t.vm.InitErrorObject(errors.HTTPError, sourceLine, couldNotCompleteRequest, err)
			}
			if resp.StatusCode != http.StatusOK {
				return t.vm.InitErrorObject(errors.HTTPError, sourceLine, non200Response, resp.Status, resp.StatusCode)
			}

			content, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()

			if err != nil {
				return t.vm.InitErrorObject(errors.InternalError, sourceLine, err.Error())
			}

			return t.vm.InitStringObject(string(content))

		},
	}, {
		// Sends a HEAD request to the target with type header and body. Returns the HTTP headers as a map[string]string. Will error on non-200 responses, for more control over http requests look at the `start` method.
		Name: "head",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			arg0, ok := args[0].(*StringObject)
			if !ok {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongArgumentTypeFormatNum, 0, "String", args[0].Class().Name)
			}

			uri, err := url.Parse(arg0.value)

			if len(args) > 1 {
				var arr []string

				for i, v := range args[1:] {
					argn, ok := v.(*StringObject)
					if !ok {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, invalidSplatArgument, v.Class().Name, i)
					}
					arr = append(arr, argn.value)
				}

				uri.Path = path.Join(arr...)
			}

			resp, err := http.Head(uri.String())
			if err != nil {
				return t.vm.InitErrorObject(errors.HTTPError, sourceLine, couldNotCompleteRequest, err)
			}
			if resp.StatusCode != http.StatusOK {
				return t.vm.InitErrorObject(errors.HTTPError, sourceLine, non200Response, resp.Status, resp.StatusCode)
			}

			ret := t.vm.InitHashObject(map[string]Object{})

			for k, v := range resp.Header {
				ret.Pairs[k] = t.vm.InitStringObject(strings.Join(v, " "))
			}

			return ret

		},
	}, {
		// Starts an HTTP client. This method requires a block which takes a Net::HTTP::Client object. The return value of this method is the last evaluated value of the provided block.
		Name: "start",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 0 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
			}

			goobyClient := httpClientClass.initializeInstance()

			result := t.builtinMethodYield(blockFrame, goobyClient)

			if err, ok := result.Target.(*Error); ok {
				return err //an Error object
			}

			return result.Target

		},
	},
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func initHTTPClass(vm *VM) {
	net := vm.loadConstant("Net", true)
	http := vm.initializeClass("HTTP")
	http.setBuiltinMethods(builtinHTTPClassMethods, true)
	initRequestClass(vm, http)
	initResponseClass(vm, http)
	initClientClass(vm, http)

	net.setClassConstant(http)

	// Use Gooby code to extend request and response classes.
	vm.mainThread.execGoobyLib("net/http/response.gb")
	vm.mainThread.execGoobyLib("net/http/request.gb")
}

func initRequestClass(vm *VM, hc *RClass) *RClass {
	requestClass := vm.initializeClass("Request")
	hc.setClassConstant(requestClass)
	builtinHTTPRequestInstanceMethods := []*BuiltinMethodObject{}

	requestClass.setBuiltinMethods(builtinHTTPRequestInstanceMethods, false)

	httpRequestClass = requestClass
	return requestClass
}

func initResponseClass(vm *VM, hc *RClass) *RClass {
	responseClass := vm.initializeClass("Response")
	hc.setClassConstant(responseClass)
	builtinHTTPResponseInstanceMethods := []*BuiltinMethodObject{}

	responseClass.setBuiltinMethods(builtinHTTPResponseInstanceMethods, false)

	httpResponseClass = responseClass
	return responseClass
}
