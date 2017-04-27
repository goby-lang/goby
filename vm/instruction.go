package vm

import (
	"bytes"
	"fmt"
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

var labelTypes = map[string]labelType{
	"Def":          LabelDef,
	"DefClass":     LabelDefClass,
	"ProgramStart": Program,
	"Block":        Block,
}

type instructionSet struct {
	label        *label
	instructions []*instruction
}

type operationType string

const (
	// label types
	LabelDef      = "DefMethod"
	LabelDefClass = "DefClass"
	Block         = "Block"
	Program       = "Program"

	// instruction actions
	GetLocal            = "getlocal"
	GetConstant         = "getconstant"
	GetInstanceVariable = "getinstancevariable"
	SetLocal            = "setlocal"
	SetConstant         = "setconstant"
	SetInstanceVariable = "setinstancevariable"
	PutString           = "putstring"
	PutSelf             = "putself"
	PutObject           = "putobject"
	PutNull             = "putnil"
	NewArray            = "newarray"
	NewHash             = "newhash"
	BranchUnless        = "branchunless"
	Jump                = "jump"
	DefMethod           = "def_method"
	DefSingletonMethod  = "def_singleton_method"
	DefClass            = "def_class"
	Send                = "send"
	InvokeBlock         = "invokeblock"
	Pop                 = "pop"
	Leave               = "leave"

	RequireRelative = "require_relative"
)

var builtInActions = map[operationType]*action{
	Pop: {
		name: Pop,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			vm.stack.pop()
		},
	},
	PutObject: {
		name: PutObject,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			object := initializeObject(args[0])
			vm.stack.push(&Pointer{Target: object})
		},
	},
	GetConstant: {
		name: GetConstant,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			constName := args[0].(string)
			constant, ok := vm.constants[constName]

			if !ok {
				panic(fmt.Sprintf("Can't find constant: %s", constName))
			}
			vm.stack.push(constant)
		},
	},
	GetLocal: {
		name: GetLocal,
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
	GetInstanceVariable: {
		name: GetInstanceVariable,
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
	SetInstanceVariable: {
		name: SetInstanceVariable,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			variableName := args[0].(string)
			p := vm.stack.pop()
			cf.self.(*RObject).InstanceVariables.set(variableName, p.Target)
		},
	},
	SetLocal: {
		name: SetLocal,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			v := vm.stack.pop()
			depth := 0

			if len(args) >= 2 {
				depth = args[1].(int)
			}
			cf.insertLCL(args[0].(int), depth, v.Target)
		},
	},
	SetConstant: {
		name: SetConstant,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			constName := args[0].(string)
			v := vm.stack.pop()
			vm.constants[constName] = v
		},
	},
	NewArray: {
		name: NewArray,
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
	NewHash: {
		name: NewHash,
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
	BranchUnless: {
		name: BranchUnless,
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
	Jump: {
		name: Jump,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			cf.pc = args[0].(int)
		},
	},
	PutSelf: {
		name: PutSelf,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			vm.stack.push(&Pointer{cf.self})
		},
	},
	PutString: {
		name: PutString,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			object := initializeObject(args[0])
			vm.stack.push(&Pointer{object})
		},
	},
	PutNull: {
		name: PutNull,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			vm.stack.push(&Pointer{NULL})
		},
	},
	DefMethod: {
		name: DefMethod,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			argCount := args[0].(int)
			methodName := vm.stack.pop().Target.(*StringObject).Value
			is, _ := vm.getMethodIS(methodName)
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
	DefSingletonMethod: {
		name: DefSingletonMethod,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			argCount := args[0].(int)
			methodName := vm.stack.pop().Target.(*StringObject).Value
			is, _ := vm.getMethodIS(methodName)
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
	DefClass: {
		name: DefClass,
		operation: func(vm *VM, cf *callFrame, args ...interface{}) {
			class := initializeClass(args[0].(string))
			classPr := &Pointer{Target: class}
			vm.constants[class.Name] = classPr

			is, ok := vm.getClassIS(class.Name)

			if !ok {
				panic(fmt.Sprintf("Can't find class %s's instructions", class.Name))
			}

			if len(args) >= 2 {
				constantName := args[1].(string)
				constant := vm.constants[constantName]
				inheritedClass, ok := constant.Target.(*RClass)
				if !ok {
					newError("Constant %s is not a class. got=%T", constantName, constant)
				}

				class.SuperClass = inheritedClass
			}

			vm.stack.pop()
			c := newCallFrame(is)
			c.self = class
			vm.callFrameStack.push(c)
			vm.start()

			vm.stack.push(classPr)
		},
	},
	Send: {
		name: Send,
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

			error := newError("undefined method `%s' for %s", methodName, receiver.Inspect())

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
				block, ok := vm.getBlock(blockName)

				if !ok {
					panic(fmt.Sprintf("Can't find block %s", blockName))
				}

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
	InvokeBlock: {
		name: InvokeBlock,
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
			vm.start()

			setReturnValueAndSP(vm, receiverPr, vm.stack.top())
		},
	},
	Leave: {
		name: Leave,
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
	setReturnValueAndSP(vm, receiverPr, &Pointer{evaluated})
}

func evalMethodObject(vm *VM, receiver BaseObject, method *Method, receiverPr, argC, argPr int, blockFrame *callFrame) {
	c := newCallFrame(method.instructionSet)
	c.self = receiver

	for i := 0; i < argC; i++ {
		c.insertLCL(i, 0, vm.stack.Data[argPr+i].Target)
	}

	c.blockFrame = blockFrame
	vm.callFrameStack.push(c)
	vm.start()

	setReturnValueAndSP(vm, receiverPr, vm.stack.top())
}

func setReturnValueAndSP(vm *VM, receiverPr int, value *Pointer) {
	vm.stack.Data[receiverPr] = value
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
	return fmt.Sprintf("%s: %s \n", i.action.name, strings.Join(params, ", "))
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
