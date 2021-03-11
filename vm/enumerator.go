package vm

import (
	"fmt"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

// EnumeratorObject represents a generic enumerable object..
type EnumeratorObject struct {
	*BaseObj
	EnumerableObject interface{}
	DefaultMethod    string
	Block            *BlockObject
}

// Class methods --------------------------------------------------------
var builtinEnumeratorClassMethods = []*BuiltinMethodObject{
	{
		Name: "new",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			return t.vm.InitNoMethodError(sourceLine, "new", receiver)
		},
	},
}

// Instance methods -----------------------------------------------------
var builtinEnumeratorInstanceMethods = []*BuiltinMethodObject{
	{
		// Loops through each element in the array, with the given block.
		// Returns self.
		// A block literal is required.
		//
		// ```ruby
		// a = ["a", "b", "c"]
		//
		// b = a.each do |e|
		//   puts(e + e)
		// end
		// #=> "aa"
		// #=> "bb"
		// #=> "cc"
		// puts b
		// #=> ["a", "b", "c"]
		// ```
		//
		// @param block literal
		// @return [Array]
		Name: "each",
		Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
			if len(args) != 0 {
				return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
			}

			if blockFrame == nil {
				return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
			}

			if blockIsEmpty(blockFrame) {
				return receiver.(*EnumeratorObject)
			}

			// just delegates to "each" on Array or Hash
			t.sendMethod("each", len(args)-1, blockFrame, sourceLine)

			return t.Stack.top().Target
		},
	},
	//// Works like #each, but passes the index of the element instead of the element itself.
	//// Returns self.
	//// A block literal is required.
	////
	//// ```ruby
	//// a = [:apple, :orange, :grape, :melon]
	////
	//// b = a.each_index do |i|
	////   puts(i*i)
	//// end
	//// #=> 0
	//// #=> 1
	//// #=> 4
	//// #=> 9
	//// puts b
	//// #=> ["a", "b", "c"]
	//// ```
	////
	//// @param block literal
	//// @return [Array]
	//{
	//	Name: "each_index",
	//	Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
	//		if len(args) != 0 {
	//			return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 0, len(args))
	//		}
	//
	//		if blockFrame == nil {
	//			return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
	//		}
	//
	//		arr := receiver.(*EnumeratorObject)
	//		if blockIsEmpty(blockFrame) {
	//			return arr
	//		}
	//
	//		// If it's an empty array, pop the block's call frame
	//		if len(arr.Elements) == 0 {
	//			t.callFrameStack.pop()
	//		}
	//
	//		for i := range arr.Elements {
	//			t.builtinMethodYield(blockFrame, t.vm.InitIntegerObject(i))
	//		}
	//		return arr
	//
	//	},
	//},
	//{
	//	// Returns a new hash from the element of the receiver (array) as keys, and generates respective values of hash from the keys by using the block provided.
	//	// The method can take a default value, and a block is required.
	//	// `index_with` is equivalent to `receiver.map do |e| e, e._do_something end.to_h`
	//	// Ref: https://github.com/rails/rails/pull/32523
	//	//
	//	// ```ruby
	//	// ary = [:Mon, :Tue, :Wed, :Thu, :Fri, :Sat, :Sun]
	//	// ary.index_with("weekday") do |d|
	//	//   if d == :Sat || d == :Sun
	//	//     "off day"
	//	//   end
	//	// end
	//	// #=> {Mon: "weekday",
	//	//      Tue: "weekday"
	//	//      Wed: "weekday"
	//	//      Thu: "weekday"
	//	//      Fri: "weekday"
	//	//      Sat: "off day"
	//	//      Sun: "off day"
	//	// }
	//	// ```
	//	//
	//	// @param optional default value [Object], block
	//	// @return [Hash]
	//	Name: "index_with",
	//	Fn: func(receiver Object, sourceLine int, t *Thread, args []Object, blockFrame *normalCallFrame) Object {
	//		if blockFrame == nil {
	//			return t.vm.InitErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
	//		}
	//
	//		a := receiver.(*EnumeratorObject)
	//		// If it's an empty array, pop the block's call frame
	//		if len(a.Elements) == 0 {
	//			t.callFrameStack.pop()
	//		}
	//
	//		hash := make(map[string]Object)
	//		switch len(args) {
	//		case 0:
	//			for _, obj := range a.Elements {
	//				hash[obj.ToString()] = t.builtinMethodYield(blockFrame, obj).Target
	//			}
	//		case 1:
	//			arg := args[0]
	//			for _, obj := range a.Elements {
	//				switch b := t.builtinMethodYield(blockFrame, obj).Target; b.(type) {
	//				case *NullObject:
	//					hash[obj.ToString()] = arg
	//				default:
	//					hash[obj.ToString()] = b
	//				}
	//			}
	//		default:
	//			return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentLess, 1, len(args))
	//		}
	//
	//		return t.vm.InitHashObject(hash)
	//
	//	},
	//},
}

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

