package vm

import (
	"fmt"
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
		if msg, yes := t.hasError(); yes {
			fmt.Println(msg)
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

	defer func() {
		if p := recover(); p != nil {
			if t.vm.stackTraceCount == 0 {
				fmt.Printf("Internal Error: %s\n", p)
			}
			t.vm.stackTraceCount++
			panic(p)
		}
	}()

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

func (t *thread) evalBuiltInMethod(receiver Object, method *BuiltInMethodObject, receiverPr, argCount, argPr int, blockFrame *callFrame) {
	methodBody := method.Fn(receiver)
	args := []Object{}

	for i := 0; i < argCount; i++ {
		args = append(args, t.stack.Data[argPr+i].Target)
	}

	evaluated := methodBody(t, args, blockFrame)

	_, ok := receiver.(*RClass)
	if method.Name == "new" && ok {
		instance, ok := evaluated.(*RObject)
		if ok && instance.InitializeMethod != nil {
			t.evalMethodObject(instance, instance.InitializeMethod, receiverPr, argCount, argPr, blockFrame)
		}
	}
	t.stack.Data[receiverPr] = &Pointer{Target: evaluated}
	t.sp = receiverPr + 1
}

func (t *thread) evalMethodObject(receiver Object, method *MethodObject, receiverPr, argC, argPr int, blockFrame *callFrame) {
	var normalArgCount int

	c := newCallFrame(method.instructionSet)
	c.self = receiver

	for _, at := range method.instructionSet.argTypes {
		if at == bytecode.NormalArg {
			normalArgCount++
		}
	}

	if argC < normalArgCount {
		e := t.vm.initErrorObject(ArgumentError, "Expect at least %d args for method '%s'. got: %d", normalArgCount, method.Name, argC)
		t.stack.push(&Pointer{Target:e})
	} else if argC > method.argc {
		e := t.vm.initErrorObject(ArgumentError, "Expect at most %d args for method '%s'. got: %d", method.argc, method.Name, argC)
		t.stack.push(&Pointer{Target: e})
	} else {
		for i := 0; i < argC; i++ {
			c.insertLCL(i, 0, t.stack.Data[argPr+i].Target)
		}

		c.blockFrame = blockFrame
		t.callFrameStack.push(c)
		t.startFromTopFrame()
	}

	t.stack.Data[receiverPr] = t.stack.top()
	t.sp = receiverPr + 1
}

// TODO: Use this method to replace unnecessary panics
func (t *thread) returnError(msg string) {
	err := &Error{Message: msg}
	t.stack.push(&Pointer{Target: err})
}

func (t *thread) UndefinedMethodError(methodName string, receiver Object) {
	err := t.vm.initErrorObject(UndefinedMethodError, "Undefined Method '%+v' for %+v", methodName, receiver.toString())
	t.stack.push(&Pointer{Target:err})
}

func (t *thread) UnsupportedMethodError(methodName string, receiver Object) *Error {
	return t.vm.initErrorObject(UnsupportedMethodError, "Unsupported Method %s for %+v", methodName, receiver.toString())
}
