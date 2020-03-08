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

		res = Net::HTTP.start do |client|
			res = client.get("http://127.0.0.1:3000/index")
		end

		res.body
		`, "GET Hello World"},
		{`
		require "net/http"

		res = Net::HTTP.start do |client|
			client.post("http://127.0.0.1:3000/index", "text/plain", "Hi Again")
		end

		res.body
		`, "POST Hi Again"},
		{`
		require "net/http"

		res = Net::HTTP.start do |client|
			r = client.request()
			r.url = "http://127.0.0.1:3000/index"
			r.method = "POST"
			r.body = "Another way of doing it"
			client.exec(r)
		end

		res.body
		`, "POST Another way of doing it"},
		{`
		require "net/http"

		res = Net::HTTP.start do |client|
			client.head("http://127.0.0.1:3000/index")
		end

		res.status_code
		`, 200},
		{`
		require "net/http"

		res = Net::HTTP.start do |client|
			client.get("http://127.0.0.1:3000/error")
		end

		res.status_code
		`, 404},
	}

	//block until server is ready
	<-c

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHTTPClientObjectFail(t *testing.T) {

	testsFail := []errorTestCase{
		{`
		require "net/http"

		res = Net::HTTP.start do |client|
			client.get("http://127.0.0.1:3001")
		end

		res
		`, "HTTPError: Could not complete request, Get \"http://127.0.0.1:3001\": dial tcp 127.0.0.1:3001: connect: connection refused", 4},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 2)
	}
}
