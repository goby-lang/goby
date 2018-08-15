package vm

import (
	"testing"

	"github.com/goby-lang/goby/compiler"
	"github.com/goby-lang/goby/compiler/parser"
)

func runBench(b *testing.B, input string) {
	b.Helper()
	iss, err := compiler.CompileToInstructions(input, parser.NormalMode)

	if err != nil {
		b.Errorf("Error when compiling input: %s", input)
		b.Fatal(err.Error())
	}
	v := initTestVM()
	filepath := getFilename()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v.ExecInstructions(iss, filepath)
	}
}

func BenchmarkBasicMath(b *testing.B) {
	b.Run("add", func(b *testing.B) {
		runBench(b, `
			a = 1
			b = 2
			c = a + b
		`)
	})
	b.Run("subtract", func(b *testing.B) {
		runBench(b, `
			a = 1
			b = 2
			c = a - b
		`)
	})
	b.Run("multiply", func(b *testing.B) {
		runBench(b, `
			a = 1
			b = 2
			c = a * b
		`)
	})
	b.Run("divide", func(b *testing.B) {
		runBench(b, `
			a = 1
			b = 2
			c = a % b
		`)
	})
}
