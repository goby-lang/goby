package vm

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var (
	stringClass *RString
)

// RString is the built in string class
type RString struct {
	*BaseClass
}

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
// Currently, manipulations are based upon Golang's Unicode manipulations.
//
// - Currently, UTF-8 encoding is assumed based upon Golang's string manipulation, but the encoding is not actually specified(TBD).
// - `String.new` is not supported.
type StringObject struct {
	Class *RString
	Value string
}

func (s *StringObject) toString() string {
	return "\"" + s.Value + "\""
}

func (s *StringObject) toJSON() string {
	return "\"" + s.Value + "\""
}

func (s *StringObject) returnClass() Class {
	if s.Class == nil {
		panic(fmt.Sprintf("String %s doesn't have class.", s.toString()))
	}

	return s.Class
}

func (s *StringObject) equal(e *StringObject) bool {
	return s.Value == e.Value
}

func initStringObject(value string) *StringObject {
	return &StringObject{Value: value, Class: stringClass}
}

var builtinStringInstanceMethods = []*BuiltInMethodObject{
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
					return wrongTypeError(stringClass)
				}

				rightValue := right.Value
				return &StringObject{Value: leftValue + rightValue, Class: stringClass}
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
					return wrongTypeError(stringClass)
				}

				if right.Value < 0 {
					return newError("Second argument must be greater than or equal to 0 String#*")
				}

				var result string

				for i := 0; i < right.Value; i++ {
					result += leftValue
				}

				return &StringObject{Value: result, Class: stringClass}
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
					return wrongTypeError(stringClass)
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
					return wrongTypeError(stringClass)
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
					return wrongTypeError(stringClass)
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
					return wrongTypeError(stringClass)
				}

				rightValue := right.Value

				if leftValue < rightValue {
					return initIntegerObject(-1)
				}
				if leftValue > rightValue {
					return initIntegerObject(1)
				}

				return initIntegerObject(0)
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
					return wrongTypeError(stringClass)
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
		// "test".capitalize # => "Test"
		// "tEST".capitalize # => "Test"
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

				return initStringObject(result)
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

				return initStringObject(strings.ToUpper(str))
			}
		},
	},
	{
		// Returns a new String with all characters is lowercase
		//
		// ```Ruby
		// "erROR".downcase # => "error"
		// ```
		//
		// @return [String]
		Name: "downcase",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				return initStringObject(strings.ToLower(str))
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

				return initIntegerObject(len(str))
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

				return initIntegerObject(len(str))
			}
		},
	},
	{
		// Returns a new String with reverse order of self
		// **Note:** the length is currently byte-based, instead of charcode-based.
		//
		// ```Ruby
		// "reverse".reverse # => "esrever"
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

				return initStringObject(revert)
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

				return initStringObject(str)
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
					return initIntegerObject(int(parsedStr))
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
					return initIntegerObject(int(parsedStr))
				}

				return initIntegerObject(0)
			}
		},
	},
}

func initStringClass() {
	bc := &BaseClass{Name: "String", Methods: newEnvironment(), ClassMethods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass, superClass: objectClass}
	sc := &RString{BaseClass: bc}
	sc.setBuiltInMethods(builtinStringInstanceMethods, false)
	stringClass = sc
}
