package vm

import (
	"strconv"
	"net/http"
	"errors"
	"strings"
	"fmt"
	"io/ioutil"
)

func builtinHTTPClientClassMethods() []*BuiltInMethodObject {
	goClient := http.DefaultClient

	return []*BuiltInMethodObject{
		{
			// Sends a GET request to the target and returns the HTTP response as a string.
			Name: "send",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					if len(args) != 0 {
						return t.vm.initErrorObject(ArgumentError, "Expect 0 arguments. got=%v", strconv.Itoa(len(args)))

					}

					req := httpRequestClass.initializeInstance()

					result := t.builtInMethodYield(blockFrame, req)

					if err, ok := result.Target.(*Error); ok {
						fmt.Printf("Error: %s", err.Message)
						return err //a Error object
					}


					goReq, err := requestGobyToGo(req)
					if err != nil {
						return t.vm.initErrorObject(ArgumentError, "Could not parse request object %v", req)
					}

					resp, err := goClient.Do(goReq)
					if err != nil {
						return t.vm.initErrorObject(HTTPResponseError, "Could not get response: %s", err.Error())
					}

					gobyResp, err := responseGoToGoby(t, resp)
					if err != nil {
						return t.vm.initErrorObject(HTTPError, "Could not read response: %s", err.Error())
					}

					return gobyResp
				}
			},
		},
	}
}

func requestGobyToGo(gobyReq *RObject) (*http.Request, error) {
	//:method, :protocol, :body, :content_length, :transfer_encoding, :host, :path, :url, :params
	uObj, ok := gobyReq.instanceVariableGet("url")
Er:
	if !ok {
		return nil, errors.New("could not parse paramehters")
	}

	u := uObj.(*StringObject).value

	methodObj, ok := gobyReq.instanceVariableGet("method")
	if !ok {
		goto Er
	}

	method := methodObj.(*StringObject).value

	var body string
	if method == "GET" || method== "HEAD" {
		bodyObj, ok := gobyReq.instanceVariableGet("body")
		if !ok {
			goto Er
		}

		body = bodyObj.(*StringObject).value
	}

	return http.NewRequest(method, u, strings.NewReader(body))

}

func responseGoToGoby(t *thread, goResp *http.Response) (*RObject, error) {
	gobyResp := httpResponseClass.initializeInstance()

	//attr_accessor :body, :status, :status_code, :protocol, :transfer_encoding, :http_version, :request_http_version, :request
//attr_reader :headers

	body, err := ioutil.ReadAll(goResp.Body)
	if err != nil {
		return nil, err
	}

	gobyResp.instanceVariableSet("body", t.vm.initStringObject(string(body)))
	gobyResp.instanceVariableSet("status_code", t.vm.initObjectFromGoType(goResp.StatusCode))
	gobyResp.instanceVariableSet("status", t.vm.initObjectFromGoType(goResp.Status))
	gobyResp.instanceVariableSet("protocol", t.vm.initObjectFromGoType(goResp.Proto))
	gobyResp.instanceVariableSet("transfer_encoding", t.vm.initObjectFromGoType(goResp.TransferEncoding))


	underHeaders := map[string]Object{}

	for k, v := range goResp.Header {
		underHeaders[k] = t.vm.initObjectFromGoType(v)
	}

	gobyResp.instanceVariableSet("headers", t.vm.initHashObject(underHeaders))


	return gobyResp, nil
}