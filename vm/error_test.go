package vm

import (
  "testing"
)

func TestUndefinedMethod(t *testing.T) {
  expectError(t, "Undefined method", "a")
}
