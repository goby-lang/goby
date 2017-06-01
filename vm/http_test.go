package vm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPResponse(t *testing.T) {
	script := `
	require("net/http")

	res = Net::HTTP::Response.new

	res.body = "test"
	res.status = 200

	res.body
	`

	evaluated := testEval(t, script)
	testStringObject(t, evaluated, "test")
}

func TestNormalGet(t *testing.T) {
	expected := "Hello, client"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	}))

	defer ts.Close()

	testScript := fmt.Sprintf(`
require("net/http")

Net::HTTP.get("%s")
`, ts.URL)

	evaluated := testEval(t, testScript)
	testStringObject(t, evaluated, expected)
}

func TestNormalGetWithPath(t *testing.T) {
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
require("net/http")

Net::HTTP.get("%s", "path")
`, ts.URL)

	evaluated := testEval(t, testScript)
	testStringObject(t, evaluated, expected)
}
