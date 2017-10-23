package vm

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// Implementation of thread-safe associative arrays (Hash).
//
// The implementation internally uses Go's `sync.Map` type, with some advantages and disadvantages:
//
// - it is highly performant and predictable for a certain pattern of usage (`concurrent loops with keys that are stable over time, and either few steady-state stores, or stores localized to one goroutine per key.`); performance and predictability in other conditions are unspecified;
// - iterations are non-deterministic; during iterations, keys may not be included;
// - size can't be retrieved;
// - for the reasons above, the Hash APIs implemented are minimal.
//
// For details, see https://golang.org/pkg/sync/#Map.
//
// Concurrent hashes are instantiated via `new()`:
//
//     ConcurrentHash.new()
//     ConcurrentHash.new({"a": 1, "b": 2})
//
// ```ruby
// hash = ConcurrentHash.new({ "a": 1, "b": 2 })
// has["a"]  # => 1
// ```
//
type ConcurrentHashObject struct {
	*baseObj
	internalMap sync.Map
}

// Class methods --------------------------------------------------------
func builtinConcurrentHashClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) > 1 {
						return t.vm.initErrorObject(errors.ArgumentError, instruction, "Expect 0 or 1 arguments, got %d", len(args))
					}

					if len(args) == 0 {
						return t.vm.initConcurrentHashObject(make(map[string]Object))
					}

					arg := args[0]
					hashArg, ok := arg.(*HashObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, instruction, errors.WrongArgumentTypeFormat, classes.HashClass, arg.Class().Name)
					}

					return t.vm.initConcurrentHashObject(hashArg.Pairs)
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinConcurrentHashInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// Retrieves the value (object) that corresponds to the key specified.
			// When a key doesn't exist, `nil` is returned, or the default, if set.
			//
			// ```Ruby
			// h = ConcurrentHash.new({ a: 1, b: "2" })
			// h['a'] #=> 1
			// h['b'] #=> "2"
			// h['c'] #=> nil
			// ```
			//
			// @return [Object]
			Name: "[]",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, instruction, "Expect 1 argument. got: %d", len(args))
					}

					i := args[0]
					key, ok := i.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, instruction, errors.WrongArgumentTypeFormat, classes.StringClass, i.Class().Name)
					}

					h := receiver.(*ConcurrentHashObject)

					value, ok := h.internalMap.Load(key.value)

					if !ok {
						return NULL
					}

					return value.(Object)
				}
			},
		},
		{
			// Associates the value given by `value` with the key given by `key`.
			// Returns the `value`.
			//
			// ```Ruby
			// h = ConcurrentHash.new{ a: 1, b: "2" })
			// h['a'] = 2          #=> 2
			// h                   #=> { a: 2, b: "2" }
			// ```
			//
			// @return [Object] The value
			Name: "[]=",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					// First arg is index
					// Second arg is assigned value
					if len(args) != 2 {
						return t.vm.initErrorObject(errors.ArgumentError, instruction, "Expect 2 arguments. got: %d", len(args))
					}

					k := args[0]
					key, ok := k.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, instruction, errors.WrongArgumentTypeFormat, classes.StringClass, k.Class().Name)
					}

					h := receiver.(*ConcurrentHashObject)
					h.internalMap.Store(key.value, args[1])

					return args[1]
				}
			},
		},
		{
			// Remove the key from the hash if key exist.
			//
			// ```Ruby
			// h = ConcurrentHash.new({ a: 1, b: 2, c: 3 })
			// h.delete("b") # => NULL
			// h             # => { a: 1, c: 3 }
			// ```
			//
			// @return [NULL]
			Name: "delete",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, instruction, "Expect 1 argument. got: %d", len(args))
					}

					h := receiver.(*ConcurrentHashObject)
					d := args[0]
					deleteKeyObject, ok := d.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, instruction, errors.WrongArgumentTypeFormat, classes.StringClass, d.Class().Name)
					}

					h.internalMap.Delete(deleteKeyObject.value)

					return NULL
				}
			},
		},
		{
			// Calls block once for each key in the hash (in sorted key order), passing the
			// key-value pair as parameters.
			// Note that iteration is not deterministic under all circumstances; see
			// https://golang.org/pkg/sync/#Map.
			//
			// ```Ruby
			// h = ConcurrentHash.new({ b: "2", a: 1 })
			// h.each do |k, v|
			//   puts k.to_s + "->" + v.to_s
			// end
			// # => a->1
			// # => b->2
			// ```
			//
			// @return [Hash] self
			Name: "each",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, instruction, errors.CantYieldWithoutBlockFormat)
					}

					if len(args) != 0 {
						t.callFrameStack.pop()
						return t.vm.initErrorObject(errors.ArgumentError, instruction, "Expect 0 arguments. got: %d", len(args))
					}

					hash := receiver.(*ConcurrentHashObject)
					framePopped := false

					iterator := func(key, value interface{}) bool {
						keyObject := t.vm.initStringObject(key.(string))

						t.builtinMethodYield(blockFrame, keyObject, value.(Object))

						framePopped = true

						return true
					}

					hash.internalMap.Range(iterator)

					if !framePopped {
						t.callFrameStack.pop()
					}

					return hash
				}
			},
		},
		{
			// Returns true if the key exist in the hash.
			//
			// ```Ruby
			// h = ConcurrentHash.new({ a: 1, b: "2" })
			// h.has_key?("a") # => true
			// h.has_key?("e") # => false
			// ```
			//
			// @return [Boolean]
			Name: "has_key?",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, instruction, "Expect 1 argument. got: %d", len(args))
					}

					h := receiver.(*ConcurrentHashObject)
					i := args[0]
					input, ok := i.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, instruction, errors.WrongArgumentTypeFormat, classes.StringClass, i.Class().Name)
					}

					if _, ok := h.internalMap.Load(input.value); ok {
						return TRUE
					}

					return FALSE
				}
			},
		},
		{
			// Returns json that is corresponding to the hash.
			// Basically just like Hash#to_json in Rails but currently doesn't support options.
			//
			// ```Ruby
			// h = ConcurrentHash.new({ a: 1, b: 2 })
			// h.to_json #=> {"a":1,"b":2}
			// ```
			//
			// @return [String]
			Name: "to_json",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, instruction, "Expect 0 argument. got: %d", len(args))
					}

					r := receiver.(*ConcurrentHashObject)
					return t.vm.initStringObject(r.toJSON())
				}
			},
		},
		{
			// Returns json that is corresponding to the hash.
			// Basically just like Hash#to_json in Rails but currently doesn't support options.
			//
			// ```Ruby
			// h = { a: 1, b: "2"}
			// h.to_s #=> "{ a: 1, b: \"2\" }"
			// ```
			//
			// @return [String]
			Name: "to_s",
			Fn: func(receiver Object, instruction *instruction) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, instruction, "Expect 0 argument. got: %d", len(args))
					}

					h := receiver.(*ConcurrentHashObject)
					return t.vm.initStringObject(h.toString())
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initConcurrentHashObject(pairs map[string]Object) *ConcurrentHashObject {
	var internalMap sync.Map

	for key, value := range pairs {
		internalMap.Store(key, value)
	}

	return &ConcurrentHashObject{
		baseObj: &baseObj{class: vm.topLevelClass(classes.ConcurrentHashClass)},
		internalMap: internalMap,
	}
}

