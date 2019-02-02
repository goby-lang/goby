package main

import (
	"fmt"
	"github.com/gooby-lang/gooby/test_fixtures/import_test/plugin/lib"
)

var ReturnLibName = lib.ReturnLibName

func Bar() string {
	return "Bar"
}

// Foo ...
func Foo(s string) {
	fmt.Println("Foo" + s)
}

// Baz ...
func Baz() {
	fmt.Println("Baz")
}

func main() {
	fmt.Println("Main")
}
