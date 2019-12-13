package main

import (
	"fmt"

	"github.com/goby-lang/goby/test_fixtures/import_test/plugin/lib"
)

// ReturnLibName is an alias type of the target "lib'
var ReturnLibName = lib.ReturnLibName

// Bar ...
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
