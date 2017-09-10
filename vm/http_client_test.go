package vm

import "testing"

func TestHTTPClientObject(t *testing.T) {

	//blocking channel
	c := make(chan bool, 1)

	//server to test off of
	go startTestServer(c)

	tests := []struct {
		input    string
		expected interface{}
	}{
		//test get request
		{`
		require "net/http"

		c = Net::HTTP::Client.new

		res = c.send do |req|
			req.url = "http://127.0.0.1:3000/index"
			req.method = "GET"
		end

		res.body
		`, "GET Hello World"},
		{`
		require "net/http"

		c = Net::HTTP::Client.new

		res = c.send do |req|
			req.url = "http://127.0.0.1:3000/index"
			req.method = "POST"
			req.body = "Hi Again"
		end

		res.body
		`, "POST Hi Again"},
		{`
		require "net/http"

		c = Net::HTTP::Client.new

		res = c.send do |req|
			req.url = "http://127.0.0.1:3000/error"
			req.method = "GET"
		end

		res.status_code
		`, 404},
	}

	//block until server is ready
	<-c

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHTTPClientObjectFail(t *testing.T) {

	testsFail := []errorTestCase{
		{`
		require "net/http"

		c = Net::HTTP::Client.new

		res = c.send do |req|
			req.url = "http://127.0.0.1:3001"
			req.method = "GET"
		end

		res
		`, "HTTPError: Could not complete request, Get http://127.0.0.1:3001: dial tcp 127.0.0.1:3001: getsockopt: connection refused", 6},
	}

	//block until server is ready
	//<-c

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}
