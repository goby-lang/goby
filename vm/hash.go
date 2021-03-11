package vm

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// HashObject represents hash instances.
// Hash is a collection of key-value pair, which works like a dictionary.
// Hash literal is represented with curly brackets `{ }` like `{ key: value }`.
// Each key of the hash is unique and cannot be duplicate within the hash.
// Adding a leading space and a trailing space within curly brackets are preferable.
//
// - **Key:** an alphanumeric word that starts with alphabet, without containing space and punctuations.
// Underscore `_` can also be used within the key.
// In hash literals, only a symbol literals such as `symbol:` can be used as a key.
// String literal like "mickey mouse" cannot be used as a key in hash literals.
// (String and symbol are equivalent in Goby)
//
// Retrieving a value via `[]`, you can use both symbol literals or string literals as keys.
//
// ```ruby
// a = { balthazar1: 100 } # valid
// b = { 2melchior: 200 }  # invalid
// b = { "casper": 200 }   # invalid
// x = 'balthazar1'
//
// a["balthazar1"]  # => 100
// a[balthazar1:]   # => 100
// a[x]             # => 100
// a[balthazar1]    # => error
// ```
//
// - **value:** String literals and objects (Integer, String, Array, Hash, nil, etc) can be used.
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
var builtinHashClassMethods = []*BuiltinMethodObject{
	{
		Name: "new",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			return t.vm.InitNoMethodError(sourceLine, "new", receiver)

		},
	},
}

