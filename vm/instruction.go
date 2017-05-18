package vm

import (
	"bytes"
	"fmt"
	"github.com/goby-lang/goby/bytecode"
	"strings"
)

type operation func(vm *VM, cf *callFrame, args ...interface{})

type action struct {
	name      string
	operation operation
}

type instruction struct {
	action *action
	Params []interface{}
	Line   int
}

type label struct {
	name string
	Type labelType
}

type labelType string

type instructionSet struct {
	label        *label
	instructions []*instruction
	filename     filename
}

type operationType string

var builtInActions = map[operationType]*action{
	bytecode.Pop: {
		name: bytecode.Pop,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			vm.stack.pop()
		},
	},
	bytecode.PutObject: {
		name: bytecode.PutObject,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			object := initializeObject(args[0])
			vm.stack.push(&Pointer{Target: object})
		},
	},
	bytecode.GetConstant: {
		name: bytecode.GetConstant,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			constName := args[0].(string)
			c := vm.lookupConstant(cf, constName)

			if c == nil {
				msg := "Can't find constant: " + constName
				vm.returnError(msg)
			}

			vm.stack.push(c)
		},
	},
	bytecode.GetLocal: {
		name: bytecode.GetLocal,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			index := args[0].(int)
			depth := 0

			if len(args) >= 2 {
				depth = args[1].(int)
			}

			p := cf.getLCL(index, depth)

			if p == nil {
				panic(fmt.Sprintf("locals index: %d is nil. Callframe: %s", index, cf.instructionSet.label.name))
			}
			vm.stack.push(p)
		},
	},
	bytecode.GetInstanceVariable: {
		name: bytecode.GetInstanceVariable,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			variableName := args[0].(string)
			v, ok := cf.self.(*RObject).InstanceVariables.get(variableName)
			if !ok {
				vm.stack.push(&Pointer{Target: NULL})
				return
			}

			p := &Pointer{Target: v}
			vm.stack.push(p)
		},
	},
	bytecode.SetInstanceVariable: {
		name: bytecode.SetInstanceVariable,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			variableName := args[0].(string)
			p := vm.stack.pop()
			cf.self.(*RObject).InstanceVariables.set(variableName, p.Target)
		},
	},
	bytecode.SetLocal: {
		name: bytecode.SetLocal,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			v := vm.stack.pop()
			depth := 0

			if len(args) >= 2 {
				depth = args[1].(int)
			}
			cf.insertLCL(args[0].(int), depth, v.Target)
		},
	},
	bytecode.SetConstant: {
		name: bytecode.SetConstant,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			constName := args[0].(string)
			v := vm.stack.pop()

			cf.storeConstant(constName, v)
		},
	},
	bytecode.NewArray: {
		name: bytecode.NewArray,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			argCount := args[0].(int)
			elems := []Object{}

			for i := 0; i < argCount; i++ {
				v := vm.stack.pop()
				elems = append([]Object{v.Target}, elems...)
			}

			arr := initializeArray(elems)
			vm.stack.push(&Pointer{arr})
		},
	},
	bytecode.NewHash: {
		name: bytecode.NewHash,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			argCount := args[0].(int)
			pairs := map[string]Object{}

			for i := 0; i < argCount/2; i++ {
				v := vm.stack.pop()
				k := vm.stack.pop()
				pairs[k.Target.(*StringObject).Value] = v.Target
			}

			hash := initializeHash(pairs)
			vm.stack.push(&Pointer{hash})
		},
	},
	bytecode.BranchUnless: {
		name: bytecode.BranchUnless,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			v := vm.stack.pop()
			bool, isBool := v.Target.(*BooleanObject)

			if isBool {
				if bool.Value {
					return
				}

				line := args[0].(int)
				cf.pc = line
				return
			}

			_, isNull := v.Target.(*Null)

			if isNull {
				line := args[0].(int)
				cf.pc = line
				return
			}
		},
	},
	bytecode.BranchIf: {
		name: bytecode.BranchIf,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			v := vm.stack.pop()
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
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			cf.pc = args[0].(int)
		},
	},
	bytecode.PutSelf: {
		name: bytecode.PutSelf,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			vm.stack.push(&Pointer{cf.self})
		},
	},
	bytecode.PutString: {
		name: bytecode.PutString,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			object := initializeObject(args[0])
			vm.stack.push(&Pointer{object})
		},
	},
	bytecode.PutNull: {
		name: bytecode.PutNull,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			vm.stack.push(&Pointer{NULL})
		},
	},
	bytecode.DefMethod: {
		name: bytecode.DefMethod,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			argCount := args[0].(int)
			methodName := vm.stack.pop().Target.(*StringObject).Value
			is, _ := vm.getMethodIS(methodName, cf.instructionSet.filename)
			method := &Method{Name: methodName, argc: argCount, instructionSet: is}

			v := vm.stack.pop().Target
			switch self := v.(type) {
			case *RClass:
				self.Methods.set(methodName, method)
			case BaseObject:
				self.returnClass().(*RClass).Methods.set(methodName, method)
			default:
				panic(fmt.Sprintf("Can't define method on %T", self))
			}
		},
	},
	bytecode.DefSingletonMethod: {
		name: bytecode.DefSingletonMethod,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			argCount := args[0].(int)
			methodName := vm.stack.pop().Target.(*StringObject).Value
			is, _ := vm.getMethodIS(methodName, cf.instructionSet.filename)
			method := &Method{Name: methodName, argc: argCount, instructionSet: is}

			v := vm.stack.pop().Target

			switch self := v.(type) {
			case *RClass:
				self.setSingletonMethod(methodName, method)
			case BaseObject:
				self.returnClass().(*RClass).setSingletonMethod(methodName, method)
			default:
				panic(fmt.Sprintf("Can't define singleton method on %T", self))
			}
		},
	},
	bytecode.DefClass: {
		name: bytecode.DefClass,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			subject := strings.Split(args[0].(string), ":")
			subjectType, subjectName := subject[0], subject[1]
			class := initializeClass(subjectName, subjectType == "module")
			classPr := cf.storeConstant(class.Name, class)

			is := vm.getClassIS(class.Name, cf.instructionSet.filename)

			if len(args) >= 2 {
				superClassName := args[1].(string)
				superClass := vm.lookupConstant(cf, superClassName)
				inheritedClass, ok := superClass.Target.(*RClass)

				if !ok {
					panic("Constant " + superClassName + " is not a class. got=" + string(superClass.Target.objectType()))
				}

				class.pseudoSuperClass = inheritedClass
				class.superClass = inheritedClass
			}

			vm.stack.pop()
			c := newCallFrame(is)
			c.self = class
			vm.callFrameStack.push(c)
			vm.startFromTopFrame()

			vm.stack.push(classPr)
		},
	},
	bytecode.Send: {
		name: bytecode.Send,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			methodName := args[0].(string)
			argCount := args[1].(int)
			var blockName string
			var hasBlock bool

			if len(args) > 2 {
				hasBlock = true
				blockFlag := args[2].(string)
				blockName = strings.Split(blockFlag, ":")[1]
			} else {
				hasBlock = false
			}

			argPr := vm.sp - argCount
			receiverPr := argPr - 1
			receiver := vm.stack.Data[receiverPr].Target.(BaseObject)

			error := &Error{Message: "undefined method `" + methodName + "' for " + receiver.Inspect()}

			var method Object

			switch receiver := receiver.(type) {
			case Class:
				method = receiver.lookupClassMethod(methodName)
			case BaseObject:
				method = receiver.returnClass().lookupInstanceMethod(methodName)
			case *Error:
				panic(receiver.Inspect())
			default:
				panic(fmt.Sprintf("not a valid receiver: %s", receiver.Inspect()))
			}

			if method == nil {
				panic(error.Message)
			}

			var blockFrame *callFrame

			if hasBlock {
				block := vm.getBlock(blockName, cf.instructionSet.filename)

				c := newCallFrame(block)
				c.isBlock = true
				c.ep = cf
				c.self = cf.self

				vm.callFrameStack.push(c)
				blockFrame = c
			}

			switch m := method.(type) {
			case *Method:
				evalMethodObject(vm, receiver, m, receiverPr, argCount, argPr, blockFrame)
			case *BuiltInMethod:
				evalBuiltInMethod(vm, receiver, m, receiverPr, argCount, argPr, blockFrame)
			case *Error:
				panic(m.Inspect())
			default:
				panic(fmt.Sprintf("unknown instance method type: %T", m))
			}
		},
	},
	bytecode.InvokeBlock: {
		name: bytecode.InvokeBlock,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			argCount := args[0].(int)
			argPr := vm.sp - argCount
			receiverPr := argPr - 1
			receiver := vm.stack.Data[receiverPr].Target.(BaseObject)

			if cf.blockFrame == nil {
				panic("Can't yield without a block")
			}

			c := newCallFrame(cf.blockFrame.instructionSet)
			c.blockFrame = cf.blockFrame
			c.ep = cf.blockFrame.ep
			c.self = receiver

			for i := 0; i < argCount; i++ {
				c.locals[i] = vm.stack.Data[argPr+i]
			}

			vm.callFrameStack.push(c)
			vm.startFromTopFrame()

			vm.stack.Data[receiverPr] = vm.stack.top()
			vm.sp = receiverPr + 1
		},
	},
	bytecode.Leave: {
		name: bytecode.Leave,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			cf = vm.callFrameStack.pop()
			cf.pc = len(cf.instructionSet.instructions)
		},
	},
}

