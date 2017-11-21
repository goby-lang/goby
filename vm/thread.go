package vm

import (
	"fmt"
	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/vm/errors"
	"os"
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
	defer func() {
		if r := recover(); r != nil {
			t.reportErrorAndStop()
		}
	}()

	switch cf := cf.(type) {
	case *normalCallFrame:
		for cf.pc < cf.instructionsCount() {
			i := cf.instructionSet.instructions[cf.pc]
			t.execInstruction(cf, i)
		}
	case *goMethodCallFrame:
		args := []Object{}

		for _, obj := range cf.locals {
			if obj != nil {
				args = append(args, obj.Target)
			}
		}
		//fmt.Println("-----------------------")
		//fmt.Println(t.callFrameStack.inspect())
		result := cf.method(t, args, cf.blockFrame)
		t.stack.push(&Pointer{Target: result})
		//fmt.Println(t.callFrameStack.inspect())
		//fmt.Println("-----------------------")
		t.callFrameStack.pop()
	}

	t.removeUselessBlockFrame(cf)
}

/*
	Remove top frame if it's a block frame

	Block execution frame <- This was popped after callframe is executed
	---------------------
	Block frame           <- So this frame is useless
	---------------------
	Main frame
*/

func (t *thread) removeUselessBlockFrame(frame callFrame) {
	topFrame := t.callFrameStack.top()

	if topFrame != nil && topFrame.IsSourceBlock() {
		t.callFrameStack.pop().stopExecution()
	}
}

func (t *thread) reportErrorAndStop() {
	cf := t.callFrameStack.top()

	if cf != nil {
		cf.stopExecution()
	}

	top := t.stack.top().Target
	err := top.(*Error)

	if !err.storedTraces {
		for i := t.cfp - 1; i > 0; i-- {
			frame := t.callFrameStack.callFrames[i]

			if frame.IsBlock() {
				continue
			}

			msg := fmt.Sprintf("from %s:%d", frame.FileName(), frame.SourceLine())
			err.stackTraces = append(err.stackTraces, msg)
		}

		err.storedTraces = true
	}

	if t.vm.mode == REPLMode {
		return
	}

	panic(err)

	if t.vm.mode == NormalMode {
		if t.isMainThread() {
			os.Exit(1)
		}
	}
}

func (t *thread) execInstruction(cf *normalCallFrame, i *instruction) {
	cf.pc++

	//fmt.Println(t.callFrameStack.inspect())
	//fmt.Println(i.inspect())
	i.action.operation(t, i.sourceLine, cf, i.Params...)
	//fmt.Println("============================")
	//fmt.Println(t.callFrameStack.inspect())
}

func (t *thread) builtinMethodYield(blockFrame *normalCallFrame, args ...Object) *Pointer {
	c := newNormalCallFrame(blockFrame.instructionSet, blockFrame.FileName(), blockFrame.sourceLine)
	c.blockFrame = blockFrame
	c.ep = blockFrame.ep
	c.self = blockFrame.self
	c.sourceLine = blockFrame.SourceLine()
	c.isBlock = true

	for i := 0; i < len(args); i++ {
		c.insertLCL(i, 0, args[i])
	}

	t.callFrameStack.push(c)
	t.startFromTopFrame()

	return t.stack.top()
}

func (t *thread) retrieveBlock(fileName, blockFlag string, sourceLine int) (blockFrame *normalCallFrame) {
	var blockName string
	var hasBlock bool

	if len(blockFlag) != 0 {
		hasBlock = true
		blockName = strings.Split(blockFlag, ":")[1]
	}

	if hasBlock {
		block := t.getBlock(blockName, fileName)

		c := newNormalCallFrame(block, fileName, sourceLine)
		c.isSourceBlock = true
		c.isBlock = true
		blockFrame = c
	}

	return
}

