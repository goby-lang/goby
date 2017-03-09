package object

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
