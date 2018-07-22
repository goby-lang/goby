package vm

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/goby-lang/goby/compiler"
	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/compiler/parser"
	"github.com/goby-lang/goby/vm/errors"
)

const mainThreadID = 0

// Thread is the context needed for a single thread of execution
type Thread struct {
	// a stack that holds call frames
	callFrameStack callFrameStack

	// The acall frame currently being executed
	currentFrame callFrame

	// data Stack
	Stack Stack

	// theads have an id so they can be looked up in the vm. The main thread is always 0
	id int64

	vm *VM
}

func (t *Thread) VM() *VM {
	return t.vm
}

func (t *Thread) isMainThread() bool {
	return t.id == mainThreadID
}

func (t *Thread) getBlock(name string, filename filename) *instructionSet {
	// The "name" here is actually an index of block
	// for example <Block:1>'s name is "1"
	is, ok := t.vm.blockTables[filename][name]

	if !ok {
		panic(fmt.Sprintf("Can't find block %s", name))
	}

	return is
}

func (t *Thread) getMethodIS(name string, filename filename) (*instructionSet, bool) {
	iss, ok := t.vm.isTables[bytecode.MethodDef][name]

	if !ok {
		return nil, false
	}

	is := iss[t.vm.methodISIndexTables[filename].Data[name]]

	t.vm.methodISIndexTables[filename].Data[name]++

	return is, ok
}

func (t *Thread) getClassIS(name string, filename filename) *instructionSet {
	iss, ok := t.vm.isTables[bytecode.ClassDef][name]

	if !ok {
		panic(fmt.Sprintf("Can't find class %s's instructions", name))
	}

	is := iss[t.vm.classISIndexTables[filename].Data[name]]

	t.vm.classISIndexTables[filename].Data[name]++

	return is
}

func (t *Thread) execGobyLib(libName string) (err error) {
	libPath := filepath.Join(t.vm.libPath, libName)
	err = t.execFile(libPath)
	return
}

func (t *Thread) execFile(fpath string) (err error) {
	file, err := ioutil.ReadFile(fpath)

	if err != nil {
		return
	}

	instructionSets, err := compiler.CompileToInstructions(string(file), parser.NormalMode)

	if err != nil {
		return
	}

	oldMethodTable := isTable{}
	oldClassTable := isTable{}

	// Copy current file's instruction sets.
	for name, is := range t.vm.isTables[bytecode.MethodDef] {
		oldMethodTable[name] = is
	}

	for name, is := range t.vm.isTables[bytecode.ClassDef] {
		oldClassTable[name] = is
	}

	// This creates new execution environments for required file, including new instruction set table.
	// So we need to copy old instruction sets and restore them later, otherwise current program's instruction set would be overwrite.
	t.vm.ExecInstructions(instructionSets, fpath)

	// Restore instruction sets.
	t.vm.isTables[bytecode.MethodDef] = oldMethodTable
	t.vm.isTables[bytecode.ClassDef] = oldClassTable
	return
}

func (t *Thread) captureAndHandlePanic() {
	switch e := recover().(type) {
	case *Error:
		if t.vm.mode == NormalMode {
			fmt.Println(e.Message())
			if t.isMainThread() {
				os.Exit(1)
			}
		} else {
			log.Println(e)
			panic(e)
		}
	case error:
		log.Println(e)
		fmt.Println(e.Error())
	case nil:
		return
	default:
		log.Println(e)
		panic(e)
	}
}

func (t *Thread) startFromTopFrame() {
	defer t.captureAndHandlePanic()
	cf := t.callFrameStack.top()
	t.evalCallFrame(cf)
}

