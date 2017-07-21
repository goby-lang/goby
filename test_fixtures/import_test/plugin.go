package main

import "fmt"

// Bar ...
type Bar struct {
	name string
}

// Name ...
func (b *Bar) Name() string {
	return b.name
}

// NewBar ...
func NewBar(name string) (*Bar, error) {
	return &Bar{name: name}, nil
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
