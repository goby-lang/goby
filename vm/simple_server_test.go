package vm

import (
	"net/http/httptest"
	"strings"
	"testing"
	"net/http"
	"io/ioutil"
	"time"
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

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestServerSetupResponse(t *testing.T) {
	tests := []struct {
		input    string
		expectedBody string
		expectedStatus int
	}{
		{`
		require "net/simple_server"

		server = Net::SimpleServer.new(4000)
		server.get("/") do |req, res|
		  res.body = req.method + " Hello World"
		  res.status = 200
		end
		puts("foo")
		server.start
		`,
			"GET Hello World",
			200},
	}

	for _, tt := range tests {
		go func() {
			v := initTestVM()
			v.testEval(t, tt.input, getFilename())
		}()

		time.Sleep(1 * time.Second)
		resp, err := http.Get("http://localhost:4000")

		if err != nil {
			t.Fatal(err.Error())
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if string(body) != tt.expectedBody {
			t.Fatalf("Expect response body to be: \n %s, got \n %s", tt.expectedBody, string(body))
		}

		if resp.StatusCode != tt.expectedStatus {
			t.Fatalf("Expect response status to be %d, got %d", tt.expectedStatus, resp.StatusCode)
		}
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

func TestServerRequestInitialization(t *testing.T) {
	v := initTestVM()
	reader := strings.NewReader("Hello World")
	r := initRequest(v.mainThread, httptest.NewRecorder(), httptest.NewRequest("GET", "https://google.com/path", reader))

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
		{
			"@protocol",
			"HTTP/1.1",
		},
		{
			"@content_length",
			11, // Length of the body: "Hello World"
		},
		//{
		//	"@transfer_encoding",
		//	0,
		//},
		//{
		//	"@headers",
		//	123,
		//},
	}

	for i, tt := range tests {
		v, ok := r.InstanceVariables.get(tt.varName)

		if !ok {
			t.Fatalf("Expect request object to have %s attribute.", tt.varName)
		}

		verifyExpected(t, i, v, tt.expected)
	}

}
