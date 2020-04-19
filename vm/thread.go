package vm

import (
	"fmt"
	"io/ioutil"
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

// VM returns the vm of the thread
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

func (t *Thread) startFromTopFrame() (err *Error) {
	cf := t.callFrameStack.top()
	err = t.evalCallFrame(cf)

	if err != nil {
		cf := t.callFrameStack.top()

		if cf != nil {
			cf.stopExecution()
		}
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

		return err
	}

	return
}

func (t *Thread) evalCallFrame(cf callFrame) (err *Error) {
	t.currentFrame = cf

	switch cf := cf.(type) {
	case *normalCallFrame:
		for cf.pc < cf.instructionsCount() {
			i := cf.instructionSet.instructions[cf.pc]
			err = t.execInstruction(cf, i)

			if err != nil {
				return
			}
		}
	case *goMethodCallFrame:
		args := []Object{}

		for i := 0; i < cf.argCount; i++ {
			args = append(args, t.Stack.data[cf.argPtr+i].Target)
		}
		//fmt.Println("-----------------------")
		//fmt.Println(t.callFrameStack.inspect())
		result := cf.method(cf.receiver, cf.sourceLine, t, args, cf.blockFrame)
		err, ok := result.(*Error)

		if ok {
			return err
		}

		t.Stack.Push(&Pointer{Target: result})
		//fmt.Println(t.callFrameStack.inspect())
		//fmt.Println("-----------------------")

		// this check is to prevent popping out a call frame that should be kept
		// usually the top call frame would be the `cf` here
		// but in some rare cases it could be other call frame, and if we pop it here we'll have some troubles in the later execution
		//
		// one example of the problem would be https://github.com/goby-lang/goby/issues/584
		// the cause of this issue is that the "break" instruction always pops 3 call frames
		// (this is necessary, please read the comments inside `Break` instruction for more detail)
		// and because we already popped the `cf` during the `Break` instruction, what we pop here would be the top-level call frame, which causes the program to crash
		if t.callFrameStack.top() == cf {
			t.callFrameStack.pop()
		}
	}

	t.removeUselessBlockFrame(cf)

	return
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

func (t *Thread) execInstruction(cf *normalCallFrame, i *bytecode.Instruction) (err *Error) {
	cf.pc++

	//fmt.Println(t.callFrameStack.inspect())
	//fmt.Println(i.inspect())
	ins := operations[i.Opcode]
	err = ins(t, i.SourceLine(), cf, i.Params...)
	//fmt.Println("============================")
	//fmt.Println(t.callFrameStack.inspect())
	return
}

// Yield to a call frame
func (t *Thread) Yield(args ...Object) (*Pointer, *Error) {
	return t.builtinMethodYield(t.currentFrame.BlockFrame(), args...)
}

// BlockGiven returns whethe or not we have a block frame below us in the stack
func (t *Thread) BlockGiven() bool {
	return t.currentFrame.BlockFrame() != nil
}

func (t *Thread) builtinMethodYield(blockFrame *normalCallFrame, args ...Object) (p *Pointer, err *Error) {
	if blockFrame.IsRemoved() {
		return &Pointer{Target: NULL}, nil
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
	err = t.startFromTopFrame()

	if err != nil {
		return nil, err
	}

	if blockFrame.IsRemoved() {
		return &Pointer{Target: NULL}, nil
	}

	return t.Stack.top(), nil
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

func (t *Thread) findMethod(receiver Object, methodName string, receiverPr int, argCount int, argPr int, sourceLine int) (method Object, argC int, err *Error) {
	method = receiver.findMethod(methodName)

	if method == nil {
		mm := receiver.findMethodMissing(receiver.Class().inheritsMethodMissing)

		if mm == nil {
			err = t.setErrorObject(receiverPr, argPr, errors.NoMethodError, sourceLine, errors.UndefinedMethod, methodName, receiver.Inspect())
			return
		} else {
			// Move up args for missed method's name
			// before: | arg 1       | arg 2 |
			// after:  | method name | arg 1 | arg 2 |
			// TODO: Improve this
			t.Stack.Push(nil)

			for i := argCount - 1; i >= 0; i-- {
				position := argPr + i
				arg := t.Stack.data[argPr+i]
				t.Stack.Set(position+1, arg)
			}

			t.Stack.Set(argPr, &Pointer{Target: t.vm.InitStringObject(methodName)})
			argCount++

			method = mm
		}
	}

	return method, argCount, nil
}

func (t *Thread) findAndCallMethod(receiver Object, methodName string, receiverPr int, argSet *bytecode.ArgSet, argCount int, argPr int, sourceLine int, blockFrame *normalCallFrame, fileName string) (err *Error) {
	// argCount change if we ended up calling method_missing
	method, argCount, err := t.findMethod(receiver, methodName, receiverPr, argCount, argPr, sourceLine)

	if err != nil {
		return err
	}

	switch m := method.(type) {
	case *MethodObject:
		callObj := newCallObject(receiver, m, receiverPr, argCount, argSet, blockFrame, sourceLine)
		err = t.evalMethodObject(callObj)
	case *BuiltinMethodObject:
		err = t.evalBuiltinMethod(receiver, m, receiverPr, argCount, argSet, blockFrame, sourceLine, fileName)
	}

	return
}

func (t *Thread) sendMethod(methodName string, argCount int, blockFrame *normalCallFrame, sourceLine int) (err *Error) {
	if arr, ok := t.Stack.top().Target.(*ArrayObject); ok && arr.splat {
		// Pop array
		t.Stack.Pop()
		// Can't count array self, only the number of array elements
		argCount += len(arr.Elements)
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

	sendCallFrame := t.callFrameStack.top()

	return t.findAndCallMethod(receiver, methodName, receiverPr, &bytecode.ArgSet{}, argCount, argPr, sourceLine, blockFrame, sendCallFrame.FileName())
}

func (t *Thread) evalBuiltinMethod(receiver Object, method *BuiltinMethodObject, receiverPtr, argCount int, argSet *bytecode.ArgSet, blockFrame *normalCallFrame, sourceLine int, fileName string) (err *Error) {
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

	err = t.startFromTopFrame()

	if err != nil {
		return err
	}
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
		return err
	}

	return nil
}

// TODO: Move instruction into call object
func (t *Thread) evalMethodObject(call *callObject) (err *Error) {
	normalParamsCount := call.normalParamsCount()
	paramTypes := call.paramTypes()
	paramsCount := len(call.paramTypes())
	stack := t.Stack.data
	sourceLine := call.sourceLine

	if call.argCount > paramsCount && !call.method.isSplatArgIncluded() {
		return t.reportArgumentError(sourceLine, paramsCount, call.methodName(), call.argCount, call.receiverPtr)
	}

	if normalParamsCount > call.argCount {
		return t.reportArgumentError(sourceLine, normalParamsCount, call.methodName(), call.argCount, call.receiverPtr)
	}

	// Check if arguments include all the required keys before assign keyword arguments
	for paramIndex, paramType := range paramTypes {
		switch paramType {
		case bytecode.RequiredKeywordArg:
			paramName := call.paramNames()[paramIndex]
			if _, ok := call.hasKeywordArgument(paramName); !ok {
				return t.setErrorObject(call.receiverPtr, call.argPtr(), errors.ArgumentError, sourceLine, "Method %s requires key argument %s", call.methodName(), paramName)
			}
		}
	}

	e := call.assignKeywordArguments(stack)

	if e != nil {
		return t.setErrorObject(call.receiverPtr, call.argPtr(), errors.ArgumentError, sourceLine, e.Error())
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
	err = t.startFromTopFrame()

	if err != nil {
		return err
	}

	t.Stack.Set(call.receiverPtr, t.Stack.top())
	t.Stack.pointer = call.argPtr()

	return nil
}

func (t *Thread) reportArgumentError(sourceLine, idealArgNumber int, methodName string, exactArgNumber int, receiverPtr int) (err *Error) {
	var message string

	if idealArgNumber > exactArgNumber {
		message = "Expect at least %d args for method '%s'. got: %d"
	} else {
		message = "Expect at most %d args for method '%s'. got: %d"
	}

	return t.setErrorObject(receiverPtr, receiverPtr+1, errors.ArgumentError, sourceLine, message, idealArgNumber, methodName, exactArgNumber)
}

// pushErrorObject pushes the Error object to the stack
func (t *Thread) pushErrorObject(errorType string, sourceLine int, format string, args ...interface{}) (err *Error) {
	return t.vm.InitErrorObject(errorType, sourceLine, format, args...)
}

// setErrorObject replaces a certain stack element with the Error object
func (t *Thread) setErrorObject(receiverPtr, sp int, errorType string, sourceLine int, format string, args ...interface{}) (err *Error) {
	return t.vm.InitErrorObject(errorType, sourceLine, format, args...)
}

// Other helper functions  ----------------------------------------------

// blockIsEmpty returns true if the block is empty
func blockIsEmpty(blockFrame *normalCallFrame) bool {
	if blockFrame.instructionSet.instructions[0].ActionName() == "leave" {
		return true
	}
	return false
}
