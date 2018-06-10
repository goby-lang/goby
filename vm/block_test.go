package vm

import "testing"

func TestBlockInitialize(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
b = Block.new do
  100
end

b.call
`, 100},

		{`
def baz
  1000
end

class Foo
  def exec_block(block)
	block.call
  end

  def baz
    100
  end
end

b = Block.new do
  baz
end

f = Foo.new
f.exec_block(b)
`, 1000},
		{`
bl = Block.new do |array|
  array.reduce do |sum, i|
    sum + i
  end
end
bl.call([1, 2, 3, 4])`, 10},
		{`
p = Block.new do |i, j, k|
  [i, j, k]
end
p.call(1, 2, 3, 4, 5)`, []interface{}{1, 2, 3}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
