package vm

import (
	"bytes"
	"fmt"
	"strings"
)

type Operation func(vm *VM, cf *CallFrame, args ...Object)

type Action struct {
	Name      string
	Operation Operation
}

type Instruction struct {
	Action *Action
	Params []Object
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
	SET_INSTANCE_VARIABLE = "setinstancevariable"
	PUT_STRING            = "putstring"
	PUT_SELF              = "putself"
	PUT_OBJECT            = "putobject"
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
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			vm.Stack.pop()
		},
	},
	PUT_OBJECT: {
		Name: PUT_OBJECT,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			vm.Stack.push(args[0])
		},
	},
	GET_CONSTANT: {
		Name: GET_CONSTANT,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			constName := args[0].(*StringObject).Value
			constant, ok := vm.Constants[constName]

			if !ok {
				panic(fmt.Sprintf("Can't find constant: %s", constName))
			}
			vm.Stack.push(constant)
		},
	},
	GET_LOCAL: {
		Name: GET_LOCAL,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			i := args[0].(*IntegerObject)
			v := cf.Local[i.Value]

			if v == nil {
				panic(fmt.Sprintf("Local index: %d is nil. Callframe: %s", i.Value, cf.InstructionSet.Label.Name))
			}
			vm.Stack.push(v)
		},
	},
	GET_INSTANCE_VARIABLE: {
		Name: GET_INSTANCE_VARIABLE,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			variableName := args[0].(*StringObject).Value
			v, ok := cf.Self.(*RObject).InstanceVariables.Get(variableName)
			if !ok {
				vm.Stack.push(NULL)
				return
			}
			vm.Stack.push(v)
		},
	},
	SET_INSTANCE_VARIABLE: {
		Name: SET_INSTANCE_VARIABLE,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			variableName := args[0].(*StringObject).Value
			value := vm.Stack.pop()
			cf.Self.(*RObject).InstanceVariables.Set(variableName, value)
		},
	},
	SET_LOCAL: {
		Name: SET_LOCAL,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			v := vm.Stack.pop()
			cf.insertLCL(args[0].(*IntegerObject).Value, v)
		},
	},
	NEW_ARRAY: {
		Name: NEW_ARRAY,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			argCount := args[0].(*IntegerObject).Value
			elems := []Object{}

			for i := 0; i < argCount; i++ {
				v := vm.Stack.pop()
				elems = append([]Object{v}, elems...)
			}

			arr := InitializeArray(elems)
			vm.Stack.push(arr)
		},
	},
	NEW_HASH: {
		Name: NEW_HASH,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			argCount := args[0].(*IntegerObject).Value
			pairs := map[string]Object{}

			for i := 0; i < argCount/2; i++ {
				v := vm.Stack.pop()
				k := vm.Stack.pop()
				pairs[k.(*StringObject).Value] = v
			}

			hash := InitializeHash(pairs)
			vm.Stack.push(hash)
		},
	},
	PLUS: {
		Name: PLUS,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			second := vm.Stack.Top().(*IntegerObject).Value
			vm.SP -= 1
			first := vm.Stack.Top().(*IntegerObject).Value
			vm.SP -= 1
			vm.Stack.push(InitilaizeInteger(first + second))
		},
	},
	MINUS: {
		Name: MINUS,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			second := vm.Stack.Top().(*IntegerObject).Value
			vm.SP -= 1
			first := vm.Stack.Top().(*IntegerObject).Value
			vm.SP -= 1
			vm.Stack.push(InitilaizeInteger(first - second))
		},
	},
	MULT: {
		Name: MULT,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			second := vm.Stack.Top().(*IntegerObject).Value
			vm.SP -= 1
			first := vm.Stack.Top().(*IntegerObject).Value
			vm.SP -= 1
			vm.Stack.push(InitilaizeInteger(first * second))
		},
	},
	DIV: {
		Name: DIV,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			second := vm.Stack.Top().(*IntegerObject).Value
			vm.SP -= 1
			first := vm.Stack.Top().(*IntegerObject).Value
			vm.SP -= 1
			vm.Stack.push(InitilaizeInteger(first / second))
		},
	},
	GT: {
		Name: GT,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			second := vm.Stack.Top().(*IntegerObject).Value
			vm.SP -= 1
			first := vm.Stack.Top().(*IntegerObject).Value
			vm.SP -= 1
			result := first > second

			if result {
				vm.Stack.push(TRUE)
			} else {
				vm.Stack.push(FALSE)
			}
		},
	},
	LT: {
		Name: LT,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			second := vm.Stack.Top().(*IntegerObject).Value
			vm.SP -= 1
			first := vm.Stack.Top().(*IntegerObject).Value
			vm.SP -= 1
			result := first < second

			if result {
				vm.Stack.push(TRUE)
			} else {
				vm.Stack.push(FALSE)
			}
		},
	},
	GE: {
		Name: GE,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			second := vm.Stack.Top().(*IntegerObject).Value
			vm.SP -= 1
			first := vm.Stack.Top().(*IntegerObject).Value
			vm.SP -= 1
			result := first >= second

			if result {
				vm.Stack.push(TRUE)
			} else {
				vm.Stack.push(FALSE)
			}
		},
	},
	LE: {
		Name: LE,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			second := vm.Stack.Top().(*IntegerObject).Value
			vm.SP -= 1
			first := vm.Stack.Top().(*IntegerObject).Value
			vm.SP -= 1
			result := first <= second

			if result {
				vm.Stack.push(TRUE)
			} else {
				vm.Stack.push(FALSE)
			}
		},
	},
	BRANCH_UNLESS: {
		Name: BRANCH_UNLESS,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			cond := vm.Stack.pop().(*BooleanObject).Value
			if cond {
				return
			}

			line := args[0].(*IntegerObject).Value
			cf.PC = line
		},
	},
	JUMP: {
		Name: JUMP,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			cf.PC = args[0].(*IntegerObject).Value
		},
	},
	PUT_SELF: {
		Name: PUT_SELF,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			vm.Stack.push(cf.Self)
		},
	},
	PUT_STRING: {
		Name: PUT_STRING,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			vm.Stack.push(args[0])
		},
	},
	DEF_METHOD: {
		Name: DEF_METHOD,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			argCount := args[0].(*IntegerObject).Value
			methodName := vm.Stack.pop().(*StringObject).Value
			is, _ := vm.getMethodIS(methodName)
			method := &Method{Name: methodName, Argc: argCount, InstructionSet: is}

			v := vm.Stack.pop()
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
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			argCount := args[0].(*IntegerObject).Value
			methodName := vm.Stack.pop().(*StringObject).Value
			is, _ := vm.getMethodIS(methodName)
			method := &Method{Name: methodName, Argc: argCount, InstructionSet: is}

			v := vm.Stack.pop()

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
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			class := InitializeClass(args[0].(*StringObject).Value)
			vm.Constants[class.Name] = class

			is, ok := vm.getClassIS(class.Name)

			if !ok {
				panic(fmt.Sprintf("Can't find class %s's instructions", class.Name))
			}

			if len(args) >= 2 {
				constantName := args[1].(*StringObject).Value
				constant := vm.Constants[constantName]
				inheritedClass, ok := constant.(*RClass)
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

			vm.Stack.push(class)
		},
	},
	SEND: {
		Name: SEND,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			methodName := args[0].(*StringObject).Value
			argCount := args[1].(*IntegerObject).Value
			argPr := vm.SP - argCount
			receiverPr := argPr - 1
			receiver := vm.Stack.Data[receiverPr].(BaseObject)

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

			switch m := method.(type) {
			case *Method:
				evalMethodObject(vm, receiver, m, receiverPr, argCount, argPr)
			case *BuiltInMethod:
				evalBuiltInMethod(vm, receiver, m, receiverPr, argCount, argPr)
			case *Error:
				panic(m.Inspect())
			default:
				panic(fmt.Sprintf("unknown instance method type: %T", m))
			}
		},
	},
	INVOKE_BLOCK: {
		Name: INVOKE_BLOCK,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			argCount := args[0].(*IntegerObject).Value
			argPr := vm.SP - argCount
			receiverPr := argPr - 1
			receiver := vm.Stack.Data[receiverPr].(BaseObject)

			block, _ := vm.getBlock()

			c := NewCallFrame(block)
			c.Self = receiver
			c.ArgPr = argPr
			copy(c.Local, cf.Local)

			vm.CallFrameStack.Push(c)
			vm.Exec()

			setReturnValueAndSP(vm, receiverPr, vm.Stack.Top())
		},
	},
	LEAVE: {
		Name: LEAVE,
		Operation: func(vm *VM, cf *CallFrame, args ...Object) {
			vm.CallFrameStack.Pop()
		},
	},
}

