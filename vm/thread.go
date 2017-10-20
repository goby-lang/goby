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

func (t *thread) evalCallFrame(cf callFrame) {
	switch cf := cf.(type) {
	case *normalCallFrame:
		for cf.pc < len(cf.instructionSet.instructions) {
			i := cf.instructionSet.instructions[cf.pc]
			t.execInstruction(cf, i)
			if _, yes := t.hasError(); yes {
				return
			}
		}
	case *goMethodCallFrame:
		args := []Object{}

		for _, obj := range cf.locals {
			if obj != nil {
				args = append(args, obj.Target)
			}
		}
		result := cf.method(t, args, cf.blockFrame)
		t.stack.push(&Pointer{Target: result})
		t.callFrameStack.pop()
	}
}

func (t *thread) hasError() (string, bool) {
	var hasError bool
	var msg string
	if t.stack.top() != nil {
		top := t.stack.top().Target
		err, ok := top.(*Error)
		if ok {
			hasError = true
			msg = err.Message
		}
	}
	return msg, hasError
}

func (t *thread) execInstruction(cf *normalCallFrame, i *instruction) {
	cf.pc++

	i.action.operation(t, i, cf, i.Params...)
}

func (t *thread) builtinMethodYield(blockFrame *normalCallFrame, args ...Object) *Pointer {
	c := newNormalCallFrame(blockFrame.instructionSet, blockFrame.FileName())
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

func (t *thread) retrieveBlock(cf *normalCallFrame, args []interface{}) (blockFrame *normalCallFrame) {
	var blockName string
	var hasBlock bool

	blockFlag := args[2].(string)

	if len(blockFlag) != 0 {
		hasBlock = true
		blockName = strings.Split(blockFlag, ":")[1]
	}

	if hasBlock {
		block := t.getBlock(blockName, cf.instructionSet.filename)

		c := newNormalCallFrame(block, cf.instructionSet.filename)
		c.isBlock = true
		c.ep = cf
		c.self = cf.self

		t.callFrameStack.push(c)

		blockFrame = c
	}

	return
}

func (t *thread) sendMethod(methodName string, argCount int, blockFrame *normalCallFrame, instruction *instruction) {
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
		err := t.vm.initErrorObject(errors.UndefinedMethodError, instruction, "Undefined Method '%+v' for %+v", methodName, receiver.toString())
		t.stack.set(receiverPr, &Pointer{Target: err})
		t.sp = argPr
		return
	}

	sendCallFrame := t.callFrameStack.top()

	switch m := method.(type) {
	case *MethodObject:
		callObj := newCallObject(receiver, m, receiverPr, argCount, &bytecode.ArgSet{}, blockFrame, sendCallFrame.SourceLine(), sendCallFrame.FileName())
		t.evalMethodObject(callObj, instruction)
	case *BuiltinMethodObject:
		t.evalBuiltinMethod(receiver, m, receiverPr, argCount, &bytecode.ArgSet{}, blockFrame, instruction, sendCallFrame.FileName())
	case *Error:
		t.pushErrorObject(errors.InternalError, instruction, m.toString())
	}
}

func (t *thread) evalBuiltinMethod(receiver Object, method *BuiltinMethodObject, receiverPtr, argCount int, argSet *bytecode.ArgSet, blockFrame *normalCallFrame, instruction *instruction, fileName string) {
	cf := newGoMethodCallFrame(method.Fn(receiver, instruction), method.Name, fileName)
	cf.sourceLine = instruction.sourceLine
	cf.blockFrame = blockFrame
	argPtr := receiverPtr + 1

	for i := 0; i < argCount; i++ {
		cf.locals = append(cf.locals, t.stack.Data[argPtr+i])
	}

	t.callFrameStack.push(cf)
	t.startFromTopFrame()
	evaluated := t.stack.top()

	_, ok := receiver.(*RClass)
	if method.Name == "new" && ok {
		instance, ok := evaluated.Target.(*RObject)
		if ok && instance.InitializeMethod != nil {
			callObj := newCallObject(instance, instance.InitializeMethod, receiverPtr, argCount, argSet, blockFrame, instruction.sourceLine, fileName)
			t.evalMethodObject(callObj, instruction)
		}
	}

	t.stack.set(receiverPtr, evaluated)
	t.sp = argPtr
}

func (t *thread) reportArgumentError(instruction *instruction, idealArgNumber int, methodName string, exactArgNumber int, receiverPtr int) {
	var message string

	if idealArgNumber > exactArgNumber {
		message = "Expect at least %d args for method '%s'. got: %d"
	} else {
		message = "Expect at most %d args for method '%s'. got: %d"
	}

	e := t.vm.initErrorObject(errors.ArgumentError, instruction, message, idealArgNumber, methodName, exactArgNumber)
	t.stack.set(receiverPtr, &Pointer{Target: e})
	t.sp = receiverPtr + 1
}

// TODO: Move instruction into call object
func (t *thread) evalMethodObject(call *callObject, instruction *instruction) {
	normalParamsCount := call.normalParamsCount()
	paramTypes := call.paramTypes()
	paramsCount := len(call.paramTypes())
	stack := t.stack.Data

	if call.argCount > paramsCount && !call.method.isSplatArgIncluded() {
		t.reportArgumentError(instruction, paramsCount, call.methodName(), call.argCount, call.receiverPtr)
		return
	}

	if normalParamsCount > call.argCount {
		t.reportArgumentError(instruction, normalParamsCount, call.methodName(), call.argCount, call.receiverPtr)
		return
	}

	// Check if arguments include all the required keys before assign keyword arguments
	for paramIndex, paramType := range paramTypes {
		switch paramType {
		case bytecode.RequiredKeywordArg:
			paramName := call.paramNames()[paramIndex]
			if _, ok := call.hasKeywordArgument(paramName); !ok {
				e := t.vm.initErrorObject(errors.ArgumentError, instruction, "Method %s requires key argument %s", call.methodName(), paramName)
				t.stack.set(call.receiverPtr, &Pointer{Target: e})
				t.sp = call.argPtr()
				return
			}
		}
	}

	err := call.assignKeywordArguments(stack)

	if err != nil {
		e := t.vm.initErrorObject(errors.ArgumentError, instruction, err.Error())
		t.stack.set(call.receiverPtr, &Pointer{Target: e})
		t.sp = call.argPtr()
		return
	}

	// If given arguments is more than the normal arguments.
	// It might mean we have optioned argument been override.
	// Or we have some keyword arguments
	if normalParamsCount < call.argCount {
		for paramIndex, paramType := range paramTypes {
			switch paramType {
			case bytecode.NormalArg, bytecode.OptionedArg:
				call.assignNormalAndOptionedArguments(paramIndex, stack)
			case bytecode.SplatArg:
				call.argIndex = paramIndex
				call.assignSplatArgument(stack, t.vm.initArrayObject([]Object{}))
			}
		}
	} else {
		call.assignNormalArguments(stack)
	}

	t.callFrameStack.push(call.callFrame)
	t.startFromTopFrame()

	t.stack.set(call.receiverPtr, t.stack.top())
	t.sp = call.argPtr()
}

func (t *thread) pushErrorObject(errorType string, instruction *instruction, format string, args ...interface{}) {
	err := t.vm.initErrorObject(errorType, instruction, format, args...)
	t.stack.push(&Pointer{Target: err})
}

func (t *thread) initUnsupportedMethodError(instruction *instruction, methodName string, receiver Object) *Error {
	return t.vm.initErrorObject(errors.UnsupportedMethodError, instruction, "Unsupported Method %s for %+v", methodName, receiver.toString())
}