func evalBuiltInMethod(vm *VM, receiver BaseObject, method *BuiltInMethod, receiverPr, argCount, argPr int, blockFrame *callFrame) {
	methodBody := method.Fn(receiver)
	args := []Object{}

	for i := 0; i < argCount; i++ {
		args = append(args, vm.stack.Data[argPr+i].Target)
	}

	evaluated := methodBody(vm, args, blockFrame)

	_, ok := receiver.(*RClass)
	if method.Name == "new" && ok {
		instance := evaluated.(*RObject)
		if instance.InitializeMethod != nil {
			evalMethodObject(vm, instance, instance.InitializeMethod, receiverPr, argCount, argPr, blockFrame)
		}
	}
	vm.stack.Data[receiverPr] = &Pointer{evaluated}
	vm.sp = receiverPr + 1
}

func evalMethodObject(vm *VM, receiver BaseObject, method *Method, receiverPr, argC, argPr int, blockFrame *callFrame) {
	c := newCallFrame(method.instructionSet)
	c.self = receiver

	for i := 0; i < argC; i++ {
		c.insertLCL(i, 0, vm.stack.Data[argPr+i].Target)
	}

	c.blockFrame = blockFrame
	vm.callFrameStack.push(c)
	vm.startFromTopFrame()

	vm.stack.Data[receiverPr] = vm.stack.top()
	vm.sp = receiverPr + 1
}

func (is *instructionSet) define(line int, a *action, params ...interface{}) {
	i := &instruction{action: a, Params: params, Line: line}
	is.instructions = append(is.instructions, i)
}

func (is *instructionSet) inspect() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("<%s>\n", is.label.name))
	for _, i := range is.instructions {
		out.WriteString(i.inspect())
		out.WriteString("\n")
	}

	return out.String()
}

func (i *instruction) inspect() string {
	var params []string

	for _, param := range i.Params {
		params = append(params, fmt.Sprint(param))
	}
	return fmt.Sprintf("%s: %s", i.action.name, strings.Join(params, ", "))
}

func initializeObject(value interface{}) Object {
	switch v := value.(type) {
	case int:
		return initilaizeInteger(int(v))
	case int64:
		return initilaizeInteger(int(v))
	case string:
		switch v {
		case "true":
			return TRUE
		case "false":
			return FALSE
		case "nil":
			return NULL
		default:
			return initializeString(v)
		}
	default:
		panic(fmt.Sprintf("Unknown data type: %T", v))
	}
}
