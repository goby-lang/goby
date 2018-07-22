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
	b.Run("thread_and_channel", func(b *testing.B) {
		script := `
THREAD_COUNT = 503
RES = Channel.new

class Receiver
  def initialize(name)
    @name = name
    @mailbox = Channel.new
  end

  def next=(n)
    @next = n
  end

  def put(msg)
    @mailbox.deliver(msg)
  end

  def message_loop
    while true do
      msg = @mailbox.receive
      if msg == 0
        RES.deliver(@name)
      else
        @next.put(msg - 1)
      end
    end
  end

  def run
    thread do
      message_loop
    end
  end
end

receivers = []
THREAD_COUNT.times do |i|
  r = Receiver.new(i + 1)
  receivers[i] = r
  if i > 0
    receivers[i-1].next = r
  end
end

receivers[THREAD_COUNT - 1].next = receivers[0]

THREAD_COUNT.times do |i|
  receivers[i].run
end

receivers[0].put(1000)
RES.receive
`
		runBench(b, script)
	})
}