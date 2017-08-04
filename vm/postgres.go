package vm

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func initPGClass(vm *VM) *RClass {
	pg := vm.initializeClass("PG", false)
	pg.setBuiltInMethods(builtInPGClassMethods(), true)
	pg.setBuiltInMethods(builtInPGInstanceMethods(), false)
	vm.objectClass.setClassConstant(pg)

	vm.execGobyLib("db_drivers/pg.gb")

	return pg
}

func builtInPGClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					dataSource, ok := args[1].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(ArgumentError, "Expect postgres' data source to be a String object. got: %s", args[1].Class().Name)
					}

					conn, err := sql.Open("postgres", dataSource.Value)

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					connObj := t.vm.initObjectFromGoType(conn)
					pgDriver := receiver.(*RClass).initializeInstance()
					pgDriver.instanceVariableSet("@connection", connObj)

					return pgDriver
				}
			},
		},
	}
}

func builtInPGInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{}
}