func evalBuiltInMethod(vm *VM, receiver BaseObject, method *BuiltInMethod, receiverPr, argCount, argPr int) {
	methodBody := method.Fn(receiver)
	args := []Object{}

	for i := 0; i < argCount; i++ {
		args = append(args, vm.Stack.Data[argPr+i])
	}

	evaluated := methodBody(args, nil)

	_, ok := receiver.(*RClass)
	if method.Name == "new" && ok {
		instance := evaluated.(*RObject)
		if instance.InitializeMethod != nil {
			evalMethodObject(vm, instance, instance.InitializeMethod, receiverPr, argCount, argPr)
		}
	}
	setReturnValueAndSP(vm, receiverPr, evaluated)
}

func evalMethodObject(vm *VM, receiver BaseObject, method *Method, receiverPr, argC, argPr int) {
	c := NewCallFrame(method.InstructionSet)
	c.Self = receiver

	for i := 0; i < argC; i++ {
		c.insertLCL(i, vm.Stack.Data[argPr+i])
	}

	vm.CallFrameStack.Push(c)
	vm.Exec()

	setReturnValueAndSP(vm, receiverPr, vm.Stack.Top())
}

func setReturnValueAndSP(vm *VM, receiverPr int, value Object) {
	vm.Stack.Data[receiverPr] = value
	vm.SP = receiverPr + 1
}

func (is *InstructionSet) Define(line int, action *Action, params ...interface{}) {
	ps := []Object{}

	for _, param := range params {
		var p Object

		switch param := param.(type) {
		case int:
			p = InitilaizeInteger(int(param))
		case int64:
			p = InitilaizeInteger(int(param))
		case string:
			switch param {
			case "true":
				p = TRUE
			case "false":
				p = FALSE
			case "nil":
				p = NULL
			default:
				p = InitializeString(param)
			}
		}

		ps = append(ps, p)
	}

	i := &Instruction{Action: action, Params: ps, Line: line}
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
		params = append(params, param.Inspect())
	}
	return fmt.Sprintf("%s: %s \n", i.Action.Name, strings.Join(params, ", "))
}
