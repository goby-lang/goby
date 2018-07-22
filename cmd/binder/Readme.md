# binder

A helper cli for generating goby class bindings for go structures.

## Usage
binder -in file_name.go -type MyGoType

This will create a file named `bindings.go` which contains wrapper functions and an init function which will load those bindings into the vm at runtime.

To ensure your package is loaded in the vm, include a null import to your package in the main file of your finial binary.

```go
import _ "github.com/path/to/your/package"
```

## Auto generate using go generate

Adding a go generate comment in the file your define your go structure will allow you to automatically generate updated bindings easily by running `go generate ./...` in the root of your project.

example in
`github.com/goby-lang/native/result/result.go`
```go
package result
//go:generate binder -in result.go -type Result
```

## Binding rules

* Methods with a named receiver will be instance methods.
* Methods without a named receiver will be class methods.
* Types must be exported.
* Camel case names will be converted to snake case names.

examples.
```go
func (t *MyType) func MyFunc() vm.Object
```

will generate the equivalent instance method in goby.
```ruby
class MyType
    def my_func()
    end
end
```


```go
func (MyType) func MyFunc() vm.Object
```

will generate the equivalent class method in goby.
```ruby
class MyType
    def self.my_func()
    end
end
```

## Current Limitations

* Only one type can have generated bindings per package.
* Only functions that return `vm.Object` will have bindings generated.
* Function names cannot contain special characters like `?`.
