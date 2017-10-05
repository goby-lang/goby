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
		t.evalMethodObject(receiver, m, receiverPr, argCount, blockFrame)
	case *BuiltinMethodObject:
		t.evalBuiltinMethod(receiver, m, receiverPr, argCount, blockFrame)
	case *Error:
		t.returnError(errors.InternalError, m.toString())
	}
}

func (t *thread) evalBuiltinMethod(receiver Object, method *BuiltinMethodObject, receiverPr, argCount int, blockFrame *callFrame) {
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
			t.evalMethodObject(instance, instance.InitializeMethod, receiverPr, argCount, blockFrame)
		}
	}
	t.stack.set(receiverPr, &Pointer{Target: evaluated})
	t.sp = argPr
}

func (t *thread) evalMethodObject(receiver Object, method *MethodObject, receiverPr, argC int, blockFrame *callFrame) {
	c := newCallFrame(method.instructionSet)
	c.self = receiver
	argPr := receiverPr + 1
	minimumArgNumber := 0
	argTypesCount := len(method.argTypes())

	for _, at := range method.argTypes() {
		if at == bytecode.NormalArg || at == bytecode.RequiredKeywordArg {
			minimumArgNumber++
		}
	}

	if argC > method.argc && !method.isSplatArgIncluded() {
		e := t.vm.initErrorObject(errors.ArgumentError, "Expect at most %d args for method '%s'. got: %d", method.argc, method.Name, argC)
		t.stack.set(receiverPr, &Pointer{Target: e})
		t.sp = argPr
		return
	}

	if minimumArgNumber > argC {
		e := t.vm.initErrorObject(errors.ArgumentError, "Expect at least %d args for method '%s'. got: %d", minimumArgNumber, method.Name, argC)
		t.stack.set(receiverPr, &Pointer{Target: e})
		t.sp = argPr
		return
	}

	argIndex := 0

	for i, argType := range method.argTypes() {
		if argType == bytecode.NormalArg {
			c.insertLCL(i, 0, t.stack.Data[argPr+argIndex].Target)
			argIndex++
		}
	}

	/*
	 def foo(a = 10, b = 11, c); end

	 foo(1, 2)

	 In the above example, method foo's minimum argument number is 1 (`c`).
	 And the given argument number is 2.

	 So we first assign arguments to those doesn't have a default value (`c` in this example).
	 And then we assign the rest of given values from first parameter, so `a` would be assigned 2.

	 Result:

	 a == 2
	 b == 11
	 c == 1
	*/

	if minimumArgNumber < argC {
		// Fill arguments with default value from beginning
		for i, argType := range method.argTypes() {
			if argType == bytecode.OptionedArg {
				c.insertLCL(i, 0, t.stack.Data[argPr+argIndex].Target)
				argIndex++
			}

			if argType == bytecode.RequiredKeywordArg || argType == bytecode.OptionalKeywordArg {
				h := t.stack.Data[argPr+argIndex].Target.(*HashObject)

				for _, data := range h.Pairs {
					c.insertLCL(i, 0, data)
					argIndex++
				}
			}

			// If argument index equals argument count means we already assigned all arguments
			if argIndex == argC || argType == bytecode.SplatArg {
				break
			}
		}
	}

	if argTypesCount > 0 && method.isSplatArgIncluded() && !method.isKeywordArgIncluded() {
		elems := []Object{}
		for argIndex < argC {
			elems = append(elems, t.stack.Data[argPr+argIndex].Target)
			argIndex++
		}

		c.insertLCL(len(method.argTypes())-1, 0, t.vm.initArrayObject(elems))
	}

	// TODO: Implement this
	if argTypesCount > 0 && method.isSplatArgIncluded() && method.isKeywordArgIncluded() {

	}

	c.blockFrame = blockFrame
	t.callFrameStack.push(c)
	t.startFromTopFrame()

	t.stack.set(receiverPr, t.stack.top())
	t.sp = argPr
}

func (t *thread) returnError(errorType, format string, args ...interface{}) {
	err := t.vm.initErrorObject(errorType, format, args...)
	t.stack.push(&Pointer{Target: err})
}

func (t *thread) unsupportedMethodError(methodName string, receiver Object) *Error {
	return t.vm.initErrorObject(errors.UnsupportedMethodError, "Unsupported Method %s for %+v", methodName, receiver.toString())
}
