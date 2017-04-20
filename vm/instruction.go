package vm

import (
	"bytes"
	"fmt"
	"strings"
)

type operation func(vm *VM, cf *CallFrame, args ...interface{})

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
)

var builtInActions = map[operationType]*action{
	Pop: {
		name: Pop,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			vm.Stack.pop()
		},
	},
	PutObject: {
		name: PutObject,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			object := initializeObject(args[0])
			vm.Stack.push(&Pointer{Target: object})
		},
	},
	GetConstant: {
		name: GetConstant,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			constName := args[0].(string)
			constant, ok := vm.Constants[constName]

			if !ok {
				panic(fmt.Sprintf("Can't find constant: %s", constName))
			}
			vm.Stack.push(constant)
		},
	},
	GetLocal: {
		name: GetLocal,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			index := args[0].(int)
			depth := 0

			if len(args) >= 2 {
				depth = args[1].(int)
			}

			p := cf.getLCL(index, depth)

			if p == nil {
				panic(fmt.Sprintf("Local index: %d is nil. Callframe: %s", index, cf.instructionSet.label.name))
			}
			vm.Stack.push(p)
		},
	},
	GetInstanceVariable: {
		name: GetInstanceVariable,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			variableName := args[0].(string)
			v, ok := cf.Self.(*RObject).InstanceVariables.Get(variableName)
			if !ok {
				vm.Stack.push(&Pointer{Target: NULL})
				return
			}

			p := &Pointer{Target: v}
			vm.Stack.push(p)
		},
	},
	SetInstanceVariable: {
		name: SetInstanceVariable,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			variableName := args[0].(string)
			p := vm.Stack.pop()
			cf.Self.(*RObject).InstanceVariables.Set(variableName, p.Target)
		},
	},
	SetLocal: {
		name: SetLocal,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			v := vm.Stack.pop()
			depth := 0

			if len(args) >= 2 {
				depth = args[1].(int)
			}
			cf.insertLCL(args[0].(int), depth, v.Target)
		},
	},
	SetConstant: {
		name: SetConstant,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			constName := args[0].(string)
			v := vm.Stack.pop()
			vm.Constants[constName] = v
		},
	},
	NewArray: {
		name: NewArray,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			argCount := args[0].(int)
			elems := []Object{}

			for i := 0; i < argCount; i++ {
				v := vm.Stack.pop()
				elems = append([]Object{v.Target}, elems...)
			}

			arr := InitializeArray(elems)
			vm.Stack.push(&Pointer{arr})
		},
	},
	NewHash: {
		name: NewHash,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			argCount := args[0].(int)
			pairs := map[string]Object{}

			for i := 0; i < argCount/2; i++ {
				v := vm.Stack.pop()
				k := vm.Stack.pop()
				pairs[k.Target.(*StringObject).Value] = v.Target
			}

			hash := InitializeHash(pairs)
			vm.Stack.push(&Pointer{hash})
		},
	},
	BranchUnless: {
		name: BranchUnless,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			v := vm.Stack.pop()
			bool, isBool := v.Target.(*BooleanObject)

			if isBool {
				if bool.Value {
					return
				}

				line := args[0].(int)
				cf.PC = line
				return
			}

			_, isNull := v.Target.(*Null)

			if isNull {
				line := args[0].(int)
				cf.PC = line
				return
			}
		},
	},
	Jump: {
		name: Jump,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			cf.PC = args[0].(int)
		},
	},
	PutSelf: {
		name: PutSelf,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			vm.Stack.push(&Pointer{cf.Self})
		},
	},
	PutString: {
		name: PutString,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			object := initializeObject(args[0])
			vm.Stack.push(&Pointer{object})
		},
	},
	PutNull: {
		name: PutNull,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			vm.Stack.push(&Pointer{NULL})
		},
	},
	DefMethod: {
		name: DefMethod,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			argCount := args[0].(int)
			methodName := vm.Stack.pop().Target.(*StringObject).Value
			is, _ := vm.getMethodIS(methodName)
			method := &Method{Name: methodName, Argc: argCount, instructionSet: is}

			v := vm.Stack.pop().Target
			switch self := v.(type) {
			case *RClass:
				self.Methods.Set(methodName, method)
			case BaseObject:
				self.ReturnClass().(*RClass).Methods.Set(methodName, method)
			default:
				panic(fmt.Sprintf("Can't define method on %T", self))
			}
		},
	},
	DefSingletonMethod: {
		name: DefSingletonMethod,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			argCount := args[0].(int)
			methodName := vm.Stack.pop().Target.(*StringObject).Value
			is, _ := vm.getMethodIS(methodName)
			method := &Method{Name: methodName, Argc: argCount, instructionSet: is}

			v := vm.Stack.pop().Target

			switch self := v.(type) {
			case *RClass:
				self.SetSingletonMethod(methodName, method)
			case BaseObject:
				self.ReturnClass().(*RClass).SetSingletonMethod(methodName, method)
			default:
				panic(fmt.Sprintf("Can't define singleton method on %T", self))
			}
		},
	},
	DefClass: {
		name: DefClass,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			class := InitializeClass(args[0].(string))
			classPr := &Pointer{Target: class}
			vm.Constants[class.Name] = classPr

			is, ok := vm.getClassIS(class.Name)

			if !ok {
				panic(fmt.Sprintf("Can't find class %s's instructions", class.Name))
			}

			if len(args) >= 2 {
				constantName := args[1].(string)
				constant := vm.Constants[constantName]
				inheritedClass, ok := constant.Target.(*RClass)
				if !ok {
					newError("Constant %s is not a class. got=%T", constantName, constant)
				}

				class.SuperClass = inheritedClass
			}

			vm.Stack.pop()
			c := NewCallFrame(is)
			c.Self = class
			vm.CallFrameStack.Push(c)
			vm.Exec()

			vm.Stack.push(classPr)
		},
	},
	Send: {
		name: Send,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
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

			argPr := vm.SP - argCount
			receiverPr := argPr - 1
			receiver := vm.Stack.Data[receiverPr].Target.(BaseObject)

			error := newError("undefined method `%s' for %s", methodName, receiver.Inspect())

			var method Object

			switch receiver := receiver.(type) {
			case Class:
				method = receiver.LookupClassMethod(methodName)
			case BaseObject:
				method = receiver.ReturnClass().LookupInstanceMethod(methodName)
			case *Error:
				panic(receiver.Inspect())
			default:
				panic(fmt.Sprintf("not a valid receiver: %s", receiver.Inspect()))
			}

			if method == nil {
				panic(error.Message)
			}

			var blockFrame *CallFrame

			if hasBlock {
				block, ok := vm.getBlock(blockName)

				if !ok {
					panic(fmt.Sprintf("Can't find block %s", blockName))
				}

				c := NewCallFrame(block)
				c.IsBlock = true
				c.EP = cf
				vm.CallFrameStack.Push(c)
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
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			argCount := args[0].(int)
			argPr := vm.SP - argCount
			receiverPr := argPr - 1
			receiver := vm.Stack.Data[receiverPr].Target.(BaseObject)

			if cf.BlockFrame == nil {
				panic("Can't yield without a block")
			}

			c := NewCallFrame(cf.BlockFrame.instructionSet)
			c.BlockFrame = cf.BlockFrame
			c.EP = cf.BlockFrame.EP
			c.Self = receiver

			for i := 0; i < argCount; i++ {
				c.Local[i] = vm.Stack.Data[argPr+i]
			}

			vm.CallFrameStack.Push(c)
			vm.Exec()

			setReturnValueAndSP(vm, receiverPr, vm.Stack.Top())
		},
	},
	Leave: {
		name: Leave,
		operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			cf = vm.CallFrameStack.Pop()
			cf.PC = len(cf.instructionSet.instructions)
		},
	},
}