// InitEnumeratorObject returns a new object with the given elemnts
func (vm *VM) InitEnumeratorObject(e Object, m string, b *BlockObject) *EnumeratorObject {
	if m == "" {
		m = "each"
	}
	return &EnumeratorObject{
		BaseObj:          NewBaseObject(vm.TopLevelClass(classes.EnumeratorClass)),
		EnumerableObject: e,
		DefaultMethod:    m,
		Block:            b,
	}
}

func (vm *VM) initEnumeratorClass() *RClass {
	ec := vm.initializeClass(classes.EnumeratorClass)
	ec.setBuiltinMethods(builtinEnumeratorInstanceMethods, false)
	ec.setBuiltinMethods(builtinEnumeratorClassMethods, true)
	return ec
}

// Polymorphic helper functions -----------------------------------------

// Value returns the elements from the object
func (e *EnumeratorObject) Value() interface{} {
	return e.EnumerableObject
}

// ToString returns the object's elements as the string format
func (e *EnumeratorObject) ToString() string {
	var eo string
	switch e.EnumerableObject.(type) {
	case *ArrayObject:
		eo = e.EnumerableObject.(*ArrayObject).Inspect()
	case *HashObject:
		eo = e.EnumerableObject.(*HashObject).Inspect()
	default:
		fmt.Println("debug: enum? ", eo)
	}
	return fmt.Sprintf("<Enumerator: %s:%s>", eo, e.DefaultMethod)
}

// Inspect delegates to ToString
func (e *EnumeratorObject) Inspect() string {
	return e.ToString()
}

// ToJSON returns the object's elements as the JSON string format
func (e *EnumeratorObject) ToJSON(t *Thread) string {
	var out string
	switch e.EnumerableObject.(type) {
	case *ArrayObject:
		out = e.EnumerableObject.(*ArrayObject).ToJSON(t)
	case *HashObject:
		out = e.EnumerableObject.(*HashObject).ToJSON(t)
	default:
		fmt.Println("debug: enum? ", e)
	}
	return out
}

