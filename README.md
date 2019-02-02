![Gooby](http://i.imgur.com/ElGAzRn.png?3)
=========

[![](https://gooby-lang-slackin.herokuapp.com/badge.svg)](https://gooby-lang-slackin.herokuapp.com)
[![Build Status](https://travis-ci.org/gooby-lang/gooby.svg?branch=master)](https://travis-ci.org/gooby-lang/gooby)
[![GoDoc](https://godoc.org/github.com/gooby-lang/gooby?status.svg)](https://godoc.org/github.com/gooby-lang/gooby)
[![Go Report Card](https://goreportcard.com/badge/github.com/gooby-lang/gooby)](https://goreportcard.com/report/github.com/gooby-lang/gooby)
[![codecov](https://codecov.io/gh/gooby-lang/gooby/branch/master/graph/badge.svg)](https://codecov.io/gh/gooby-lang/gooby)
[![Readme Score](http://readme-score-api.herokuapp.com/score.svg?url=gooby-lang/gooby)](http://clayallsopp.github.io/readme-score?url=gooby-lang/gooby)
[![Snap Status](https://build.snapcraft.io/badge/gooby-lang/gooby.svg)](https://build.snapcraft.io/user/gooby-lang/gooby)
[![Open Source Helpers](https://www.codetriage.com/gooby-lang/gooby/badges/users.svg)](https://www.codetriage.com/gooby-lang/gooby)
[![Reviewed by Hound](https://img.shields.io/badge/Reviewed_by-Hound-8E64B0.svg)](https://houndci.com)

**Gooby** is an object-oriented interpreter language deeply inspired by **Ruby** as well as its core implementation by 100% pure **Go**. Moreover, it has standard libraries to provide several features such as the Plugin system. Note that we do not intend to reproduce whole of the honorable works of Ruby syntax/implementation/libraries.

One of our goal is to provide web developers a sort of small and handy environment that mainly focusing on creating **API servers or microservices**. For this, Gooby includes the following native features:

- Robust thread/channel mechanism powered by Go's goroutine
- Builtin high-performance HTTP server
- Builtin database library (currently only support PostgreSQL adapter)
- JSON support
- [Plugin system](https://gooby-lang.gitbooks.io/gooby/content/plugin-system.html) that can load existing Go packages dynamically (Only for Linux and MacOS right now)
- Accessing Go objects from Gooby directly

> Note: Gooby had formerly been known as "Rooby", which was renamed in May 2017.

## Table of contents

- [Demo and sample Gooby app](#demo_and_sample_app)
- [Structure](#structure)
- [Aspects](#aspects)
- [Features](#major-features)
- [Current roadmap](#current_roadmap)
- [Installation](#installation)
- [Usage](#usage)
- [Sample codes](#sample_codes)
- [Documentation](https://gooby-lang.org/docs/introduction.html)
- [Joining to Gooby](#joining-to-gooby)
- [Maintainers](#maintainers)
- [Support us](#support-us)
- [References](#references)

## Demo screen and sample Gooby app

<img src="https://i.imgur.com/1Le7nTe.gif" width="60%">

**New!** Check-out our [sample app](http://sample.gooby-lang.org) built with Gooby. Source code is also available [here](https://github.com/gooby-lang/sample-web-app).

## Structure

![](https://github.com/gooby-lang/gooby/blob/master/wiki/gooby_structure.png)

## Aspects

Gooby has several aspects: language specification, design of compiler and vm, implementation (just one for now), library, and the whole of them. [See more](https://github.com/gooby-lang/gooby/wiki/Aspects)

We are optimizing and expanding Gooby all the time. Toward the first release, we've been focusing on implementing Gooby first.

## Major Features

- Plugin system
    - Allows using Go libraries (packages) dynamically
    - Allows calling Go's methods from Gooby directly (only on Linux for now)
- Builtin multi-threaded server and DB library
- REPL (run `gooby -i`)

Here's a [complete list](https://github.com/gooby-lang/gooby/wiki/Features) of all the features.

## Current roadmap

See wiki: [Current roadmap](https://github.com/gooby-lang/gooby/wiki/Current-Roadmap)

## Installation

Confirmed Gooby runs on Mac OS and Linux for now. Try Gooby on Windows and let us know the result.

### A. Via Homebrew (binary installation for Mac OS)

**Note: Please check the [latest release](https://github.com/gooby-lang/gooby/releases) before installing Gooby via Homebrew**

```
brew tap gooby-lang/gooby
brew install gooby
```

In the case, `$GOBY_ROOT` is automatically configured.

### B. From Source

Try this if you'd like to contribute Gooby! Skip 1 if you already have Golang in your environment.

1. Prepare Golang environment
    - Install Golang >= 1.10
    - Make sure `$GOPATH` in your shell's config file( like .bashrc) is correct
    - Add you `$GOPATH/bin` to `$PATH`
2. Run `go get github.com/gooby-lang/gooby`
3. Set the Gooby project's exact root path `$GOBY_ROOT` manually, which should be:

```
$GOPATH/src/github.com/gooby-lang/gooby
```

### C. Installation on a Linux system

In order to install Go, Gooby and PostgreSQL on a Linux system, see the [wiki page](https://github.com/gooby-lang/gooby/wiki/Setup-Go,-Gooby-and-PostgreSQL-on-a-Linux-system).

### Verifying Gooby installation

1. Run `gooby -v` to see the version.
2. Run `gooby -i` to launch igb REPL.
3. Type `require "uri"` in igb.

FYI: You can just run `brew test gooby` to check Homebrew installation.

**If you have any issue installing Gooby, please let us know via [GitHub issues](https://github.com/gooby-lang/gooby/issues)**

### Using Docker

Gooby has official [docker image](https://hub.docker.com/r/goobylang/gooby/) as well. You can try the [Plugin System](https://gooby-lang.gitbooks.io/gooby/content/plugin-system.html) using docker.

## Syntax highlighting

The Gooby syntax is currently a subset of the Ruby one, with an exception (`get_block`), therefore, it's possible to attain syntax highlighting on any platform/editor by simply switching it to Ruby for the currently opened file.

### Sublime Text 3

Sublime Text 3 users can use the `Only Gooby` package, by typing the following in a terminal:

```sh
git clone git@github.com:saveriomiroddi/only-gooby-for-sublime-text "$HOME/.config/sublime-text-3/Packages/only-gooby-for-sublime-text"
```

this will automatically apply the Gooby syntax highlighting to the `.gb` files.

### Vim

Vim users can use the `vim-gooby-syntax-highlighting` definition, by typing the following in a terminal:

```sh
mkdir -p "$HOME/.vim/syntax"
wget -O "$HOME/.vim/syntax/gooby.vim" https://raw.githubusercontent.com/saveriomiroddi/vim-gooby-syntax-highlighting/master/gooby.vim
echo 'au BufNewFile,BufRead *.gb    setf gooby' >> "$HOME/.vim/filetype.vim"
```

this will automatically apply the Gooby syntax highlighting to the `.gb` files.

## Sample codes

- [Built a stack data structure using Gooby](https://github.com/gooby-lang/gooby/blob/master/samples/stack.gb)
- [Running a "Hello World" app with built in server library](https://github.com/gooby-lang/gooby/blob/master/samples/server/server.gb)
- [Sending request using http library](https://github.com/gooby-lang/gooby/blob/master/samples/http.gb)
- [Running load test on blocking server](https://github.com/gooby-lang/gooby/blob/master/samples/server/blocking_server.gb) (This shows `Gooby`'s simple server is very performant and can handle requests concurrently)
- [One thousand threads](https://github.com/gooby-lang/gooby/blob/master/samples/one_thousand_threads.gb)

More sample Gooby codes can be found in [sample directory](https://github.com/gooby-lang/gooby/tree/master/samples).

## Joining to Gooby

**Join us on Slack!** [![](https://gooby-lang-slackin.herokuapp.com/badge.svg)](https://gooby-lang-slackin.herokuapp.com)

See the [guideline](https://github.com/gooby-lang/gooby/blob/master/CONTRIBUTING.md).

## Maintainers

- @st0012
- @hachi8833
- @saveriomiroddi
- @ear7h

## Designer
- [steward379](https://dribbble.com/steward379)

## Support Us

### Donations

Support us with a monthly donation and help us continue our activities. [[Become a backer](https://opencollective.com/gooby#backer)]

<a href="https://opencollective.com/gooby/backer/0/website" target="_blank"><img src="https://opencollective.com/gooby/backer/0/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/1/website" target="_blank"><img src="https://opencollective.com/gooby/backer/1/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/2/website" target="_blank"><img src="https://opencollective.com/gooby/backer/2/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/3/website" target="_blank"><img src="https://opencollective.com/gooby/backer/3/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/4/website" target="_blank"><img src="https://opencollective.com/gooby/backer/4/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/5/website" target="_blank"><img src="https://opencollective.com/gooby/backer/5/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/6/website" target="_blank"><img src="https://opencollective.com/gooby/backer/6/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/7/website" target="_blank"><img src="https://opencollective.com/gooby/backer/7/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/8/website" target="_blank"><img src="https://opencollective.com/gooby/backer/8/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/9/website" target="_blank"><img src="https://opencollective.com/gooby/backer/9/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/10/website" target="_blank"><img src="https://opencollective.com/gooby/backer/10/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/11/website" target="_blank"><img src="https://opencollective.com/gooby/backer/11/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/12/website" target="_blank"><img src="https://opencollective.com/gooby/backer/12/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/13/website" target="_blank"><img src="https://opencollective.com/gooby/backer/13/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/14/website" target="_blank"><img src="https://opencollective.com/gooby/backer/14/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/15/website" target="_blank"><img src="https://opencollective.com/gooby/backer/15/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/16/website" target="_blank"><img src="https://opencollective.com/gooby/backer/16/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/17/website" target="_blank"><img src="https://opencollective.com/gooby/backer/17/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/18/website" target="_blank"><img src="https://opencollective.com/gooby/backer/18/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/19/website" target="_blank"><img src="https://opencollective.com/gooby/backer/19/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/20/website" target="_blank"><img src="https://opencollective.com/gooby/backer/20/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/21/website" target="_blank"><img src="https://opencollective.com/gooby/backer/21/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/22/website" target="_blank"><img src="https://opencollective.com/gooby/backer/22/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/23/website" target="_blank"><img src="https://opencollective.com/gooby/backer/23/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/24/website" target="_blank"><img src="https://opencollective.com/gooby/backer/24/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/25/website" target="_blank"><img src="https://opencollective.com/gooby/backer/25/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/26/website" target="_blank"><img src="https://opencollective.com/gooby/backer/26/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/27/website" target="_blank"><img src="https://opencollective.com/gooby/backer/27/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/28/website" target="_blank"><img src="https://opencollective.com/gooby/backer/28/avatar.svg"></a>
<a href="https://opencollective.com/gooby/backer/29/website" target="_blank"><img src="https://opencollective.com/gooby/backer/29/avatar.svg"></a>

### Sponsors

<a href="https://opencollective.com/gooby/sponsor/0/website" target="_blank"><img src="https://opencollective.com/gooby/sponsor/0/avatar.svg"></a>
<a href="https://opencollective.com/gooby/sponsor/1/website" target="_blank"><img src="https://opencollective.com/gooby/sponsor/1/avatar.svg"></a>
<a href="https://opencollective.com/gooby/sponsor/2/website" target="_blank"><img src="https://opencollective.com/gooby/sponsor/2/avatar.svg"></a>
<a href="https://opencollective.com/gooby/sponsor/3/website" target="_blank"><img src="https://opencollective.com/gooby/sponsor/3/avatar.svg"></a>
<a href="https://opencollective.com/gooby/sponsor/4/website" target="_blank"><img src="https://opencollective.com/gooby/sponsor/4/avatar.svg"></a>
<a href="https://opencollective.com/gooby/sponsor/5/website" target="_blank"><img src="https://opencollective.com/gooby/sponsor/5/avatar.svg"></a>
<a href="https://opencollective.com/gooby/sponsor/6/website" target="_blank"><img src="https://opencollective.com/gooby/sponsor/6/avatar.svg"></a>
<a href="https://opencollective.com/gooby/sponsor/7/website" target="_blank"><img src="https://opencollective.com/gooby/sponsor/7/avatar.svg"></a>
<a href="https://opencollective.com/gooby/sponsor/8/website" target="_blank"><img src="https://opencollective.com/gooby/sponsor/8/avatar.svg"></a>
<a href="https://opencollective.com/gooby/sponsor/9/website" target="_blank"><img src="https://opencollective.com/gooby/sponsor/9/avatar.svg"></a>

### Powered by

* JetBrains [Goland IDE](https://www.jetbrains.com/go/)

[![JetBrains Goland](https://github.com/gooby-lang/gooby/blob/master/wiki/goland_logo-text.png)](https://www.jetbrains.com/go/)

**Supporting Gooby by sending your first PR! See [contribution guideline](https://github.com/gooby-lang/gooby/blob/master/CONTRIBUTING.md)**

**Or [support us on opencollective](https://opencollective.com/gooby)**

## References

The followings are the essential resources to create Gooby; I highly recommend you to check them first if you'd be interested in building your own languages:

- [Write An Interpreter In Go](https://interpreterbook.com)
- [Nand2Tetris II](https://www.coursera.org/learn/nand2tetris2/home/welcome)
- [Ruby under a microscope](http://patshaughnessy.net/ruby-under-a-microscope)
- [YARV's instruction table](http://www.atdot.net/yarv/insnstbl.html)
