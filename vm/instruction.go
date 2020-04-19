package vm

import (
	"fmt"
	"strings"

	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

type operation func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) *Error

type operationType = uint8

type setType = string

type instructionSet struct {
	name         string
	instructions []*bytecode.Instruction
	filename     filename
	paramTypes   *bytecode.ArgSet
}

var operations [bytecode.InstructionCount]operation

// This is for avoiding initialization loop
func init() {
	operations = [bytecode.InstructionCount]operation{
		bytecode.Pop:                 pop,
		bytecode.Dup:                 dup,

		bytecode.PutSelf:             putSelf,
		bytecode.PutBoolean:          putBoolean,
		bytecode.PutString:           putString,
		bytecode.PutFloat:            putFloat,
		bytecode.PutNull:             putNull,
		bytecode.PutObject:           putObject,

		bytecode.NewRange:            newRange,
		bytecode.NewHash:             newHash,
		bytecode.NewArray:            newArray,
		bytecode.ExpandArray:         expandArray,
		bytecode.SplatArray:          splatArray,

		bytecode.DefMethod:           defMethod,
		bytecode.DefSingletonMethod:  defSingletonMethod,
		bytecode.DefClass:            defClass,

		bytecode.SetConstant:         setConstant,
		bytecode.GetConstant:         getConstant,
		bytecode.SetLocal:            setLocal,
		bytecode.GetLocal:            getLocal,
		bytecode.SetInstanceVariable: setInstanceVariable,
		bytecode.GetInstanceVariable: getInstanceVariable,

		bytecode.Send:         send,

		bytecode.InvokeBlock:  invokeBlock,
		bytecode.GetBlock:     getBlock,

		bytecode.BranchIf:     branchIf,
		bytecode.BranchUnless: branchUnless,
		bytecode.Jump:         jump,
		bytecode.Break:        breakInc,
		bytecode.Leave:        leave,
	}
}

// InitObjectFromGoType creates an object based on Go's type
func (v *VM) InitObjectFromGoType(value interface{}) Object {
	switch val := value.(type) {
	case nil:
		return NULL
	case int:
		return v.InitIntegerObject(val)
	case int64:
		return v.InitIntegerObject(int(val))
	case int32:
		return v.InitIntegerObject(int(val))
	case float64:
		return v.initFloatObject(val)
	case []uint8:
		var bytes []byte

		for _, i := range val {
			bytes = append(bytes, byte(i))
		}

		return v.InitStringObject(string(bytes))
	case string:
		return v.InitStringObject(val)
	case bool:
		return toBooleanObject(val)
	case []interface{}:
		var objects []Object

		for _, elem := range val {
			objects = append(objects, v.InitObjectFromGoType(elem))
		}

		return v.InitArrayObject(objects)
	default:
		return v.initGoObject(value)
	}
}

func pop(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	t.Stack.Pop()
	return
}

func dup(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	obj := t.Stack.top().Target
	t.Stack.Push(&Pointer{Target: obj})
	return
}

func getBlock(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	if cf.blockFrame == nil {
		t.pushErrorObject(errors.InternalError, sourceLine, "Can't get block without a block argument")
	}

	blockFrame := cf.blockFrame

	if cf.blockFrame.ep == cf.ep {
		blockFrame = cf.blockFrame.ep.blockFrame
	}

	blockObject := t.vm.initBlockObject(blockFrame.instructionSet, blockFrame.ep, t.Stack.data[t.Stack.pointer-1].Target)

	t.Stack.Push(&Pointer{Target: blockObject})

	return
}

