package main

import "fmt"

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
