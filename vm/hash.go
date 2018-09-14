package vm

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// HashObject represents hash instances
// Hash is a collection of key-value pair, which works like a dictionary.
// Hash literal is represented with curly brackets `{ }` like `{ key: value }`.
// Each key of the hash is unique and cannot be duplicate within the hash.
// Adding a leading space and a trailing space within curly brackets are preferable.
//
// - **Key:** an alphanumeric word that starts with alphabet, without containing space and punctuations.
// Underscore `_` can also be used within the key.
// String literal like "mickey mouse" cannot be used as a hash key.
// The internal key is actually a String and **not a Symbol** for now (TBD).
// Thus only a String object or a string literal should be used when referencing with `[ ]`.
//
// ```ruby
// a = { balthazar1: 100 } # valid
// b = { 2melchior: 200 }  # invalid
// x = 'balthazar1'
//
// a["balthazar1"]  # => 100
// a[x]             # => 100
// a[balthazar1]    # => error
// ```
//
// - **value:** String literal and objects (Integer, String, Array, Hash, nil, etc) can be used.
//
// **Note:**
// - The order of key-value pairs are **not** preserved.
// - Operator `=>` is not supported.
// - `Hash.new` is not supported.
type HashObject struct {
	*BaseObj
	Pairs map[string]Object

	// See `[]` and `[]=` for the operational explanation of the default value.
	Default Object
}

