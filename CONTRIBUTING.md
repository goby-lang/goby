First of all, thank you for trying to contribute rooby, any contribution will be appreciated 😁

The following is the guideline (not rules) for contributing rooby, I suggest you read them all before you start your contribution.

## What to contribute

- Any ticket listed in [GitHub Projects](https://github.com/rooby-lang/rooby/projects). And if you think the ticket is confusing, please open an issue or ask me on [gitter](https://gitter.im/Rooby-lang/Lobby).
- Any issues that has `bug` label.
- Any grammar error in readme, wiki, and code comments...etc.
- Any issues litsted in rooby's [codeclimate](https://codeclimate.com/github/rooby-lang/rooby/issues).
- Play around rooby and report any bug you find.
- Write benchmarks for rooby (we really need this and really haven't have time to do it yet 😢)


**If you're interested in lexeing/parsing, please check `token`, `lexer`, `ast` and `parser` packages**

**If you're interested in compiler, check `bytecode` package and [its tests](https://github.com/rooby-lang/rooby/blob/master/bytecode/generator_test.go) for some compiling examples.**

**If yor're interested in VM, please contact with me directly I will tell you about how it works and bytecode specifications.**

**If you're a Ruby developer, you can start with adding methods to built in classes like [`Array`](https://github.com/rooby-lang/rooby/blob/master/vm/array.go) or [`Hash`](https://github.com/rooby-lang/rooby/blob/master/vm/hash.go).**

## If you want to propose a feature

Open an issue with `[feature request]` prefix on title.

## How add new method to `Array` (or another Object)

First of all we need to choose some method. For example `index`:

```
a = ["a", "b", "c", "d", 2] # create an array
a.index("a") # get index of "a" it's will be 0

# or index can get block

c = a.index do |x|
    x == "c"
end

c # will be 2
```

Second we need to add this method to [`Array`](https://github.com/rooby-lang/rooby/blob/master/vm/array.go) to `builtinArrayMethods`.

```
{
    // receiver it's our Array a in example.
    Fn: func(receiver Object) builtinMethodBody {
        // vm - it's pointer to VM
        // args - it's array of arguments in "()": a.index("c") args will be:
        //
        // []Object{
        //   0: StringObject{
        //       Class: *RString
        //       Value: "c"
        //   }
        // }
        //
        // blockFrame it's our block, it will be nil if block empty.
        return func(vm *VM, args []Object, blockFrame *callFrame) Object {
            arr := receiver.(*ArrayObject) // get our Array ["a", "b", "c", "d", 2]

            arg = args[0] // get object what we are looking for
            // now we need check type of object
            elInt, isInt := arg.(*IntegerObject)
            elStr, isStr := arg.(*StringObject)

            // 'index' searches the same object in array, and after finding it returns index
            // i - index of element, o - object to compare with arg
            for i, o := range arr.Elements {
                switch o.(type) {
                case *IntegerObject:
                    el := o.(*IntegerObject)
                    if isInt && el.equal(elInt) { // if object from args is IntegerObject and the same as object in array then returns IntegerObject with i
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

After implementation this method, we need to add tests in [`ArrayTest`](https://github.com/rooby-lang/rooby/blob/master/vm/array_test.go)
Create new function `TestIndexMethod`:

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




