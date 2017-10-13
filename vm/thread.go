package vm

import (
	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/vm/errors"
	"strings"
)

type thread struct {
	// a stack that holds call frames
	callFrameStack *callFrameStack
	// call frame pointer
	cfp int
	// data stack
	stack *stack
	// stack pointer
	sp int

	vm *VM
}

func (t *thread) isMainThread() bool {
	return t == t.vm.mainThread
}

func (t *thread) getBlock(name string, filename filename) *instructionSet {
	return t.vm.getBlock(name, filename)
}

func (t *thread) getMethodIS(name string, filename filename) (*instructionSet, bool) {
	return t.vm.getMethodIS(name, filename)
}

func (t *thread) getClassIS(name string, filename filename) *instructionSet {
	return t.vm.getClassIS(name, filename)
}

func (t *thread) startFromTopFrame() {
	cf := t.callFrameStack.top()
	t.evalCallFrame(cf)
}

func (t *thread) evalCallFrame(cf *callFrame) {
	for cf.pc < len(cf.instructionSet.instructions) {
		i := cf.instructionSet.instructions[cf.pc]
		t.execInstruction(cf, i)
		if _, yes := t.hasError(); yes {
			return
		}
	}
}

func (t *thread) hasError() (string, bool) {
	var hasError bool
	var msg string
	if t.stack.top() != nil {
		if err, ok := t.stack.top().Target.(*Error); ok {
			hasError = true
			msg = err.Message
		}
	}
	return msg, hasError
}

func (t *thread) execInstruction(cf *callFrame, i *instruction) {
	cf.pc++

	i.action.operation(t, cf, i.Params...)
}

func (t *thread) builtinMethodYield(blockFrame *callFrame, args ...Object) *Pointer {
	c := newCallFrame(blockFrame.instructionSet)
	c.blockFrame = blockFrame
	c.ep = blockFrame.ep
	c.self = blockFrame.self

	for i := 0; i < len(args); i++ {
		c.insertLCL(i, 0, args[i])
	}

	t.callFrameStack.push(c)
	t.startFromTopFrame()

	return t.stack.top()
}

func (t *thread) retrieveBlock(cf *callFrame, args []interface{}) (blockFrame *callFrame) {
	var blockName string
	var hasBlock bool

	blockFlag := args[2].(string)

	if len(blockFlag) != 0 {
		hasBlock = true
		blockName = strings.Split(blockFlag, ":")[1]
	}

	if hasBlock {
		block := t.getBlock(blockName, cf.instructionSet.filename)

		c := newCallFrame(block)
		c.isBlock = true
		c.ep = cf
		c.self = cf.self

		t.callFrameStack.push(c)

		blockFrame = c
	}

	return
}

func (t *thread) sendMethod(methodName string, argCount int, blockFrame *callFrame) {
	var method Object

	if arr, ok := t.stack.top().Target.(*ArrayObject); ok && arr.splat {
		// Pop array
		t.stack.pop()
		// Can't count array self, only the number of array elements
		argCount = argCount - 1 + len(arr.Elements)
		for _, elem := range arr.Elements {
			t.stack.push(&Pointer{Target: elem})
		}
	}

	argPr := t.sp - argCount
	receiverPr := argPr - 1
	receiver := t.stack.Data[receiverPr].Target

	/*
		Because send method adds additional object (method name) to the stack.
		So we need to move down real arguments like

		---------------
		Foo (*vm.RClass) 0
		bar (*vm.StringObject) 1
		5 (*vm.IntegerObject) 2
		---------------

		To

		-----------
		Foo (*vm.RClass) 0
		5 (*vm.IntegerObject) 1
		---------

		This also means we need to minus one on argument count and stack pointer
	*/
	for i := 0; i < argCount-1; i++ {
		t.stack.Data[argPr+i] = t.stack.Data[argPr+i+1]
	}

	argCount--
	t.sp--

	method = receiver.findMethod(methodName)

	if method == nil {
		err := t.vm.initErrorObject(errors.UndefinedMethodError, "Undefined Method '%+v' for %+v", methodName, receiver.toString())
		t.stack.set(receiverPr, &Pointer{Target: err})
		t.sp = argPr
		return
	}

	switch m := method.(type) {
	case *MethodObject:
		callObj := newCallObject(receiver, m, receiverPr, argCount, &bytecode.ArgSet{}, blockFrame)
		t.evalMethodObject(callObj)
	case *BuiltinMethodObject:
		t.evalBuiltinMethod(receiver, m, receiverPr, argCount, &bytecode.ArgSet{}, blockFrame)
	case *Error:
		t.returnError(errors.InternalError, m.toString())
	}
}

