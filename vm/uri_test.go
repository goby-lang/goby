package vm

import "testing"

func TestURIParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// Scheme
		{`
		require "uri"

		u = URI.parse("http://example.com")
		u.scheme
		`, "http"},
		{`
		require "uri"

		u = URI.parse("https://example.com")
		u.scheme
		`, "https"},
		// Host
		{`
		require "uri"

		u = URI.parse("http://example.com")
		u.host
		`, "example.com"},
		// Port
		{`
		require "uri"

		u = URI.parse("http://example.com")
		u.port
		`, 80},
		{`
		require "uri"

		u = URI.parse("https://example.com")
		u.port
		`, 443},
		// Path
		{`
		require "uri"

		u = URI.parse("https://example.com/posts/1")
		u.path
		`, "/posts/1"},
		{`
		require "uri"

		u = URI.parse("https://example.com")
		u.path
		`, "/"},
		// Query
		{`
		require "uri"

		u = URI.parse("https://example.com?foo=bar&a=b")
		u.query
		`, "foo=bar&a=b"},
		{`
		require "uri"

		u = URI.parse("https://example.com")
		u.query
		`, nil},
		// User
		{`
		require "uri"

		u = URI.parse("https://example.com?foo=bar&a=b")
		u.user
		`, nil},
		// Password
		{`
		require "uri"

		u = URI.parse("https://example.com")
		u.password
		`, nil},
	}

	for i, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}
