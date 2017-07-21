package main

import "fmt"

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
