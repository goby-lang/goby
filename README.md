# rooby

[![Build Status](https://travis-ci.org/st0012/rooby.svg?branch=master)](https://travis-ci.org/st0012/rooby)

Rooby is a new object oriented language written in Go.

##  Sample snippet.
```
class Foo {
  def set(x) {
    @x = x;
  }
  def get() {
    @x
  }
}
class Bar {
  def set(x) {
    @x = x;
  }
  def get() {
    @x
  }
}
f1 = Foo.new;
f1.set(10);
f2 = Foo.new;
f2.set(21);
b = Bar.new;
b.set(9)
f2.get() + f1.get() + b.get(); #=> 40
```

## TODO

- ~~Customize initialization method~~
- ~~Inheritance~~
- Definable class methods
- REPL(Hard to be done)
- Advanced data structues (Array/Hash)
- ~~Execution command~~
- Basic IO
- More documentation
- Primitive type class
