package vm

import (
	"strconv"
	"strings"
	"unicode"
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
	Value string
}

func (vm *VM) initStringObject(value string) *StringObject {
	replacer := strings.NewReplacer("\\n", "\n", "\\r", "\r", "\\t", "\t", "\\v", "\v", "\\\\", "\\")
	return &StringObject{Value: replacer.Replace(value), baseObj: &baseObj{class: vm.topLevelClass(stringClass)}}
}

func (vm *VM) initStringClass() *RClass {
	sc := vm.initializeClass(stringClass, false)
	sc.setBuiltInMethods(builtinStringInstanceMethods(), false)
	sc.setBuiltInMethods(builtInStringClassMethods(), true)
	return sc
}

func builtinStringInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{

		{
			// Returns the concatenation of self and another String
			//
			// ```Ruby
			// "first" + "-second" # => "first-second"
			// ```
			//
			// @return [String]
			Name: "+",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*StringObject).Value
					right, ok := args[0].(*StringObject)

					if !ok {
						return wrongTypeError(receiver.Class())
					}

					rightValue := right.Value
					return t.vm.initStringObject(leftValue + rightValue)
				}
			},
		},
		{
			// Returns self multiplying another Integer
			//
			// ```Ruby
			// "string " * 2 # => "string string string "
			// ```
			//
			// @return [String]
			Name: "*",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*StringObject).Value
					right, ok := args[0].(*IntegerObject)

					if !ok {
						return wrongTypeError(receiver.Class())
					}

					if right.Value < 0 {
						return newError("Second argument must be greater than or equal to 0 String#*")
					}

					var result string

					for i := 0; i < right.Value; i++ {
						result += leftValue
					}

					return t.vm.initStringObject(result)
				}
			},
		},
		{
			// Returns a Boolean if first string greater than second string
			//
			// ```Ruby
			// "a" < "b" # => true
			// ```
			//
			// @return [Boolean]
			Name: ">",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*StringObject).Value
					right, ok := args[0].(*StringObject)

					if !ok {
						return wrongTypeError(receiver.Class())
					}

					rightValue := right.Value

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
			// ```Ruby
			// "a" < "b" # => true
			// ```
			//
			// @return [Boolean]
			Name: "<",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*StringObject).Value
					right, ok := args[0].(*StringObject)

					if !ok {
						return wrongTypeError(receiver.Class())
					}

					rightValue := right.Value

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
			// ```Ruby
			// "first" == "second" # => false
			// "two" == "two" # => true
			// ```
			//
			// @return [Boolean]
			Name: "==",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*StringObject).Value
					right, ok := args[0].(*StringObject)

					if !ok {
						return wrongTypeError(receiver.Class())
					}

					rightValue := right.Value

					if leftValue == rightValue {
						return TRUE
					}

					return FALSE
				}
			},
		},
		{
			// Returns a Integer. If first string is less than second string returns -1, if equal to returns 0, if greater returns 1
			//
			//
			// ```Ruby
			// "abc" <=> "abcd" # => -1
			// "abc" <=> "abc" # => 0
			// "abcd" <=> "abc" # => 1
			// ```
			//
			// @return [Integer]
			Name: "<=>",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*StringObject).Value
					right, ok := args[0].(*StringObject)

					if !ok {
						return wrongTypeError(receiver.Class())
					}

					rightValue := right.Value

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
			// ```Ruby
			// "first" != "second" # => true
			// "two" != "two" # => false
			// ```
			//
			// @return [Boolean]
			Name: "!=",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					leftValue := receiver.(*StringObject).Value
					right, ok := args[0].(*StringObject)

					if !ok {
						return wrongTypeError(receiver.Class())
					}

					rightValue := right.Value

					if leftValue != rightValue {
						return TRUE
					}

					return FALSE
				}
			},
		},
		{
			// Return a new String with the first character converted to uppercase but the rest of string converted to lowercase.
			//
			// ```Ruby
			// "test".capitalize         # => "Test"
			// "tEST".capitalize         # => "Test"
			// "heLlo\nWoRLd".capitalize # => "Hello\nworld"
			// ```
			//
			// @return [String]
			Name: "capitalize",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					str := []byte(receiver.(*StringObject).Value)
					start := string(str[0])
					rest := string(str[1:])
					result := strings.ToUpper(start) + strings.ToLower(rest)

					return t.vm.initStringObject(result)
				}
			},
		},
		{
			// Returns a new String with all characters is upcase
			//
			// ```Ruby
			// "very big".upcase # => "VERY BIG"
			// ```
			//
			// @return [String]
			Name: "upcase",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					str := receiver.(*StringObject).Value

					return t.vm.initStringObject(strings.ToUpper(str))
				}
			},
		},
		{
			// Returns a new String with all characters is lowercase
			//
			// ```Ruby
			// "erROR".downcase        # => "error"
			// "HeLlO\tWorLD".downcase # => "hello\tworld"
			// ```
			//
			// @return [String]
			Name: "downcase",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					str := receiver.(*StringObject).Value

					return t.vm.initStringObject(strings.ToLower(str))
				}
			},
		},
		{
			// Returns the character length of self
			// **Note:** the length is currently byte-based, instead of charcode-based.
			//
			// ```Ruby
			// "zero".size # => 4
			// "".size # => 0
			// ```
			//
			// @return [Integer]
			Name: "size",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					str := receiver.(*StringObject).Value

					return t.vm.initIntegerObject(len(str))
				}
			},
		},
		{
			// Returns the character length of self
			// **Note:** the length is currently byte-based, instead of charcode-based.
			//
			// ```Ruby
			// "zero".size # => 4
			// "".size # => 0
			// ```
			//
			// @return [Integer]
			Name: "length",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					str := receiver.(*StringObject).Value

					return t.vm.initIntegerObject(len(str))
				}
			},
		},
		{
			// Returns a new String with reverse order of self
			// **Note:** the length is currently byte-based, instead of charcode-based.
			//
			// ```Ruby
			// "reverse".reverse      # => "esrever"
			// "Hello\nWorld".reverse # => "dlroW\nolleH"
			// ```
			//
			// @return [String]
			Name: "reverse",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					str := receiver.(*StringObject).Value

					var revert string
					for i := len(str) - 1; i >= 0; i-- {
						revert += string(str[i])
					}

					return t.vm.initStringObject(revert)
				}
			},
		},
		{
			// Returns a new String with self value
			//
			// ```Ruby
			// "string".to_s # => "string"
			// ```
			//
			// @return [String]
			Name: "to_s",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					str := receiver.(*StringObject).Value

					return t.vm.initStringObject(str)
				}
			},
		},
		{
			// Returns the result of converting self to Integer
			//
			// ```Ruby
			// "123".to_i # => 123
			// "3d print".to_i # => 3
			// "some text".to_i # => 0
			// ```
			//
			// @return [Integer]
			Name: "to_i",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					str := receiver.(*StringObject).Value
					parsedStr, err := strconv.ParseInt(str, 10, 0)

					if err == nil {
						return t.vm.initIntegerObject(int(parsedStr))
					}

					var digits string
					for _, char := range str {
						if unicode.IsDigit(char) {
							digits += string(char)
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
			// Checks if the specified string is included in the receiver
			//
			// ```Ruby
			// "Hello\nWorld".include("\n") # => true
			// ```
			//
			// @return [Bool]
			Name: "include",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {

					rcv := receiver.(*StringObject).Value
					arg := args[0].(*StringObject).Value

					if strings.Contains(rcv, arg) {
						return TRUE
					}

					return FALSE
				}
			},
		},
	}
}

func builtInStringClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					return t.UnsupportedMethodError("#new", receiver)
				}
			},
		},
	}
}

// Polymorphic helper functions -----------------------------------------

// toString just returns the value of string.
func (s *StringObject) toString() string {
	return s.Value
}

// toJSON converts the receiver into JSON string.
func (s *StringObject) toJSON() string {
	return "\"" + s.Value + "\""
}

func (s *StringObject) equal(e *StringObject) bool {
	return s.Value == e.Value
}
