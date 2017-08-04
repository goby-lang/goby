package vm

func initDBClass(vm *VM) {
	pg := vm.initializeClass("DB", false)
	pg.setBuiltInMethods(builtInDBClassMethods(), true)
	pg.setBuiltInMethods(builtInDBInstanceMethods(), false)
	vm.objectClass.setClassConstant(pg)

	vm.execGobyLib("db.gb")
}

var driverTable = map[string]func(*VM) *RClass{
	"postgres": initPGClass,
}

func builtInDBClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "init_driver",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 2 {
						return t.vm.initErrorObject(ArgumentError, WrongNumberOfArgumentFormat, 2, len(args))
					}

					driverName, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(ArgumentError, "Expect database's driver name to be a String object. got: %s", args[0].Class().Name)
					}

					_, ok = args[1].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(ArgumentError, "Expect database's data source to be a String object. got: %s", args[1].Class().Name)
					}

					driverInitFunc, ok := driverTable[driverName.Value]

					if !ok {
						return t.vm.initErrorObject(InternalError, "Can't find specified driver: %s", driverName.Value)
					}

					driverClass := driverInitFunc(t.vm)
					initMethod := driverClass.singletonClass.lookupMethod("new").(*BuiltInMethodObject)
					driver := initMethod.Fn(driverClass)(t, args, blockFrame)

					return driver
				}
			},
		},
	}
}

func builtInDBInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{}
}