// Instance methods -----------------------------------------------------
var builtinHashInstanceMethods = []*BuiltinMethodObject{
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
		// @param key [String]
		// @return [Object]
		Name: "[]",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
			}

			typeErr := t.vm.checkArgTypes(args, sourceLine, classes.StringClass)

			if typeErr != nil {
				return typeErr
			}

			h := receiver.(*HashObject)

			value, ok := h.Pairs[args[0].Value().(string)]

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
		// @param key [String]
		// @return [Object] The value
		Name: "[]=",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			// First arg is index
			// Second arg is assigned value
			if len(args) != 2 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 2, len(args))
			}

			typeErr := t.vm.checkArgTypes(args, sourceLine, classes.StringClass)

			if typeErr != nil {
				return typeErr
			}

			h := receiver.(*HashObject)
			h.Pairs[args[0].Value().(string)] = args[1]

			return args[1]

		},
	},
	{
		// Passes each (key, value) pair of the collection to the given block.
		// The method returns true if any of the results by the block is true.
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
		//
		// @return [Boolean]
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

				if result.isTruthy() {
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
		// @return [Hash]
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
		// Returns the configured default value of the Hash.
		// If no default value has been specified, nil is returned.
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
		// Sets the default value of this Hash for the missing keys, and returns the default value.
		// Note that Arrays/Hashes are not accepted because they're unsafe.
		//
		// ```Ruby
		// h = { a: 1 }
		// h['c']         #=> nil
		// h.default = 2
		// h['c']         #=> 2
		// h.default = [] #=> ArgumentError
		// ```
		//
		// @param default value [Object]
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
		// @param key [String]
		// @return [Hash]
		Name: "delete",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
			}

			typeErr := t.vm.checkArgTypes(args, sourceLine, classes.StringClass)

			if typeErr != nil {
				return typeErr
			}

			deleteKeyValue := args[0].Value().(string)

			h := receiver.(*HashObject)

			if _, ok := h.Pairs[deleteKeyValue]; ok {
				delete(h.Pairs, deleteKeyValue)
			}
			return h

		},
	},
	{
		// Deletes every key-value pair from the hash for which block evaluates to anything except false and nil.
		//
		// Returns the modified hash.
		//
		// ```Ruby
		// { a: 1, b: 2}.delete_if do |k, v| v == 1 end # =>  { b: 2 }
		// { a: 1, b: 2}.delete_if do |k, v| 5 end      # =>  { }
		// { a: 1, b: 2}.delete_if do |k, v| false end  # =>  { a: 1, b: 2}
		// { a: 1, b: 2}.delete_if do |k, v| nil end    # =>  { a: 1, b: 2}
		// ```
		//
		// @param block
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

				booleanResult, isResultBoolean := result.(*BooleanObject)

				if isResultBoolean {
					if booleanResult.value {
						delete(hash.Pairs, stringKey)
					}
				} else if result != NULL {
					delete(hash.Pairs, stringKey)
				}
			}

			return hash

		},
	},
	{
		// Extracts the nested value specified by the sequence of idx objects by calling `dig` at each step,
		// Returns nil if any intermediate step is nil.
		//
		// ```Ruby
		// { a: 1 , b: 2 }.dig(:a)         # => 1
		// { a: {}, b: 2 }.dig(:a, :b)     # => nil
		// { a: {}, b: 2 }.dig(:a, :b, :c) # => nil
		// { a: 1, b: 2 }.dig(:a, :b)      # => TypeError: Expect target to be Diggable
		// ```
		//
		// @param key [String]
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
		// Performs a copy of the hash, including the keys and values, and returns it.
		// Any arguments are ignored.
		// The object_id of the returned object is different from the one of the receiver.

		// Caveat: any keys of hash ARE also copied with different object ids for now.
		// This comes from the fact that the string objects are NOT frozen in current Goby.
		//
		// See also `Object#dup`, `String#dup`, `Array#dup`.
		//
		// ```ruby
		// h = { k1: :key1, k2: :key2 }
		// h.object_id           #» 824633779744
		// h.each do |k, v|
		//   print "key:   "
		//   puts k.object_id
		//   print "value: "
		//   puts v.object_id
		// end
		// key:   824636231680
		// value: 824635528224
		// key:   824636232480
		// value: 824635528448
		//
		// b = h.dup
		// b.object_id           #» 824633779904
		// b.each do |k, v|
		//   print "key:   "
		//   puts k.object_id
		//   print "value: "
		//   puts v.object_id
		// end
		// key:   824638121536
		// value: 824635528224
		// key:   824638122336
		// value: 824635528448
		// ```
		//
		// @return [Hash]
		Name: "dup",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			return receiver.(*HashObject).copy()
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
		// @param block
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
		// Loops through keys of the hash with given block frame.
		// Then returns an array of keys in alphabetical order.
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
		// @param block
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
		// Loops through values of the hash with given block frame.
		// Then returns an array of values of the hash in the alphabetical order of the keys.
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
		// @param block
		// @return [Array]
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
		// @param object [Object]
		// @return [Boolean]
		Name: "eql?",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
			}

			h := receiver.(*HashObject)
			c := args[0]
			compare, ok := c.(*HashObject)

			if ok && h.equalTo(compare) {
				return TRUE
			}
			return FALSE

		},
	},
	{
		// Returns a value from the hash for the given key.
		// If the key can’t be found, there are several options:
		//
		// - With no other arguments, it will raise an ArgumentError.
		// - If a default value is given as a second argument, then that will be returned.
		// - If an optional code block is specified, then runs the block and returns the result.
		// - If a block and a second argument is given together, it raises an ArgumentError.
		//
		// ```Ruby
		// h = { spaghetti: "eat" }
		// h.fetch("spaghetti")                     #=> "eat"
		// h.fetch("pizza")                         #=> ArgumentError
		// h.fetch("pizza", "not eat")              #=> "not eat"
		// h.fetch("pizza") do |el| "eat " + el end #=> "eat pizza"
		// ```
		//
		// @param key [String], default value [Object]
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
				return t.builtinMethodYield(blockFrame, key)
			}
			return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, "The value was not found, and no block has been provided")
		},
	},
	{
		// Returns an array containing the values associated with the given keys.
		// When even one of keys can’t be found, it raises an ArgumentError.
		//
		// ```Ruby
		// h = { cat: "feline", dog: "canine", cow: "bovine" }
		//
		// h.fetch_values("cow", "cat")                      #=> ["bovine", "feline"]
		// h.fetch_values("cow", "bird")                     # raises ArgumentError
		// h.fetch_values("cow", "bird") do |k| k.upcase end #=> ["bovine", "BIRD"]
		// ```
		//
		// @param key [String]...
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
						value = t.builtinMethodYield(blockFrame, objectKey)
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
		// Returns true if the specified key exists in the hash
		// Currently, only string can be taken.
		// type object.
		//
		// ```Ruby
		// h = { a: 1, b: "2", c: [1, 2, 3], d: { k: "v" } }
		// h.has_key?("a") # => true
		// h.has_key?("e") # => false
		// h.has_key?(:b)  # => true
		// h.has_key?(:f)  # => false
		// ```
		//
		// @param key [String]
		// @return [Boolean]
		Name: "has_key?",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
			}

			typeErr := t.vm.checkArgTypes(args, sourceLine, classes.StringClass)

			if typeErr != nil {
				return typeErr
			}

			if _, ok := receiver.(*HashObject).Pairs[args[0].Value().(string)]; ok {
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
		// @param value [Object]
		// @return [Boolean]
		Name: "has_value?",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 1 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 1, len(args))
			}

			h := receiver.(*HashObject)

			for _, v := range h.Pairs {
				if v.equalTo(args[0]) {
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
		// @param block
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
				result[k] = t.builtinMethodYield(blockFrame, v)
			}
			return t.vm.InitHashObject(result)

		},
	},
	{
		// Returns a newly merged hash. One or more hashes can be taken.
		// If keys are duplicate between the receiver and the argument, the last ones in the argument are prioritized.
		//
		// ```Ruby
		// h = { a: 1, b: "2", c: [1, 2, 3] }
		// h.merge({ b: "Hello", d: "World" })
		// # => { a: 1, b: "Hello", c: [1, 2, 3], d: "World" }
		//
		// { a: "Hello"}.merge({a: 0}, {a: 99})
		// # => { a: 99 }
		// ```
		//
		// @param hash [Hash]
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
		//
		// @param block
		// @return [Hash]
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

				if result.isTruthy() {
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
		// @return [Array]
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
		// @param sorting [Boolean]
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
				typeErr := t.vm.checkArgTypes(args, sourceLine, classes.BooleanClass)

				if typeErr != nil {
					return typeErr
				}

				sorted = args[0].Value().(bool)
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
		// @return [Enumerator]
		Name: "to_enum",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			return t.vm.InitEnumeratorObject(receiver, "", nil)
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
		// This method does not change the keys. Unlike Hash#map_values, it does not
		// change the receiver's hash values.
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
		// @param block
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
				resultHash[k] = t.builtinMethodYield(blockFrame, v)
			}
			return t.vm.InitHashObject(resultHash)

		},
	},
	{
		// Returns an array of values.
		// The order of the returned values are indeterminable.
		//
		// ```Ruby
		// { a: 1, b: "2", c: [3, true, "Hello"] }.keys
		// # =>  [1, "2", [3, true, "Hello"]] or ["2", [3, true, "Hello"], 1] ... etc
		// ```
		//
		// @return [Array]
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
		// @param key [String]...
		// @return [Array]
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

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

// InitHashObject creates a HashObject
func (vm *VM) InitHashObject(pairs map[string]Object) *HashObject {
	return &HashObject{
		BaseObj: NewBaseObject(vm.TopLevelClass(classes.HashClass)),
		Pairs:   pairs,
	}
}

func (vm *VM) initHashClass() *RClass {
	hc := vm.initializeClass(classes.HashClass)
	hc.setBuiltinMethods(builtinHashInstanceMethods, false)
	hc.setBuiltinMethods(builtinHashClassMethods, true)
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
		pairs = append(pairs, fmt.Sprintf("%s: %s", key, h.Pairs[key].Inspect()))
	}

	out.WriteString("{ ")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString(" }")

	return out.String()
}

