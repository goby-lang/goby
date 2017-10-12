![Goby](http://i.imgur.com/ElGAzRn.png?3)
=========

[![Build Status](https://travis-ci.org/goby-lang/goby.svg?branch=master)](https://travis-ci.org/goby-lang/goby)
[![Code Climate](https://codeclimate.com/github/goby-lang/goby/badges/gpa.svg)](https://codeclimate.com/github/goby-lang/goby)
[![GoDoc](https://godoc.org/github.com/goby-lang/goby?status.svg)](https://godoc.org/github.com/goby-lang/goby)
[![Go Report Card](https://goreportcard.com/badge/github.com/goby-lang/goby)](https://goreportcard.com/report/github.com/goby-lang/goby)
[![codecov](https://codecov.io/gh/goby-lang/goby/branch/master/graph/badge.svg)](https://codecov.io/gh/goby-lang/goby)
[![Readme Score](http://readme-score-api.herokuapp.com/score.svg?url=goby-lang/goby)](http://clayallsopp.github.io/readme-score?url=goby-lang/goby)

**Goby** is an object-oriented interpreter language deeply inspired by **Ruby** as well as its core implementation by 100% pure **Go**. Moreover, it has standard libraries to provide several features such as the Plugin system. Note that we do not intend to reproduce whole of the honorable works of Ruby syntax/implementation/libraries. 

One of our goal is to provide web developers a sort of small and handy environment that mainly focusing on creating **API servers or microservices**. For this, Goby includes the following native features:

- Robust thread/channel mechanism powered by Go's goroutine
- Builtin high-performance HTTP server
- Builtin database library (currently only support PostgreSQL adapter)
- JSON support
- [Plugin system](https://goby-lang.gitbooks.io/goby/content/plugin-system.html) that can load existing Go packages dynamically (Only for Linux by now)
- Accessing Go objects from Goby directly

> Note: Goby had formerly been known as "Rooby", which was renamed in May 2017.

## Table of contents

- [Demo and sample Goby app](#demo_and_sample_app)
- [Aspects](#aspects)
    - [Features](#features)
    - [Language](#language)
    - [Native class](#native_class)
    - [Standard class](#standard_class)
- [Installation](#installation)
- [Usage](#usage)
- [Sample codes](#sample_codes)
- [Documentation](#documentation)
- [Joining to Goby](#joining-to-goby)
- [Maintainers](#maintainers)
- [Support us](#support-us)
- [References](#references)

## Demo screen and sample Goby app

<img src="http://i.imgur.com/9YrDZOR.gif" width="60%">

**New!** Check-out our [sample app](http://sample.goby-lang.org) built with Goby. Source code is also available [here](https://github.com/goby-lang/sample-web-app).

## Aspects

Goby has several aspects: language specification, design of compiler and vm, implementation (just one for now), library, and the whole of them. 

----------

**Language**: Class-based, straight-ahead object-oriented script language. Syntax is influenced by Ruby language (and by Go a bit), but has been **condensed and simplified** (and slightly modified) to keep Goby VM simple and concise. Several aspects of Ruby such as meta-programming (known as 'magic'), special variables with `$`, have been dropped for now, but note that we might resurrect some of them with a different form or implementation in the future.

**Class**: Single inheritance. Module is supported for mixin with `#include` or `#extend`. Defining singleton class and singleton method are also supported. Goby has two kinds of class internally: native class and standard class. **Native class** (or builtin class) provides fundamental classes such as `Array` or `String`. `Object` class is a superclass of any other native/standard classes including `Class` class. `Class` class contains most common methods such as `#puts`. **Standard class** (or standard library) can be loaded via `require` and provides additional methods. Standard classes are often split internal Go code and external Goby code in order to make implementation easier. Both kinds of class are transparent to Goby developers and can be overridden by child classes. Any classes including `Class` class are under `Object` class. 

**Compiler**: Consists of **AST**, **lexer**, **parser**, and **token**, which of the structure is pretty conventional and should be familiar to language creators. These components are all written in 100% pure Go, instead of using conventional static yacc/lex/bison conversion with a mess of ad-hoc macros. This makes Goby's codes far smaller, concise, and legible. You can inspect, maintain, or improve Goby codes more easily, being free from pains like C/C++ era. 

**VM**: YARV-conscious, including **stack** and **call_frame**, as well as containing Goby's native classes and some standard library and additional components. All are written in Go as well.

**Implementation**: Built-in monolithic Go binary executable which equips several native features such as a robust **thread/channel** mechanism powered by goroutine, a very new experimental [**Plugin system**](https://goby-lang.gitbooks.io/goby/content/plugin-system.html) to manage existing Go packages dynamically from Goby codes, **igb** (REPL) powered by readline package. Goby contains some standard or third-party Go packages, but the dependency to them is not high. These packages contain **no CGO** codes (at least by now) thus cross-compile for any OS environments that Go supports should work fine. 

**Library**: Provides some lean but sufficient standard libraries to support developers, including **threaded high-performance HTTP server**, **DB adapter**, **file** or **JSON**. Curiously, most of them are split into Go and Goby codes, and Goby codes are not within Goby executable but placed under lib directory as Goby script files. Of course you can create custom libraries and include them to your codes. Thanks to the flexibility of **Plugin system**, we expect that you can quickly import most of the existing Go packages to your Goby scripts without creating additional libraries from scratch in almost all cases. 

-----------

**Let's improve Goby together!**: We are optimizing and expanding Goby all the time. Toward the first release, we've been focusing on implementing Goby first. 

### Features

- Plugin system
    - Allows to use Go libraries (packages) dynamically
    - Allows to call Go's methods from Goby directly (only on Linux for now)
- Builtin multi-threaded server and DB library
- REPL (run `goby -i`)

### Language

Perhaps Goby should be far easier for Rubyists to comprehend. You can use Ruby's syntax highlighting for Goby as well😀
 
- Everything is object
- Object and Class
    - Top level main object
    - Constructor
    - Class/instance method
    - Class
        - Can be inherited with `<`
        - Singleton class
        - `#send` **new!**
    - `self`
- Module for supporting mixin
    - `#include` for instance methods
    - `#extend` for class methods
    - `::` for delimiting namespaces
- Variable: starts with lowercase letter like `var`
    - Local variable
    - Instance variable
- Constant
    - Starts with uppercase like `Var` or `VAR`
    - Global if defined on top-level
    - **not reentrant** by assignment, but still permits redefining class/module
    - (special variables with `$` are unsupported)
- Methods
    - Definition: order of parameter is determined:
        1. normal params (ex: `a`, `b`)
        2. opt params (ex: `ary=[]`, `hs={}`)
        3. splat params (ex: `*sp`) for compatibility with Go functions
    - Evaluation with/without arguments
    - Evaluation with a block (closure)
    - Defining singleton methods
- Block
    - `do` - `end`
- Flow control
    - `if`, `else`, `elsif`
    - `while`
- IO
    - `#puts`
    - `ARGV`, `STDIN`, `STDOUT`, `STDERR`, `ENV` constants
- Import files
    - `require` (Just for standard libraries by now)
    - `require_relative`
- Thread (not a class!)
    - Goroutine-based `thread` method to create a new thread
    - Works with `Channel` class for passing objects between threads, like `chan` in Go
    - See this sample: [One thousand threads](https://github.com/goby-lang/goby/blob/master/samples/one_thousand_threads.gb)

### Native class
 
Written in Go.

- `Class`
- `Integer`
- `String`
- `Boolean`
- `Null` (`nil`)
- `Hash`
- `Array`
- `Range`
- `URI`
- `Channel`
- `File` (Changed from loadable class)
- `GoObject` (provides `#go_func` that wraps pure Go objects or pointers for interaction)
- `Regexp`

### Standard library

written in Go and Goby.

- Loadable class
    - `DB` (only for PostgreSQL by now)
    - `Plugin`
- Loadable module
    - NET
        - `Net::HTTP:Request`
        - `Net::HTTP:Response`
        - `Net::SimpleServer` (try [sample Goby app](http://sample.goby-lang.org) and [source](https://github.com/goby-lang/sample-web-app), or [sample code](https://github.com/goby-lang/goby/blob/master/samples/server.gb)!)

## Installation

Confirmed Goby runs on Mac OS and Linux for now. Try Goby on Windows and let us know the result.

### A. Via Homebrew (binary installation for Mac OS)

**Note: Please check the [latest release](https://github.com/goby-lang/goby/releases) before installing Goby via Homebrew**

```
brew tap goby-lang/goby
brew install goby
```

In the case, `$GOBY_ROOT` is automatically configured. 

### B. From Source

Try this if you'd like to contribute Goby! Skip 1 if you already have Golang in your environment.

1. Prepare Golang environment
    - Install Golang >= 1.9
    - Make sure `$GOPATH` in your shell's config file( like .bashrc) is correct
    - Add you `$GOPATH/bin` to `$PATH`
2. Run `go get github.com/goby-lang/goby`
3. Set the Goby project's exact root path `$GOBY_ROOT` manually, which should be:

```
$GOPATH/src/github.com/goby-lang/goby
```

### C. Installation on a clean Linux environment

For installing both Go and Goby on a clean Linux environment, see the [wiki page](https://github.com/goby-lang/goby/wiki/Linux-Go-and-Goby-setup).

### Verifying Goby installation

1. Run `goby -v` to see the version.
2. Run `goby -i` to launch igb REPL.
3. Type `require "uri"` in igb.

FYI: You can just run `brew test goby` to check Homebrew installation.

**If you have any issue installing Goby, please let us know via [Github issues](https://github.com/goby-lang/goby/issues)**

### Using Docker

Goby has official [docker image](https://cloud.docker.com/app/gobylang/repository/docker/gobylang/goby/general) as well. You can try the [Plugin System](https://goby-lang.gitbooks.io/goby/content/plugin-system.html) using docker.

## Sample codes

- [Built a stack data structure using Goby](https://github.com/goby-lang/goby/blob/master/samples/stack.gb)
- [Running a "Hello World" app with built in server library](https://github.com/goby-lang/goby/blob/master/samples/server/server.gb)
- [Sending request using http library](https://github.com/goby-lang/goby/blob/master/samples/http.gb)
- [Running load test on blocking server](https://github.com/goby-lang/goby/blob/master/samples/server/blocking_server.gb) (This shows `Goby`'s simple server is very performant and can handle requests concurrently)
- [One thousand threads](https://github.com/goby-lang/goby/blob/master/samples/one_thousand_threads.gb)

More sample Goby codes can be found in [sample directory](https://github.com/goby-lang/goby/tree/master/samples).

## Documentation

- [**User Manual (WIP)**](https://goby-lang.gitbooks.io/goby/content/)(Gitbooks)
- [API Documentation](https://goby-lang.github.io/api.doc/) -- needs update the build script. See the [guide for API doc](https://github.com/goby-lang/api.doc#documenting-goby-code) if you'd like to contribute. 

## Joining to Goby

**Join us on Slack!** [![](https://goby-lang-slackin.herokuapp.com/badge.svg)](https://goby-lang-slackin.herokuapp.com)

See the [guideline](https://github.com/goby-lang/goby/blob/master/CONTRIBUTING.md).

## Maintainers

- @st0012
- @hachi8833
- @Maxwell-Alexius

## Designer
- [steward379](https://dribbble.com/steward379)

## Support Us

### Donations

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

[![](http://i.imgur.com/dsKTzXZ.png?1)](https://5xruby.tw/en)

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

**Supporting Goby by sending your first PR! See [contribution guideline](https://github.com/goby-lang/goby/blob/master/CONTRIBUTING.md)**

**Or [support us on opencollective](https://opencollective.com/goby) (I quit my job to develop `Goby` in full-time, so financial support are needed 😢)**

## References

The followings are the essential resources to create Goby; I highly recommend you to check them first if you'd be interested in building your own languages:

- [Write An Interpreter In Go](https://interpreterbook.com)
- [Nand2Tetris II](https://www.coursera.org/learn/nand2tetris2/home/welcome)
- [Ruby under a microscope](http://patshaughnessy.net/ruby-under-a-microscope)
- [YARV's instruction table](http://www.atdot.net/yarv/insnstbl.html)
