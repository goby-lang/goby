package vm

import (
	"bytes"
	"fmt"
	"strings"
)

var (
	hashClass *RHash
)

// RHash is the class of hash objects
type RHash struct {
	*BaseClass
}

// HashObject represents hash instances
type HashObject struct {
	Class *RHash
	Pairs map[string]Object
}

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

func (h *HashObject) returnClass() Class {
	return h.Class
}

func (h *HashObject) length() int {
	return len(h.Pairs)
}

func initializeHash(pairs map[string]Object) *HashObject {
	return &HashObject{Pairs: pairs, Class: hashClass}
}

func generateJSONFromPair(key string, v Object) string {
	var data string
	var out bytes.Buffer

	out.WriteString(data)
	out.WriteString("\"" + key + "\"")
	out.WriteString(":")
	out.WriteString(v.toJSON())

	return out.String()
}

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

var builtinHashInstanceMethods = []*BuiltInMethodObject{
	{
		// Returns json that is corresponding to the hash.
		// Basically just like Hash#to_json in Rails but currently doesn't support options.
		//
		// ```Ruby
		// h = { a: 1, b: [1, "2", [4, 5, nil], { foo: "bar" }]}.to_json
		// puts(h) #=> {"a":1,"b":[1, "2", [4, 5, null], {"foo":"bar"}]}
		// ```
		// @return [String]
		Name: "to_json",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				r := receiver.(*HashObject)
				return initializeString(r.toJSON())
			}
		},
	},
	{
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
		Name: "length",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				if len(args) != 0 {
					return newError("Expect 0 argument. got=%d", len(args))
				}

				hash := receiver.(*HashObject)
				return initilaizeInteger(hash.length())
			}
		},
	},
}

func initHash() {
	bc := &BaseClass{Name: "Hash", ClassMethods: newEnvironment(), Methods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass, superClass: objectClass}
	hc := &RHash{BaseClass: bc}
	hc.setBuiltInMethods(builtinHashInstanceMethods, false)
	hashClass = hc
}
