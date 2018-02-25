![Goby](http://i.imgur.com/ElGAzRn.png?3)
=========

[![](https://goby-lang-slackin.herokuapp.com/badge.svg)](https://goby-lang-slackin.herokuapp.com)
[![Build Status](https://travis-ci.org/goby-lang/goby.svg?branch=master)](https://travis-ci.org/goby-lang/goby)
[![GoDoc](https://godoc.org/github.com/goby-lang/goby?status.svg)](https://godoc.org/github.com/goby-lang/goby)
[![Go Report Card](https://goreportcard.com/badge/github.com/goby-lang/goby)](https://goreportcard.com/report/github.com/goby-lang/goby)
[![codecov](https://codecov.io/gh/goby-lang/goby/branch/master/graph/badge.svg)](https://codecov.io/gh/goby-lang/goby)
[![Readme Score](http://readme-score-api.herokuapp.com/score.svg?url=goby-lang/goby)](http://clayallsopp.github.io/readme-score?url=goby-lang/goby)
[![Snap Status](https://build.snapcraft.io/badge/goby-lang/goby.svg)](https://build.snapcraft.io/user/goby-lang/goby)
[![Open Source Helpers](https://www.codetriage.com/goby-lang/goby/badges/users.svg)](https://www.codetriage.com/goby-lang/goby)

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
- [Structure](#structure)
- [Aspects](#aspects)
- [Features](#major-features)
- [Current roadmap](#current_roadmap)
- [Installation](#installation)
- [Usage](#usage)
- [Sample codes](#sample_codes)
- [Documentation](#documentation)
- [Joining to Goby](#joining-to-goby)
- [Maintainers](#maintainers)
- [Support us](#support-us)
- [References](#references)

## Demo screen and sample Goby app

<img src="https://i.imgur.com/1Le7nTe.gif" width="60%">

**New!** Check-out our [sample app](http://sample.goby-lang.org) built with Goby. Source code is also available [here](https://github.com/goby-lang/sample-web-app).

## Structure

![](https://github.com/goby-lang/goby/blob/master/wiki/goby_structure.png)

## Aspects

Goby has several aspects: language specification, design of compiler and vm, implementation (just one for now), library, and the whole of them. [See more](https://github.com/goby-lang/goby/wiki/Aspects)

We are optimizing and expanding Goby all the time. Toward the first release, we've been focusing on implementing Goby first.

## Major Features

- Plugin system
    - Allows to use Go libraries (packages) dynamically
    - Allows to call Go's methods from Goby directly (only on Linux for now)
- Builtin multi-threaded server and DB library
- REPL (run `goby -i`)

Here's a [complete list](https://github.com/goby-lang/goby/wiki/Features) of all the features.

## Current roadmap

- [ ] create functions for testing framework like `rescue`, `begin` -- by 2018 March
- [ ] testing framework -- by 2018 April
- [ ] third-party library support (like `rubygem`)

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
    - Install Golang >= 1.10
    - Make sure `$GOPATH` in your shell's config file( like .bashrc) is correct
    - Add you `$GOPATH/bin` to `$PATH`
2. Run `go get github.com/goby-lang/goby`
3. Set the Goby project's exact root path `$GOBY_ROOT` manually, which should be:

```
$GOPATH/src/github.com/goby-lang/goby
```

### C. Installation on a Linux system

In order to install Go, Goby and PostgreSQL on a Linux system, see the [wiki page](https://github.com/goby-lang/goby/wiki/Setup-Go,-Goby-and-PostgreSQL-on-a-Linux-system).

### Verifying Goby installation

1. Run `goby -v` to see the version.
2. Run `goby -i` to launch igb REPL.
3. Type `require "uri"` in igb.

FYI: You can just run `brew test goby` to check Homebrew installation.

**If you have any issue installing Goby, please let us know via [Github issues](https://github.com/goby-lang/goby/issues)**

### Using Docker

Goby has official [docker image](https://hub.docker.com/r/gobylang/goby/) as well. You can try the [Plugin System](https://goby-lang.gitbooks.io/goby/content/plugin-system.html) using docker.

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
- @saveriomiroddi
- @ear7h

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

### Powered by

* JetBrains [Goland IDE](https://www.jetbrains.com/go/)

[![JetBrains Goland](https://github.com/goby-lang/goby/blob/master/wiki/goland_logo-text.png)](https://www.jetbrains.com/go/)

**Supporting Goby by sending your first PR! See [contribution guideline](https://github.com/goby-lang/goby/blob/master/CONTRIBUTING.md)**

**Or [support us on opencollective](https://opencollective.com/goby)**

## References

The followings are the essential resources to create Goby; I highly recommend you to check them first if you'd be interested in building your own languages:

- [Write An Interpreter In Go](https://interpreterbook.com)
- [Nand2Tetris II](https://www.coursera.org/learn/nand2tetris2/home/welcome)
- [Ruby under a microscope](http://patshaughnessy.net/ruby-under-a-microscope)
- [YARV's instruction table](http://www.atdot.net/yarv/insnstbl.html)