// Class methods --------------------------------------------------------
func builtinHashClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				return t.vm.InitNoMethodError(sourceLine, "new", receiver)

			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinHashInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Retrieves the value (object) that corresponds to the key specified.
			// When a key doesn't exist, `nil` is returned, or the default, if set.
			//
			// ```Ruby
			// h = { a: 1, b: "2", c: [1, 2, 3], d: { k: 'v' } }
			// h['a'] #=> 1
			// h['b'] #=> "2"
			// h['c'] #=> [1, 2, 3]
			// h['d'] #=> { k: 'v' }
			//
			// h = { a: 1 }
			// h['c']        #=> nil
			// h.default = 0
			// h['c']        #=> 0
			// h             #=> { a: 1 }
			// h['d'] += 2
			// h             #=> { a: 1, d: 2 }
			// ```
			//
			// @return [Object]
			Name: "[]",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
				}

				i := args[0]
				key, ok := i.(*StringObject)

				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, i.Class().Name)
				}

				h := receiver.(*HashObject)

				value, ok := h.Pairs[key.value]

				if !ok {
					if h.Default != nil {
						return h.Default
					}

					return NULL
				}

				return value

			},
		},
		{
			// Associates the value given by `value` with the key given by `key`.
			// Returns the `value`.
			//
			// ```Ruby
			// h = { a: 1, b: "2", c: [1, 2, 3], d: { k: 'v' } }
			// h['a'] = 1          #=> 1
			// h['b'] = "2"        #=> "2"
			// h['c'] = [1, 2, 3]  #=> [1, 2, 3]
			// h['d'] = { k: 'v' } #=> { k: 'v' }
			// ```
			//
			// @return [Object] The value
			Name: "[]=",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				// First arg is index
				// Second arg is assigned value
				if len(args) != 2 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 2, len(args))
				}
				k := args[0]
				key, ok := k.(*StringObject)

				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, k.Class().Name)
				}

				h := receiver.(*HashObject)
				h.Pairs[key.value] = args[1]

				return args[1]

			},
		},
		{
			// Passes each (key, value) pair  of the collection to the given block. The method returns
			// true if the block ever returns a value other than false or nil.
			//
			// ```ruby
			// a = { a: 1, b: 2 }
			//
			// a.any? do |k, v|
			//   v == 2
			// end            # => true
			// a.any? do |k, v|
			//   v
			// end            # => true
			// a.any? do |k, v|
			//   v == 5
			// end            # => false
			// a.any? do |k, v|
			//   nil
			// end            # => false
			//
			// a = {}
			//
			// a.any? do |k, v|
			//   true
			// end            # => false
			// ```
			Name: "any?",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
				}

				if blockFrame == nil {
					return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
				}

				hash := receiver.(*HashObject)
				if blockIsEmpty(blockFrame) {
					return FALSE
				}

				if len(hash.Pairs) == 0 {
					t.callFrameStack.pop()
				}

				for stringKey, value := range hash.Pairs {
					objectKey := t.vm.InitStringObject(stringKey)
					result := t.builtinMethodYield(blockFrame, objectKey, value)

					/*
						TODO: Discuss this behavior

						```ruby
						{ key: "foo", bar: "baz" }.any? do |k, v|
						  true
						  break
						end
						```

						The block returns nil because of the break.
						But in Ruby the final result is nil, which means the block's result is completely ignored
					*/
					if blockFrame.IsRemoved() {
						return NULL
					}

					if result.Target.isTruthy() {
						return TRUE
					}
				}

				return FALSE

			},
		},
		{
			// Returns empty hash (no key-value pairs)
			//
			// ```Ruby
			// { a: "Hello", b: "World" }.clear # => {}
			// {}.clear                         # => {}
			// ```
			//
			// @return [Boolean]
			Name: "clear",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
				}

				h := receiver.(*HashObject)

				h.Pairs = make(map[string]Object)

				return h

			},
		},
		{
			// Return the default value of this Hash.
			//
			// ```Ruby
			// h = { a: 1 }
			// h.default     #=> nil
			// h.default = 2
			// h.default     #=> 2
			// ```
			//
			// @return [Object]
			Name: "default",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
				}

				hash := receiver.(*HashObject)

				if hash.Default == nil {
					return NULL
				}

				return hash.Default

			},
		},
		{
			// Set the default value of this Hash.
			// Arrays/Hashes are not accepted, since they're unsafe.
			//
			// ```Ruby
			// h = { a: 1 }
			// h['c']         #=> nil
			// h.default = 2
			// h['c']         #=> 2
			// h.default = [] #=> ArgumentError
			// ```
			//
			// @return [Object]
			Name: "default=",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
				}

				// Arrays and Hashes are generally a mistake, since a single instance would be used for all the accesses
				// via default.
				switch args[0].(type) {
				case *HashObject, *ArrayObject:
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "Arrays and Hashes are not accepted as default values")
				}

				hash := receiver.(*HashObject)
				hashDefault := args[0]

				hash.Default = hashDefault

				return hashDefault

			},
		},
		{
			// Remove the key from the hash if key exist
			//
			// ```Ruby
			// h = { a: 1, b: 2, c: 3 }
			// h.delete("b") # =>  { a: 1, c: 3 }
			// ```
			//
			// @return [Hash]
			Name: "delete",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
				}

				h := receiver.(*HashObject)
				d := args[0]
				deleteKey, ok := d.(*StringObject)

				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, d.Class().Name)
				}

				deleteKeyValue := deleteKey.value
				if _, ok := h.Pairs[deleteKeyValue]; ok {
					delete(h.Pairs, deleteKeyValue)
				}
				return h

			},
		},
		{
			// Deletes every key-value pair from the hash for which block evaluates to anything except
			// false and nil.
			//
			// Returns the hash.
			//
			// ```Ruby
			// { a: 1, b: 2}.delete_if do |k, v| v == 1 end # =>  { b: 2 }
			// { a: 1, b: 2}.delete_if do |k, v| 5 end      # =>  { }
			// { a: 1, b: 2}.delete_if do |k, v| false end  # =>  { a: 1, b: 2}
			// { a: 1, b: 2}.delete_if do |k, v| nil end    # =>  { a: 1, b: 2}
			// ```
			//
			// @return [Hash]
			Name: "delete_if",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
				}

				if blockFrame == nil {
					return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
				}

				hash := receiver.(*HashObject)
				if blockIsEmpty(blockFrame) {
					return hash
				}

				if len(hash.Pairs) == 0 {
					t.callFrameStack.pop()
				}

				// Note that from the Go specification, https://golang.org/ref/spec#For_statements,
				// it's safe to delete elements from a Map, while iterating it.
				for stringKey, value := range hash.Pairs {
					objectKey := t.vm.InitStringObject(stringKey)
					result := t.builtinMethodYield(blockFrame, objectKey, value)

					booleanResult, isResultBoolean := result.Target.(*BooleanObject)

					if isResultBoolean {
						if booleanResult.value {
							delete(hash.Pairs, stringKey)
						}
					} else if result.Target != NULL {
						delete(hash.Pairs, stringKey)
					}
				}

				return hash

			},
		},
		{
			// Extracts the nested value specified by the sequence of idx objects by calling `dig` at
			// each step, returning nil if any intermediate step is nil.
			//
			// ```Ruby
			// { a: 1 , b: 2 }.dig(:a)         # => 1
			// { a: {}, b: 2 }.dig(:a, :b)     # => nil
			// { a: {}, b: 2 }.dig(:a, :b, :c) # => nil
			// { a: 1, b: 2 }.dig(:a, :b)      # => TypeError: Expect target to be Diggable
			// ```
			//
			// @return [Object]
			Name: "dig",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) < 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentMore, 1, len(args))
				}

				hash := receiver.(*HashObject)
				value := hash.dig(t, args, sourceLine)

				return value

			},
		},
		{
			// Calls block once for each key in the hash (in sorted key order), passing the
			// key-value pair as parameters.
			// Returns `self`.
			//
			// ```Ruby
			// h = { b: "2", a: 1 }
			// h.each do |k, v|
			//   puts k.to_s + "->" + v.to_s
			// end
			// # => a->1
			// # => b->2
			// ```
			//
			// @return [Hash]
			Name: "each",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
				}

				if blockFrame == nil {
					return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
				}

				h := receiver.(*HashObject)

				if len(h.Pairs) == 0 {
					t.callFrameStack.pop()
				} else {
					keys := h.sortedKeys()

					for _, k := range keys {
						v := h.Pairs[k]
						strK := t.vm.InitStringObject(k)

						t.builtinMethodYield(blockFrame, strK, v)
					}
				}

				return h

			},
		},
		{
			// Loop through keys of the hash with given block frame. It also returns array of
			// keys in alphabetical order.
			//
			// ```Ruby
			// h = { a: 1, b: "2", c: [1, 2, 3], d: { k: 'v' } }
			// h.each_key do |k|
			//   puts k
			// end
			// # => a
			// # => b
			// # => c
			// # => d
			// ```
			//
			// @return [Array]
			Name: "each_key",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
				}

				if blockFrame == nil {
					return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
				}

				h := receiver.(*HashObject)

				if len(h.Pairs) == 0 {
					t.callFrameStack.pop()
				}

				keys := h.sortedKeys()
				var arrOfKeys []Object

				for _, k := range keys {
					obj := t.vm.InitStringObject(k)
					arrOfKeys = append(arrOfKeys, obj)
					t.builtinMethodYield(blockFrame, obj)
				}

				return t.vm.InitArrayObject(arrOfKeys)

			},
		},
		{
			// Loop through values of the hash with given block frame. It also returns array of
			// values of the hash in the alphabetical order of its key
			//
			// ```Ruby
			// h = { a: 1, b: "2", c: [1, 2, 3], d: { k: "v" } }
			// h.each_value do |v|
			//   puts v
			// end
			// # => 1
			// # => "2"
			// # => [1, 2, 3]
			// # => { k: "v" }
			// ```
			//
			Name: "each_value",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
				}

				if blockFrame == nil {
					return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
				}

				h := receiver.(*HashObject)

				if len(h.Pairs) == 0 {
					t.callFrameStack.pop()
				}

				keys := h.sortedKeys()
				var arrOfValues []Object

				for _, k := range keys {
					value := h.Pairs[k]
					arrOfValues = append(arrOfValues, value)
					t.builtinMethodYield(blockFrame, value)
				}

				return t.vm.InitArrayObject(arrOfValues)

			},
		},
		{
			// Returns true if hash has no key-value pairs
			//
			// ```Ruby
			// {}.empty?       # => true
			// { a: 1 }.empty? # => false
			// ```
			//
			// @return [Boolean]
			Name: "empty?",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
				}

				h := receiver.(*HashObject)
				if h.length() == 0 {
					return TRUE
				}
				return FALSE

			},
		},
		{
			// Returns true if hash is exactly equal to another hash
			//
			// ```Ruby
			// { a: "Hello", b: "World" }.eql?(1) # => false
			// ```
			//
			// @return [Boolean]
			Name: "eql?",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
				}

				h := receiver.(*HashObject)
				c := args[0]
				compare, ok := c.(*HashObject)

				if ok && reflect.DeepEqual(h, compare) {
					return TRUE
				}
				return FALSE

			},
		},
		{
			// Returns a value from the hash for the given key. If the key can’t be found, there are several
			// options: With no other arguments, it will raise an ArgumentError exception; if default is
			// given, then that will be returned; if the optional code block is specified, then that will be
			// run and its result returned.
			//
			// ```Ruby
			// h = { "spaghetti" => "eat" }
			// h.fetch("spaghetti")                     #=> "eat"
			// h.fetch("pizza")                         #=> ArgumentError
			// h.fetch("pizza", "not eat")              #=> "not eat"
			// h.fetch("pizza") do |el| "eat " + el end #=> "eat pizza"
			// ```
			//
			// @return [Object]
			Name: "fetch",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				aLen := len(args)
				if aLen < 1 || aLen > 2 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentRange, 1, 2, aLen)
				}

				key, ok := args[0].(*StringObject)
				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, key.Class().Name)
				}

				if aLen == 2 {
					if blockFrame != nil {
						return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "The default argument can't be passed along with a block")
					}
					return args[1]
				}

				hash := receiver.(*HashObject)
				value, ok := hash.Pairs[key.value]

				if ok {
					if blockFrame != nil {
						t.callFrameStack.pop()
					}
					return value
				}

				if blockFrame != nil {
					return t.builtinMethodYield(blockFrame, key).Target
				}
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "The value was not found, and no block has been provided")

			},
		},
		{
			// Returns an array containing the values associated with the given keys but also raises
			// ArgumentError when one of keys can’t be found.
			//
			// ```Ruby
			// h = { cat: "feline", dog: "canine", cow: "bovine" }
			//
			// h.fetch_values("cow", "cat")                      #=> ["bovine", "feline"]
			// h.fetch_values("cow", "bird")                     # raises ArgumentError
			// h.fetch_values("cow", "bird") do |k| k.upcase end #=> ["bovine", "BIRD"]
			// ```
			//
			// @return [ArrayObject]
			Name: "fetch_values",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				aLen := len(args)
				if aLen < 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentMore, 1, aLen)
				}

				values := make([]Object, aLen)

				hash := receiver.(*HashObject)
				blockFramePopped := false

				for index, objectKey := range args {
					stringKey, ok := objectKey.(*StringObject)

					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, objectKey.Class().Name)
					}

					value, ok := hash.Pairs[stringKey.value]

					if !ok {
						if blockFrame != nil {
							value = t.builtinMethodYield(blockFrame, objectKey).Target
							blockFramePopped = true
						} else {
							return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "There is no value for the key `%s`, and no block has been provided", stringKey.value)
						}
					}

					values[index] = value
				}

				if blockFrame != nil && !blockFramePopped {
					t.callFrameStack.pop()
				}

				return t.vm.InitArrayObject(values)

			},
		},
		{
			// Returns true if the key exist in the hash. Currently, it can only input string
			// type object.
			//
			// ```Ruby
			// h = { a: 1, b: "2", c: [1, 2, 3], d: { k: "v" } }
			// h.has_key?("a") # => true
			// h.has_key?("e") # => false
			// # TODO: Support Symbol Type Key Input
			// h.has_key?(:b)  # => true
			// h.has_key?(:f)  # => false
			// ```
			//
			// @return [Boolean]
			Name: "has_key?",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
				}

				h := receiver.(*HashObject)
				i := args[0]
				input, ok := i.(*StringObject)

				if !ok {
					return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, i.Class().Name)
				}

				if _, ok := h.Pairs[input.value]; ok {
					return TRUE
				}
				return FALSE

			},
		},
		{
			// Returns true if the value exist in the hash.
			//
			// ```Ruby
			// h = { a: 1, b: "2", c: [1, 2, 3], d: { k: "v" } }
			// h.has_value?(1)          # => true
			// h.has_value?(2)          # => false
			// h.has_value?("2")        # => true
			// h.has_value?([1, 2, 3])  # => true
			// h.has_value?({ k: "v" }) # => true
			// ```
			//
			// @return [Boolean]
			Name: "has_value?",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
				}

				h := receiver.(*HashObject)

				for _, v := range h.Pairs {
					if reflect.DeepEqual(v, args[0]) {
						return TRUE
					}
				}
				return FALSE

			},
		},
		{
			// Returns an array of keys (in arbitrary order)
			//
			// ```Ruby
			// { a: 1, b: "2", c: [3, true, "Hello"] }.keys
			// # =>  ["c", "b", "a"] or ["b", "a", "c"] ... etc
			// ```
			//
			// @return [Boolean]
			Name: "keys",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
				}

				h := receiver.(*HashObject)
				var keys []Object
				for k := range h.Pairs {
					keys = append(keys, t.vm.InitStringObject(k))
				}
				return t.vm.InitArrayObject(keys)

			},
		},
		{
			// Returns the number of key-value pairs of the hash.
			//
			// ```Ruby
			// h = { a: 1, b: "2", c: [1, 2, 3], d: { k: 'v' } }
			// h.length  #=> 4
			// ```
			//
			// @return [Integer]
			Name: "length",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
				}

				h := receiver.(*HashObject)
				return t.vm.InitIntegerObject(h.length())

			},
		},
		{
			// Returns a new hash with the results of running the block once for every value.
			// This method does not change the keys and the receiver hash values.
			//
			// ```Ruby
			// h = { a: 1, b: 2, c: 3 }
			// result = h.map_values do |v|
			//   v * 3
			// end
			// h      # => { a: 1, b: 2, c: 3 }
			// result # => { a: 3, b: 6, c: 9 }
			// ```
			//
			// @return [Boolean]
			Name: "map_values",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
				}

				if blockFrame == nil {
					return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
				}

				h := receiver.(*HashObject)
				if blockIsEmpty(blockFrame) {
					return h
				}

				result := make(map[string]Object)

				if len(h.Pairs) == 0 {
					t.callFrameStack.pop()
				}

				for k, v := range h.Pairs {
					result[k] = t.builtinMethodYield(blockFrame, v).Target
				}
				return t.vm.InitHashObject(result)

			},
		},
		{
			// Returns the number of key-value pairs of the hash.
			//
			// ```Ruby
			// h = { a: 1, b: "2", c: [1, 2, 3] }
			// h.merge({ b: "Hello", d: "World" })
			// # => { a: 1, b: "Hello", c: [1, 2, 3], d: "World" }
			// ```
			//
			// @return [Hash]
			Name: "merge",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) < 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentMore, 1, len(args))
				}

				h := receiver.(*HashObject)
				result := make(map[string]Object)
				for k, v := range h.Pairs {
					result[k] = v
				}

				for _, obj := range args {
					hashObj, ok := obj.(*HashObject)
					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.HashClass, obj.Class().Name)
					}
					for k, v := range hashObj.Pairs {
						result[k] = v
					}
				}

				return t.vm.InitHashObject(result)

			},
		},
		{
			// Returns a new hash consisting of entries for which the block does not return false
			// or nil.
			//
			// ```ruby
			// a = { a: 1, b: 2 }
			//
			// a.select do |k, v|
			//   v == 2
			// end            # => { a: 1 }
			// a.select do |k, v|
			//   5
			// end            # => { a: 1, b: 2 }
			// a.select do |k, v|
			//   nil
			// end            # => { }
			// a.select do |k, v|
			//   false
			// end            # => { }
			// ```
			Name: "select",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
				}

				if blockFrame == nil {
					return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
				}

				destinationPairs := map[string]Object{}
				if blockIsEmpty(blockFrame) {
					return t.vm.InitHashObject(destinationPairs)
				}

				sourceHash := receiver.(*HashObject)

				if len(sourceHash.Pairs) == 0 {
					t.callFrameStack.pop()
				}

				for stringKey, value := range sourceHash.Pairs {
					objectKey := t.vm.InitStringObject(stringKey)
					result := t.builtinMethodYield(blockFrame, objectKey, value)

					if result.Target.isTruthy() {
						destinationPairs[stringKey] = value
					}
				}

				return t.vm.InitHashObject(destinationPairs)

			},
		},
		{
			// Returns an array of keys (in arbitrary order)
			//
			// ```Ruby
			// { a: 1, b: "2", c: [3, true, "Hello"] }.sorted_keys
			// # =>  ["a", "b", "c"]
			// { c: 1, b: "2", a: [3, true, "Hello"] }.sorted_keys
			// # =>  ["a", "b", "c"]
			// { b: 1, c: "2", a: [3, true, "Hello"] }.sorted_keys
			// # =>  ["a", "b", "c"]
			// { b: 1, c: "2", b: [3, true, "Hello"] }.sorted_keys
			// # =>  ["b", "c"]
			// ```
			//
			// @return [Boolean]
			Name: "sorted_keys",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
				}

				h := receiver.(*HashObject)
				sortedKeys := h.sortedKeys()
				var keys []Object
				for _, k := range sortedKeys {
					keys = append(keys, t.vm.InitStringObject(k))
				}
				return t.vm.InitArrayObject(keys)

			},
		},
		{
			// Returns two-dimensional array with the key-value pairs of hash. If specified true
			// then it will return sorted key value pairs array
			//
			// ```Ruby
			// { a: 1, b: 2, c: 3 }.to_a
			// # => [["a", 1], ["c", 3], ["b", 2]] or [["b", 2], ["c", 3], ["a", 1]] ... etc
			// { a: 1, b: 2, c: 3 }.to_a(true)
			// # => [["a", 1], ["b", 2], ["c", 3]]
			// { b: 1, a: 2, c: 3 }.to_a(true)
			// # => [["a", 2], ["b", 1], ["c", 3]]
			// { b: 1, a: 2, a: 3 }.to_a(true)
			// # => [["a", 3], ["b", 1]]
			// ```
			//
			// @return [Array]
			Name: "to_a",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				aLen := len(args)
				if aLen > 1 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentLess, 1, aLen)
				}

				var sorted bool
				if aLen == 0 {
					sorted = false
				} else {
					s := args[0]
					st, ok := s.(*BooleanObject)
					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.BooleanClass, s.Class().Name)
					}
					sorted = st.value
				}

				h := receiver.(*HashObject)
				var resultArr []Object
				if sorted {
					for _, k := range h.sortedKeys() {
						var pairArr []Object
						pairArr = append(pairArr, t.vm.InitStringObject(k))
						pairArr = append(pairArr, h.Pairs[k])
						resultArr = append(resultArr, t.vm.InitArrayObject(pairArr))
					}
				} else {
					for k, v := range h.Pairs {
						var pairArr []Object
						pairArr = append(pairArr, t.vm.InitStringObject(k))
						pairArr = append(pairArr, v)
						resultArr = append(resultArr, t.vm.InitArrayObject(pairArr))
					}
				}
				return t.vm.InitArrayObject(resultArr)

			},
		},
		{
			// Returns json that is corresponding to the hash.
			// Basically just like Hash#to_json in Rails but currently doesn't support options.
			//
			// ```Ruby
			// h = { a: 1, b: [1, "2", [4, 5, nil], { foo: "bar" }]}.to_json
			// puts(h) #=> {"a":1,"b":[1, "2", [4, 5, null], {"foo":"bar"}]}
			// ```
			//
			// @return [String]
			Name: "to_json",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
				}

				r := receiver.(*HashObject)
				return t.vm.InitStringObject(r.ToJSON(t))

			},
		},
		{
			// Returns json that is corresponding to the hash.
			// Basically just like Hash#to_json in Rails but currently doesn't support options.
			//
			// ```Ruby
			// h = { a: 1, b: [1, "2", [4, 5, nil], { foo: "bar" }]}.to_s
			// puts(h) #=> "{ a: 1, b: [1, \"2\", [4, 5, null], { foo: \"bar \" }] }"
			// ```
			//
			// @return [String]
			Name: "to_s",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
				}

				h := receiver.(*HashObject)
				return t.vm.InitStringObject(h.ToString())

			},
		},
		{
			// Returns a new hash with the results of running the block once for every value.
			// This method does not change the keys and unlike Hash#map_values, it does not
			// change the receiver hash values.
			//
			// ```Ruby
			// h = { a: 1, b: 2, c: 3 }
			// result = h.transform_values do |v|
			//   v * 3
			// end
			// h      # => { a: 1, b: 2, c: 3 }
			// result # => { a: 3, b: 6, c: 9 }
			// ```
			//
			// @return [Boolean]
			Name: "transform_values",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
				}

				if blockFrame == nil {
					return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
				}

				h := receiver.(*HashObject)

				if len(h.Pairs) == 0 {
					t.callFrameStack.pop()
				}

				resultHash := make(map[string]Object)
				for k, v := range h.Pairs {
					result := t.builtinMethodYield(blockFrame, v)
					resultHash[k] = result.Target
				}
				return t.vm.InitHashObject(resultHash)

			},
		},
		{
			// Returns an array of values (in arbitrary order)
			//
			// ```Ruby
			// { a: 1, b: "2", c: [3, true, "Hello"] }.keys
			// # =>  [1, "2", [3, true, "Hello"]] or ["2", [3, true, "Hello"], 1] ... etc
			// ```
			//
			// @return [Boolean]
			Name: "values",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				if len(args) != 0 {
					return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
				}

				h := receiver.(*HashObject)
				var keys []Object
				for _, v := range h.Pairs {
					keys = append(keys, v)
				}
				return t.vm.InitArrayObject(keys)

			},
		},
		{
			// Return an array containing the values associated with the given keys.
			//
			// ```Ruby
			// { a: 1, b: "2" }.values_at("a", "c") # => [1, nil]
			// ```
			//
			// @return [Boolean]
			Name: "values_at",
			Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
				hash := receiver.(*HashObject)
				var result []Object

				for _, objectKey := range args {
					stringObjectKey, ok := objectKey.(*StringObject)

					if !ok {
						return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, objectKey.Class().Name)
					}

					value, ok := hash.Pairs[stringObjectKey.value]

					if !ok {
						value = NULL
					}

					result = append(result, value)
				}

				return t.vm.InitArrayObject(result)

			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) InitHashObject(pairs map[string]Object) *HashObject {
	return &HashObject{
		BaseObj: &BaseObj{class: vm.TopLevelClass(classes.HashClass)},
		Pairs:   pairs,
	}
}

