package main

import "fmt"

type Bar struct {
	name string
}

func (b *Bar) Name() {
	fmt.Println(b.name)
}

func NewBar(name string) *Bar {
	return &Bar{name: name}
}

func Foo(s string) {
	fmt.Println("Foo" + s)
}

func main() {
	fmt.Println("Main")
}
