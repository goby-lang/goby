package vm

import (
	"bytes"
	"fmt"
	"strings"
)

type Operation func(vm *VM, cf *CallFrame, args ...interface{})

type Action struct {
	Name      string
	Operation Operation
}

type Instruction struct {
	Action *Action
	Params []interface{}
	Line   int
}

type Label struct {
	Name string
	Type LabelType
}

type LabelType string

var labelTypes = map[string]LabelType{
	"Def":          LABEL_DEF,
	"DefClass":     LABEL_DEFCLASS,
	"ProgramStart": PROGRAM,
	"Block":        BLOCK,
}

type InstructionSet struct {
	Label        *Label
	Instructions []*Instruction
}

type OperationType string

const (
	// Label types
	LABEL_DEF      = "DefMethod"
	LABEL_DEFCLASS = "DefClass"
	BLOCK          = "Block"
	PROGRAM        = "Program"

	// Instruction actions
	GET_LOCAL             = "getlocal"
	GET_CONSTANT          = "getconstant"
	GET_INSTANCE_VARIABLE = "getinstancevariable"
	SET_LOCAL             = "setlocal"
	SET_CONSTANT          = "setconstant"
	SET_INSTANCE_VARIABLE = "setinstancevariable"
	PUT_STRING            = "putstring"
	PUT_SELF              = "putself"
	PUT_OBJECT            = "putobject"
	PUT_NULL              = "putnil"
	NEW_ARRAY             = "newarray"
	NEW_HASH              = "newhash"
	PLUS                  = "opt_plus"
	MINUS                 = "opt_minus"
	MULT                  = "opt_mult"
	DIV                   = "opt_div"
	GT                    = "opt_gt"
	GE                    = "opt_ge"
	LT                    = "opt_lt"
	LE                    = "opt_le"
	BRANCH_UNLESS         = "branchunless"
	JUMP                  = "jump"
	DEF_METHOD            = "def_method"
	DEF_SINGLETON_METHOD  = "def_singleton_method"
	DEF_CLASS             = "def_class"
	SEND                  = "send"
	INVOKE_BLOCK          = "invokeblock"
	POP                   = "pop"
	LEAVE                 = "leave"
)

var BuiltInActions = map[OperationType]*Action{
	POP: {
		Name: POP,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			vm.Stack.pop()
		},
	},
	PUT_OBJECT: {
		Name: PUT_OBJECT,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			object := initializeObject(args[0])
			vm.Stack.push(&Pointer{Target: object})
		},
	},
	GET_CONSTANT: {
		Name: GET_CONSTANT,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			constName := args[0].(string)
			constant, ok := vm.Constants[constName]

			if !ok {
				panic(fmt.Sprintf("Can't find constant: %s", constName))
			}
			vm.Stack.push(constant)
		},
	},
	GET_LOCAL: {
		Name: GET_LOCAL,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			i := args[0].(int)
			p := cf.getLCL(i)

			if p == nil {
				panic(fmt.Sprintf("Local index: %d is nil. Callframe: %s", i, cf.InstructionSet.Label.Name))
			}
			vm.Stack.push(p)
		},
	},
	GET_INSTANCE_VARIABLE: {
		Name: GET_INSTANCE_VARIABLE,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
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
	SET_INSTANCE_VARIABLE: {
		Name: SET_INSTANCE_VARIABLE,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			variableName := args[0].(string)
			p := vm.Stack.pop()
			cf.Self.(*RObject).InstanceVariables.Set(variableName, p.Target)
		},
	},
	SET_LOCAL: {
		Name: SET_LOCAL,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			v := vm.Stack.pop()
			cf.insertLCL(args[0].(int), v.Target)
		},
	},
	SET_CONSTANT: {
		Name: SET_CONSTANT,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			constName := args[0].(string)
			v := vm.Stack.pop()
			vm.Constants[constName] = v
		},
	},
	NEW_ARRAY: {
		Name: NEW_ARRAY,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
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
	NEW_HASH: {
		Name: NEW_HASH,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
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
	BRANCH_UNLESS: {
		Name: BRANCH_UNLESS,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
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
	JUMP: {
		Name: JUMP,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			cf.PC = args[0].(int)
		},
	},
	PUT_SELF: {
		Name: PUT_SELF,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			vm.Stack.push(&Pointer{cf.Self})
		},
	},
	PUT_STRING: {
		Name: PUT_STRING,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			object := initializeObject(args[0])
			vm.Stack.push(&Pointer{object})
		},
	},
	PUT_NULL: {
		Name: PUT_NULL,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			vm.Stack.push(&Pointer{NULL})
		},
	},
	DEF_METHOD: {
		Name: DEF_METHOD,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			argCount := args[0].(int)
			methodName := vm.Stack.pop().Target.(*StringObject).Value
			is, _ := vm.getMethodIS(methodName)
			method := &Method{Name: methodName, Argc: argCount, InstructionSet: is}

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
	DEF_SINGLETON_METHOD: {
		Name: DEF_SINGLETON_METHOD,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			argCount := args[0].(int)
			methodName := vm.Stack.pop().Target.(*StringObject).Value
			is, _ := vm.getMethodIS(methodName)
			method := &Method{Name: methodName, Argc: argCount, InstructionSet: is}

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
	DEF_CLASS: {
		Name: DEF_CLASS,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
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
	SEND: {
		Name: SEND,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
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
	INVOKE_BLOCK: {
		Name: INVOKE_BLOCK,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			argCount := args[0].(int)
			argPr := vm.SP - argCount
			receiverPr := argPr - 1
			receiver := vm.Stack.Data[receiverPr].Target.(BaseObject)

			if cf.BlockFrame == nil {
				panic("Can't yield without a block")
			}

			c := NewCallFrame(cf.BlockFrame.InstructionSet)
			c.BlockFrame = cf.BlockFrame
			c.EP = cf.BlockFrame.EP
			c.Self = receiver

			for i := 0; i < argCount; i++ {
				c.Local[i] = vm.Stack.Data[argPr+i]
			}

			fmt.Println(c.inspect())
			fmt.Println(c.Local)
			fmt.Println(vm.Stack.Data[argPr].Target)

			vm.CallFrameStack.Push(c)
			vm.Exec()

			setReturnValueAndSP(vm, receiverPr, vm.Stack.Top())
		},
	},
	LEAVE: {
		Name: LEAVE,
		Operation: func(vm *VM, cf *CallFrame, args ...interface{}) {
			cf = vm.CallFrameStack.Pop()
			cf.PC = len(cf.InstructionSet.Instructions)
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
	c := NewCallFrame(method.InstructionSet)
	c.Self = receiver

	for i := 0; i < argC; i++ {
		c.insertLCL(i, vm.Stack.Data[argPr+i].Target)
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

func (is *InstructionSet) Define(line int, action *Action, params ...interface{}) {
	i := &Instruction{Action: action, Params: params, Line: line}
	is.Instructions = append(is.Instructions, i)
}

func (is *InstructionSet) Inspect() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("<%s>\n", is.Label.Name))
	for _, i := range is.Instructions {
		out.WriteString(i.Inspect())
		out.WriteString("\n")
	}

	return out.String()
}

func (i *Instruction) Inspect() string {
	var params []string

	for _, param := range i.Params {
		params = append(params, fmt.Sprint(param))
	}
	return fmt.Sprintf("%s: %s \n", i.Action.Name, strings.Join(params, ", "))
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
