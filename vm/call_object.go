package vm

import "github.com/goby-lang/goby/compiler/bytecode"

type callObject struct {
	method       *MethodObject
	receiverPtr  int
	argCount     int
	argSet       *bytecode.ArgSet
	argIndex     int
	lastArgIndex int
	callFrame    *callFrame
}

func newCallObject(receiver Object, method *MethodObject, receiverPtr, argCount int, argSet *bytecode.ArgSet, blockFrame *callFrame) *callObject {
	cf := newCallFrame(method.instructionSet)
	cf.self = receiver
	cf.blockFrame = blockFrame

	return &callObject{
		method:      method,
		receiverPtr: receiverPtr,
		argCount:    argCount,
		argSet:      argSet,
		// This is only for normal/optioned arguments
		lastArgIndex: -1,
		callFrame:    cf,
	}
}

func (co *callObject) instructionSet() *instructionSet {
	return co.method.instructionSet
}

func (co *callObject) paramTypes() []int {
	return co.instructionSet().paramTypes.Types()
}

func (co *callObject) paramNames() []string {
	return co.instructionSet().paramTypes.Names()
}

func (co *callObject) methodName() string {
	return co.method.Name
}

func (co *callObject) argTypes() []int {
	return co.argSet.Types()
}

func (co *callObject) argPtr() int {
	return co.receiverPtr + 1
}

func (co *callObject) argPosition() int {
	return co.argPtr() + co.argIndex
}

func (co *callObject) assignNormalArguments(stack []*Pointer) {
	for i, paramType := range co.paramTypes() {
		if paramType == bytecode.NormalArg {
			co.callFrame.insertLCL(i, 0, stack[co.argPosition()].Target)
			co.argIndex++
		}
	}
}

func (co *callObject) assignNormalAndOptionedArguments(paramIndex int, stack []*Pointer) {
	/*
		Find first usable value as normal argument, for example:

		```ruby
		  def foo(x, y:); end

		  foo(y: 100, 10)
		```

		In the example we can see that 'x' is the first parameter,
		but in the method call it's the second argument.

		This loop is for skipping other types of arguments and get the correct argument index.
	*/
	for argIndex, at := range co.argTypes() {
		if co.lastArgIndex < argIndex && (at == bytecode.NormalArg || at == bytecode.OptionedArg) {
			co.callFrame.insertLCL(paramIndex, 0, stack[co.argPtr()+argIndex].Target)

			// Store latest index value (and compare them to current argument index)
			// This is to make sure we won't get same argument's index twice.
			co.lastArgIndex = argIndex
			break
		}
	}
}

func (co *callObject) assignKeywordArguments(paramIndex int, stack []*Pointer) (paramName string, success bool) {
	paramName = co.paramNames()[paramIndex]
	argIndex := co.argSet.FindIndex(paramName)

	if argIndex != -1 {
		co.callFrame.insertLCL(paramIndex, 0, stack[co.argPtr()+argIndex].Target)
		success = true
	}

	return
}

func (co *callObject) assignSplatArgument(stack []*Pointer, arr *ArrayObject) {
	index := len(co.paramTypes()) - 1

	for co.argIndex < co.argCount {
		arr.Elements = append(arr.Elements, stack[co.argPosition()].Target)
		co.argIndex++
	}

	co.callFrame.insertLCL(index, 0, arr)
}

func (co *callObject) minimumArgNumber() (n int) {
	for _, at := range co.paramTypes() {
		if at == bytecode.NormalArg {
			n++
		}
	}

	return
}
