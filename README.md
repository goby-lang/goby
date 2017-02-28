# Rooby

[![Build Status](https://travis-ci.org/st0012/rooby.svg?branch=master)](https://travis-ci.org/st0012/rooby)

Rooby is a new object oriented language written in Go.

## Install

```
$ git clone git@github.com:st0012/rooby.git
$ cd rooby
$ make install
```

## Try it!
```
$ rooby ./samples/sample-1.ro
$ rooby ./samples/sample-2.ro
$ rooby ./samples/sample-3.ro
$ rooby ./samples/sample-4.ro
```

##  Sample snippet.
```ruby
class Foo {
  def set(x) {
    @x = x
  }
  def get {
    @x
  }
}
class Bar < Foo {}
class Baz < Foo {}

bar = Bar.new
baz = Baz.new
foo = Foo.new
bar.set(10)
baz.set(1)
foo.set(5)

puts(bar.get + baz.get + foo.get) #=> 16
```

```ruby
class User {
  def initialize(name, age) {
      @name = name
      @age = age
  }

  def name {
      @name
  }

  def age {
      @age
  }

  def say_hi(user) {
      puts(@name + " says hi to " + user.name)
  }

  def self.sum_age(user1, user2) {
      user1.age + user2.age
  }
}

stan = User.new("Stan", 22)
john = User.new("John", 40)
puts(User.sum_age(stan, john)) #=> 62
stan.say_hi(john) #=> Stan says hi to John
```

```ruby
class JobPosition {
  def initialize(name) {
    @name = name
  }

  def name {
    @name
  }
    
  def self.engineer {
    new("Engineer")
  }
}

job = JobPosition.engineer
puts(job.name) #=> "Engineer"
```


```ruby
puts("123".class.name) #=> String
puts(123.class.name) #=> Integer
puts(true.class.name) #=> Boolean
```

## TODO

- ~~Customize initialization method~~
- ~~Inheritance~~
- ~~Definable class methods~~
- ~~Execution command~~
- ~~Improve built in method's self implementation~~ 
- Improve syntax
    - ~~method call without self~~
    - ~~remove semicolon~~
- ~~Makefile~~
- for loop support
- ~~Comment support~~
- Advanced data structures (Array/Hash)
- Basic IO
- More documentation
    - Samples
    - Feature list
- ~~Primitive type class~~
    - ~~String~~
    - ~~Integer~~
    - ~~Boolean~~
- REPL(Hard to be done)
