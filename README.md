# rooby

[![Build Status](https://travis-ci.org/st0012/rooby.svg?branch=master)](https://travis-ci.org/st0012/rooby)

Rooby is a new object oriented language written in Go.

##  Sample snippet.
```
class Foo {
  def set(x) {
    @x = x;
  }
  def get {
    @x
  }
}
class Bar {
  def set(x) {
    @x = x;
  }
  def get {
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

```
class User {
    def initialize(name, age) {
        @name = name;
        @age = age;
    }

    def name {
        @name
    }

    def age {
        @age
    }

    def say_hi(user) {
        @name + " says hi to " + user.name
    }

    def self.sum_age(user1, user2) {
        user1.age() + user2.age()
    }
}

stan = User.new("Stan", 22);
john = User.new("John", 40);

self.puts(User.sum_age(stan, john)); #=> 62
self.puts(stan.say_hi(john)); #=> Stan says hi to John
```

## TODO

- ~~Customize initialization method~~
- ~~Inheritance~~
- ~~Definable class methods~~
- ~~Execution command~~
- Makefile
- Improve syntax
    - method call without self
    - remove semicolon
- for loop support
- Advanced data structures (Array/Hash)
- Basic IO
- More documentation
    - Samples
    - Feature list
- Primitive type class
    - String
    - Integer
    - Boolean
- REPL(Hard to be done)
