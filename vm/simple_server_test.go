package vm

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServerInitialization(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require "net/simple_server"

		s = Net::SimpleServer.new(4000)
		s.port
		`, 4000},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, evaluated, tt.expected)
	}
}

func TestSetupResponseDefaultValue(t *testing.T) {
	reader := strings.NewReader("")
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "https://google.com/path", reader)

	res := httpResponseClass.initializeInstance()

	setupResponse(recorder, req, res)

	if recorder.Code != 200 {
		t.Fatalf("Expect response code to be 200. got=%d", recorder.Code)
	}

	if recorder.HeaderMap.Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Fatalf("Expect content type to be \"text/plain; charset=utf-8\". got=%s", recorder.HeaderMap.Get("Content-Type"))
	}

	if recorder.Body.String() != "" {
		t.Fatalf("Expect response body to be empty. got=%s", recorder.Body.String())
	}
}

func TestInitRequest(t *testing.T) {
	reader := strings.NewReader("Hello World")
	r := initRequest(httptest.NewRecorder(), httptest.NewRequest("GET", "https://google.com/path", reader))

	tests := []struct {
		varName  string
		expected interface{}
	}{
		{
			"@method",
			"GET",
		},
		{
			"@path",
			"/path",
		},
		{
			"@url",
			"https://google.com/path",
		},
		{
			"@host",
			"google.com",
		},
		{
			"@body",
			"Hello World",
		},
	}

	for _, tt := range tests {
		v, ok := r.InstanceVariables.get(tt.varName)

		if !ok {
			t.Fatalf("Expect request object to have %s attribute.", tt.varName)
		}

		checkExpected(t, v, tt.expected)
	}

}