func initConcurrentHashClass(vm *VM) {
	chc := vm.initializeClass(classes.ConcurrentHashClass, false)
	chc.setBuiltinMethods(builtinConcurrentHashInstanceMethods(), false)
	chc.setBuiltinMethods(builtinConcurrentHashClassMethods(), true)
	vm.objectClass.setClassConstant(chc)
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (h *ConcurrentHashObject) Value() interface{} {
	return h.internalMap
}

// toString returns the object's name as the string format
func (h *ConcurrentHashObject) toString() string {
	var out bytes.Buffer
	var pairs []string

	iterator := func(key, value interface{}) bool {
		var template string

		switch value.(type) {
		case *StringObject:
			template = "%s: \"%s\""
		default:
			template = "%s: %s"
		}

		pairs = append(pairs, fmt.Sprintf(template, key, value.(Object).toString()))

		return true
	}

	h.internalMap.Range(iterator)

	out.WriteString("{ ")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString(" }")

	return out.String()
}

// toJSON returns the object's name as the JSON string format
func (h *ConcurrentHashObject) toJSON() string {
	var out bytes.Buffer
	var values []string
	out.WriteString("{")

	iterator := func(key, value interface{}) bool {
		values = append(values, generateJSONFromPair(key.(string), value.(Object)))

		return true
	}

	h.internalMap.Range(iterator)

	out.WriteString(strings.Join(values, ","))
	out.WriteString("}")
	return out.String()
}
