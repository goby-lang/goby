package vm

import (
	"testing"
	//"net/http/httptest"
	//"net/http"
	"fmt"
	"io/ioutil"
	"net/http"
)

func TestHTTPObject(t *testing.T) {

	c := make(chan bool, 1)

	//server to test off of
	go func() {
		m := http.NewServeMux()

		m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			if r.Method == http.MethodPost {
				b, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				fmt.Fprintf(w, "POST %s", b)
			} else {
				fmt.Fprint(w, "GET Hello World")
			}

		})

		c <- true

		http.ListenAndServe(":3000", m)
	}()

	tests := []struct {
		input    string
		expected interface{}
	}{
		//test get request
		{`
		require "net/http"

		Net::HTTP.get("http://127.0.0.1:3000")
		`, "GET Hello World"},
		{`
		require "net/http"

		Net::HTTP.post("http://127.0.0.1:3000", "text/plain", "Hi Again")
		`, "POST Hi Again"},
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