func invokeBlock(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	argCount := args[0].(int)
	argPr := t.Stack.pointer - argCount
	receiverPr := argPr - 1
	receiver := t.Stack.data[receiverPr].Target

	if cf.blockFrame == nil {
		t.pushErrorObject(errors.InternalError, sourceLine, errors.CantYieldWithoutBlockFormat)
	}

	blockFrame := cf.blockFrame

	/*
		This is for such condition:

		```ruby
		def foo(x)
		  yield(x + 10)
		end

		def bar(y)
		  foo(y) do |f|
			yield(f) # <------- here
		  end
		end

		bar(100) do |b|
		  puts(b) #=> 110
		end
		```

		In this case the target frame is not first block frame we meet. It should be `bar`'s block.
		And bar's frame is foo block frame's ep, so our target frame is ep's block frame.
	*/
	if cf.blockFrame.ep == cf.ep {
		blockFrame = cf.blockFrame.ep.blockFrame
	}

	c := newNormalCallFrame(blockFrame.instructionSet, blockFrame.instructionSet.filename, sourceLine)
	c.blockFrame = blockFrame
	c.ep = blockFrame.ep
	c.self = receiver
	c.isBlock = true

	for i := 0; i < argCount; i++ {
		c.locals[i] = t.Stack.data[argPr+i]
	}

	t.callFrameStack.push(c)
	t.startFromTopFrame()

	t.Stack.Set(receiverPr, t.Stack.top())
	t.Stack.pointer = receiverPr + 1

	return
}

func leave(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	t.callFrameStack.pop()
	cf.stopExecution()
	return
}

func send(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	var blockFlag string

	methodName := args[0].(string)
	argCount := args[1].(int)
	blockFlag, ok := args[2].(string)

	if !ok {
		blockFlag = ""
	}

	argSet := args[3].(*bytecode.ArgSet)

	// Deal with splat arguments
	if arr, ok := t.Stack.top().Target.(*ArrayObject); ok && arr.splat {
		// Pop array
		t.Stack.Pop()
		// Can't count array itself, only the number of array elements
		argCount = argCount - 1 + len(arr.Elements)
		for _, elem := range arr.Elements {
			t.Stack.Push(&Pointer{Target: elem})
		}
	}

	argPr := t.Stack.pointer - argCount
	receiverPr := argPr - 1
	receiver := t.Stack.data[receiverPr].Target

	// Find Block
	blockFrame := t.retrieveBlock(cf.FileName(), blockFlag, cf.SourceLine())

	if blockFrame != nil {
		blockFrame.ep = cf
		blockFrame.self = cf.self
		blockFrame.sourceLine = sourceLine
		t.callFrameStack.push(blockFrame)
	}

	t.findAndCallMethod(receiver, methodName, receiverPr, argSet, argCount, argPr, sourceLine, blockFrame, cf.fileName)
	return
}

func defClass(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	subject := strings.Split(args[0].(string), ":")
	subjectType, subjectName := subject[0], subject[1]

	classPtr := cf.lookupConstantUnderAllScope(subjectName)

	if classPtr == nil {
		var class *RClass
		if subjectType == "module" {
			class = t.vm.initializeModule(subjectName)
		} else {
			class = t.vm.initializeClass(subjectName)
		}

		classPtr = cf.storeConstant(class.Name, class)

		if len(args) >= 2 {
			superClassName := args[1].(string)
			superClass := t.vm.lookupConstant(cf, superClassName)
			inheritedClass, ok := superClass.Target.(*RClass)

			if !ok {
				t.pushErrorObject(errors.InternalError, sourceLine, "Constant %s is not a class. got: %s", superClassName, string(superClass.Target.Class().ReturnName()))
			}

			if inheritedClass.isModule {
				t.pushErrorObject(errors.InternalError, sourceLine, "Module inheritance is not supported: %s", inheritedClass.Name)
			}

			class.inherits(inheritedClass)
		}
	}

	is := t.getClassIS(subjectName, cf.FileName())

	t.Stack.Pop()
	c := newNormalCallFrame(is, cf.FileName(), sourceLine)
	c.self = classPtr.Target
	t.callFrameStack.push(c)
	t.startFromTopFrame()

	t.Stack.Push(classPtr)
	return
}

