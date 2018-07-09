package plugin

import (
	"strings"
	"testing"
)

func TestCompilePluginTemplate(t *testing.T) {
	pkgs := []*pkg{
		{
			Prefix: "",
			Name:   "database/sql",
		},
		{
			Prefix: "_",
			Name:   "github.com/lib/pq",
		},
	}

	funcs := []*function{
		{
			Prefix: "sql",
			Name:   "Open",
		},
	}

	result := strings.TrimSpace(compilePluginTemplate(pkgs, funcs))
	expected := strings.TrimSpace(`
package main


import(
	 "database/sql"

	_ "github.com/lib/pq"

)

var Open = sql.Open

func main() {}
`)

	if result != expected {
		t.Errorf("Expect result template:\n `%q`.\n got:\n `%q`", expected, result)
	}
}
