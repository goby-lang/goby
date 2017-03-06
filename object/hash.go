package object

import (
	"bytes"
	"fmt"
	"strings"
)

type HashClass struct {
	*BaseClass
}

type HashObject struct {
	Class *HashClass
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