func defSingletonMethod(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	argCount := args[0].(int)
	methodName := t.Stack.Pop().Target.(*StringObject).value
	is, _ := t.getMethodIS(methodName, cf.FileName())
	method := &MethodObject{Name: methodName, argc: argCount, instructionSet: is, BaseObj: NewBaseObject(t.vm.TopLevelClass(classes.MethodClass))}

	v := t.Stack.Pop().Target

	switch v := v.(type) {
	case *RClass:
		v.SingletonClass().Methods.set(methodName, method)
	default:
		singletonClass := t.vm.createRClass(fmt.Sprintf("#<Class:#<%s:%d>>", v.Class().Name, v.ID()))
		singletonClass.Methods.set(methodName, method)
		singletonClass.isSingleton = true
		v.SetSingletonClass(singletonClass)
	}
	return
}

func defMethod(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	argCount := args[0].(int)
	methodName := t.Stack.Pop().Target.(*StringObject).value
	is, ok := t.getMethodIS(methodName, cf.FileName())

	if !ok {
		t.pushErrorObject(errors.InternalError, sourceLine, "Can't get method %s's instruction set.", methodName)
	}

	method := &MethodObject{Name: methodName, argc: argCount, instructionSet: is, BaseObj: NewBaseObject(t.vm.TopLevelClass(classes.MethodClass))}

	v := t.Stack.Pop().Target
	switch self := v.(type) {
	case *RClass:
		self.Methods.set(methodName, method)
	default:
		self.Class().Methods.set(methodName, method)
	}

	return
}

func putNull(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	t.Stack.Push(&Pointer{Target: NULL})

	return
}

func putFloat(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	value := args[0].(float64)
	object := t.vm.initFloatObject(value)
	t.Stack.Push(&Pointer{Target: object})
	return

}

func putString(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	object := t.vm.InitObjectFromGoType(args[0])
	t.Stack.Push(&Pointer{Target: object})
	return

}

func putSelf(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	t.Stack.Push(&Pointer{Target: cf.self})
	return

}

func breakInc(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	/*
		Normal frame. IS name: ProgramStart. is block: false. source line: 1
		Normal frame. IS name: 0. is block: true. ep: 17. source line: 5 <- The block source
		Go method frame. Method name: each. <- The method call with block
		Normal frame. IS name: 0. is block: true. ep: 17. source line: 5 <- The block execution
	*/

	if cf.IsBlock() {
		/*
		  1. Remove block execution frame
		  2. Remove method call frame
		  3. Remove block source frame
		*/
		for i := 0; i < 3; i++ {
			frame := t.callFrameStack.pop()
			frame.stopExecution()
			frame.setAsRemoved()
		}
	}

	return
}

func jump(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	cf.pc = args[0].(int)
	return

}

func branchIf(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	v := t.Stack.Pop()
	bo, isBool := v.Target.(*BooleanObject)

	if isBool && !bo.value {
		return
	}

	_, isNull := v.Target.(*NullObject)

	if isNull {
		return
	}

	line := args[0].(int)
	cf.pc = line
	return

}

func branchUnless(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	v := t.Stack.Pop()
	bo, isBool := v.Target.(*BooleanObject)

	if isBool {
		if bo.value {
			return
		}

		line := args[0].(int)
		cf.pc = line
		return
	}

	_, isNull := v.Target.(*NullObject)

	if isNull {
		line := args[0].(int)
		cf.pc = line
		return
	}
	return

}

func newHash(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	argCount := args[0].(int)
	pairs := map[string]Object{}

	for i := 0; i < argCount/2; i++ {
		v := t.Stack.Pop()
		k := t.Stack.Pop()
		pairs[k.Target.(*StringObject).value] = v.Target
	}

	hash := t.vm.InitHashObject(pairs)
	t.Stack.Push(&Pointer{Target: hash})
	return

}

func splatArray(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	obj := t.Stack.top().Target
	arr, ok := obj.(*ArrayObject)

	if !ok {
		return
	}

	arr.splat = true
	return

}

