First of all, thank you for trying to contribute goby, any contribution will be appreciated 😁

The following is the guideline (not rules) for contributing goby, I suggest you read them all before you start your contribution.

If you are very interested in `Goby` or planning contribute `Goby` frequently, please contact me directly.

## What to contribute

- Any issues you see, and if you think the ticket is confusing, please open an issue or ask me on [slack](https://goby-lang-slackin.herokuapp.com).
- Any grammar error in readme, wiki, and code comments...etc.
- Any issues litsted in goby's [codeclimate](https://codeclimate.com/github/goby-lang/goby/issues).
- Play around goby and report any bug you find.
- Write benchmarks for goby (we really need this and really haven't have time to do it yet 😢)
- Help us document built in class and libraries' api, see the [guideline](https://github.com/goby-lang/goby/wiki/Documenting-Goby-Code)


#### If you're interested in lexeing/parsing, please check `token`, `lexer`, `ast` and `parser` packages

#### If you're interested in compiler, check [bytecode specifications](https://github.com/goby-lang/goby/wiki/Bytecode-Instruction-specs) and bytecode package's [tests](https://github.com/goby-lang/goby/blob/master/bytecode/generator_test.go) for some compiled examples.

#### If yor're interested in VM's structure, please contact me directly since a lot things haven't been documented yet.

#### If you're a Ruby developer:
  - you can start with adding methods to built in classes like [`Array`](https://github.com/goby-lang/goby/blob/master/vm/array.go) or [`Hash`](https://github.com/goby-lang/goby/blob/master/vm/hash.go) using `Golang`. And here's a [guideline](https://github.com/goby-lang/goby/wiki/Contibuting-a-Method) for contributing built in methods.
  - you can also porting Ruby's standard lib using `Goby` (not Go), see [lib directory](https://github.com/goby-lang/goby/tree/master/lib/net). You'll feel like you're just writing plain Ruby 😄

#### If you want to propose a feature, just open an issue with `[feature request]` prefix on title.

**Note**: Before sending PR, you should perform `make test` on the root directory of the project to perform all tests (`go test` works only against goby.go file and will be incomplete for the test).










