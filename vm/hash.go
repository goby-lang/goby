package vm

import (
	"bytes"
	"fmt"
	"strings"
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
// - **Value:** String literal and objects (Integer, String, Array, Hash, nil, etc) can be used.
//
// **Note:**
// - The order of key-value pairs are **not** preserved.
// - Operator `=>` is not supported.
// - `Hash.new` is not supported.
type HashObject struct {
	*baseObj
	Pairs map[string]Object
}

func (vm *VM) initHashObject(pairs map[string]Object) *HashObject {
	return &HashObject{Pairs: pairs, baseObj: &baseObj{class: vm.topLevelClass(hashClass)}}
}

func (vm *VM) initHashClass() *RClass {
	hc := vm.initializeClass(hashClass, false)
	hc.setBuiltInMethods(builtinHashInstanceMethods(), false)
	hc.setBuiltInMethods(builtInHashClassMethods(), true)
	return hc
}

func builtInHashClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					return t.UnsupportedMethodError("#new", receiver)
				}
			},
		},
	}
}

func builtinHashInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			// Retrieves the value (object) that corresponds to the key specified.
			// Returns `nil` when specifying a nonexistent key.
			//
			// ```Ruby
			// h = { a: 1, b: "2", c: [1, 2, 3], d: { k: 'v' } }
			// h['a'] #=> 1
			// h['b'] #=> "2"
			// h['c'] #=> [1, 2, 3]
			// h['d'] #=> { k: 'v' }
			// ```
			//
			// @return [Object]
			Name: "[]",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					if len(args) != 1 {
						return newError("Expect 1 arguments. got=%d", len(args))
					}

					i := args[0]
					key, ok := i.(*StringObject)

					if !ok {
						return newError("Expect index argument to be String. got=%T", i)
					}

					hash := receiver.(*HashObject)

					if len(hash.Pairs) == 0 {
						return NULL
					}

					value, ok := hash.Pairs[key.Value]

					if !ok {
						return NULL
					}

					return value

				}
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					// First arg is index
					// Second arg is assigned value
					if len(args) != 2 {
						return newError("Expect 2 arguments. got=%d", len(args))
					}

					k := args[0]
					key, ok := k.(*StringObject)

					if !ok {
						return newError("Expect index argument to be String. got=%T", k)
					}

					hash := receiver.(*HashObject)
					hash.Pairs[key.Value] = args[1]

					return args[1]
				}
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					if len(args) != 0 {
						return newError("Expect 0 argument. got=%d", len(args))
					}

					hash := receiver.(*HashObject)
					return t.vm.initIntegerObject(hash.length())
				}
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
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*HashObject)
					return t.vm.initStringObject(r.toJSON())
				}
			},
		},
	}
}

// Polymorphic helper functions -----------------------------------------

// toString converts the receiver into string.
func (h *HashObject) toString() string {
	var out bytes.Buffer
	var pairs []string

	for key, value := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", key, value.toString()))
	}

	out.WriteString("{ ")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString(" }")

	return out.String()
}

// toJSON converts the receiver into JSON string.
func (h *HashObject) toJSON() string {
	var out bytes.Buffer
	var values []string
	pairs := h.Pairs
	out.WriteString("{")

	for key, value := range pairs {
		values = append(values, generateJSONFromPair(key, value))
	}

	out.WriteString(strings.Join(values, ","))
	out.WriteString("}")
	return out.String()
}

func (h *HashObject) length() int {
	return len(h.Pairs)
}

// Other helper functions ----------------------------------------------

func generateJSONFromPair(key string, v Object) string {
	var data string
	var out bytes.Buffer

	out.WriteString(data)
	out.WriteString("\"" + key + "\"")
	out.WriteString(":")
	out.WriteString(v.toJSON())

	return out.String()
}
