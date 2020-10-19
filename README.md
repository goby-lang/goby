![Goby](http://i.imgur.com/ElGAzRn.png?3)
=========

[![](https://goby-slack-invite.herokuapp.com/badge.svg)](https://goby-slack-invite.herokuapp.com)
[![Discord](https://img.shields.io/discord/678892628955103232?label=discord)](https://discord.gg/SS5HbYN)
[![Build Status](https://travis-ci.org/goby-lang/goby.svg?branch=master)](https://travis-ci.org/goby-lang/goby)
[![GoDoc](https://godoc.org/github.com/goby-lang/goby?status.svg)](https://godoc.org/github.com/goby-lang/goby)
[![Go Report Card](https://goreportcard.com/badge/github.com/goby-lang/goby)](https://goreportcard.com/report/github.com/goby-lang/goby)
[![codecov](https://codecov.io/gh/goby-lang/goby/branch/master/graph/badge.svg)](https://codecov.io/gh/goby-lang/goby)
[![Readme Score](http://readme-score-api.herokuapp.com/score.svg?url=goby-lang/goby)](http://clayallsopp.github.io/readme-score?url=goby-lang/goby)
[![Snap Status](https://build.snapcraft.io/badge/goby-lang/goby.svg)](https://build.snapcraft.io/user/goby-lang/goby)
[![Open Source Helpers](https://www.codetriage.com/goby-lang/goby/badges/users.svg)](https://www.codetriage.com/goby-lang/goby)
[![Reviewed by Hound](https://img.shields.io/badge/Reviewed_by-Hound-8E64B0.svg)](https://houndci.com)

**Goby** is an object-oriented interpreter language deeply inspired by **Ruby** as well as its core implementation by 100% pure **Go**. Moreover, it has standard libraries to provide several features such as the Plugin system. Note that we do not intend to reproduce whole of the honorable works of Ruby syntax/implementation/libraries.

The expected use case for Goby would be backend development. With this goal, it equips (but not limited to) the following features:

- thread/channel mechanism powered by Go's goroutine
- Builtin database library (currently only support PostgreSQL adapter)
- JSON support
- [Plugin system](https://goby-lang.gitbooks.io/goby/content/plugin-system.html) that can load existing Go packages dynamically (Only for Linux and MacOS right now)
- Accessing Go objects from Goby directly

> Note: Goby had formerly been known as "Rooby", which was renamed in May 2017.

## Table of contents

- [!Goby](#img-srchttpiimgurcomelgazrnpng3-altgoby)
  - [Table of contents](#table-of-contents)
  - [Demo screen and sample Goby app](#demo-screen-and-sample-goby-app)
  - [Structure](#structure)
  - [3D Visualization](#3d-visualization)
  - [Major Features](#major-features)
  - [Installation](#installation)
    - [A. Via Homebrew (binary installation for Mac OS)](#a-via-homebrew-binary-installation-for-mac-os)
    - [B. From Source](#b-from-source)
    - [C. Installation on a Linux system](#c-installation-on-a-linux-system)
    - [Verifying Goby installation](#verifying-goby-installation)
    - [Using Docker](#using-docker)
  - [Syntax highlighting](#syntax-highlighting)
    - [Sublime Text 3](#sublime-text-3)
    - [Vim](#vim)
    - [SpaceVim](#spacevim)
  - [Sample codes](#sample-codes)
  - [Joining to Goby](#joining-to-goby)
  - [Maintainers](#maintainers)
  - [Designer](#designer)
  - [Supporters](#supporters)
    - [Sponsors](#sponsors)
    - [Powered by](#powered-by)
  - [References](#references)

## Demo screen and sample Goby app

Click to see the demo below (powered by [asciinema](https://asciinema.org) and [GIPHY](https://giphy.com/)).

![](https://github.com/goby-lang/animation-gif/blob/master/goby_demo_large.gif)

**New!** Check-out our [sample app](http://sample.goby-lang.org) built with Goby. Source code is also available [here](https://github.com/goby-lang/sample-web-app).

## Structure

![](https://github.com/goby-lang/goby/blob/master/wiki/goby_structure.png)

## 3D Visualization

A 3D visualization of Goby codebase, powered by [GoCity](https://go-city.github.io/)

[![Goby 3D Visualization](https://github.com/goby-lang/goby/blob/master/wiki/goby_codebase_gocity-min.png)](https://go-city.github.io/#/github.com/goby-lang/goby)

## Major Features

- Plugin system
    - Allows using Go libraries (packages) dynamically
    - Allows calling Go's methods from Goby directly (only on Linux for now)
- Builtin multi-threaded server and DB library
- REPL (run `goby -i`)

Here's a [complete list](https://github.com/goby-lang/goby/wiki/Features) of all the features.

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
    - Install Golang >= 1.14
    - Make sure `$GOPATH` in your shell's config file( like .bashrc) is correct
    - Add your `$GOPATH/bin` to `$PATH`
    - Add `export GO111MODULE=on` to your shell profile
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

**If you have any issue installing Goby, please let us know via [GitHub issues](https://github.com/goby-lang/goby/issues)**

### Using Docker

Goby has official [docker image](https://hub.docker.com/r/gobylang/goby/) as well. You can try the [Plugin System](https://goby-lang.gitbooks.io/goby/content/plugin-system.html) using docker.

## Syntax highlighting

The Goby syntax is currently a subset of the Ruby one, with an exception (`get_block`), therefore, it's possible to attain syntax highlighting on any platform/editor by simply switching it to Ruby for the currently opened file.

### Sublime Text 3

Sublime Text 3 users can use the `Only Goby` package, by typing the following in a terminal:

```sh
git clone git@github.com:saveriomiroddi/only-goby-for-sublime-text "$HOME/.config/sublime-text-3/Packages/only-goby-for-sublime-text"
```

this will automatically apply the Goby syntax highlighting to the `.gb` files.

### Vim

Vim users can use the `vim-goby-syntax-highlighting` definition, by typing the following in a terminal:

```sh
mkdir -p "$HOME/.vim/syntax"
wget -O "$HOME/.vim/syntax/goby.vim" https://raw.githubusercontent.com/saveriomiroddi/vim-goby-syntax-highlighting/master/goby.vim
echo 'au BufNewFile,BufRead *.gb    setf goby' >> "$HOME/.vim/filetype.vim"
```

this will automatically apply the Goby syntax highlighting to the `.gb` files.

### SpaceVim

SpaceVim users can load the [`lang#goby`](https://spacevim.org/layers/lang/goby/) layer by adding following configuration:

```toml
[[layers]]
  name = "lang#goby"
```

## Sample codes

- [Built a stack data structure using Goby](https://github.com/goby-lang/goby/blob/master/samples/stack.gb)
- [Running a "Hello World" app with built in server library](https://github.com/goby-lang/goby/blob/master/samples/server/server.gb)
- [Sending request using http library](https://github.com/goby-lang/goby/blob/master/samples/http.gb)
- [Running load test on blocking server](https://github.com/goby-lang/goby/blob/master/samples/server/blocking_server.gb) (This shows `Goby`'s simple server is very performant and can handle requests concurrently)
- [One thousand threads](https://github.com/goby-lang/goby/blob/master/samples/one_thousand_threads.gb)

More sample Goby codes can be found in [sample directory](https://github.com/goby-lang/goby/tree/master/samples).

## Joining to Goby

**Join us on Slack!** [![](https://goby-slack-invite.herokuapp.com/badge.svg)](https://goby-slack-invite.herokuapp.com)

See the [guideline](https://github.com/goby-lang/goby/blob/master/CONTRIBUTING.md).

## Maintainers

- @st0012
- @hachi8833
- @saveriomiroddi

## Designer
- [steward379](https://dribbble.com/steward379)

## Supporters

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

* JetBrains [Goland IDE](https://www.jetbrains.com/go/?from=goby)

[![JetBrains Goland](https://github.com/goby-lang/goby/blob/master/wiki/goland_logo-text.png)](https://www.jetbrains.com/go/?from=goby)

**Supporting Goby by sending your first PR! See [contribution guideline](https://github.com/goby-lang/goby/blob/master/CONTRIBUTING.md)**


## References

The followings are the essential resources to create Goby; I highly recommend you to check them first if you'd be interested in building your own languages:

- [Write An Interpreter In Go](https://interpreterbook.com)
- [Nand2Tetris II](https://www.coursera.org/learn/nand2tetris2/home/welcome)
- [Ruby under a microscope](http://patshaughnessy.net/ruby-under-a-microscope)
- [YARV's instruction table](http://www.atdot.net/yarv/insnstbl.html)
