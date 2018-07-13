package result

import (
	"testing"

	"github.com/goby-lang/goby/vm"
)

func TestEmptyResult(t *testing.T) {
	r := vm.ExecAndReturn(t, `
		require("result")
		r = Result.empty()
	`)

	if r, ok := r.(*Result); ok {
		if !r.empty {
			t.Error("Should be empty")
		}

	} else {
		t.Fatalf("Wrong type %T", r)
	}
}

func TestNewResult(t *testing.T) {
	r := vm.ExecAndReturn(t, `
		require("result")
		r = Result.new(:hello, "World")
	`)

	if r, ok := r.(*Result); ok {
		if r.empty {
			t.Error("Should not be empty")
		}

	} else {
		t.Fatalf("Wrong type %T", r)
	}
}

func TestResultUsage(t *testing.T) {
	r := vm.ExecAndReturn(t, `
		require("result")
		r = Result.new(:hello, "World")
		out = 100
		r.hello do |x|
		  out = true
		end.or do |name, type|
		  out = false
		end
		return out
	`)
	vm.VerifyExpected(t, 0, r, true)
}
