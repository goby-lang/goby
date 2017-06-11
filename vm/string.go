package vm

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"sync"
)

var (
	stringClass *RString
)

// RString is the built in string class
type RString struct {
	*BaseClass
}

// StringObject represents string instances
type StringObject struct {
	Class *RString
	Value string
}

func (s *StringObject) Inspect() string {
	return s.Value
}

func (s *StringObject) returnClass() Class {
	if s.Class == nil {
		panic(fmt.Sprintf("String %s doesn't have class.", s.Inspect()))
	}

	return s.Class
}

func (s *StringObject) equal(e *StringObject) bool {
	return s.Value == e.Value
}

var (
	stringTable = make(map[string]*StringObject)
	mutex = &sync.Mutex{}
)

func initializeString(value string) *StringObject {
	addr, ok := stringTable[value]

	if !ok {
		s := &StringObject{Value: value, Class: stringClass}

		mutex.Lock()
		stringTable[value] = s
		mutex.Unlock()

		return s
	}

	return addr
}

var builtinStringInstanceMethods = []*BuiltInMethodObject{
	{
		// Returns the concatineted of self and another String
		//
		// ```Ruby
		// "first" + "-second" # => "first-second"
		// ```
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
					return initilaizeInteger(-1)
				}
				if leftValue > rightValue {
					return initilaizeInteger(1)
				}

				return initilaizeInteger(0)
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
		// @return [String]
		Name: "capitalize",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := []byte(receiver.(*StringObject).Value)
				start := string(str[0])
				rest := string(str[1:])
				result := strings.ToUpper(start) + strings.ToLower(rest)

				return initializeString(result)
			}
		},
	},
	{
		// Returns a new String with all characters is upcase
		//
		// ```Ruby
		// "very big".upcase # => "VERY BIG"
		// ```
		// @return [String]
		Name: "upcase",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				return initializeString(strings.ToUpper(str))
			}
		},
	},
	{
		// Returns a new String with all characters is lowercase
		//
		// ```Ruby
		// "erROR".downcase # => "error"
		// ```
		// @return [String]
		Name: "downcase",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				return initializeString(strings.ToLower(str))
			}
		},
	},
	{
		// Returns the character length of self
		//
		// ```Ruby
		// "zero".size # => 4
		// "".size # => 0
		// ```
		// @return [Integer]
		Name: "size",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				return initilaizeInteger(len(str))
			}
		},
	},
	{
		// Returns the character length of self
		//
		// ```Ruby
		// "zero".size # => 4
		// "".size # => 0
		// ```
		// @return [Integer]
		Name: "length",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				return initilaizeInteger(len(str))
			}
		},
	},
	{
		// Returns a new String with reverse order of self
		//
		// ```Ruby
		// "reverse".reverse # => "esrever"
		// ```
		// @return [String]
		Name: "reverse",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value
				var revert string
				for i := len(str) - 1; i >= 0; i-- {
					revert += string(str[i])
				}

				return initializeString(revert)
			}
		},
	},
	{
		// Returns a new String with self value
		//
		// ```Ruby
		// "string".to_s # => "string"
		// ```
		// @return [String]
		Name: "to_s",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value

				return initializeString(str)
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
		// @return [Integer]
		Name: "to_i",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {

				str := receiver.(*StringObject).Value
				parsedStr, err := strconv.ParseInt(str, 10, 0)

				if err == nil {
					return initilaizeInteger(int(parsedStr))
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
					return initilaizeInteger(int(parsedStr))
				}

				return initilaizeInteger(0)
			}
		},
	},
}

func initString() {
	bc := &BaseClass{Name: "String", Methods: newEnvironment(), ClassMethods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass, superClass: objectClass}
	sc := &RString{BaseClass: bc}
	sc.setBuiltInMethods(builtinStringInstanceMethods, false)
	stringClass = sc
}
