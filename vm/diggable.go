package vm

// Diggable represents a class that support the #dig method.
type Diggable interface {
	dig(t *Thread, keys []Object, sourceLine int) Object
}
