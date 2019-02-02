First of all, thank you for trying to contribute gooby, any contribution will be appreciated 😁

The following is the guideline (not rules) for contributing gooby, I suggest you read them all before you start your contribution.

If you are very interested in `Gooby` or planning contribute `Gooby` frequently, please contact me directly.

## What to contribute

- Any issues you see, and if you think the ticket is confusing, please open an issue or ask me on [slack](https://gooby-lang-slackin.herokuapp.com).
- Any grammar error in readme, wiki, and code comments...etc.
- Any issues litsted in gooby's [codeclimate](https://codeclimate.com/github/gooby-lang/gooby/issues).
- Play around gooby and report any bug you find.
- Write benchmarks for gooby (we really need this and really haven't have time to do it yet 😢)
- Help us document built in class and libraries' api, see the [guideline](https://github.com/gooby-lang/gooby/wiki/Documenting-Gooby-Code)


#### If you're interested in lexeing/parsing, please check `token`, `lexer`, `ast` and `parser` packages

#### If you're interested in compiler, check [bytecode specifications](https://github.com/gooby-lang/gooby/wiki/Bytecode-Instruction-specs) and bytecode package's [tests](https://github.com/gooby-lang/gooby/blob/master/bytecode/generator_test.go) for some compiled examples.

#### If yor're interested in VM's structure, please contact me directly since a lot things haven't been documented yet.

#### If you're a Ruby developer:
  - you can start with adding methods to built in classes like [`Array`](https://github.com/gooby-lang/gooby/blob/master/vm/array.go) or [`Hash`](https://github.com/gooby-lang/gooby/blob/master/vm/hash.go) using `Golang`. And here's a [guideline](https://github.com/gooby-lang/gooby/wiki/Contibuting-a-Method) for contributing built in methods.
  - you can also porting Ruby's standard lib using `Gooby` (not Go), see [lib directory](https://github.com/gooby-lang/gooby/tree/master/lib/net). You'll feel like you're just writing plain Ruby 😄

#### If you want to propose a feature, just open an issue with `[feature request]` prefix on title.

**Note**:
  - Before sending PR, you should perform `make test` on the root directory of the project to perform all tests (`go test` works only against gooby.go file and will be incomplete for the test).
  - DB library tests requires Postgresql to be opened and export port `5432`


## Setup Environment


### `$GOBY_ROOT`

By default Gooby finds standard libs in `/usr/local/gooby` when you install it via homebrew.
But if you want to develop Gooby or you installed Gooby from source, you might want to set `$GOBY_ROOT` to Gooby's project root so you can use latest libs.
Add the following line to your shell config file, either `~/.bashrc`, `~/.bash_profile`, or `~/.zshrc` if you're using zsh.

```
export GOBY_ROOT=$GOPATH/src/github.com/gooby-lang/gooby
```

The most common messages you'll see if you do not set `$GOBY_ROOT` right are 'library not found'. For example:

```ruby
require 'net/http'
# => Internal Error: open lib/net/http/response.gb: no such file or directory
```







