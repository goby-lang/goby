package main

import "fmt"

type Bar struct {
	name string
}

func Foo(s string) {
	fmt.Println("Foo" + s)
}

func Baz() {
	fmt.Println("Baz")
}

func main() {
	fmt.Println("Main")
}
