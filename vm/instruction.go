package vm

import (
	"fmt"
	"github.com/goby-lang/goby/bytecode"
	"strings"
)

type operation func(t *thread, cf *callFrame, args ...interface{})

type operationType string

type label struct {
	name string
	Type labelType
}

type labelType string

type action struct {
	name      string
	operation operation
}

type instruction struct {
	action *action
	Params []interface{}
	Line   int
}

type instructionSet struct {
	label        *label
	instructions []*instruction
	filename     filename
}

func (is *instructionSet) define(line int, a *action, params ...interface{}) {
	i := &instruction{action: a, Params: params, Line: line}
	is.instructions = append(is.instructions, i)
}

var builtInActions = map[operationType]*action{
	bytecode.Pop: {
		name: bytecode.Pop,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			t.stack.pop()
		},
	},
	bytecode.PutObject: {
		name: bytecode.PutObject,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			object := initializeObjectFromInstruction(args[0])
			t.stack.push(&Pointer{Target: object})
		},
	},
	bytecode.GetConstant: {
		name: bytecode.GetConstant,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			constName := args[0].(string)
			c := t.vm.lookupConstant(cf, constName)

			if c == nil {
				msg := "Can't find constant: " + constName
				t.returnError(msg)
			}

			t.stack.push(c)
		},
	},
	bytecode.GetLocal: {
		name: bytecode.GetLocal,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			index := args[0].(int)
			depth := 0

			if len(args) >= 2 {
				depth = args[1].(int)
			}

			p := cf.getLCL(index, depth)

			if p == nil {
				t.stack.push(&Pointer{NULL})
				return
			}

			t.stack.push(p)
		},
	},
	bytecode.GetInstanceVariable: {
		name: bytecode.GetInstanceVariable,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			variableName := args[0].(string)
			v, ok := cf.self.(*RObject).InstanceVariables.get(variableName)
			if !ok {
				t.stack.push(&Pointer{Target: NULL})
				return
			}

			p := &Pointer{Target: v}
			t.stack.push(p)
		},
	},
	bytecode.SetInstanceVariable: {
		name: bytecode.SetInstanceVariable,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			variableName := args[0].(string)
			p := t.stack.pop()
			cf.self.(*RObject).InstanceVariables.set(variableName, p.Target)
		},
	},
	bytecode.SetLocal: {
		name: bytecode.SetLocal,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			v := t.stack.pop()
			depth := 0

			if len(args) >= 2 {
				depth = args[1].(int)
			}
			cf.insertLCL(args[0].(int), depth, v.Target)
		},
	},
	bytecode.SetConstant: {
		name: bytecode.SetConstant,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			constName := args[0].(string)
			v := t.stack.pop()

			cf.storeConstant(constName, v)
		},
	},
	bytecode.NewArray: {
		name: bytecode.NewArray,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			argCount := args[0].(int)
			elems := []Object{}

			for i := 0; i < argCount; i++ {
				v := t.stack.pop()
				elems = append([]Object{v.Target}, elems...)
			}

			arr := initializeArray(elems)
			t.stack.push(&Pointer{arr})
		},
	},
	bytecode.NewHash: {
		name: bytecode.NewHash,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			argCount := args[0].(int)
			pairs := map[string]Object{}

			for i := 0; i < argCount/2; i++ {
				v := t.stack.pop()
				k := t.stack.pop()
				pairs[k.Target.(*StringObject).Value] = v.Target
			}

			hash := initializeHash(pairs)
			t.stack.push(&Pointer{hash})
		},
	},
	bytecode.BranchUnless: {
		name: bytecode.BranchUnless,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			v := t.stack.pop()
			bool, isBool := v.Target.(*BooleanObject)

			if isBool {
				if bool.Value {
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
		},
	},
	bytecode.BranchIf: {
		name: bytecode.BranchIf,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			v := t.stack.pop()
			bool, isBool := v.Target.(*BooleanObject)

			if isBool {
				if !bool.Value {
					return
				}

				line := args[0].(int)
				cf.pc = line
				return
			}
		},
	},
	bytecode.Jump: {
		name: bytecode.Jump,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			cf.pc = args[0].(int)
		},
	},
	bytecode.PutSelf: {
		name: bytecode.PutSelf,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			t.stack.push(&Pointer{cf.self})
		},
	},
	bytecode.PutString: {
		name: bytecode.PutString,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			object := initializeObjectFromInstruction(args[0])
			t.stack.push(&Pointer{object})
		},
	},
	bytecode.PutNull: {
		name: bytecode.PutNull,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			t.stack.push(&Pointer{NULL})
		},
	},
	bytecode.DefMethod: {
		name: bytecode.DefMethod,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			argCount := args[0].(int)
			methodName := t.stack.pop().Target.(*StringObject).Value
			is, _ := t.getMethodIS(methodName, cf.instructionSet.filename)
			method := &MethodObject{Name: methodName, argc: argCount, instructionSet: is, class: methodClass}

			v := t.stack.pop().Target
			switch self := v.(type) {
			case *RClass:
				self.Methods.set(methodName, method)
			default:
				self.returnClass().(*RClass).Methods.set(methodName, method)
			}
		},
	},
	bytecode.DefSingletonMethod: {
		name: bytecode.DefSingletonMethod,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			argCount := args[0].(int)
			methodName := t.stack.pop().Target.(*StringObject).Value
			is, _ := t.getMethodIS(methodName, cf.instructionSet.filename)
			method := &MethodObject{Name: methodName, argc: argCount, instructionSet: is, class: methodClass}

			v := t.stack.pop().Target

			switch self := v.(type) {
			case *RClass:
				self.setSingletonMethod(methodName, method)
			}
			// TODO: Support something like:
			// ```
			// f = Foo.new
			// def f.bar
			//   10
			// end
			// ```
		},
	},
	bytecode.DefClass: {
		name: bytecode.DefClass,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			subject := strings.Split(args[0].(string), ":")
			subjectType, subjectName := subject[0], subject[1]

			classPtr, ok := cf.lookupConstant(subjectName)

			if !ok {
				class := initializeClass(subjectName, subjectType == "module")
				classPtr = cf.storeConstant(class.Name, class)

				if len(args) >= 2 {
					superClassName := args[1].(string)
					superClass := t.vm.lookupConstant(cf, superClassName)
					inheritedClass, ok := superClass.Target.(*RClass)

					if !ok {
						t.returnError("Constant " + superClassName + " is not a class. got=" + string(superClass.Target.returnClass().ReturnName()))
					}

					class.pseudoSuperClass = inheritedClass
					class.superClass = inheritedClass
				}
			}

			is := t.getClassIS(subjectName, cf.instructionSet.filename)

			t.stack.pop()
			c := newCallFrame(is)
			c.self = classPtr.Target
			t.callFrameStack.push(c)
			t.startFromTopFrame()

			t.stack.push(classPtr)
		},
	},
	bytecode.Send: {
		name: bytecode.Send,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			var method Object

			methodName := args[0].(string)
			argCount := args[1].(int)
			argPr := t.sp - argCount
			receiverPr := argPr - 1
			receiver := t.stack.Data[receiverPr].Target

			switch r := receiver.(type) {
			case Class:
				method = r.lookupClassMethod(methodName)
			default:
				method = r.returnClass().lookupInstanceMethod(methodName)
			}

			if method == nil {
				t.UndefinedMethodError(methodName, receiver)
				return
			}

			blockFrame := t.retrieveBlock(cf, args)

			switch m := method.(type) {
			case *MethodObject:
				t.evalMethodObject(receiver, m, receiverPr, argCount, argPr, blockFrame)
			case *BuiltInMethodObject:
				t.evalBuiltInMethod(receiver, m, receiverPr, argCount, argPr, blockFrame)
			case *Error:
				t.returnError(m.toString())
			}
		},
	},
	bytecode.InvokeBlock: {
		name: bytecode.InvokeBlock,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			argCount := args[0].(int)
			argPr := t.sp - argCount
			receiverPr := argPr - 1
			receiver := t.stack.Data[receiverPr].Target

			if cf.blockFrame == nil {
				t.returnError("Can't yield without a block")
			}

			c := newCallFrame(cf.blockFrame.instructionSet)
			c.blockFrame = cf.blockFrame
			c.ep = cf.blockFrame.ep
			c.self = receiver

			for i := 0; i < argCount; i++ {
				c.locals[i] = t.stack.Data[argPr+i]
			}

			t.callFrameStack.push(c)
			t.startFromTopFrame()

			t.stack.Data[receiverPr] = t.stack.top()
			t.sp = receiverPr + 1
		},
	},
	bytecode.Leave: {
		name: bytecode.Leave,
		operation: func(t *thread, cf *callFrame, args ...interface{}) {
			cf = t.callFrameStack.pop()
			cf.pc = len(cf.instructionSet.instructions)
		},
	},
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
	t.stack.Data[receiverPr] = &Pointer{evaluated}
	t.sp = receiverPr + 1
}

func (t *thread) evalMethodObject(receiver Object, method *MethodObject, receiverPr, argC, argPr int, blockFrame *callFrame) {
	c := newCallFrame(method.instructionSet)
	c.self = receiver

	for i := 0; i < argC; i++ {
		c.insertLCL(i, 0, t.stack.Data[argPr+i].Target)
	}

	c.blockFrame = blockFrame
	t.callFrameStack.push(c)
	t.startFromTopFrame()

	t.stack.Data[receiverPr] = t.stack.top()
	t.sp = receiverPr + 1
}

func initializeObjectFromInstruction(value interface{}) Object {
	switch v := value.(type) {
	case int:
		return initilaizeInteger(int(v))
	case string:
		switch v {
		case "true":
			return TRUE
		case "false":
			return FALSE
		default:
			return initializeString(v)
		}
	}

	return newError(fmt.Sprintf("Unknow data type %T", value))
}
