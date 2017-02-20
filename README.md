# rooby

[![Build Status](https://travis-ci.org/st0012/rooby.svg?branch=master)](https://travis-ci.org/st0012/rooby)

Rooby is a new object oriented language written in Go.

##  Sample snippet.
```
class Foo {
  def set(x) {
    let @x = x;
  }
  def get() {
    @x
  }
}
class Bar {
  def set(x) {
    let @x = x;
  }
  def get() {
    @x
  }
}
let f1 = Foo.new;
f1.set(10);
let f2 = Foo.new;
f2.set(21);
let b = Bar.new;
b.set(9)
f2.get() + f1.get() + b.get(); #=> 40
```

## TODO

1. Customize initialization method
2. Inheritance
3. REPL
4. Advanced data structues (Array/Hash)
5. Execution command
6. Basic IO
7. More documentation
8. Primitive type class
