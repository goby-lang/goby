package result_test

import (
	"testing"

	"github.com/goby-lang/goby/vm"
)

func TestNewResult(t *testing.T) {
	r := vm.ExecAndReturn(t, `
		require("result")
		r = Result.new()

		`)
	t.Error(r)

}
