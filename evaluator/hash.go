package evaluator

import (
	"bytes"
	"fmt"
	"strings"
)

var (
	HashClass *RHash
)

type RHash struct {
	*BaseClass
}

type HashObject struct {
	Class *RHash
	Pairs map[string]Object
}

func (h *HashObject) Type() ObjectType {
	return HASH_OBJ
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

func (h *HashObject) ReturnClass() Class {
	return h.Class
}

func (h *HashObject) Length() int {
	return len(h.Pairs)
}

func InitializeHash(pairs map[string]Object) *HashObject {
	return &HashObject{Pairs: pairs, Class: HashClass}
}

var builtinHashMethods = []*BuiltInMethod{
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
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
		Name: "[]",
	},
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
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
		Name: "[]=",
	},
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				if len(args) != 0 {
					return newError("Expect 0 argument. got=%d", len(args))
				}

				hash := receiver.(*HashObject)
				return InitilaizeInteger(hash.Length())
			}
		},
		Name: "length",
	},
}

func init() {
	methods := NewEnvironment()

	for _, m := range builtinHashMethods {
		methods.Set(m.Name, m)
	}

	bc := &BaseClass{Name: "Hash", Methods: methods, Class: ClassClass, SuperClass: ObjectClass}
	hc := &RHash{BaseClass: bc}
	HashClass = hc
}