func (t *Thread) evalCallFrame(cf callFrame) {
	defer func() {
		if r := recover(); r != nil {
			t.reportErrorAndStop()
		}
	}()

	t.currentFrame = cf

	switch cf := cf.(type) {
	case *normalCallFrame:
		for cf.pc < cf.instructionsCount() {
			i := cf.instructionSet.instructions[cf.pc]
			t.execInstruction(cf, i)
		}
	case *goMethodCallFrame:
		args := []Object{}

		for i := 0; i < cf.argCount; i++ {
			args = append(args, t.Stack.data[cf.argPtr+i].Target)
		}
		//fmt.Println("-----------------------")
		//fmt.Println(t.callFrameStack.inspect())
		result := cf.method(cf.receiver, cf.sourceLine, t, args, cf.blockFrame)
		t.Stack.Push(&Pointer{Target: result})
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

func (t *Thread) removeUselessBlockFrame(frame callFrame) {
	topFrame := t.callFrameStack.top()

	if topFrame != nil && topFrame.IsSourceBlock() {
		t.callFrameStack.pop().stopExecution()
	}
}

func (t *Thread) reportErrorAndStop() {
	cf := t.callFrameStack.top()

	if cf != nil {
		cf.stopExecution()
	}

	top := t.Stack.top().Target
	err := top.(*Error)

	if !err.storedTraces {
		for i := t.callFrameStack.pointer - 1; i > 0; i-- {
			frame := t.callFrameStack.callFrames[i]

			if frame.IsBlock() {
				continue
			}

			msg := fmt.Sprintf("from %s:%d", frame.FileName(), frame.SourceLine())
			err.stackTraces = append(err.stackTraces, msg)
		}

		err.storedTraces = true
	}

	panic(err)

	if t.vm.mode == NormalMode {
		if t.isMainThread() {
			os.Exit(1)
		}
	}
}

func (t *Thread) execInstruction(cf *normalCallFrame, i *instruction) {
	cf.pc++

	//fmt.Println(t.callFrameStack.inspect())
	//fmt.Println(i.inspect())
	i.action.operation(t, i.sourceLine, cf, i.Params...)
	//fmt.Println("============================")
	//fmt.Println(t.callFrameStack.inspect())
}

// Yield to a call frame
func (t *Thread) Yield(args ...Object) *Pointer {
	return t.builtinMethodYield(t.currentFrame.BlockFrame(), args...)
}

// BlockGiven returns whethe or not we have a block frame below us in the stack
func (t *Thread) BlockGiven() bool {
	return t.currentFrame.BlockFrame() != nil
}

func (t *Thread) builtinMethodYield(blockFrame *normalCallFrame, args ...Object) *Pointer {
	if blockFrame.IsRemoved() {
		return &Pointer{Target: NULL}
	}

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

	if blockFrame.IsRemoved() {
		return &Pointer{Target: NULL}
	}

	return t.Stack.top()
}

func (t *Thread) retrieveBlock(fileName, blockFlag string, sourceLine int) (blockFrame *normalCallFrame) {
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

func (t *Thread) sendMethod(methodName string, argCount int, blockFrame *normalCallFrame, sourceLine int) {
	var method Object

	if arr, ok := t.Stack.top().Target.(*ArrayObject); ok && arr.splat {
		// Pop array
		t.Stack.Pop()
		// Can't count array self, only the number of array elements
		argCount = argCount + len(arr.Elements)
		for _, elem := range arr.Elements {
			t.Stack.Push(&Pointer{Target: elem})
		}
	}

	argPr := t.Stack.pointer - argCount - 1
	receiverPr := argPr - 1
	receiver := t.Stack.data[receiverPr].Target

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
		t.Stack.data[argPr+i] = t.Stack.data[argPr+i+1]
	}

	t.Stack.pointer--

	method = receiver.findMethod(methodName)

	if method == nil {
		t.setErrorObject(receiverPr, argPr, errors.UndefinedMethodError, sourceLine, "Undefined Method '%+v' for %+v", methodName, receiver.ToString())
	}

	sendCallFrame := t.callFrameStack.top()

	switch m := method.(type) {
	case *MethodObject:
		callObj := newCallObject(receiver, m, receiverPr, argCount, &bytecode.ArgSet{}, blockFrame, 1)
		t.evalMethodObject(callObj)
	case *BuiltinMethodObject:
		t.evalBuiltinMethod(receiver, m, receiverPr, argCount, &bytecode.ArgSet{}, blockFrame, sourceLine, sendCallFrame.FileName())
	case *Error:
		t.pushErrorObject(errors.InternalError, sourceLine, m.ToString())
	}
}

func (t *Thread) evalBuiltinMethod(receiver Object, method *BuiltinMethodObject, receiverPtr, argCount int, argSet *bytecode.ArgSet, blockFrame *normalCallFrame, sourceLine int, fileName string) {
	argPtr := receiverPtr + 1

	cf := newGoMethodCallFrame(
		method.Fn,
		receiver,
		argCount,
		argPtr,
		method.Name,
		fileName,
		sourceLine,
		blockFrame,
	)

	t.callFrameStack.push(cf)
	t.startFromTopFrame()
	evaluated := t.Stack.top()

	_, ok := receiver.(*RClass)
	if method.Name == "new" && ok {
		instance, ok := evaluated.Target.(*RObject)
		if ok && instance.InitializeMethod != nil {
			callObj := newCallObject(instance, instance.InitializeMethod, receiverPtr, argCount, argSet, blockFrame, sourceLine)
			t.evalMethodObject(callObj)
		}
	}

	t.Stack.Set(receiverPtr, evaluated)
	t.Stack.pointer = cf.argPtr

	if err, ok := evaluated.Target.(*Error); ok {
		panic(err.Message())
	}
}

// TODO: Move instruction into call object
func (t *Thread) evalMethodObject(call *callObject) {
	normalParamsCount := call.normalParamsCount()
	paramTypes := call.paramTypes()
	paramsCount := len(call.paramTypes())
	stack := t.Stack.data
	sourceLine := call.sourceLine

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
				call.assignSplatArgument(stack, t.vm.InitArrayObject([]Object{}))
			}
		}
	} else {
		call.assignNormalArguments(stack)
	}

	t.callFrameStack.push(call.callFrame)
	t.startFromTopFrame()

	t.Stack.Set(call.receiverPtr, t.Stack.top())
	t.Stack.pointer = call.argPtr()
}

func (t *Thread) reportArgumentError(sourceLine, idealArgNumber int, methodName string, exactArgNumber int, receiverPtr int) {
	var message string

	if idealArgNumber > exactArgNumber {
		message = "Expect at least %d args for method '%s'. got: %d"
	} else {
		message = "Expect at most %d args for method '%s'. got: %d"
	}

	t.setErrorObject(receiverPtr, receiverPtr+1, errors.ArgumentError, sourceLine, message, idealArgNumber, methodName, exactArgNumber)
}

func (t *Thread) pushErrorObject(errorType string, sourceLine int, format string, args ...interface{}) {
	err := t.vm.InitErrorObject(errorType, sourceLine, format, args...)
	t.Stack.Push(&Pointer{Target: err})
	panic(err.Message())
}

func (t *Thread) setErrorObject(receiverPtr, sp int, errorType string, sourceLine int, format string, args ...interface{}) {
	err := t.vm.InitErrorObject(errorType, sourceLine, format, args...)
	t.Stack.Set(receiverPtr, &Pointer{Target: err})
	t.Stack.pointer = sp
	panic(err.Message())
}

// Other helper functions  ----------------------------------------------

// blockIsEmpty returns true if the block is empty
func blockIsEmpty(blockFrame *normalCallFrame) bool {
	if blockFrame.instructionSet.instructions[0].action.name == bytecode.Leave {
		return true
	}
	return false
}
