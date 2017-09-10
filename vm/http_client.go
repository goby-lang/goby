package vm

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// Instance methods --------------------------------------------------------
func builtinHTTPClientInstanceMethods() []*BuiltinMethodObject {
	//TODO: cookie jar
	goClient := http.DefaultClient

	return []*BuiltinMethodObject{
		{
			// Sends a GET request to the target and returns the HTTP response as a string.
			Name: "get",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, errors.WrongNumberOfArgumentFormat, 1, len(args))
					}

					u, ok := args[0].(*StringObject)
					if !ok {
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.StringClass, u.Class().Name)
					}

					resp, err := goClient.Get(u.value)
					if err != nil {
						return t.vm.initErrorObject(errors.HTTPError, "Could not complete request, %s", err)
					}

					gobyResp, err := responseGoToGoby(t, resp)
					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
					}

					return gobyResp
				}
			},
		}, {
			// Sends a GET request to the target and returns the HTTP response as a string.
			Name: "post",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 3 {
						return t.vm.initErrorObject(errors.ArgumentError, errors.WrongNumberOfArgumentFormat, 3, len(args))
					}

					u, ok := args[0].(*StringObject)
					if !ok {
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.StringClass, u.Class().Name)
					}

					contentType, ok := args[1].(*StringObject)
					if !ok {
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.StringClass, u.Class().Name)
					}

					body, ok := args[2].(*StringObject)
					if !ok {
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.StringClass, u.Class().Name)
					}

					bodyR := strings.NewReader(body.value)

					resp, err := goClient.Post(u.value, contentType.value, bodyR)
					if err != nil {
						return t.vm.initErrorObject(errors.HTTPError, "Could not complete request, %s", err)
					}

					gobyResp, err := responseGoToGoby(t, resp)
					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
					}

					return gobyResp
				}
			},
		}, {
			// Sends a GET request to the target and returns the HTTP response as a string.
			Name: "head",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, errors.WrongNumberOfArgumentFormat, 1, len(args))
					}

					u, ok := args[0].(*StringObject)
					if !ok {
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.StringClass, u.Class().Name)
					}

					resp, err := goClient.Head(u.value)
					if err != nil {
						return t.vm.initErrorObject(errors.HTTPError, "Could not complete request, %s", err)
					}

					gobyResp, err := responseGoToGoby(t, resp)
					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
					}

					return gobyResp
				}
			},
		}, {
			// Sends a GET request to the target and returns the HTTP response as a string.
			Name: "request",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					return httpRequestClass.initializeInstance()
				}
			},
		}, {
			Name: "exec",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, errors.WrongNumberOfArgumentFormat, 1, len(args))
					}

					if args[0].Class().Name != httpRequestClass.Name {
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, "HTTP Response", args[0].Class().Name)
					}

					goReq, err := requestGobyToGo(args[0])
					if err != nil {
						return t.vm.initErrorObject(errors.ArgumentError, err.Error())
					}

					goResp, err := goClient.Do(goReq)
					if err != nil {
						return t.vm.initErrorObject(errors.HTTPError, "Could not complete request, %s", err)
					}

					gobyResp, err := responseGoToGoby(t, goResp)

					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
					}

					return gobyResp
				}
			},
		},
	}
}

func requestGobyToGo(gobyReq Object) (*http.Request, error) {
	//:method, :protocol, :body, :content_length, :transfer_encoding, :host, :path, :url, :params
	uObj, ok := gobyReq.instanceVariableGet("@url")
	if !ok {
		return nil, fmt.Errorf("could not get url")
	}

	u := uObj.(*StringObject).value

	methodObj, ok := gobyReq.instanceVariableGet("@method")
	if !ok {
		return nil, fmt.Errorf("could not get method")
	}

	method := methodObj.(*StringObject).value

	var body string
	if !(method == "GET" || method == "HEAD") {
		bodyObj, ok := gobyReq.instanceVariableGet("@body")
		if !ok {
			return nil, fmt.Errorf("could not get body")
		}

		body = bodyObj.(*StringObject).value
	}

	return http.NewRequest(method, u, strings.NewReader(body))

}

// Other helper functions -----------------------------------------------

func responseGoToGoby(t *thread, goResp *http.Response) (Object, error) {
	gobyResp := httpResponseClass.initializeInstance()

	//attr_accessor :body, :status, :status_code, :protocol, :transfer_encoding, :http_version, :request_http_version, :request
	//attr_reader :headers

	body, err := ioutil.ReadAll(goResp.Body)
	if err != nil {
		return nil, err
	}

	gobyResp.instanceVariableSet("@body", t.vm.initStringObject(string(body)))
	gobyResp.instanceVariableSet("@status_code", t.vm.initObjectFromGoType(goResp.StatusCode))
	gobyResp.instanceVariableSet("@status", t.vm.initObjectFromGoType(goResp.Status))
	gobyResp.instanceVariableSet("@protocol", t.vm.initObjectFromGoType(goResp.Proto))
	gobyResp.instanceVariableSet("@transfer_encoding", t.vm.initObjectFromGoType(goResp.TransferEncoding))

	underHeaders := map[string]Object{}

	for k, v := range goResp.Header {
		underHeaders[k] = t.vm.initObjectFromGoType(v)
	}

	gobyResp.instanceVariableSet("@headers", t.vm.initHashObject(underHeaders))

	return gobyResp, nil
}
