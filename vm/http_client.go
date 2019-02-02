package vm

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gooby-lang/gooby/vm/classes"
	"github.com/gooby-lang/gooby/vm/errors"
)

// Instance methods --------------------------------------------------------

func builtinHTTPClientInstanceMethods() []*BuiltinMethodObject {
	//TODO: cookie jar and mutable client
	goClient := http.DefaultClient

	return []*BuiltinMethodObject{
		{
			// Sends a GET request to the target and returns a `Net::HTTP::Response` object.
			Name: "get",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
				}

				u, ok := args[0].(*StringObject)
				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, u.Class().Name)
				}

				resp, err := goClient.Get(u.value)
				if err != nil {
					return t.vm.InitErrorObject(errors.HTTPError, sourceLine, couldNotCompleteRequest, err)
				}

				goobyResp, err := responseGoToGooby(t, resp)
				if err != nil {
					return t.vm.InitErrorObject(errors.InternalError, sourceLine, err.Error())
				}

				return goobyResp

			},
		}, {
			// Sends a POST request to the target and returns a `Net::HTTP::Response` object.
			Name: "post",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 3 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 3, len(args))
				}

				u, ok := args[0].(*StringObject)
				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, u.Class().Name)
				}

				contentType, ok := args[1].(*StringObject)
				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, u.Class().Name)
				}

				body, ok := args[2].(*StringObject)
				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, u.Class().Name)
				}

				bodyR := strings.NewReader(body.value)

				resp, err := goClient.Post(u.value, contentType.value, bodyR)
				if err != nil {
					return t.vm.InitErrorObject(errors.HTTPError, sourceLine, "Could not complete request, %s", err)
				}

				goobyResp, err := responseGoToGooby(t, resp)
				if err != nil {
					return t.vm.InitErrorObject(errors.InternalError, sourceLine, err.Error())
				}

				return goobyResp

			},
		}, {
			// Sends a HEAD request to the target and returns a `Net::HTTP::Response` object.
			Name: "head",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
				}

				u, ok := args[0].(*StringObject)
				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, u.Class().Name)
				}

				resp, err := goClient.Head(u.value)
				if err != nil {
					return t.vm.InitErrorObject(errors.HTTPError, sourceLine, couldNotCompleteRequest, err)
				}

				goobyResp, err := responseGoToGooby(t, resp)
				if err != nil {
					return t.vm.InitErrorObject(errors.InternalError, sourceLine, err.Error())
				}

				return goobyResp

			},
		}, {
			// Returns a blank `Net::HTTP::Request` object to be sent with the`exec` method
			Name: "request",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				return httpRequestClass.initializeInstance()

			},
		}, {
			// Sends a passed `Net::HTTP::Request` object and returns a `Net::HTTP::Response` object
			Name: "exec",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
				}

				if args[0].Class().Name != httpRequestClass.Name {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, "HTTP Response", args[0].Class().Name)
				}

				goReq, err := requestGoobyToGo(args[0])
				if err != nil {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, err.Error())
				}

				goResp, err := goClient.Do(goReq)
				if err != nil {
					return t.vm.InitErrorObject(errors.HTTPError, sourceLine, couldNotCompleteRequest, err)
				}

				goobyResp, err := responseGoToGooby(t, goResp)

				if err != nil {
					return t.vm.InitErrorObject(errors.InternalError, sourceLine, err.Error())
				}

				return goobyResp

			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func initClientClass(vm *VM, hc *RClass) *RClass {
	clientClass := vm.initializeClass("Client")
	hc.setClassConstant(clientClass)

	clientClass.setBuiltinMethods(builtinHTTPClientInstanceMethods(), false)

	httpClientClass = clientClass
	return clientClass
}

// Other helper functions -----------------------------------------------

func requestGoobyToGo(goobyReq Object) (*http.Request, error) {
	//:method, :protocol, :body, :content_length, :transfer_encoding, :host, :path, :url, :params
	uObj, ok := goobyReq.InstanceVariableGet("@url")
	if !ok {
		return nil, fmt.Errorf("could not get url")
	}

	u := uObj.(*StringObject).value

	methodObj, ok := goobyReq.InstanceVariableGet("@method")
	if !ok {
		return nil, fmt.Errorf("could not get method")
	}

	method := methodObj.(*StringObject).value

	var body string
	if !(method == "GET" || method == "HEAD") {
		bodyObj, ok := goobyReq.InstanceVariableGet("@body")
		if !ok {
			return nil, fmt.Errorf("could not get body")
		}

		body = bodyObj.(*StringObject).value
	}

	return http.NewRequest(method, u, strings.NewReader(body))

}

func responseGoToGooby(t *Thread, goResp *http.Response) (Object, error) {
	goobyResp := httpResponseClass.initializeInstance()

	//attr_accessor :body, :status, :status_code, :protocol, :transfer_encoding, :http_version, :request_http_version, :request
	//attr_reader :headers

	body, err := ioutil.ReadAll(goResp.Body)
	if err != nil {
		return nil, err
	}

	goobyResp.InstanceVariableSet("@body", t.vm.InitStringObject(string(body)))
	goobyResp.InstanceVariableSet("@status_code", t.vm.InitObjectFromGoType(goResp.StatusCode))
	goobyResp.InstanceVariableSet("@status", t.vm.InitObjectFromGoType(goResp.Status))
	goobyResp.InstanceVariableSet("@protocol", t.vm.InitObjectFromGoType(goResp.Proto))
	goobyResp.InstanceVariableSet("@transfer_encoding", t.vm.InitObjectFromGoType(goResp.TransferEncoding))

	underHeaders := map[string]Object{}

	for k, v := range goResp.Header {
		underHeaders[k] = t.vm.InitObjectFromGoType(v)
	}

	goobyResp.InstanceVariableSet("@headers", t.vm.InitHashObject(underHeaders))

	return goobyResp, nil
}
