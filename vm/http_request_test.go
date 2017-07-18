package vm

import (
	"testing"
	//"net/http/httptest"
	//"net/http"
	//"fmt"
)

func TestHTTPRequestObject(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require "net/http"

		req = Net::HTTP::Request.new
		req.method = "GET"

		req.method
		`, "GET"},
		{`
		require "net/http"

		req = Net::HTTP::Request.new
		req.set_header("Content-Type", "text/plain")

		req.headers["Content-Type"]
		`, "text/plain"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}