func (vm *VM) initHashClass() *RClass {
	hc := vm.initializeClass(classes.HashClass)
	hc.setBuiltinMethods(builtinHashInstanceMethods(), false)
	hc.setBuiltinMethods(builtinHashClassMethods(), true)
	return hc
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (h *HashObject) Value() interface{} {
	return h.Pairs
}

// ToString returns the object's name as the string format
func (h *HashObject) ToString() string {
	var out bytes.Buffer
	var pairs []string

	for _, key := range h.sortedKeys() {
		// TODO: Improve this conditional statement
		if _, isString := h.Pairs[key].(*StringObject); isString {
			pairs = append(pairs, fmt.Sprintf("%s: \"%s\"", key, h.Pairs[key].ToString()))
		} else {
			pairs = append(pairs, fmt.Sprintf("%s: %s", key, h.Pairs[key].ToString()))
		}
	}

	out.WriteString("{ ")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString(" }")

	return out.String()
}

// ToJSON returns the object's name as the JSON string format
func (h *HashObject) ToJSON(t *Thread) string {
	var out bytes.Buffer
	var values []string
	pairs := h.Pairs
	out.WriteString("{")

	for key, value := range pairs {
		values = append(values, generateJSONFromPair(key, value, t))
	}

	out.WriteString(strings.Join(values, ","))
	out.WriteString("}")
	return out.String()
}

// Returns the length of the hash
func (h *HashObject) length() int {
	return len(h.Pairs)
}

// Returns the sorted keys of the hash
func (h *HashObject) sortedKeys() []string {
	var arr []string
	for k := range h.Pairs {
		arr = append(arr, k)
	}
	sort.Strings(arr)
	return arr
}

// Returns the duplicate of the Hash object
func (h *HashObject) copy() Object {
	elems := map[string]Object{}

	for k, v := range h.Pairs {
		elems[k] = v
	}

	newHash := &HashObject{
		BaseObj: &BaseObj{class: h.class},
		Pairs:   elems,
	}

	return newHash
}

// recursive indexed access - see ArrayObject#dig documentation.
func (h *HashObject) dig(t *Thread, keys []Object, sourceLine int) Object {
	currentKey := keys[0]
	stringCurrentKey, ok := currentKey.(*StringObject)

	if !ok {
		return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, currentKey.Class().Name)
	}

	nextKeys := keys[1:]
	currentValue, ok := h.Pairs[stringCurrentKey.value]

	if !ok {
		return NULL
	}

	if len(nextKeys) == 0 {
		return currentValue
	}

	diggableCurrentValue, ok := currentValue.(Diggable)

	if !ok {
		return t.vm.InitErrorObject(errors.TypeError, sourceLine, "Expect target to be Diggable, got %s", currentValue.Class().Name)
	}

	return diggableCurrentValue.dig(t, nextKeys, sourceLine)
}

// Other helper functions ----------------------------------------------

// Return the JSON style strings of the Hash object
func generateJSONFromPair(key string, v Object, t *Thread) string {
	var data string
	var out bytes.Buffer

	out.WriteString(data)
	out.WriteString("\"" + key + "\"")
	out.WriteString(":")
	out.WriteString(v.ToJSON(t))

	return out.String()
}
