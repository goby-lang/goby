package vm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPResponseObject(t *testing.T) {
	input := `
	require "net/http"

	res = Net::HTTP::Response.new

	res.body = "test"
	res.status = 200

	res.body
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	VerifyExpected(t, 0, evaluated, "test")
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestNormalGetResponse(t *testing.T) {
	expected := "Hello, client"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	}))

	defer ts.Close()

	testScript := fmt.Sprintf(`
require "net/http"

Net::HTTP.get("%s")
`, ts.URL)

	v := initTestVM()
	evaluated := v.testEval(t, testScript, getFilename())
	VerifyExpected(t, 0, evaluated, expected)
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestNormalGetResponseWithPath(t *testing.T) {
	expected := "Hello, client"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/path":
			fmt.Fprint(w, expected)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))

	defer ts.Close()

	testScript := fmt.Sprintf(`
require "net/http"

Net::HTTP.get("%s", "path")
`, ts.URL)

	v := initTestVM()
	evaluated := v.testEval(t, testScript, getFilename())
	VerifyExpected(t, 0, evaluated, expected)
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}
