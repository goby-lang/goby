package vm

import "testing"

func TestURIParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// Scheme
		{`
		u = URI.parse("http://example.com")
		u.scheme
		`, "http"},
		{`
		u = URI.parse("https://example.com")
		u.scheme
		`, "https"},
		// Host
		{`
		u = URI.parse("http://example.com")
		u.host
		`, "example.com"},
		// Port
		{`
		u = URI.parse("http://example.com")
		u.port
		`, 80},
		{`
		u = URI.parse("https://example.com")
		u.port
		`, 443},
		// Path
		{`
		u = URI.parse("https://example.com/posts/1")
		u.path
		`, "/posts/1"},
		{`
		u = URI.parse("https://example.com")
		u.path
		`, "/"},
		// Query
		{`
		u = URI.parse("https://example.com?foo=bar&a=b")
		u.query
		`, "foo=bar&a=b"},
		{`
		u = URI.parse("https://example.com")
		u.query
		`, nil},
		// User
		{`
		u = URI.parse("https://example.com?foo=bar&a=b")
		u.user
		`, nil},
		// Password
		{`
		u = URI.parse("https://example.com")
		u.password
		`, nil},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEvalWithRequire(t, tt.input, getFilename(), "uri")
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