func expandArray(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	arrLength := args[0].(int)
	arr, ok := t.Stack.Pop().Target.(*ArrayObject)

	if !ok {
		t.pushErrorObject(errors.TypeError, sourceLine, "Expect stack top's value to be an Array when executing 'expandarray' instruction.")
	}

	var elems []Object

	for i := 0; i < arrLength; i++ {
		var elem Object
		if i < len(arr.Elements) {
			elem = arr.Elements[i]
		} else {
			elem = NULL
		}

		elems = append([]Object{elem}, elems...)
	}

	for _, elem := range elems {
		t.Stack.Push(&Pointer{Target: elem})
	}
	return

}

func newArray(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	argCount := args[0].(int)
	var elems []Object

	for i := 0; i < argCount; i++ {
		v := t.Stack.Pop()
		elems = append([]Object{v.Target}, elems...)
	}

	arr := t.vm.InitArrayObject(elems)
	t.Stack.Push(&Pointer{Target: arr})
	return

}

func newRange(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	rangeEnd := t.Stack.Pop().Target.(*IntegerObject).value
	rangeStart := t.Stack.Pop().Target.(*IntegerObject).value

	t.Stack.Push(&Pointer{Target: t.vm.initRangeObject(rangeStart, rangeEnd)})
	return

}

func setConstant(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	constName := args[0].(string)
	c := cf.lookupConstantInCurrentScope(constName)
	v := t.Stack.Pop()

	if c != nil {
		t.pushErrorObject(errors.ConstantAlreadyInitializedError, sourceLine, "Constant %s already been initialized. Can't assign value to a constant twice.", constName)
	}

	cf.storeConstant(constName, v)
	return

}

func setLocal(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	var optioned bool
	p := t.Stack.Pop()
	depth := args[0].(int)
	index := args[1].(int)

	if len(args) > 2 && args[2].(int) == 1 {
		optioned = true
	}

	if optioned {
		if cf.getLCL(index, depth) == nil {
			cf.insertLCL(index, depth, p.Target)
		}

		return
	}

	cf.insertLCL(index, depth, p.Target)

	var obj Object

	switch v := p.Target.(type) {
	case *HashObject:
		obj = v.copy()
	case *ArrayObject:
		obj = v.copy()
	case *ChannelObject:
		obj = v.copy()
	default:
		obj = v
	}

	t.Stack.Push(&Pointer{Target: obj})
	return
}

func setInstanceVariable(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	variableName := args[0].(string)
	p := t.Stack.Pop()
	cf.self.InstanceVariableSet(variableName, p.Target)

	var obj Object

	switch v := p.Target.(type) {
	case *HashObject:
		obj = v.copy()
	case *ArrayObject:
		obj = v.copy()
	case *ChannelObject:
		obj = v.copy()
	default:
		obj = v
	}

	t.Stack.Push(&Pointer{Target: obj})
	return
}

func getInstanceVariable(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	variableName := args[0].(string)
	v, ok := cf.self.InstanceVariableGet(variableName)
	if !ok {
		t.Stack.Push(&Pointer{Target: NULL})
		return
	}

	p := &Pointer{Target: v}
	t.Stack.Push(p)
	return
}

func getLocal(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	depth := args[0].(int)
	index := args[1].(int)

	p := cf.getLCL(index, depth)

	if p == nil {
		t.Stack.Push(&Pointer{Target: NULL})
		return
	}

	t.Stack.Push(p)
	return
}

func getConstant(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	constName := args[0].(string)
	c := t.vm.lookupConstant(cf, constName)

	if c == nil {
		t.pushErrorObject(errors.NameError, sourceLine, "uninitialized constant %s", constName)
	}

	c.isNamespace = args[1].(bool)

	if t.Stack.top() != nil && t.Stack.top().isNamespace {
		t.Stack.Pop()
	}

	t.Stack.Push(c)
	return
}

func putObject(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	object := t.vm.InitObjectFromGoType(args[0])
	t.Stack.Push(&Pointer{Target: object})
	return
}

func putBoolean(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) (err *Error) {
	object := t.vm.InitObjectFromGoType(args[0])
	t.Stack.Push(&Pointer{Target: object})
	return
}
