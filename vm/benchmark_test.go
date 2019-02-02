package vm

import (
	"testing"

	"github.com/gooby-lang/gooby/compiler"
	"github.com/gooby-lang/gooby/compiler/parser"
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

func BenchmarkConcurrency(b *testing.B) {
	b.Run("concurrency", func(b *testing.B) {
		script := `
		c = Channel.new

		1001.times do |i| # i start from 0 to 1000
		  thread do
		  	c.deliver(i)
		  end
		end

		r = 0
		1001.times do
		  r = r + c.receive
		end

		r
`
		runBench(b, script)
	})
}

func BenchmarkContextSwitch(b *testing.B) {
	b.Run("fib", func(b *testing.B) {
		script := `
		def fib(n)
			if n <= 1
				return n
			else
				return fib(n - 1) + fib(n - 2)
			end
		end

		25.times do |i|
			fib(i/2)
		end
`
		runBench(b, script)
	})

	b.Run("quicksort", func(b *testing.B) {
		script := `
		def quicksort(arr, l, r)
			if l >= r
				return
			end

			pivot = arr[l]
			i = l - 1
			j = r + 1

			while true do
				i += 1
				while arr[i] < pivot do
					i += 1
				end

				j -= 1
				while arr[j] > pivot do
					j -= 1
				end

				if i >= j
					break
				end

				# swap
				tmp = arr[i]
				arr[i] = arr[j]
				arr[j] = tmp
			end


			quicksort(arr, l, j)
			quicksort(arr, j + 1, r)
		end

		arr = [ 0, 5, 3, 2, 5, 7, 3, 5, 6, 9] * 10
		quicksort(arr, 0, arr.length - 1)
`
		runBench(b, script)
	})
}
