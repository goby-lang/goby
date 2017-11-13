package vm

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// StringObject represents string instances
// String object holds and manipulates a sequence of characters.
// String objects may be created using as string literals.
// Double or single quotations can be used for representation.
//
// ```ruby
// a = "Three"
// b = 'zero'
// c = '漢'
// d = 'Tiếng Việt'
// e = "😏️️"
// ```
//
// **Note:**
//
// - Currently, manipulations are based upon Golang's Unicode manipulations.
// - Currently, UTF-8 encoding is assumed based upon Golang's string manipulation, but the encoding is not actually specified(TBD).
// - `String.new` is not supported.
type StringObject struct {
	*baseObj
	value string
}

// Class methods --------------------------------------------------------
func builtinStringClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			// The String.fmt implements formatted I/O with functions analogous to C's printf and scanf
			// Currently only support plain "%s" formatting
			// TODO: Support other kind of formatting such as %f, %v ... etc
			//
			// ```ruby
			// String.fmt("Hello! %s Lang!", "Goby")                    # => "Hello! Goby Lang!"
			// String.fmt("I love to eat %s and %s!", "Sushi", "Ramen") # => "I love to eat Sushi and Ramen"
			// ```
			//
			// @return [String]
			Name: "fmt",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) < 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect at least 1 argument. got=%v", strconv.Itoa(len(args)))
					}

					formatObj, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
					}

					format := formatObj.value
					arguments := []interface{}{}

					for _, arg := range args[1:] {
						arguments = append(arguments, arg.toString())
					}

					count := strings.Count(format, "%s")

					if len(args[1:]) != count {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect %d string arguments. got=%d", count, len(args[1:]))
					}

					return t.vm.initStringObject(fmt.Sprintf(format, arguments...))
				}
			},
		},
		{
			Name: "new",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					return t.vm.initUnsupportedMethodError(sourceLine, "#new", receiver)
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinStringInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{

		{
			// Returns the concatenation of self and another String
			//
			// ```ruby
			// "first" + "-second" # => "first-second"
			// ```
			//
			// @return [String]
			Name: "+",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					leftValue := receiver.(*StringObject).value
					r := args[0]
					right, ok := r.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, r.Class().Name)
					}

					rightValue := right.value
					return t.vm.initStringObject(leftValue + rightValue)
				}
			},
		},
		{
			// Returns self multiplying another Integer
			//
			// ```ruby
			// "string " * 2 # => "string string string "
			// ```
			//
			// @return [String]
			Name: "*",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					leftValue := receiver.(*StringObject).value
					r := args[0]
					right, ok := r.(*IntegerObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, r.Class().Name)
					}

					if right.value < 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Second argument must be greater than or equal to 0. got=%v", right.value)
					}

					var result string

					for i := 0; i < right.value; i++ {
						result += leftValue
					}

					return t.vm.initStringObject(result)
				}
			},
		},
		{
			// Returns a Boolean if first string greater than second string
			//
			// ```ruby
			// "a" < "b" # => true
			// ```
			//
			// @return [Boolean]
			Name: ">",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					leftValue := receiver.(*StringObject).value
					r := args[0]
					right, ok := r.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, r.Class().Name)
					}

					rightValue := right.value

					if leftValue > rightValue {
						return TRUE
					}

					return FALSE
				}
			},
		},
		{
			// Returns a Boolean if first string less than second string
			//
			// ```ruby
			// "a" < "b" # => true
			// ```
			//
			// @return [Boolean]
			Name: "<",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					leftValue := receiver.(*StringObject).value
					r := args[0]
					right, ok := r.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, r.Class().Name)
					}

					rightValue := right.value

					if leftValue < rightValue {
						return TRUE
					}

					return FALSE
				}
			},
		},
		{
			// Returns a Boolean of compared two strings
			//
			// ```ruby
			// "first" == "second" # => false
			// "two" == "two" # => true
			// ```
			//
			// @return [Boolean]
			Name: "==",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					leftValue := receiver.(*StringObject).value
					r := args[0]
					right, ok := r.(*StringObject)

					if !ok {
						return FALSE
					}

					rightValue := right.value

					if leftValue == rightValue {
						return TRUE
					}

					return FALSE
				}
			},
		},
		{
			// Matches the receiver with a Regexp
			//
			// ```ruby
			// "pizza" =~ Regex.new("zz")  # => 2
			// "pizza" =~ Regex.new("OH!") # => nil
			// ```
			//
			// @return [Integer]
			Name: "=~",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%d", len(args))
					}

					arg := args[0]

					regexp, ok := arg.(*RegexpObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.RegexpClass, arg.Class().Name)
					}

					text := receiver.(*StringObject).value

					match, _ := regexp.Regexp.FindStringMatch(text)

					if match == nil {
						return NULL
					}

					position := match.Groups()[0].Captures[0].Index

					return t.vm.initIntegerObject(position)
				}
			},
		},
		{
			// Returns a Integer. If first string is less than second string returns -1, if equal to returns 0, if greater returns 1
			//
			//
			// ```ruby
			// "abc" <=> "abcd" # => -1
			// "abc" <=> "abc" # => 0
			// "abcd" <=> "abc" # => 1
			// ```
			//
			// @return [Integer]
			Name: "<=>",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					leftValue := receiver.(*StringObject).value
					r := args[0]
					right, ok := r.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, r.Class().Name)
					}

					rightValue := right.value

					if leftValue < rightValue {
						return t.vm.initIntegerObject(-1)
					}
					if leftValue > rightValue {
						return t.vm.initIntegerObject(1)
					}

					return t.vm.initIntegerObject(0)
				}
			},
		},
		{
			// Returns a Boolean of compared two strings
			//
			// ```ruby
			// "first" != "second" # => true
			// "two" != "two" # => false
			// ```
			//
			// @return [Boolean]
			Name: "!=",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					leftValue := receiver.(*StringObject).value
					right, ok := args[0].(*StringObject)

					if !ok {
						return TRUE
					}

					rightValue := right.value

					if leftValue != rightValue {
						return TRUE
					}

					return FALSE
				}
			},
		},
		{
			// Returns the character of the string with specified index
			// It will raise error if the input is not an Integer type
			//
			// ```ruby
			// "Hello"[1]        # => "e"
			// "Hello"[5]        # => nil
			// "Hello\nWorld"[5] # => "\n"
			// "Hello"[-1]       # => "o"
			// "Hello"[-6]       # => nil
			// "Hello😊"[5]      # => "😊"
			// "Hello😊"[-1]     # => "😊"
			// ```
			//
			// @return [String]
			Name: "[]",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%d", len(args))
					}

					str := receiver.(*StringObject).value
					i := args[0]
					index, ok := i.(*IntegerObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, i.Class().Name)
					}

					indexValue := index.value

					if indexValue < 0 {
						strLength := utf8.RuneCountInString(str)
						if -indexValue > strLength {
							return NULL
						}
						return t.vm.initStringObject(string([]rune(str)[strLength+indexValue]))
					}

					if len(str) > indexValue {
						return t.vm.initStringObject(string([]rune(str)[indexValue]))
					}
					return NULL
				}
			},
		},
		{
			// Replace character of the string with input string
			// It will raise error if the index is not Integer type or the index value is out of
			// range of the string length
			//
			// Currently only support assign string type value
			// TODO: Support to assign type which have to_s method
			//
			// ```ruby
			// "Ruby"[1] = "oo" # => "Rooby"
			// "Go"[2] = "by"   # => "Goby"
			// "Hello\nWorld"[5] = " " # => "Hello World"
			// "Ruby"[-3] = "oo" # => "Rooby"
			// "Hello😊"[5] = "🐟" # => "Hello🐟"
			// ```
			//
			// @return [String]
			Name: "[]=",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 2 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 2 arguments. got=%v", strconv.Itoa(len(args)))
					}

					str := receiver.(*StringObject).value
					i := args[0]
					index, ok := i.(*IntegerObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, i.Class().Name)
					}

					indexValue := index.value
					strLength := utf8.RuneCountInString(str)

					if strLength < indexValue {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Index value out of range. got=%v", strconv.Itoa(indexValue))
					}

					r := args[1]
					replaceStr, ok := r.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, r.Class().Name)
					}
					replaceStrValue := replaceStr.value

					// Negative Index Case
					if indexValue < 0 {
						if -indexValue > strLength {
							return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Index value out of range. got=%v", strconv.Itoa(indexValue))
						}
						// Change to positive index to replace the string
						indexValue += strLength
					}

					if strLength == indexValue {
						return t.vm.initStringObject(str + replaceStrValue)
					}
					// Using rune type to support UTF-8 encoding to replace character
					result := string([]rune(str)[:indexValue]) + replaceStrValue + string([]rune(str)[indexValue+1:])
					return t.vm.initStringObject(result)
				}
			},
		},
		{
			// Return a new String with the first character converted to uppercase but the rest of string converted to lowercase.
			//
			// ```ruby
			// "test".capitalize         # => "Test"
			// "tEST".capitalize         # => "Test"
			// "heLlo\nWoRLd".capitalize # => "Hello\nworld"
			// "😊HeLlO🐟".capitalize    # => "😊hello🐟"
			// ```
			//
			// @return [String]
			Name: "capitalize",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					str := receiver.(*StringObject).value
					start := string([]rune(str)[0])
					rest := string([]rune(str)[1:])
					result := strings.ToUpper(start) + strings.ToLower(rest)

					return t.vm.initStringObject(result)
				}
			},
		},
		{
			// Returns a string with the last character chopped
			//
			// ```ruby
			// "Hello".chop         # => "Hell"
			// "Hello World\n".chop # => "Hello World"
			// "Hello😊".chop       # => "Hello"
			// ```
			//
			// @return [String]
			Name: "chop",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					str := receiver.(*StringObject).value
					strLength := utf8.RuneCountInString(str)

					// Support UTF-8 Encoding
					return t.vm.initStringObject(string([]rune(str)[:strLength-1]))
				}
			},
		},
		{
			// Returns a string which is concatenate with the input string or character
			//
			// ```ruby
			// "Hello ".concat("World")   # => "Hello World"
			// "Hello World".concat("😊") # => "Hello World😊"
			// ```
			//
			// @return [String]
			Name: "concat",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%v", strconv.Itoa(len(args)))
					}

					str := receiver.(*StringObject).value
					c := args[0]
					concatStr, ok := c.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, c.Class().Name)
					}

					return t.vm.initStringObject(str + concatStr.value)
				}
			},
		},
		{
			// Returns the integer that count the string chars as UTF-8
			//
			// ```ruby
			// "abcde".count          # => 5
			// "哈囉！世界！".count     # => 6
			// "Hello\nWorld".count   # => 11
			// "Hello\nWorld😊".count # => 12
			// ```
			//
			// @return [Integer]
			Name: "count",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					str := receiver.(*StringObject).value

					// Support UTF-8 Encoding
					return t.vm.initIntegerObject(utf8.RuneCountInString(str))
				}
			},
		},
		{
			// Returns a string which is being partially deleted with specified values
			//
			// ```ruby
			// "Hello hello HeLlo".delete("el")        # => "Hlo hlo HeLlo"
			// "Hello 😊 Hello 😊 Hello".delete("😊") # => "Hello  Hello  Hello"
			// # TODO: Handle delete intersection of multiple strings' input case
			// "Hello hello HeLlo".delete("el", "e") # => "Hllo hllo HLlo"
			// ```
			//
			// @return [String]
			Name: "delete",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%v", strconv.Itoa(len(args)))
					}

					str := receiver.(*StringObject).value
					d := args[0]
					deleteStr, ok := d.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, d.Class().Name)
					}

					return t.vm.initStringObject(strings.Replace(str, deleteStr.value, "", -1))
				}
			},
		},
		{
			// Returns a new String with all characters is lowercase
			//
			// ```ruby
			// "erROR".downcase        # => "error"
			// "HeLlO\tWorLD".downcase # => "hello\tworld"
			// ```
			//
			// @return [String]
			Name: "downcase",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					str := receiver.(*StringObject).value

					return t.vm.initStringObject(strings.ToLower(str))
				}
			},
		},
		{
			// Split and loop through the string byte
			//
			// ```ruby
			// "Sushi 🍣".each_byte do |byte|
			//   puts byte
			// end
			// # => 83  # "S"
			// # => 117 # "u"
			// # => 115 # "s"
			// # => 104 # "h"
			// # => 105 # "i"
			// # => 32  # " "
			// # => 240 # "\xF0"
			// # => 159 # "\x9F"
			// # => 141 # "\x8D"
			// # => 163 # "\xA3"
			// ```
			//
			// @return [String]
			Name: "each_byte",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					str := receiver.(*StringObject).value

					for _, byte := range []byte(str) {
						t.builtinMethodYield(blockFrame, t.vm.initIntegerObject(int(byte)))
					}

					return t.vm.initStringObject(str)
				}
			},
		},
		{
			// Split and loop through the string characters
			//
			// ```ruby
			// "Sushi 🍣".each_char do |char|
			//   puts char
			// end
			// # => "S"
			// # => "u"
			// # => "s"
			// # => "h"
			// # => "i"
			// # => " "
			// # => "🍣"
			// ```
			//
			// @return [String]
			Name: "each_char",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					str := receiver.(*StringObject).value

					for _, char := range []rune(str) {
						t.builtinMethodYield(blockFrame, t.vm.initStringObject(string(char)))
					}

					return t.vm.initStringObject(str)
				}
			},
		},
		{
			// Split and loop through the string segment split by the newline escaped character
			//
			// ```ruby
			// "Hello\nWorld\nGoby".each_line do |line|
			//   puts line
			// end
			// # => "Hello"
			// # => "World"
			// # => "Goby"
			// ```
			//
			// @return [String]
			Name: "each_line",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 0 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 0 argument. got=%d", len(args))
					}

					if blockFrame == nil {
						return t.vm.initErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
					}

					str := receiver.(*StringObject).value
					lineArray := strings.Split(str, "\n")

					for _, line := range lineArray {
						t.builtinMethodYield(blockFrame, t.vm.initStringObject(line))
					}

					return t.vm.initStringObject(str)
				}
			},
		},
		{
			// Returns true if string is empty value
			//
			// ```ruby
			// "".empty?      # => true
			// "Hello".empty? # => false
			// ```
			//
			// @return [Boolean]
			Name: "empty?",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					str := receiver.(*StringObject).value

					if str == "" {
						return TRUE
					}
					return FALSE
				}
			},
		},
		{
			// Returns true if receiver string end with the argument string
			//
			// ```ruby
			// "Hello".end_with?("llo")     # => true
			// "Hello".end_with?("ell")     # => false
			// "😊Hello🐟".end_with?("🐟") # => true
			// "😊Hello🐟".end_with?("😊") # => false
			// ```
			//
			// @return [Boolean]
			Name: "end_with?",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%v", strconv.Itoa(len(args)))
					}

					str := receiver.(*StringObject).value
					c := args[0]
					compareStr, ok := c.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, c.Class().Name)
					}

					compareStrValue := compareStr.value
					compareStrLength := utf8.RuneCountInString(compareStrValue)
					strLength := utf8.RuneCountInString(str)

					if compareStrLength > strLength {
						return FALSE
					}

					if compareStrValue == string([]rune(str)[strLength-compareStrLength:]) {
						return TRUE
					}
					return FALSE
				}
			},
		},
		{
			// Returns true if receiver string is equal to argument string
			//
			// ```ruby
			// "Hello".eql?("Hello")     # => true
			// "Hello".eql?("World")     # => false
			// "Hello😊".eql?("Hello😊") # => true
			// ```
			//
			// @return [Boolean]
			Name: "eql?",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%v", strconv.Itoa(len(args)))
					}

					str := receiver.(*StringObject).value
					compareStr, ok := args[0].(*StringObject)

					if !ok {
						return FALSE
					} else if compareStr.value == str {
						return TRUE
					}
					return FALSE
				}
			},
		},
		{
			// Returns a copy of str with the all occurrences of pattern substituted for the second argument.
			// The pattern is typically a String or Regexp (Not implemented yet); if given as a String, any
			// regular expression metacharacters it contains will be interpreted literally, e.g. '\\d' will
			// match a backslash followed by ‘d’, instead of a digit.
			//
			// Currently only support string version of String#gsub.
			//
			// ```ruby
			// "Ruby Lang".gsub("Ru", "Go")                # => "Goby Lang"
			// "Hello 😊 Hello 😊 Hello".gsub("😊", "🐟") # => "Hello 🐟 Hello 🐟 Hello"
			// ```
			//
			// @return [String]
			Name: "gsub",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 2 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 2 arguments. got=%v", len(args))
					}

					str := receiver.(*StringObject).value

					p := args[0]
					pattern, ok := p.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, "Expect pattern to be String. got: %s", p.Class().Name)
					}

					r := args[1]
					replacement, ok := r.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, "Expect replacement to be String. got: %s", r.Class().Name)
					}

					return t.vm.initStringObject(strings.Replace(str, pattern.value, replacement.value, -1))
				}
			},
		},
		{
			// Checks if the specified string is included in the receiver
			//
			// ```ruby
			// "Hello\nWorld".include?("\n")   # => true
			// "Hello 😊 Hello".include?("😊") # => true
			// ```
			//
			// @return [Bool]
			Name: "include?",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%v", strconv.Itoa(len(args)))
					}

					str := receiver.(*StringObject).value
					i := args[0]
					includeStr, ok := i.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, i.Class().Name)
					}

					if strings.Contains(str, includeStr.value) {
						return TRUE
					}

					return FALSE
				}
			},
		},
		{
			// Insert a string input in specified index value of the receiver string
			//
			// It will raise error if index value is not an integer or index value is out
			// of receiver string's range
			//
			// It will also raise error if the input string value is not type string
			//
			// ```ruby
			// "Hello".insert(0, "X") # => "XHello"
			// "Hello".insert(2, "X") # => "HeXllo"
			// "Hello".insert(5, "X") # => "HelloX"
			// "Hello".insert(-1, "X") # => "HelloX"
			// "Hello".insert(-3, "X") # => "HelXlo"
			// ```
			//
			// @return [String]
			Name: "insert",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 2 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 2 arguments. got=%d", len(args))
					}

					str := receiver.(*StringObject).value
					i := args[0]
					index, ok := i.(*IntegerObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, i.Class().Name)
					}

					indexValue := index.value
					ins := args[1]
					insertStr, ok := ins.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, "Expect insert string to be String. got: %s", ins.Class().Name)
					}
					strLength := utf8.RuneCountInString(str)

					if indexValue < 0 {
						if -indexValue > strLength+1 {
							return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Index value out of range. got=%v", indexValue)
						} else if -indexValue == strLength+1 {
							return t.vm.initStringObject(insertStr.value + str)
						}
						// Change it to positive index value to replace the string via index
						indexValue += strLength
					}

					if strLength < indexValue {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Index value out of range. got=%v", indexValue)
					}

					// Support UTF-8 Encoding
					return t.vm.initStringObject(string([]rune(str)[:indexValue]) + insertStr.value + string([]rune(str)[indexValue:]))
				}
			},
		},
		{
			// Returns the character length of self
			// **Note:** the length is currently byte-based, instead of charcode-based.
			//
			// ```ruby
			// "zero".length # => 4
			// "".length     # => 0
			// "😊".length   # => 1
			// ```
			//
			// @return [Integer]
			Name: "length",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					str := receiver.(*StringObject).value

					// Support UTF-8 Encoding
					return t.vm.initIntegerObject(utf8.RuneCountInString(str))
				}
			},
		},
		{
			// If input integer is greater than the length of receiver string, returns a new String of
			// length integer with receiver string left justified and padded with default " "; otherwise,
			// returns receiver string.
			//
			// It will raise error if the input string length is not integer type
			//
			// ```ruby
			// "Hello".ljust(2)           # => "Hello"
			// "Hello".ljust(7)           # => "Hello  "
			// "Hello".ljust(10, "xo")    # => "Helloxoxox"
			// "Hello".ljust(10, "😊🐟") # => "Hello😊🐟😊🐟😊"
			// ```
			//
			// @return [String]
			Name: "ljust",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 && len(args) != 2 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1..2 arguments. got=%v", strconv.Itoa(len(args)))
					}

					str := receiver.(*StringObject).value

					l := args[0]
					strLength, ok := l.(*IntegerObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, "Expect justify width to be Integer. got: %s", l.Class().Name)
					}

					strLengthValue := strLength.value

					var padStrValue string
					if len(args) == 1 {
						padStrValue = " "
					} else {
						p := args[1]
						padStr, ok := p.(*StringObject)

						if !ok {
							return t.vm.initErrorObject(errors.TypeError, sourceLine, "Expect padding string to be String. got: %s", p.Class().Name)
						}

						padStrValue = padStr.value
					}

					currentStrLength := utf8.RuneCountInString(str)
					padStrLength := utf8.RuneCountInString(padStrValue)

					if strLengthValue > currentStrLength {
						for i := currentStrLength; i < strLengthValue; i += padStrLength {
							str += padStrValue
						}
						str = string([]rune(str)[:strLengthValue])
					}

					// Support UTF-8 Encoding
					return t.vm.initStringObject(str)
				}
			},
		},
		{
			// Returns the matching data of the regex with the given string.
			//
			// ```ruby
			// 'pow'.match(Regexp.new("o")) # => #<MatchData "o">
			// 'pow'.match(Regexp.new("x")) # => nil
			// ```
			//
			// @param string [String]
			// @return [MatchData]
			Name: "match",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%d", len(args))
					}

					arg := args[0]
					regexpObj, ok := arg.(*RegexpObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.RegexpClass, arg.Class().Name)
					}

					regexp := regexpObj.Regexp
					text := receiver.(*StringObject).value

					match, _ := regexp.FindStringMatch(text)

					if match == nil {
						return NULL
					}

					return t.vm.initMatchDataObject(match, regexp.String(), text)
				}
			},
		},
		{
			// Return a string replaced by the input string
			//
			// ```ruby
			// "Hello".replace("World")          # => "World"
			// "你好"replace("再見")              # => "再見"
			// "Ruby\nLang".replace("Goby\nLang") # => "Goby Lang"
			// "Hello😊".replace("World🐟")      # => "World🐟"
			// ```
			//
			// @return [String]
			Name: "replace",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%v", strconv.Itoa(len(args)))
					}

					r := args[0]
					replaceStr, ok := r.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, r.Class().Name)
					}

					return t.vm.initStringObject(replaceStr.value)
				}
			},
		},
		{
			// Returns a new String with reverse order of self
			// **Note:** the length is currently byte-based, instead of charcode-based.
			//
			// ```ruby
			// "reverse".reverse           # => "esrever"
			// "Hello\nWorld".reverse      # => "dlroW\nolleH"
			// "Hello 😊🐟 World".reverse # => "dlroW 🐟😊 olleH"
			// ```
			//
			// @return [String]
			Name: "reverse",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					str := receiver.(*StringObject).value

					var revert string
					for i := utf8.RuneCountInString(str) - 1; i >= 0; i-- {
						revert += string([]rune(str)[i])
					}

					// Support UTF-8 Encoding
					return t.vm.initStringObject(revert)
				}
			},
		},
		{
			// If input integer is greater than the length of receiver string, returns a new String of
			// length integer with receiver string right justified and padded with default " "; otherwise,
			// returns receiver string.
			//
			// It will raise error if the input string length is not integer type
			//
			// ```ruby
			// "Hello".rjust(2)          # => "Hello"
			// "Hello".rjust(7)          # => "  Hello"
			// "Hello".rjust(10, "xo")   # => "xoxoxHello"
			// "Hello".rjust(10, "😊🐟") # => "😊🐟😊🐟😊Hello"
			// ```
			//
			// @return [String]
			Name: "rjust",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 && len(args) != 2 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1..2 arguments. got=%v", strconv.Itoa(len(args)))
					}

					str := receiver.(*StringObject).value
					l := args[0]
					strLength, ok := l.(*IntegerObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, "Expect justify width to be Integer. got: %s", l.Class().Name)
					}

					strLengthValue := strLength.value

					var padStrValue string
					if len(args) == 1 {
						padStrValue = " "
					} else {
						p := args[1]
						padStr, ok := p.(*StringObject)

						if !ok {
							return t.vm.initErrorObject(errors.TypeError, sourceLine, "Expect padding string to be String. got: %s", p.Class().Name)
						}

						padStrValue = padStr.value
					}

					padStrLength := utf8.RuneCountInString(padStrValue)

					if strLengthValue > len(str) {
						origin := str
						originStrLength := utf8.RuneCountInString(origin)
						for i := originStrLength; i < strLengthValue; i += padStrLength {
							str = padStrValue + str
						}
						currentStrLength := utf8.RuneCountInString(str)
						if currentStrLength > strLengthValue {
							chopLength := currentStrLength - strLengthValue
							str = string([]rune(str)[:currentStrLength-originStrLength-chopLength]) + origin
						}
					}

					// Support UTF-8 Encoding
					return t.vm.initStringObject(str)
				}
			},
		},
		{
			// Returns the character length of self
			// **Note:** the length is currently byte-based, instead of charcode-based.
			//
			// ```ruby
			// "zero".size  # => 4
			// "".size      # => 0
			// "😊".size   # => 1
			// ```
			//
			// @return [Integer]
			Name: "size",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					str := receiver.(*StringObject).value

					// Support UTF-8 Encoding
					return t.vm.initIntegerObject(utf8.RuneCountInString(str))
				}
			},
		},
		{
			// Returns a string sliced according to the input range
			//
			// ```ruby
			// "Hello World".slice(1..6)    # => "ello W"
			// "1234567890".slice(6..1)     # => ""
			// "1234567890".slice(11..1)    # => nil
			// "1234567890".slice(11..-1)   # => nil
			// "1234567890".slice(-10..1)   # => "12"
			// "1234567890".slice(-5..1)    # => ""
			// "1234567890".slice(-10..-1)  # => "1234567890"
			// "1234567890".slice(-10..-11) # => ""
			// "1234567890".slice(1..-1)    # => "234567890"
			// "1234567890".slice(1..-1234) # => ""
			// "1234567890".slice(-11..5)   # => nil
			// "1234567890".slice(-10..-5)  # => "123456"
			// "1234567890".slice(-5..-10)  # => ""
			// "1234567890".slice(-11..-12) # => nil
			// "1234567890".slice(-10..-12) # => ""
			// "Hello 😊🐟 World".slice(1..6)    # => "ello 😊"
			// "Hello 😊🐟 World".slice(-10..7)  # => "o 😊🐟"
			// "Hello 😊🐟 World".slice(1..-1)   # => "ello 😊🐟 World"
			// "Hello 😊🐟 World".slice(-12..-5) # => "llo 😊🐟 W"
			// "Hello World".slice(4)       # => "o"
			// "Hello\nWorld".slice(6)      # => "\n"
			// "Hello World".slice(-3)      # => "r"
			// "Hello World".slice(-11)     # => "H"
			// "Hello World".slice(-12)     # => nil
			// "Hello World".slice(11)      # => nil
			// "Hello World".slice(4)       # => "o"
			// "Hello 😊🐟 World".slice(6)      # => "😊"
			// "Hello 😊🐟 World".slice(-7)      # => "🐟"
			// "Hello 😊🐟 World".slice(-10)     # => "o"
			// "Hello 😊🐟 World".slice(-15)     # => nil
			// "Hello 😊🐟 World".slice(14)      # => nil
			// ```
			//
			// @return [String]
			Name: "slice",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%v", strconv.Itoa(len(args)))
					}

					str := receiver.(*StringObject).value
					strLength := utf8.RuneCountInString(str)

					// All Case Support UTF-8 Encoding
					switch args[0].(type) {
					case *RangeObject:
						ran := args[0].(*RangeObject)
						switch {
						case ran.Start >= 0 && ran.End >= 0:
							if ran.Start > strLength {
								return NULL
							} else if ran.Start > ran.End {
								return t.vm.initStringObject("")
							}
							return t.vm.initStringObject(string([]rune(str)[ran.Start : ran.End+1]))
						case ran.Start < 0 && ran.End >= 0:
							positiveStart := strLength + ran.Start
							if -ran.Start > strLength {
								return NULL
							} else if positiveStart > ran.End {
								return t.vm.initStringObject("")
							}
							return t.vm.initStringObject(string([]rune(str)[positiveStart : ran.End+1]))
						case ran.Start >= 0 && ran.End < 0:
							positiveEnd := strLength + ran.End
							if ran.Start > strLength {
								return NULL
							} else if positiveEnd < 0 || ran.Start > positiveEnd {
								return t.vm.initStringObject("")
							}
							return t.vm.initStringObject(string([]rune(str)[ran.Start : positiveEnd+1]))
						default:
							positiveStart := strLength + ran.Start
							positiveEnd := strLength + ran.End
							if positiveStart < 0 {
								return NULL
							} else if positiveStart > positiveEnd {
								return t.vm.initStringObject("")
							}
							return t.vm.initStringObject(string([]rune(str)[positiveStart : positiveEnd+1]))
						}

					case *IntegerObject:
						intValue := args[0].(*IntegerObject).value
						if intValue < 0 {
							if -intValue > strLength {
								return NULL
							}
							return t.vm.initStringObject(string([]rune(str)[strLength+intValue]))
						}
						if intValue > strLength-1 {
							return NULL
						}
						return t.vm.initStringObject(string([]rune(str)[intValue]))

					default:
						return t.vm.initErrorObject(errors.TypeError, sourceLine, "Expect slice range to be Range or Integer. got: %s", args[0].Class().Name)
					}
				}
			},
		},
		{
			// Returns an array of strings separated by the given separator
			//
			// ```ruby
			// "Hello World".split("o") # => ["Hell", " W", "rld"]
			// "Goby".split("")         # => ["G", "o", "b", "y"]
			// "Hello\nWorld\nGoby".split("o") # => ["Hello", "World", "Goby"]
			// "Hello🐟World🐟Goby".split("🐟") # => ["Hello", "World", "Goby"]
			// ```
			//
			// @return [Array]
			Name: "split",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%v", strconv.Itoa(len(args)))
					}

					s := args[0]
					seperator, ok := s.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, s.Class().Name)
					}

					str := receiver.(*StringObject).value
					arr := strings.Split(str, seperator.value)

					var elements []Object
					for i := 0; i < len(arr); i++ {
						elements = append(elements, t.vm.initStringObject(arr[i]))
					}

					return t.vm.initArrayObject(elements)
				}
			},
		},
		{
			// Returns true if receiver string start with the argument string
			//
			// ```ruby
			// "Hello".start_with("Hel")     # => true
			// "Hello".start_with("hel")     # => false
			// "😊Hello🐟".start_with("😊") # => true
			// "😊Hello🐟".start_with("🐟") # => false
			// ```
			//
			// @return [Boolean]
			Name: "start_with",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, sourceLine, "Expect 1 argument. got=%v", strconv.Itoa(len(args)))
					}

					str := receiver.(*StringObject).value
					c := args[0]
					compareStr, ok := c.(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.StringClass, c.Class().Name)
					}

					compareStrValue := compareStr.value
					compareStrLength := utf8.RuneCountInString(compareStrValue)
					strLength := utf8.RuneCountInString(str)

					if compareStrLength > strLength {
						return FALSE
					}

					if compareStrValue == string([]rune(str)[:compareStrLength]) {
						return TRUE
					}
					return FALSE
				}
			},
		},
		{
			// Returns a copy of str with leading and trailing whitespace removed.
			// Whitespace is defined as any of the following characters: null, horizontal tab,
			// line feed, vertical tab, form feed, carriage return, space.
			//
			// ```ruby
			// "  Goby Lang  ".strip   # => "Goby Lang"
			// "\nGoby Lang\r\t".strip # => "Goby Lang"
			// ```
			//
			// @return [String]
			Name: "strip",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					str := receiver.(*StringObject).value

					for {
						str = strings.Trim(str, " ")

						if strings.HasPrefix(str, "\n") || strings.HasPrefix(str, "\t") || strings.HasPrefix(str, "\r") || strings.HasPrefix(str, "\v") {
							str = string([]rune(str)[1:])
							continue
						}
						if strings.HasSuffix(str, "\n") || strings.HasSuffix(str, "\t") || strings.HasSuffix(str, "\r") || strings.HasSuffix(str, "\v") {
							str = string([]rune(str)[:utf8.RuneCountInString(str)-2])
							continue
						}
						break
					}
					return t.vm.initStringObject(str)
				}
			},
		},
		{
			// Returns an array of characters converted from a string
			//
			// ```ruby
			// "Goby".to_a       # => ["G", "o", "b", "y"]
			// "😊Hello🐟".to_a # => ["😊", "H", "e", "l", "l", "o", "🐟"]
			// ```
			//
			// @return [String]
			Name: "to_a",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					str := receiver.(*StringObject)
					strLength := utf8.RuneCountInString(str.value)
					elems := []Object{}

					for i := 0; i < strLength; i++ {
						elems = append(elems, t.vm.initStringObject(string([]rune(str.value)[i])))
					}

					return t.vm.initArrayObject(elems)
				}
			},
		},
		{
			// Returns the result of converting self to Float.
			// Unexpected characters will cause a 0.0 value, except trailing whitespace,
			// which is ignored.
			//
			// ```ruby
			// "123.5".to_f     # => 123.5
			// ".5".to_f      	# => 0.5
			// "  3.5".to_f     # => 3.5
			// "3.5e2".to_f     # => 350
			// "3.5ef".to_f     # => 0
			// ```
			//
			// @return [Float]
			Name: "to_f",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					str := receiver.(*StringObject).value

					for i, char := range str {
						if !unicode.IsSpace(char) {
							str = str[i:]
							break
						}
					}

					parsedStr, err := strconv.ParseFloat(str, 64)

					if err != nil {
						return t.vm.initFloatObject(0)
					}

					return t.vm.initFloatObject(parsedStr)
				}
			},
		},
		{
			// Returns the result of converting self to Integer
			//
			// ```ruby
			// "123".to_i       # => 123
			// "3d print".to_i  # => 3
			// "  321".to_i     # => 321
			// "some text".to_i # => 0
			// ```
			//
			// @return [Integer]
			Name: "to_i",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					str := receiver.(*StringObject).value
					parsedStr, err := strconv.ParseInt(str, 10, 0)

					if err == nil {
						return t.vm.initIntegerObject(int(parsedStr))
					}

					var digits string
					for _, char := range str {
						if unicode.IsDigit(char) {
							digits += string(char)
						} else if unicode.IsSpace(char) && len(digits) == 0 {
							// do nothing; allow trailing spaces
						} else {
							break
						}
					}

					if len(digits) > 0 {
						parsedStr, _ = strconv.ParseInt(digits, 10, 0)
						return t.vm.initIntegerObject(int(parsedStr))
					}

					return t.vm.initIntegerObject(0)
				}
			},
		},
		{
			// Returns a new String with self value
			//
			// ```ruby
			// "string".to_s # => "string"
			// ```
			//
			// @return [String]
			Name: "to_s",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					str := receiver.(*StringObject).value

					return t.vm.initStringObject(str)
				}
			},
		},
		{
			// Returns a new String with all characters is upcase
			//
			// ```ruby
			// "very big".upcase # => "VERY BIG"
			// ```
			//
			// @return [String]
			Name: "upcase",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {

					str := receiver.(*StringObject).value

					return t.vm.initStringObject(strings.ToUpper(str))
				}
			},
		},
		{
			Name: "to_bytes",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					r := receiver.(*StringObject)
					return t.vm.initGoObject([]byte(r.value))
				}
			},
		},
	}
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initStringObject(value string) *StringObject {
	return &StringObject{
		baseObj: &baseObj{class: vm.topLevelClass(classes.StringClass)},
		value:   value,
	}
}

func (vm *VM) initStringClass() *RClass {
	sc := vm.initializeClass(classes.StringClass, false)
	sc.setBuiltinMethods(builtinStringInstanceMethods(), false)
	sc.setBuiltinMethods(builtinStringClassMethods(), true)
	return sc
}

// Polymorphic helper functions -----------------------------------------

// Value returns the object
func (s *StringObject) Value() interface{} {
	return s.value
}

// toString returns the object's name as the string format
func (s *StringObject) toString() string {
	return s.value
}

// toJSON just delegates to toString
func (s *StringObject) toJSON() string {
	return strconv.Quote(s.value)
}

// equal returns true if the String values between receiver and parameter are equal
func (s *StringObject) equal(e *StringObject) bool {
	return s.value == e.value
}
