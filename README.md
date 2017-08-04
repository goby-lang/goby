# goby (rooby)

[![Build Status](https://travis-ci.org/goby-lang/goby.svg?branch=master)](https://travis-ci.org/goby-lang/goby)
[![Code Climate](https://codeclimate.com/github/goby-lang/goby/badges/gpa.svg)](https://codeclimate.com/github/goby-lang/goby)
[![GoDoc](https://godoc.org/github.com/goby-lang/goby?status.svg)](https://godoc.org/github.com/goby-lang/goby)
[![Go Report Card](https://goreportcard.com/badge/github.com/goby-lang/goby)](https://goreportcard.com/report/github.com/goby-lang/goby)
[![codecov](https://codecov.io/gh/goby-lang/goby/branch/master/graph/badge.svg)](https://codecov.io/gh/goby-lang/goby)
[![Readme Score](http://readme-score-api.herokuapp.com/score.svg?url=goby-lang/goby)](http://clayallsopp.github.io/readme-score?url=goby-lang/goby)

Join us on Slack! [![](https://goby-lang-slackin.herokuapp.com/badge.svg)](https://goby-lang-slackin.herokuapp.com)

Goby is an object-oriented interpreter language deeply inspired by Ruby and written in 100% pure Go. The goal of Goby is to help web developers create api servers or microservices simply and efficiently, with a help of tough thread-mechanism from Go's goroutine, see the [thread's example](https://github.com/goby-lang/goby/blob/master/samples/one_thousand_threads.gb). Howerver, We do not intend to reproduce all of works in Ruby implementation.

Goby will finally equip a reduced set of Ruby's fundamental syntax, including Ruby's common methods and libraries, but will not equip most of Ruby's meta-programming magic to make Goby VM simpler. Goby will also finally equip a built-in HTTP library and multi-threaded server that comes from Go's HTTP packages.

Goby interpreter is a monolithic binary executable, which consists of a YARV-like VM and a compiler, and a REPL (which is better than `irb`!). All components of Goby compiler, such as AST, lexer, parser, token, are written in 100% pure Go, instead of using conventional static yacc/lex/bison conversion. Goby maintainers don't need to care about C language anymore!

We are optimizing and expanding Goby all the time and need your help. One of our vision is to utilize and manages tons of Go's packages easily from Goby scripts. To bind Goby scripts and Go packages, we might introduce a type-system partially.

**Demo:**

<img src="http://i.imgur.com/5RxFgIW.gif?1" width="60%">

**Project Structure:**

![](http://i.imgur.com/WoO1EEY.jpg)

Support us with a monthly donation and help us continue our activities. [[Become a backer](https://opencollective.com/goby#backer)]

<a href="https://opencollective.com/goby/backer/0/website" target="_blank"><img src="https://opencollective.com/goby/backer/0/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/1/website" target="_blank"><img src="https://opencollective.com/goby/backer/1/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/2/website" target="_blank"><img src="https://opencollective.com/goby/backer/2/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/3/website" target="_blank"><img src="https://opencollective.com/goby/backer/3/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/4/website" target="_blank"><img src="https://opencollective.com/goby/backer/4/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/5/website" target="_blank"><img src="https://opencollective.com/goby/backer/5/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/6/website" target="_blank"><img src="https://opencollective.com/goby/backer/6/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/7/website" target="_blank"><img src="https://opencollective.com/goby/backer/7/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/8/website" target="_blank"><img src="https://opencollective.com/goby/backer/8/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/9/website" target="_blank"><img src="https://opencollective.com/goby/backer/9/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/10/website" target="_blank"><img src="https://opencollective.com/goby/backer/10/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/11/website" target="_blank"><img src="https://opencollective.com/goby/backer/11/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/12/website" target="_blank"><img src="https://opencollective.com/goby/backer/12/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/13/website" target="_blank"><img src="https://opencollective.com/goby/backer/13/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/14/website" target="_blank"><img src="https://opencollective.com/goby/backer/14/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/15/website" target="_blank"><img src="https://opencollective.com/goby/backer/15/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/16/website" target="_blank"><img src="https://opencollective.com/goby/backer/16/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/17/website" target="_blank"><img src="https://opencollective.com/goby/backer/17/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/18/website" target="_blank"><img src="https://opencollective.com/goby/backer/18/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/19/website" target="_blank"><img src="https://opencollective.com/goby/backer/19/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/20/website" target="_blank"><img src="https://opencollective.com/goby/backer/20/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/21/website" target="_blank"><img src="https://opencollective.com/goby/backer/21/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/22/website" target="_blank"><img src="https://opencollective.com/goby/backer/22/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/23/website" target="_blank"><img src="https://opencollective.com/goby/backer/23/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/24/website" target="_blank"><img src="https://opencollective.com/goby/backer/24/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/25/website" target="_blank"><img src="https://opencollective.com/goby/backer/25/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/26/website" target="_blank"><img src="https://opencollective.com/goby/backer/26/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/27/website" target="_blank"><img src="https://opencollective.com/goby/backer/27/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/28/website" target="_blank"><img src="https://opencollective.com/goby/backer/28/avatar.svg"></a>
<a href="https://opencollective.com/goby/backer/29/website" target="_blank"><img src="https://opencollective.com/goby/backer/29/avatar.svg"></a>


### Sponsors

Become a sponsor and get your logo on our README on Github with a link to your site. [[Become a sponsor](https://opencollective.com/goby#sponsor)]

<a href="https://opencollective.com/goby/sponsor/0/website" target="_blank"><img src="https://opencollective.com/goby/sponsor/0/avatar.svg"></a>
<a href="https://opencollective.com/goby/sponsor/1/website" target="_blank"><img src="https://opencollective.com/goby/sponsor/1/avatar.svg"></a>
<a href="https://opencollective.com/goby/sponsor/2/website" target="_blank"><img src="https://opencollective.com/goby/sponsor/2/avatar.svg"></a>
<a href="https://opencollective.com/goby/sponsor/3/website" target="_blank"><img src="https://opencollective.com/goby/sponsor/3/avatar.svg"></a>
<a href="https://opencollective.com/goby/sponsor/4/website" target="_blank"><img src="https://opencollective.com/goby/sponsor/4/avatar.svg"></a>
<a href="https://opencollective.com/goby/sponsor/5/website" target="_blank"><img src="https://opencollective.com/goby/sponsor/5/avatar.svg"></a>
<a href="https://opencollective.com/goby/sponsor/6/website" target="_blank"><img src="https://opencollective.com/goby/sponsor/6/avatar.svg"></a>
<a href="https://opencollective.com/goby/sponsor/7/website" target="_blank"><img src="https://opencollective.com/goby/sponsor/7/avatar.svg"></a>
<a href="https://opencollective.com/goby/sponsor/8/website" target="_blank"><img src="https://opencollective.com/goby/sponsor/8/avatar.svg"></a>
<a href="https://opencollective.com/goby/sponsor/9/website" target="_blank"><img src="https://opencollective.com/goby/sponsor/9/avatar.svg"></a>

(**We'll release first beta version in August, please checkout this [issue](https://github.com/goby-lang/goby/issues/72) for what features `Goby` will support.**)

## Table of contents

- [Supported Features](#supported-features)
- [Install](#install)
- [Usage](#usage)
- [Samples](#samples)
- [Documentations](#documentations)
- [Contribute](#contribute)
- [Maintainers](#maintainers)
- [Support us](#support-us)
- [References](#references)

## Supported Features
- Can evaluate bytecode directly
- Everything is object
- Support comment 
- Object and Class
    - Top level main object
    - Constructor
    - Support class methods
    - Support inheritance
    - Support instance variable
    - Support `self`
- Module
- Namespace
- Variables
    - Constant
    - Local variable
    - Instance variable
- Method
    - Support evaluation with arguments
    - Support evaluation without arguments
    - Support evaluation with block (closure)
- BuiltIn Data Types (All of them are classes 😀)
    - Class
    - Integer
    - String
    - Boolean
    - nil
    - Hash (with built in `to_json` method)
    - Array
- Flow control
    - If statement
    - while statement
- Import other files
    - require_relative
    - require (only for standard libraries now)
- Standard Libraries (all of them are at very early stage)
    - `URI`
    - `Net::HTTP`
    - `Net::SimpleServer` (This is very cool and quite performante, check the [sample](https://github.com/goby-lang/goby/blob/master/samples/server.gb))
    - `File`
- IO
    - `puts`
    - `ARGV`
- REPL (run `goby -i`)
- Thread (this should work but the implementation is quite naive and will be refined in the future)
    - Support `thread` method to create a new thread (like `goroutine`)
    - Has `Channel` class for passing objects between threads (like `chan` in Go)
    - See this sample: [One thousand threads](https://github.com/goby-lang/goby/blob/master/samples/one_thousand_threads.gb)

    
**(You can open an issue for any feature request)**

## Install

### From Source

1. You must have Golang installed
2. You must have set $GOPATH
3. Add your $GOPATH/bin into $PATH
4. Run following command 

    ```
    $ go get github.com/goby-lang/goby
    ```
5. Set `GOBY_ROOT` to the project's root path, which should be:

    ```
    $GOPATH/src/github.com/goby-lang/goby
    ```

### Via Homebrew

**Please checkout the [latest release](https://github.com/goby-lang/goby/releases) before using this approach**

```
brew tap goby-lang/goby
brew install goby
```

### Verify Your Installation

1. Run `goby -v` to see the version.
2. Run `goby -i` to enter interactive console.
3. Type `require "file"`.
4. If no error shows up than you have successfully installed Goby :)
5. You can also just run `brew test goby` if you install it via homebrew.

**If you have any issue installing Goby, please open an issue for it**

## Usage

**Execute goby file:**
```
$ goby ./samples/server.gb
```

**Run interactive console:**
```
$ goby -i
```

## Samples

See [sample directory](https://github.com/goby-lang/goby/tree/master/samples) for sample code snippets, like:

- [Built a stack data structure using Goby](https://github.com/goby-lang/goby/blob/master/samples/stack.gb)
- [Running a "Hello World" app with built in server library](https://github.com/goby-lang/goby/blob/master/samples/server/server.gb)
- [Sending request using http library](https://github.com/goby-lang/goby/blob/master/samples/http.gb)
- [Running load test on blocking server](https://github.com/goby-lang/goby/blob/master/samples/server/blocking_server.gb) (This shows `Goby`'s simple server is very performant and can handle requests concurrently)
- [One thousand threads](https://github.com/goby-lang/goby/blob/master/samples/one_thousand_threads.gb)

## Documentations

Check out our [API Documentation](https://goby-lang.github.io/api.doc/).

There is still a lot of document to add. Feel free to contribute following [this guide](https://github.com/goby-lang/api.doc#documenting-goby-code).

## Contribute

See the [guideline](https://github.com/goby-lang/goby/blob/master/CONTRIBUTING.md).

## Maintainers

- @st0012
- @janczer
- @adlerhsieh
- @hachi8833
- @Maxwell-Alexius
- @shes50103

## Support Us

### Backers

**Supporting Goby by sending your first PR! See [contribution guideline](https://github.com/goby-lang/goby/blob/master/CONTRIBUTING.md)**

**Or [support us on opencollective](https://opencollective.com/goby) (I quit my job to develop `Goby` in full-time, so financial support are needed 😢)**

## References

I can't build this project without these resources, and I highly recommend you to check them out if you're interested in building your own languages:

- [Write An Interpreter In Go](https://interpreterbook.com)
- [Nand2Tetris II](https://www.coursera.org/learn/nand2tetris2/home/welcome)
- [Ruby under a microscope](http://patshaughnessy.net/ruby-under-a-microscope)
- [YARV's instruction table](http://www.atdot.net/yarv/insnstbl.html)
