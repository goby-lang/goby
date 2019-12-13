package main

import "fmt"

// Bar ...
type Bar struct {
	name string
}

// Name ...
func (b *Bar) Name(s string) (string, error) {
	return b.name + s, nil
}

// Add ...
func (b *Bar) Add(x int, y int64) int64 {
	return int64(x) + y
}

// NewBar ...
func NewBar(name string) (*Bar, error) {
	return &Bar{name: name}, nil
}

// GetBarName ...
func GetBarName(b *Bar) string {
	return b.name
}

func main() {
	fmt.Println("Main")
}
