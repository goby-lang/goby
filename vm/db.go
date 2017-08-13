package vm

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func initDBClass(vm *VM) {
	pg := vm.initializeClass("DB", false)
	pg.setBuiltInMethods(builtInDBClassMethods(), true)
	pg.setBuiltInMethods(builtInDBInstanceMethods(), false)
	vm.objectClass.setClassConstant(pg)

	vm.execGobyLib("db.gb")
}

func getDBConn(t *thread, receiver Object) (*sqlx.DB, error) {
	connection, _ := receiver.instanceVariableGet("@connection")
	connObj, _ := connection.instanceVariableGet("@conn_obj")

	if connObj == NULL {
		return nil, fmt.Errorf("DB connection is nil")
	}

	conn, ok := connObj.(*GoObject).data.(*sqlx.DB)

	if !ok {
		return nil, fmt.Errorf("Connection is not *sql.DB")
	}

	return conn, nil
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

					conn, err := sqlx.Open(driverName.value, dataSource.value)

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
	return []*BuiltInMethodObject{
		{
			Name: "close",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					conn, err := getDBConn(t, receiver)

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					err = conn.Close()

					if err != nil {
						if err != nil {
							return t.vm.initErrorObject(InternalError, "Error happens when closing DB connection: %s", err.Error())
						}
					}

					return TRUE
				}
			},
		},
		{
			Name: "exec",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) < 1 {
						return t.vm.initErrorObject(ArgumentError, "Expect at least 1 argument.")
					}

					conn, err := getDBConn(t, receiver)

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					queryString := args[0].(*StringObject).value
					execArgs := []interface{}{}

					for _, arg := range args[1:] {
						execArgs = append(execArgs, arg.(builtInType).Value())
					}

					// The reason I implement this way: https://github.com/lib/pq/issues/24
					var id int

					err = conn.QueryRow(fmt.Sprintf("%s RETURNING id", queryString), execArgs...).Scan(&id)

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					return t.vm.initIntegerObject(id)
				}
			},
		},
		{
			Name: "query",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) < 1 {
						return t.vm.initErrorObject(ArgumentError, "Expect at least 1 argument.")
					}

					conn, err := getDBConn(t, receiver)

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					queryString := args[0].(*StringObject).value
					execArgs := []interface{}{}

					for _, arg := range args[1:] {
						execArgs = append(execArgs, arg.(builtInType).Value())
					}

					rows, err := conn.Queryx(queryString, execArgs...)

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					results := []Object{}

					for rows.Next() {
						row := make(map[string]interface{})

						err = rows.MapScan(row)

						if err != nil {
							return t.vm.initErrorObject(InternalError, err.Error())
						}

						data := map[string]Object{}

						for k, v := range row {
							data[k] = t.vm.initObjectFromGoType(v)
						}

						result := t.vm.initHashObject(data)
						results = append(results, result)
					}

					return t.vm.initArrayObject(results)
				}
			},
		},
	}
}
