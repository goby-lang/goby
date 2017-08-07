package vm

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func initDBClass(vm *VM) {
	pg := vm.initializeClass("DB", false)
	pg.setBuiltInMethods(builtInDBClassMethods(), true)
	pg.setBuiltInMethods(builtInDBInstanceMethods(), false)
	vm.objectClass.setClassConstant(pg)

	vm.execGobyLib("db.gb")
}

func builtInDBClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "get_connection",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 2 {
						return t.vm.initErrorObject(ArgumentError, WrongNumberOfArgumentFormat, 2, len(args))
					}

					driverName, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(ArgumentError, "Expect database's driver name to be a String object. got: %s", args[0].Class().Name)
					}

					dataSource, ok := args[1].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(ArgumentError, "Expect database's data source to be a String object. got: %s", args[1].Class().Name)
					}

					conn, err := sql.Open(driverName.value, dataSource.value)

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					connObj := t.vm.initObjectFromGoType(conn)
					return connObj
				}
			},
		},
	}
}

func builtInDBInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{}
}
