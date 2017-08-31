package vm

import (
	"github.com/goby-lang/goby/compiler/bytecode"
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

func (t *thread) builtInMethodYield(blockFrame *callFrame, args ...Object) *Pointer {
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

	if len(args) > 2 {
		hasBlock = true
		blockFlag := args[2].(string)
		blockName = strings.Split(blockFlag, ":")[1]
	} else {
		hasBlock = false
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

func (t *thread) evalBuiltInMethod(receiver Object, method *BuiltInMethodObject, receiverPr, argCount int, blockFrame *callFrame) {
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
		if at == bytecode.NormalArg {
			minimumArgNumber++
		}
	}

	if argC > method.argc && method.lastArgType() != bytecode.SplatArg {
		e := t.vm.initErrorObject(ArgumentError, "Expect at most %d args for method '%s'. got: %d", method.argc, method.Name, argC)
		t.stack.set(receiverPr, &Pointer{Target: e})
		t.sp = argPr
		return
	}

	if minimumArgNumber > argC {
		e := t.vm.initErrorObject(ArgumentError, "Expect at least %d args for method '%s'. got: %d", minimumArgNumber, method.Name, argC)
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
			if argType != bytecode.NormalArg && argType != bytecode.SplatArg {
				c.insertLCL(i, 0, t.stack.Data[argPr+argIndex].Target)
				argIndex++
			}

			// If argument index equals argument count means we already assigned all arguments
			if argIndex == argC || argType == bytecode.SplatArg {
				break
			}
		}
	}

	if argTypesCount > 0 && method.lastArgType() == bytecode.SplatArg {
		elems := []Object{}
		for argIndex < argC {
			elems = append(elems, t.stack.Data[argPr+argIndex].Target)
			argIndex++
		}

		c.insertLCL(len(method.argTypes())-1, 0, t.vm.initArrayObject(elems))
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
	return t.vm.initErrorObject(UnsupportedMethodError, "Unsupported Method %s for %+v", methodName, receiver.toString())
}