func (t *thread) sendMethod(methodName string, argCount int, blockFrame *normalCallFrame, sourceLine int) {
	var method Object

	if arr, ok := t.stack.top().Target.(*ArrayObject); ok && arr.splat {
		// Pop array
		t.stack.pop()
		// Can't count array self, only the number of array elements
		argCount = argCount + len(arr.Elements)
		for _, elem := range arr.Elements {
			t.stack.push(&Pointer{Target: elem})
		}
	}

	argPr := t.sp - argCount - 1
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
	for i := 0; i < argCount; i++ {
		t.stack.Data[argPr+i] = t.stack.Data[argPr+i+1]
	}

	t.sp--

	method = receiver.findMethod(methodName)

	if method == nil {
		t.setErrorObject(receiverPr, argPr, errors.UndefinedMethodError, sourceLine, "Undefined Method '%+v' for %+v", methodName, receiver.toString())
	}

	sendCallFrame := t.callFrameStack.top()

	switch m := method.(type) {
	case *MethodObject:
		callObj := newCallObject(receiver, m, receiverPr, argCount, &bytecode.ArgSet{}, blockFrame, sendCallFrame.SourceLine())
		t.evalMethodObject(callObj, sourceLine)
	case *BuiltinMethodObject:
		t.evalBuiltinMethod(receiver, m, receiverPr, argCount, &bytecode.ArgSet{}, blockFrame, sourceLine, sendCallFrame.FileName())
	case *Error:
		t.pushErrorObject(errors.InternalError, sourceLine, m.toString())
	}
}

func (t *thread) evalBuiltinMethod(receiver Object, method *BuiltinMethodObject, receiverPtr, argCount int, argSet *bytecode.ArgSet, blockFrame *normalCallFrame, sourceLine int, fileName string) {
	cf := newGoMethodCallFrame(method.Fn(receiver, sourceLine), method.Name, fileName, sourceLine)
	cf.sourceLine = sourceLine
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
			callObj := newCallObject(instance, instance.InitializeMethod, receiverPtr, argCount, argSet, blockFrame, sourceLine)
			t.evalMethodObject(callObj, sourceLine)
		}
	}

	t.stack.set(receiverPtr, evaluated)
	t.sp = argPtr

	if err, ok := evaluated.Target.(*Error); ok {
		panic(err.Message())
	}
}

// TODO: Move instruction into call object
func (t *thread) evalMethodObject(call *callObject, sourceLine int) {
	normalParamsCount := call.normalParamsCount()
	paramTypes := call.paramTypes()
	paramsCount := len(call.paramTypes())
	stack := t.stack.Data

	if call.argCount > paramsCount && !call.method.isSplatArgIncluded() {
		t.reportArgumentError(sourceLine, paramsCount, call.methodName(), call.argCount, call.receiverPtr)
	}

	if normalParamsCount > call.argCount {
		t.reportArgumentError(sourceLine, normalParamsCount, call.methodName(), call.argCount, call.receiverPtr)
	}

	// Check if arguments include all the required keys before assign keyword arguments
	for paramIndex, paramType := range paramTypes {
		switch paramType {
		case bytecode.RequiredKeywordArg:
			paramName := call.paramNames()[paramIndex]
			if _, ok := call.hasKeywordArgument(paramName); !ok {
				t.setErrorObject(call.receiverPtr, call.argPtr(), errors.ArgumentError, sourceLine, "Method %s requires key argument %s", call.methodName(), paramName)
			}
		}
	}

	err := call.assignKeywordArguments(stack)

	if err != nil {
		t.setErrorObject(call.receiverPtr, call.argPtr(), errors.ArgumentError, sourceLine, err.Error())
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

func (t *thread) reportArgumentError(sourceLine, idealArgNumber int, methodName string, exactArgNumber int, receiverPtr int) {
	var message string

	if idealArgNumber > exactArgNumber {
		message = "Expect at least %d args for method '%s'. got: %d"
	} else {
		message = "Expect at most %d args for method '%s'. got: %d"
	}

	t.setErrorObject(receiverPtr, receiverPtr+1, errors.ArgumentError, sourceLine, message, idealArgNumber, methodName, exactArgNumber)
}

func (t *thread) pushErrorObject(errorType string, sourceLine int, format string, args ...interface{}) {
	err := t.vm.initErrorObject(errorType, sourceLine, format, args...)
	t.stack.push(&Pointer{Target: err})
	panic(err.Message())
}

func (t *thread) setErrorObject(receiverPtr, sp int, errorType string, sourceLine int, format string, args ...interface{}) {
	err := t.vm.initErrorObject(errorType, sourceLine, format, args...)
	t.stack.set(receiverPtr, &Pointer{Target: err})
	t.sp = sp
	panic(err.Message())
}
