# metago

[![Build Status](https://travis-ci.org/st0012/metago.svg?branch=master)](https://travis-ci.org/st0012/metago)
[![codecov](https://codecov.io/gh/st0012/metago/branch/master/graph/badge.svg)](https://codecov.io/gh/st0012/metago)

`metago` is trying to provide Ruby-like meta-programming features to Go. Mostly just for fun.

## Install

```
go get github.com/st0012/metago
```

## Usage

Currently `metago` only has one function: `CallFunc`

```
CallFunc(receiver interface{}, methodName string, args ...interface{}) interface{}
```

Here's how you can use it:

```go
package main

import (
	"fmt"
	"github.com/st0012/metago"
)

type Bar struct {
	name string
}

func (b *Bar) Name() string {
	return b.name
}

func (b *Bar) SetName(n string) {
	b.name = n
}

func (b *Bar) Send(methodName string, args ...interface{}) interface{} {
	return metago.CallFunc(b, methodName, args...)
}

func main() {
	b := &Bar{}
	b.Send("SetName", "Stan")   // This is like Object#send in Ruby
	fmt.Println(b.name)         // Should be "Stan"
	fmt.Println(b.Send("Name")) // Should also be "Stan"
}
```

As you can see, you can call `*Bar` dynamically. Just like Ruby:

```ruby
b = Bar.new
b.send(:set_name, "Stan")
```

## Future work

I'll trying to add more meta-programming features in Ruby like `method_missing` or `define_method`...etc.