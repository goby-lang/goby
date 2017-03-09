package evaluator

var (
	StringClass *RString
)

type RString struct {
	*BaseClass
}

type StringObject struct {
	Class *RString
	Value string
}

func (s *StringObject) Type() ObjectType {
	return STRING_OBJ
}

func (s *StringObject) Inspect() string {
	return s.Value
}

func (s *StringObject) ReturnClass() Class {
	return s.Class
}

var (
	stringTable = make(map[string]*StringObject)
)

func InitializeString(value string) *StringObject {
	addr, ok := stringTable[value]

	if !ok {
		s := &StringObject{Value: value, Class: StringClass}
		stringTable[value] = s
		return s
	}

	return addr
}

var builtinStringMethods = []*BuiltInMethod{
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				err := checkArgumentLen(args, StringClass, "+")

				if err != nil {
					return err
				}

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(StringClass)
				}

				rightValue := right.Value
				return &StringObject{Value: leftValue + rightValue, Class: StringClass}
			}
		},
		Name: "+",
	},
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				err := checkArgumentLen(args, StringClass, ">")
				if err != nil {
					return err
				}

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(StringClass)
				}

				rightValue := right.Value

				if leftValue > rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: ">",
	},
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				err := checkArgumentLen(args, StringClass, "<")
				if err != nil {
					return err
				}

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(StringClass)
				}

				rightValue := right.Value

				if leftValue < rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: "<",
	},
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				err := checkArgumentLen(args, StringClass, "==")

				if err != nil {
					return err
				}

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(StringClass)
				}

				rightValue := right.Value

				if leftValue == rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: "==",
	},
	{
		Fn: func(receiver Object) BuiltinMethodBody {
			return func(args []Object, block *Method) Object {
				err := checkArgumentLen(args, StringClass, "!=")

				if err != nil {
					return err
				}

				leftValue := receiver.(*StringObject).Value
				right, ok := args[0].(*StringObject)

				if !ok {
					return wrongTypeError(StringClass)
				}

				rightValue := right.Value

				if leftValue != rightValue {
					return TRUE
				}

				return FALSE
			}
		},
		Name: "!=",
	},
}

func initString() {
	methods := NewEnvironment()

	for _, m := range builtinStringMethods {
		methods.Set(m.Name, m)
	}

	bc := &BaseClass{Name: "String", Methods: methods, Class: ClassClass, SuperClass: ObjectClass}
	sc := &RString{BaseClass: bc}
	StringClass = sc
}
