package igb

import (
	"bytes"
	"strings"
	"testing"
)

func TestStartIgb(t *testing.T) {
	tests := []struct {
		inputs   []string
		expected string
	}{
		{
			[]string{
				"a = 10",
				"a + 1",
				"exit",
			},
			`>> a = 10
>> a + 1
#=> 11
>> exit`},
		{
			[]string{`class Foo
  def add_ten(x)
    x + 10
  end

  def minus_two(x)
    x - 2
  end
end

f = Foo.new
a = f.add_ten(100)
f.minus_two(a)
				`,
				"exit",
			}, `>> class Foo
  def add_ten(x)
    x + 10
  end

  def minus_two(x)
    x - 2
  end
end

f = Foo.new
a = f.add_ten(100)
f.minus_two(a)

#=> 108
>> exit`},
		{
			[]string{
				"class Foo",
				"  attr_accessor :bar",
				"end",
				"f = Foo.new",
				"f.bar = 10",
				"f.bar",
				"exit",
			},
			`
>> class Foo
>>   attr_accessor :bar
>> end
#=> <Class:Foo>
>> f = Foo.new
>> f.bar = 10
#=> 10
>> f.bar
#=> 10
>> exit`},
	}

	for _, test := range tests {
		ch := make(chan string)
		out := &bytes.Buffer{}

		go func() {
			for _, input := range test.inputs {
				ch <- input
			}
		}()

		Start(ch, out)

		result := out.String()

		compareIO(t, result, test.expected)
	}

}

func compareIO(t *testing.T, value, expected string) {
	value = removeTabs(strings.TrimSpace(value))
	expected = strings.TrimSpace(expected)
	if value != expected {
		t.Fatalf(`
Expect:
"%q"

Got:
"%q"
`, expected, value)
	}
}

func removeTabs(s string) string {
	return strings.Replace(s, "\t", "", -1)
}
