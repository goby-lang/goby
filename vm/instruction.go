package vm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
)

type operation func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{})

type operationType = string

type setType = string

type action struct {
	name      string
	operation operation
}

type instruction struct {
	action     *action
	Params     []interface{}
	Line       int
	sourceLine int
	argSet     *bytecode.ArgSet
}

type instructionSet struct {
	name         string
	instructions []*instruction
	filename     filename
	paramTypes   *bytecode.ArgSet
}

func (is *instructionSet) define(line int, a *action, params ...interface{}) *instruction {
	i := &instruction{action: a, Params: params, Line: line}
	is.instructions = append(is.instructions, i)
	return i
}

var builtinActions = map[operationType]*action{
	bytecode.Pop: {
		name: bytecode.Pop,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			t.Stack.Pop()
		},
	},
	bytecode.Dup: {
		name: bytecode.Dup,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			obj := t.Stack.top().Target
			t.Stack.Push(&Pointer{Target: obj})
		},
	},
	bytecode.PutBoolean: {
		name: bytecode.PutObject,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			object := t.vm.InitObjectFromGoType(args[0])
			t.Stack.Push(&Pointer{Target: object})
		},
	},
	bytecode.PutObject: {
		name: bytecode.PutObject,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			object := t.vm.InitObjectFromGoType(args[0])
			t.Stack.Push(&Pointer{Target: object})
		},
	},
	bytecode.GetConstant: {
		name: bytecode.GetConstant,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			constName := args[0].(string)
			c := t.vm.lookupConstant(cf, constName)

			if c == nil {
				t.pushErrorObject(errors.NameError, sourceLine, "uninitialized constant %s", constName)
			}

			c.isNamespace = args[1].(string) == "true"

			if t.Stack.top() != nil && t.Stack.top().isNamespace {
				t.Stack.Pop()
			}

			t.Stack.Push(c)
		},
	},
	bytecode.GetLocal: {
		name: bytecode.GetLocal,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			depth := args[0].(int)
			index := args[1].(int)

			p := cf.getLCL(index, depth)

			if p == nil {
				t.Stack.Push(&Pointer{Target: NULL})
				return
			}

			t.Stack.Push(p)
		},
	},
	bytecode.GetInstanceVariable: {
		name: bytecode.GetInstanceVariable,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			variableName := args[0].(string)
			v, ok := cf.self.InstanceVariableGet(variableName)
			if !ok {
				t.Stack.Push(&Pointer{Target: NULL})
				return
			}

			p := &Pointer{Target: v}
			t.Stack.Push(p)
		},
	},
	bytecode.SetInstanceVariable: {
		name: bytecode.SetInstanceVariable,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
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
		},
	},
	bytecode.SetLocal: {
		name: bytecode.SetLocal,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
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
		},
	},
	bytecode.SetConstant: {
		name: bytecode.SetConstant,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			constName := args[0].(string)
			c := cf.lookupConstantInCurrentScope(constName)
			v := t.Stack.Pop()

			if c != nil {
				t.pushErrorObject(errors.ConstantAlreadyInitializedError, sourceLine, "Constant %s already been initialized. Can't assign value to a constant twice.", constName)
			}

			cf.storeConstant(constName, v)
		},
	},
	bytecode.NewRange: {
		name: bytecode.NewRange,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			rangeEnd := t.Stack.Pop().Target.(*IntegerObject).value
			rangeStart := t.Stack.Pop().Target.(*IntegerObject).value

			t.Stack.Push(&Pointer{Target: t.vm.initRangeObject(rangeStart, rangeEnd)})
		},
	},
	bytecode.NewArray: {
		name: bytecode.NewArray,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			argCount := args[0].(int)
			var elems []Object

			for i := 0; i < argCount; i++ {
				v := t.Stack.Pop()
				elems = append([]Object{v.Target}, elems...)
			}

			arr := t.vm.InitArrayObject(elems)
			t.Stack.Push(&Pointer{Target: arr})
		},
	},
	bytecode.ExpandArray: {
		name: bytecode.ExpandArray,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
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
		},
	},
	bytecode.SplatArray: {
		name: bytecode.SplatArray,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			obj := t.Stack.top().Target
			arr, ok := obj.(*ArrayObject)

			if !ok {
				return
			}

			arr.splat = true
		},
	},
	bytecode.NewHash: {
		name: bytecode.NewHash,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			argCount := args[0].(int)
			pairs := map[string]Object{}

			for i := 0; i < argCount/2; i++ {
				v := t.Stack.Pop()
				k := t.Stack.Pop()
				pairs[k.Target.(*StringObject).value] = v.Target
			}

			hash := t.vm.InitHashObject(pairs)
			t.Stack.Push(&Pointer{Target: hash})
		},
	},
	bytecode.BranchUnless: {
		name: bytecode.BranchUnless,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
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
		},
	},
	bytecode.BranchIf: {
		name: bytecode.BranchIf,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
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
		},
	},
	bytecode.Jump: {
		name: bytecode.Jump,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			cf.pc = args[0].(int)
		},
	},
	bytecode.Break: {
		name: bytecode.Break,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
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
		},
	},
	bytecode.PutSelf: {
		name: bytecode.PutSelf,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			t.Stack.Push(&Pointer{Target: cf.self})
		},
	},
	bytecode.PutString: {
		name: bytecode.PutString,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			object := t.vm.InitObjectFromGoType(args[0])
			t.Stack.Push(&Pointer{Target: object})
		},
	},
	bytecode.PutFloat: {
		name: bytecode.PutFloat,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			var value float64
			switch argValue := args[0].(type) {
			case string:
				value, _ = strconv.ParseFloat(argValue, 64)
			case int:
				value = float64(argValue)
			}
			object := t.vm.initFloatObject(value)
			t.Stack.Push(&Pointer{Target: object})
		},
	},
	bytecode.PutNull: {
		name: bytecode.PutNull,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			t.Stack.Push(&Pointer{Target: NULL})
		},
	},
	bytecode.DefMethod: {
		name: bytecode.DefMethod,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			argCount := args[0].(int)
			methodName := t.Stack.Pop().Target.(*StringObject).value
			is, ok := t.getMethodIS(methodName, cf.FileName())

			if !ok {
				t.pushErrorObject(errors.InternalError, sourceLine, "Can't get method %s's instruction set.", methodName)
			}

			method := &MethodObject{Name: methodName, argc: argCount, instructionSet: is, BaseObj: &BaseObj{class: t.vm.TopLevelClass(classes.MethodClass)}}

			v := t.Stack.Pop().Target
			switch self := v.(type) {
			case *RClass:
				self.Methods.set(methodName, method)
			default:
				self.Class().Methods.set(methodName, method)
			}
		},
	},
	bytecode.DefSingletonMethod: {
		name: bytecode.DefSingletonMethod,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			argCount := args[0].(int)
			methodName := t.Stack.Pop().Target.(*StringObject).value
			is, _ := t.getMethodIS(methodName, cf.FileName())
			method := &MethodObject{Name: methodName, argc: argCount, instructionSet: is, BaseObj: &BaseObj{class: t.vm.TopLevelClass(classes.MethodClass)}}

			v := t.Stack.Pop().Target

			switch v := v.(type) {
			case *RClass:
				v.SingletonClass().Methods.set(methodName, method)
			default:
				singletonClass := t.vm.createRClass(fmt.Sprintf("#<Class:#<%s:%d>>", v.Class().Name, v.id()))
				singletonClass.Methods.set(methodName, method)
				singletonClass.isSingleton = true
				v.SetSingletonClass(singletonClass)
			}
		},
	},
	bytecode.DefClass: {
		name: bytecode.DefClass,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
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
						t.pushErrorObject(errors.InternalError, sourceLine, "Constant %s is not a class. got=%s", superClassName, string(superClass.Target.Class().ReturnName()))
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
		},
	},
	bytecode.Send: {
		name: bytecode.Send,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			var method Object

			methodName := args[0].(string)
			argCount := args[1].(int)
			blockFlag := args[2].(string)
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

			// Find Method
			method = receiver.findMethod(methodName)

			if method == nil {
				mm := receiver.findMethodMissing(receiver.Class().inheritsMethodMissing)

				if mm == nil {
					t.setErrorObject(receiverPr, argPr, errors.NoMethodError, sourceLine, errors.UndefinedMethod, methodName, receiver.ToString())
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

			// Find Block
			blockFrame := t.retrieveBlock(cf.FileName(), blockFlag, cf.SourceLine())

			if blockFrame != nil {
				blockFrame.ep = cf
				blockFrame.self = cf.self
				blockFrame.sourceLine = sourceLine
				t.callFrameStack.push(blockFrame)
			}

			switch m := method.(type) {
			case *MethodObject:
				callObj := newCallObject(receiver, m, receiverPr, argCount, argSet, blockFrame, sourceLine)
				t.evalMethodObject(callObj)
			case *BuiltinMethodObject:
				t.evalBuiltinMethod(receiver, m, receiverPr, argCount, argSet, blockFrame, sourceLine, cf.fileName)
			case *Error:
				t.pushErrorObject(errors.InternalError, sourceLine, m.ToString())
			}
		},
	},
	bytecode.InvokeBlock: {
		name: bytecode.InvokeBlock,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			argCount := args[0].(int)
			argPr := t.Stack.pointer - argCount
			receiverPr := argPr - 1
			receiver := t.Stack.data[receiverPr].Target

			if cf.blockFrame == nil {
				t.pushErrorObject(errors.InternalError, sourceLine, "Can't yield without a block")
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
		},
	},
	bytecode.GetBlock: {
		name: bytecode.GetBlock,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			if cf.blockFrame == nil {
				t.pushErrorObject(errors.InternalError, sourceLine, "Can't get block without a block argument")
			}

			blockFrame := cf.blockFrame

			if cf.blockFrame.ep == cf.ep {
				blockFrame = cf.blockFrame.ep.blockFrame
			}

			blockObject := t.vm.initBlockObject(blockFrame.instructionSet, blockFrame.ep, t.Stack.data[t.Stack.pointer-1].Target)

			t.Stack.Push(&Pointer{Target: blockObject})
		},
	},
	bytecode.Leave: {
		name: bytecode.Leave,
		operation: func(t *Thread, sourceLine int, cf *normalCallFrame, args ...interface{}) {
			t.callFrameStack.pop()
			cf.stopExecution()
		},
	},
}

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