func evalBuiltInMethod(vm *VM, receiver BaseObject, method *BuiltInMethod, receiverPr, argCount, argPr int, blockFrame *CallFrame) {
	methodBody := method.Fn(receiver)
	args := []Object{}

	for i := 0; i < argCount; i++ {
		args = append(args, vm.Stack.Data[argPr+i].Target)
	}

	evaluated := methodBody(args, nil)

	_, ok := receiver.(*RClass)
	if method.Name == "new" && ok {
		instance := evaluated.(*RObject)
		if instance.InitializeMethod != nil {
			evalMethodObject(vm, instance, instance.InitializeMethod, receiverPr, argCount, argPr, blockFrame)
		}
	}
	setReturnValueAndSP(vm, receiverPr, &Pointer{evaluated})
}

func evalMethodObject(vm *VM, receiver BaseObject, method *Method, receiverPr, argC, argPr int, blockFrame *CallFrame) {
	c := NewCallFrame(method.instructionSet)
	c.Self = receiver

	for i := 0; i < argC; i++ {
		c.insertLCL(i, 0, vm.Stack.Data[argPr+i].Target)
	}

	c.BlockFrame = blockFrame
	vm.CallFrameStack.Push(c)
	vm.Exec()

	setReturnValueAndSP(vm, receiverPr, vm.Stack.Top())
}

func setReturnValueAndSP(vm *VM, receiverPr int, value *Pointer) {
	vm.Stack.Data[receiverPr] = value
	vm.SP = receiverPr + 1
}

func (is *instructionSet) Define(line int, a *action, params ...interface{}) {
	i := &instruction{action: a, Params: params, Line: line}
	is.instructions = append(is.instructions, i)
}

func (is *instructionSet) Inspect() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("<%s>\n", is.label.name))
	for _, i := range is.instructions {
		out.WriteString(i.Inspect())
		out.WriteString("\n")
	}

	return out.String()
}

func (i *instruction) Inspect() string {
	var params []string

	for _, param := range i.Params {
		params = append(params, fmt.Sprint(param))
	}
	return fmt.Sprintf("%s: %s \n", i.action.name, strings.Join(params, ", "))
}

func initializeObject(value interface{}) Object {
	switch v := value.(type) {
	case int:
		return InitilaizeInteger(int(v))
	case int64:
		return InitilaizeInteger(int(v))
	case string:
		switch v {
		case "true":
			return TRUE
		case "false":
			return FALSE
		case "nil":
			return NULL
		default:
			return InitializeString(v)
		}
	default:
		panic(fmt.Sprintf("Unknown data type: %T", v))
	}
}
