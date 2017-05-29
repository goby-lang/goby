First of all, thank you for trying to contribute goby, any contribution will be appreciated 😁

The following is the guideline (not rules) for contributing goby, I suggest you read them all before you start your contribution.

And I hope all contributors can join our [gitter chatroom](https://gitter.im/goby-lang/Lobby).

## What to contribute

- Any issues you see, and if you think the ticket is confusing, please open an issue or ask me on [gitter](https://gitter.im/goby-lang/Lobby).
- Any grammar error in readme, wiki, and code comments...etc.
- Any issues litsted in goby's [codeclimate](https://codeclimate.com/github/goby-lang/goby/issues).
- Play around goby and report any bug you find.
- Write benchmarks for goby (we really need this and really haven't have time to do it yet 😢)
- Help us document built in class and libraries' api, see the [guideline](https://github.com/goby-lang/goby/wiki/Documenting-Goby-Code)


**If you're interested in lexeing/parsing, please check `token`, `lexer`, `ast` and `parser` packages**

**If you're interested in compiler, check [bytecode specifications](https://github.com/goby-lang/goby/wiki/Bytecode-Instruction-specs) and bytecode package's [tests](https://github.com/goby-lang/goby/blob/master/bytecode/generator_test.go) for some compiled examples.**

**If yor're interested in VM, please contact me directly since a lot things haven't been documented yet.**

**If you're a Ruby developer, you can start with adding methods to built in classes like [`Array`](https://github.com/goby-lang/goby/blob/master/vm/array.go) or [`Hash`](https://github.com/goby-lang/goby/blob/master/vm/hash.go).
And here's a [guideline](https://github.com/goby-lang/goby/wiki/Contibuting-a-Method) for contributing built in methods**

## If you want to propose a feature

Open an issue with `[feature request]` prefix on title.

## How add new method to `Array` (or another Object)

First of all we need to choose a method. For example `index`:

```
a = ["a", "b", "c", "d", 2] # create an array
a.index("a") # get index of "a" it's will be 0

# or index can get block

c = a.index do |x|
    x == "c"
end

c # will be 2
```

Then we need to add this method to [`vm/array.go`](https://github.com/rooby-lang/rooby/blob/master/vm/array.go)'s `builtinArrayMethods`.

```
{
    // receiver it's our Array a in example.
    Fn: func(receiver Object) builtinMethodBody {
        // vm - it's a pointer to VM
        // args - it's an array of arguments in "()": a.index("c") args will be:
        //
        // []Object{
        //   0: StringObject{
        //       Class: *RString
        //       Value: "c"
        //   }
        // }
        //
        // blockFrame it's our block argument, it will be nil if there's no block.
        return func(vm *VM, args []Object, blockFrame *callFrame) Object {
            arr := receiver.(*ArrayObject) // get our Array ["a", "b", "c", "d", 2]

            arg = args[0] // get the object we are looking for
            // now we need to check the type of object
            elInt, isInt := arg.(*IntegerObject)
            elStr, isStr := arg.(*StringObject)

            // 'index' searches given object in an array, and returns it's index after finding it
            // i - index of element, o - object to compare with arg
            for i, o := range arr.Elements {
                switch o.(type) {
                case *IntegerObject:
                    el := o.(*IntegerObject)
                    if isInt && el.equal(elInt) { // if both objects are integer then returns IntegerObject with i
                        return initilaizeInteger(i)
                    }
                case *StringObject:
                    el := o.(*StringObject)
                    if isStr && el.equal(elStr) {
                        return initilaizeInteger(i)
                    }
                }
            }

            return initilaizeInteger(nil)
        }
    },
    Name: "index",
}
```

After implementating this method, we need to add tests in [`Array's Test`](https://github.com/rooby-lang/rooby/blob/master/vm/array_test.go) by creating a new function `TestIndexMethod`:

```
func TestIndexMethod(t *testing.T) {

    tests := []struct {
        input    string
        expected *IntegerObject
    }{
        {`
        a = [1, 2, "a", 3, "5", "r"]
        a.index("r")
        `, initilaizeInteger(5)},
    }

    for _, tt := range tests {
        evaluated := testEval(t, tt.input)
        testIntegerObject(t, evaluated, tt.expected.Value)
    }
}

```