func (t *thread) evalBuiltinMethod(receiver Object, method *BuiltinMethodObject, receiverPr, argCount int, argSet *bytecode.ArgSet, blockFrame *callFrame) {
	methodBody := method.Fn(receiver)
	args := []Object{}
	argPr := receiverPr + 1

	for i := 0; i < argCount; i++ {
		args = append(args, t.stack.Data[argPr+i].Target)
	}

	evaluated := methodBody(t, args, blockFrame)

	_, ok := receiver.(*RClass)
	if method.Name == "new" && ok {
		instance, ok := evaluated.(*RObject)
		if ok && instance.InitializeMethod != nil {
			callObj := newCallObject(instance, instance.InitializeMethod, receiverPr, argCount, argSet, blockFrame)
			t.evalMethodObject(callObj)
		}
	}
	t.stack.set(receiverPr, &Pointer{Target: evaluated})
	t.sp = argPr
}

func (t *thread) reportArgumentError(idealArgNumber int, methodName string, exactArgNumber int, receiverPtr int) {
	var message string

	if idealArgNumber > exactArgNumber {
		message = "Expect at least %d args for method '%s'. got: %d"
	} else {
		message = "Expect at most %d args for method '%s'. got: %d"
	}

	e := t.vm.initErrorObject(errors.ArgumentError, message, idealArgNumber, methodName, exactArgNumber)
	t.stack.set(receiverPtr, &Pointer{Target: e})
	t.sp = receiverPtr + 1
}

func (t *thread) evalMethodObject(call *callObject) {
	c := newCallFrame(call.InstructionSet())
	c.self = call.receiver
	argPr := call.receiverPtr + 1
	minimumArgNumber := 0
	paramTypes := call.ParamTypes()
	paramsCount := len(call.ParamTypes())
	argTypesCount := len(paramTypes)

	for _, at := range paramTypes {
		if at == bytecode.NormalArg {
			minimumArgNumber++
		}
	}

	if call.argCount > paramsCount && !call.method.isSplatArgIncluded() {
		t.reportArgumentError(paramsCount, call.MethodName(), call.argCount, call.receiverPtr)
		return
	}

	if minimumArgNumber > call.argCount {
		t.reportArgumentError(minimumArgNumber, call.MethodName(), call.argCount, call.receiverPtr)
		return
	}

	// argIndex + argPr == current argument's position
	argIndex := 0

	// If given arguments is more than the normal arguments.
	// It might mean we have optioned argument been override.
	// Or we have some keyword arguments
	if minimumArgNumber < call.argCount {
		// This is only for normal/optioned arguments
		lastArgIndex := -1

		for paramIndex, paramType := range paramTypes {
			// Deal with normal arguments first
			if paramType == bytecode.NormalArg || paramType == bytecode.OptionedArg {
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
				for argIndex, at := range call.ArgTypes() {
					if lastArgIndex < argIndex && (at == bytecode.NormalArg || at == bytecode.OptionedArg) {
						c.insertLCL(paramIndex, 0, t.stack.Data[argPr+argIndex].Target)

						// Store latest index value (and compare them to current argument index)
						// This is to make sure we won't get same argument's index twice.
						lastArgIndex = argIndex
						break
					}
				}
			}

			if paramType == bytecode.RequiredKeywordArg || paramType == bytecode.OptionalKeywordArg {
				paramName := call.ParamNames()[paramIndex]
				argIndex := call.argSet.FindIndex(paramName)

				if argIndex != -1 {
					c.insertLCL(paramIndex, 0, t.stack.Data[argPr+argIndex].Target)
				} else if paramType == bytecode.RequiredKeywordArg {
					e := t.vm.initErrorObject(errors.ArgumentError, "Method %s requires key argument %s", call.MethodName(), paramName)
					t.stack.set(call.receiverPtr, &Pointer{Target: e})
					t.sp = argPr
					return
				}
			}

			// If argument index equals argument count means we already assigned all arguments
			if argIndex == call.argCount || paramType == bytecode.SplatArg {
				argIndex = paramIndex
				break
			}
		}
	} else {
		for i, paramType := range paramTypes {
			if paramType == bytecode.NormalArg {
				c.insertLCL(i, 0, t.stack.Data[argPr+argIndex].Target)
				argIndex++
			}
		}
	}

	if argTypesCount > 0 && call.method.isSplatArgIncluded() {
		elems := []Object{}
		for argIndex < call.argCount {
			elems = append(elems, t.stack.Data[argPr+argIndex].Target)
			argIndex++
		}

		c.insertLCL(paramsCount-1, 0, t.vm.initArrayObject(elems))
	}

	// TODO: Implement this
	if argTypesCount > 0 && call.method.isSplatArgIncluded() && call.method.isKeywordArgIncluded() {

	}

	c.blockFrame = call.blockFrame
	t.callFrameStack.push(c)
	t.startFromTopFrame()

	t.stack.set(call.receiverPtr, t.stack.top())
	t.sp = argPr
}

func (t *thread) returnError(errorType, format string, args ...interface{}) {
	err := t.vm.initErrorObject(errorType, format, args...)
	t.stack.push(&Pointer{Target: err})
}

func (t *thread) unsupportedMethodError(methodName string, receiver Object) *Error {
	return t.vm.initErrorObject(errors.UnsupportedMethodError, "Unsupported Method %s for %+v", methodName, receiver.toString())
}
