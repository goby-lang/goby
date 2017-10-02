package vm

// Diggable represents a class that support the #dig method.
type Diggable interface {
  dig(t *thread, keys []Object) Object
}
