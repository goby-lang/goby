package vm

import (
	"testing"
)

func TestRegexpClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Regexp.class.name`, "Class"},
		{`Regexp.superclass.name`, "Object"},
		{`Regexp.ancestors.to_s`, "[Regexp, Object]"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestRegexpClassCreation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		//{`re = Regexp.new("")`, ""}, // FIXME
		{`"Hello ".concat("World")`, "Hello World"},
		//{`Regexp.new('🍣Goby🍺').class`, "Regexp"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestRegexpMatch(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`Regexp.new("Goby").match?("Hello, Goby!")`, true},
		{`Regexp.new("Python").match?("Hello, Goby!")`, false},
		{`Regexp.new("Hello Goby!").match?("Goby")`, false},
		{`Regexp.new("GOBY").match?("Hello, Goby!")`, false}, // TOFGobyIX
		{`Regexp.new("234").match?("Hello, 1234567890!")`, true},
		{`Regexp.new(" 234").match?("Hello, 1234567890!")`, false},

		// The followings are based upon Onigmo's test pattern (thanks!):
		// https://github.com/k-takata/Onigmo/blob/master/test.rb
		{`Regexp.new("").match?('')`, true},
		{`Regexp.new("a").match?('a')`, true},
		{`Regexp.new("a").match?('a')`, true},
		{`Regexp.new("b").match?('abc')`, true},
		{`Regexp.new("b").match?('abc')`, true},
		{`Regexp.new(".").match?('a')`, true},
		{`Regexp.new(".*").match?('abcde fgh')`, true},
		{`Regexp.new("a*").match?('aaabbc')`, true},
		{`Regexp.new("a+").match?('aaabbc')`, true},
		{`Regexp.new("a?").match?('bac')`, true},
		{`Regexp.new("a??").match?('bac')`, true},
		{`Regexp.new("abcde").match?('abcdeavcd')`, true},
		{`Regexp.new("\w\d\s").match?('  a2 aa $3 ')`, true},
		{`Regexp.new("[c-f]aa[x-z]").match?('3caaycaaa')`, true},
		{`Regexp.new("(?i:fG)g").match?('fGgFggFgG')`, true},
		{`Regexp.new("a|b").match?('b')`, true},
		{`Regexp.new("ab|bc|cd").match?('bcc')`, true},
		{`Regexp.new("(ffy)\1").match?('ffyffyffy')`, true},
		{`Regexp.new("|z").match?('z')`, true},
		{`Regexp.new("^az").match?('azaz')`, true},
		{`Regexp.new("az$").match?('azaz')`, true},
		{`Regexp.new("(((.a)))\3").match?('zazaaa')`, true},
		{`Regexp.new("(ac*?z)\1").match?('aacczacczacz')`, true},
		//{`Regexp.new("aaz{3, 4}").match?('bbaabbaazzzaazz')`, true},
		//{`Regexp.new("\000a").match?("b\000a")`, true},
		//{`Regexp.new("ff\xfe").match?("fff\xfe")`, true},
		//{`Regexp.new("...abcdefghijklmnopqrstuvwxyz").match?('zzzzzabcdefghijklmnopqrstuvwxyz')`, true},

		{`Regexp.new("あ").match?('あ')`, true},
		{`Regexp.new("い").match?('あいう')`, true},
		{`Regexp.new(".").match?('あ')`, true},
		{`Regexp.new(".*").match?('あいうえお かきく')`, true},
		{`Regexp.new(".*えお").match?('あいうえお かきく')`, true},
		{`Regexp.new("あ*").match?('あああいいう')`, true},
		{`Regexp.new("あ+").match?('あああいいう')`, true},
		{`Regexp.new("あ?").match?('いあう')`, true},
		{`Regexp.new("全??").match?('負全変')`, true},
		{`Regexp.new("a辺c漢e").match?('a辺c漢eavcd')`, true},
		//{`Regexp.new("(?u)\w\d\s").match?('  あ2 うう $3 ')`, true},
		{`Regexp.new("[う-お]ああ[と-ん]").match?('3うああなうあああ')`, true},
		{`Regexp.new("あ|い").match?('い')`, true},
		{`Regexp.new("あい|いう|うえ").match?('いうう')`, true},
		{`Regexp.new("(ととち)\1").match?('ととちととちととち')`, true},
		{`Regexp.new("|え").match?('え')`, true},
		{`Regexp.new("^あず").match?('あずあず')`, true},
		{`Regexp.new("あず$").match?('あずあず')`, true},
		{`Regexp.new("(((.あ)))\3").match?('zあzあああ')`, true},
		{`Regexp.new("(あう*?ん)\1").match?('ああううんあううんあうん')`, true},
		{`Regexp.new("ああん{3,4}").match?('ててああいいああんんんああんああん')`, true},
		//{`Regexp.new("\000あ").match?("い\000あ")`, true},
		//{`Regexp.new("とと\xfe\xfe").match?("ととと\xfe\xfe")`, true},
		{`Regexp.new("...あいうえおかきくけこさしすせそ").match?('zzzzzあいうえおかきくけこさしすせそ')`, true},

		{`Regexp.new("🍣").match?('🍣')`, true},
		{`Regexp.new("🍣").match?('あ🍣う')`, true},
		{`Regexp.new(".").match?('🍣')`, true},
		{`Regexp.new(".*").match?('あい🍣うえお か🍺きく')`, true},
		{`Regexp.new(".*えお").match?('あいう🍣えお🍺 かきく')`, true},
		{`Regexp.new("あ*").match?('🍣あああ🍺あいいう')`, true},
		{`Regexp.new("あ+").match?('🍣あああ🍺あいいう')`, true},
		{`Regexp.new("あ?").match?('いあ🍣う')`, true},
		{`Regexp.new("全??").match?('負🍣🍺')`, true},
		{`Regexp.new("a🍣c🍺e").match?('a🍣c🍺eavcd')`, true},
		//{`Regexp.new("(?u)\w\d\s").match?('  あ2 うう $3 ')`, true},
		{`Regexp.new("[う-お]🍣🍺[と-ん]").match?('3う🍣🍺なう🍣🍺あ')`, true},
		{`Regexp.new("🍣|🍺").match?('🍣')`, true},
		{`Regexp.new("🍣🍺|🍺😍|😍🤡").match?('🍺😍😍')`, true},
		{`Regexp.new("(ととち)\1").match?('ととちととちととち')`, true},
		{`Regexp.new("|え").match?('え')`, true},
		{`Regexp.new("^あず").match?('あずあず')`, true},
		{`Regexp.new("あず$").match?('あずあず')`, true},
		{`Regexp.new("(((.あ)))\3").match?('zあzあああ')`, true},
		{`Regexp.new("(あう*?ん)\1").match?('ああううんあううんあうん')`, true},
		{`Regexp.new("ああん{3,4}").match?('ててああいいああんんんああんああん')`, true},
		{`Regexp.new("\000あ").match?("い\000あ")`, true},
		{`Regexp.new("とと\xfe\xfe").match?("ととと\xfe\xfe")`, true},
		{`Regexp.new("...あいうえおかきくけこさしすせそ").match?('zzzzzあいうえおかきくけこさしすせそ')`, true},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
		vm.checkSP(t, i, 1)
	}
}