// Inspect delegates to ToString
func (h *HashObject) Inspect() string {
	return h.ToString()
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
		BaseObj: NewBaseObject(h.class),
		Pairs:   elems,
	}

	return newHash
}

// recursive indexed access - see ArrayObject#dig documentation.
func (h *HashObject) dig(t *Thread, keys []Object, sourceLine int) Object {
	typeErr := t.vm.checkArgTypes(keys, sourceLine, classes.StringClass)

	if typeErr != nil {
		return typeErr
	}

	nextKeys := keys[1:]
	currentValue, ok := h.Pairs[keys[0].Value().(string)]

	if !ok {
		return NULL
	}

	if len(nextKeys) == 0 {
		return currentValue
	}

	diggableCurrentValue, ok := currentValue.(Diggable)

	if !ok {
		return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.NotDiggable, currentValue.Class().Name)
	}

	return diggableCurrentValue.dig(t, nextKeys, sourceLine)
}

func (h *HashObject) equalTo(with Object) bool {
	w, ok := with.(*HashObject)

	if !ok {
		return false
	}

	if len(h.Pairs) != len(w.Pairs) {
		return false
	}

	for k, v := range h.Pairs {
		if !v.equalTo(w.Pairs[k]) {
			return false
		}
	}

	return true
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

func (h *HashObject) Enumerable() bool {
	return true
}