//// Retrieves an object in an array using Integer index; common to `[]` and `at()`.
//func (e *EnumeratorObject) index(t *Thread, args []Object, sourceLine int) Object {
//	aLen := len(args)
//	if aLen < 1 || aLen > 2 {
//		return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgumentRange, 1, 2, aLen)
//	}
//
//	i := args[0]
//	index, ok := i.(*IntegerObject)
//	arrLength := e.Len()
//
//	if !ok {
//		return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[0].Class().Name)
//	}
//
//	if index.value < 0 && index.value < -arrLength {
//		return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.TooSmallIndexValue, index.value, -arrLength)
//	}
//
//	/* Validation for the second argument if exists */
//	if aLen == 2 {
//		j := args[1]
//		count, ok := j.(*IntegerObject)
//
//		if !ok {
//			return t.vm.InitErrorObject(errors.TypeError, sourceLine, errors.WrongArgumentTypeFormat, classes.IntegerClass, args[1].Class().Name)
//		}
//		if count.value < 0 {
//			return t.vm.InitErrorObject(errors.ArgumentError, sourceLine, errors.NegativeSecondValue, count.value)
//		}
//
//		/*
//		 *  This condition meets the special case (Don't know why ~ ? Ask Ruby or try it on irb!):
//		 *
//		 *  a = [1, 2, 3, 4, 5]
//		 *  a[5, 5] #=> []
//		 */
//		if index.value > 0 && index.value == arrLength {
//			return t.vm.InitArrayObject([]Object{})
//		}
//	}
//
//	/* Start Indexing */
//	normalizedIndex := e.normalizeIndex(index)
//	if normalizedIndex == -1 {
//		return NULL
//	}
//
//	if aLen == 2 {
//		j := args[1]
//		count, _ := j.(*IntegerObject)
//
//		if normalizedIndex+count.value > arrLength {
//			return t.vm.InitArrayObject(e.Elements[normalizedIndex:])
//		}
//		return t.vm.InitArrayObject(e.Elements[normalizedIndex : normalizedIndex+count.value])
//	}
//
//	return e.Elements[normalizedIndex]
//}
//
// Len returns the length of array's elements
func (e *EnumeratorObject) Len() int {
	switch e.EnumerableObject.(type) {
	case *ArrayObject:
		return len(e.EnumerableObject.(*ArrayObject).Elements)
	case *HashObject:
		return len(e.EnumerableObject.(*HashObject).Pairs)
	default:
		fmt.Println("debug: enum? ", e)
		return -1
	}
}

//
//// Swap is one of the required method to fulfill sortable interface
//func (e *EnumeratorObject) Swap(i, j int) {
//	e.Elements[i], e.Elements[j] = e.Elements[j], e.Elements[i]
//}
//
//// Less is one of the required method to fulfill sortable interface
//func (e *EnumeratorObject) Less(i, j int) bool {
//	leftObj, rightObj := e.Elements[i], e.Elements[j]
//	switch leftObj := leftObj.(type) {
//	case Numeric:
//		return leftObj.lessThan(rightObj)
//	case *StringObject:
//		right, ok := rightObj.(*StringObject)
//
//		if ok {
//			return leftObj.value < right.value
//		}
//
//		return false
//	default:
//		return false
//	}
//}
//
//// normalizes the index to the Ruby-style:
////
//// 1. if the index is between o and the index length, returns the index
//// 2. if it's a negative value (within bounds), returns the normalized positive version
//// 3. if it's out of bounds (either positive or negative), returns -1
//func (e *EnumeratorObject) normalizeIndex(objectIndex *IntegerObject) int {
//	aLength := len(e.Elements)
//	index := objectIndex.value
//
//	// out of bounds
//
//	if index >= aLength {
//		return -1
//	}
//
//	if index < 0 && -index > aLength {
//		return -1
//	}
//
//	// within bounds
//
//	if index < 0 {
//		return aLength + index
//	}
//
//	return index
//}
//
//// pop removes the last element in the array and returns it
//func (e *EnumeratorObject) pop() Object {
//	if len(e.Elements) < 1 {
//		return NULL
//	}
//
//	value := e.Elements[len(e.Elements)-1]
//	e.Elements = e.Elements[:len(e.Elements)-1]
//	return value
//}
//
//// push appends given object into array and returns the array object
//func (e *EnumeratorObject) push(objs []Object) *EnumeratorObject {
//	e.Elements = append(e.Elements, objs...)
//	return e
//}
//
//// shift removes the first element in the array and returns it
//func (e *EnumeratorObject) shift() Object {
//	if len(e.Elements) < 1 {
//		return NULL
//	}
//
//	value := e.Elements[0]
//	e.Elements = e.Elements[1:]
//	return value
//}
//
//func (e *EnumeratorObject) equalTo(compared Object) bool {
//	c, ok := compared.(*EnumeratorObject)
//
//	if !ok {
//		return false
//	}
//
//	if len(e.Elements) != len(c.Elements) {
//		return false
//	}
//
//	for i, el := range e.Elements {
//		if !el.equalTo(c.Elements[i]) {
//			return false
//		}
//	}
//
//	return true
//}
//
//// unshift inserts an element in the first position of the array
//func (e *EnumeratorObject) unshift(objs []Object) *EnumeratorObject {
//	e.Elements = append(objs, e.Elements...)
//	return e
//}
