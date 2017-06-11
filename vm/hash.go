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

func (h *HashObject) Inspect() string {
	var out bytes.Buffer
	var pairs []string

	for key, value := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", key, value.Inspect()))
	}

	out.WriteString("{ ")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString(" }")

	return out.String()
}

func (h *HashObject) returnClass() Class {
	return h.Class
}

func (h *HashObject) Length() int {
	return len(h.Pairs)
}

func initializeHash(pairs map[string]Object) *HashObject {
	return &HashObject{Pairs: pairs, Class: hashClass}
}

func generateJson(key string, v Object) string {
	var data string
	var value string
	var out bytes.Buffer

	out.WriteString(data)
	out.WriteString("\"" + key + "\"")
	out.WriteString(":")

	switch v := v.(type) {
	case *StringObject:
		value = "\"" + v.Value + "\""
	case *IntegerObject:
		value = v.Inspect()
	}
	out.WriteString(value)

	return out.String()
}

var builtinHashInstanceMethods = []*BuiltInMethodObject{
	{
		Name: "to_json",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				var out bytes.Buffer
				var values []string

				pairs := receiver.(*HashObject).Pairs
				out.WriteString("{")

				for key, value := range pairs {
					values = append(values, generateJson(key, value))
				}

				out.WriteString(strings.Join(values, ","))
				out.WriteString("}")
				return initializeString(out.String())
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
				return initilaizeInteger(hash.Length())
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
