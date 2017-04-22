# Rooby

[![Join the chat at https://gitter.im/Rooby-lang/Lobby](https://badges.gitter.im/Rooby-lang/Lobby.svg)](https://gitter.im/Rooby-lang/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

[![Build Status](https://travis-ci.org/rooby-lang/Rooby.svg?branch=master)](https://travis-ci.org/rooby-lang/Rooby)
[![Code Climate](https://codeclimate.com/github/rooby-lang/Rooby/badges/gpa.svg)](https://codeclimate.com/github/rooby-lang/rooby)
[![GoDoc](https://godoc.org/github.com/rooby-lang/Rooby?status.svg)](https://godoc.org/github.com/rooby-lang/Rooby)
[![Go Report Card](https://goreportcard.com/badge/github.com/rooby-lang/Rooby)](https://goreportcard.com/report/github.com/rooby-lang/Rooby)
[![codecov](https://codecov.io/gh/rooby-lang/Rooby/branch/master/graph/badge.svg)](https://codecov.io/gh/rooby-lang/Rooby)
[![Readme Score](http://readme-score-api.herokuapp.com/score.svg?url=rooby-lang/rooby)](http://clayallsopp.github.io/readme-score?url=rooby-lang/rooby)

Rooby is a Ruby-like object oriented language written in Go. You can think it as a simplified, compilable Ruby for now.
   

Here's my expectation about Rooby:

- Has Ruby-like syntax and object system.
- Rooby program will be compiled at first place and then executed through VM. 
- It's compilation and evaluation are separated, it can firstly compile your program into bytecode and execute it later. This will work like python.
- **Maybe** it can gain some improvement on concurrency since it's based on Go.
- Any special idea comes to my mind or proposed by you 😄 

**Join me to build Rooby together!**

## Features
- **Can be compiled into bytecode (with `.robc` extension)**
- **Can evaluate bytecode directly**
- Everything is object
- Support comment
- Object & Class
    - Top level main object
    - Constructor
    - Support class method
    - Support inheritance
    - Support instance variable
    - Support self
- Variables
    - Constant
    - Local variable
    - Instance variable
- Method
    - Support evaluation with arguments
    - Support evaluation without arguments
    - Support evaluation with block
- BuiltIn Data Types (All of them are classes 😀)
    - Class
    - Integer
    - String
    - Boolean
    - nil (has this type internally but parser hasn't support yet)
    - Hash
    - Array
    - **Not** support symbols. Since string is already immutable, supporting symbols is not that necessary.
- Flow control
    - If statement
    - while statement
    - Haven't support `for` yet
- IO
    - Just `puts` for now
    
**(You can open an issue for any feature request)** 
    
## TODO

See [github progjects](https://github.com/rooby-lang/Rooby/projects)

## Install

1. You must have Golang installed
2. You must have set $GOPATH
3. Add your $GOPATH/bin into $PATH
4. Run following command 

```
$ go get github.com/rooby-lang/Rooby
```

## Usage

**Execute Rooby file using VM**

(might see errors on sample-6 since vm hasn't support block yet)
``` 
$ rooby ./samples/sample-1.ro
#=> 16
```

**Compile Rooby code**

```
$ rooby -c ./samples/sample-1.ro
```

You'll see `sample-1.robc` in `./samples`

**Execute bytecode**

```
$ rooby ./samples/sample-1.robc
```


## Try it!
(See sample directory)
```
$ rooby ./samples/sample-1.ro
$ rooby ./samples/sample-2.ro
$ rooby ./samples/sample-3.ro
$ rooby ./samples/sample-4.ro
$ rooby .....
```

## Development

It will be actively developed for at least a few months. Currently I'm working on building a vm that supports some basic features in Ruby (block, module...etc.).
And I will use [github project](https://github.com/rooby-lang/Rooby/projects) to manage Rooby's development progress, you can check what I'm doing and about to do there.

## Contribute

I will appreciate any feature proposal or issue report. And please contact me directly if you want to get involved. Rooby is very young so we can do a lot interesting things together.

##  Sample snippet.

```ruby
class User
  def initialize(name, age)
    @name = name
    @age = age
  end

  def name
    @name
  end

  def age
    @age
  end

  def say_hi(user)
    puts(@name + " says hi to " + user.name)
  end

  def self.sum_age(user1, user2)
    user1.age + user2.age
  end
end

stan = User.new("Stan", 22)
john = User.new("John", 40)
puts(User.sum_age(stan, john)) #=> 62
stan.say_hi(john) #=> Stan says hi to John
```

#### Build a stack using Rooby

```ruby
class Stack
  def initialize
    @data = []
  end
    
  def push(x)
    @data.push(x)
  end
    
  def pop
    @data.pop
  end
    
  def top
    @data[@data.length - 1]
  end
end

s = Stack.new
s.push(1)
s.push(2)
s.push(3)
s.push(4)
s.push(10)
puts(s.pop) #=> 10
puts(s.top) #=> 4
```

#### Block support

```ruby
class Car
  def initialize
    yield(self)
  end
  
  def color=(c)
    @color = c
  end
  
  def color
    @color
  end
  
  def doors=(ds)
    @doors = ds
  end
  
  def doors
    @doors
  end
end
 
car = Car.new do |c|
  c.color = "Red"
  c.doors = 4
end
 
puts("My car's color is " + car.color + " and it's got " + car.doors.to_s + " doors.")

```
