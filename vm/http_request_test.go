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
		req = Net::HTTP::Request.new
		req.method = "GET"

		req.method
		`, "GET"},
		{`
		req = Net::HTTP::Request.new
		req.set_header("Content-Type", "text/plain")

		req.headers["Content-Type"]
		`, "text/plain"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEvalWithRequire(t, tt.input, getFilename(), "net/http")
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
