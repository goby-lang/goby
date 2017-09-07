package vm

import (
	"fmt"
	"github.com/goby-lang/goby/vm/errors"
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
			// The get_connection method returns a connection object which requires the name of the driver
			// and the source which specifies the parameter including the name of the database and the
			// username ...etc.
			//
			// Currently supported DB driver is 'postgres'
			//
			// (The example is the DB#open class method which is implemented in db.gb file)
			//
			// ```ruby
			// class DB
			//   def self.open(driver_name, data_source)
			//	   conn_obj = get_connection(driver_name, data_source) # => Returns the Conn object
			//	   connection = Connection.new(conn_obj)
			//	   new(connection)
			//   end
			//
			//   # ... Omitted
			// ```
			//
			// @return [Object]
			//
			Name: "get_connection",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 2 {
						return t.vm.initErrorObject(errors.ArgumentError, errors.WrongNumberOfArgumentFormat, 2, len(args))
					}

					driverName, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect database's driver name to be a String object. got: %s", args[0].Class().Name)
					}

					dataSource, ok := args[1].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect database's data source to be a String object. got: %s", args[1].Class().Name)
					}

					conn, err := sqlx.Open(driverName.value, dataSource.value)

					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
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
			// The close method closes the connection of the DB instance
			//
			// ```ruby
			// require "db"
			//
			// db = DB.open("postgres", "user=postgres sslmode=disable")
			// db.ping  # => true
			//
			// db.close # Close the DB connection
			// db.ping  # => false
			// ```
			//
			Name: "close",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					conn, err := getDBConn(t, receiver)

					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
					}

					err = conn.Close()

					if err != nil {
						if err != nil {
							return t.vm.initErrorObject(errors.InternalError, "Error happens when closing DB connection: %s", err.Error())
						}
					}

					return TRUE
				}
			},
		},
		{
			Name: "run",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) < 1 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect at least 1 argument.")
					}

					conn, err := getDBConn(t, receiver)

					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
					}

					queryString := args[0].(*StringObject).value
					execArgs := []interface{}{}

					for _, arg := range args[1:] {
						execArgs = append(execArgs, arg.(builtInType).Value())
					}

					_, err = conn.Exec(queryString, execArgs...)

					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
					}

					return TRUE
				}
			},
		},
		{
			// The exec method executes the Psql and automatically returns the data's primary key value
			//
			// ```ruby
			// require "db"
			//
			// # Assume that there is a User table with name and age column
			//
			// # Create Action
			// db = DB.open("postgres", "user=postgres dbname=goby_doc sslmode=disable")
			// id = db.exec("INSERT INTO users (name, age) VALUES ('Stan', 23)")
			// puts id # => 1
			//
			// # Update Action
			//
			// id2 = db.exec("INSERT INTO users (name, age) VALUES ('Maxwell', 21)")
			// puts id2 # => 2
			// id3 = db.exec("UPDATE users SET age=18 WHERE id = $1", id)
			// puts id3 # => 2
			//
			// # Delete Action
			// id4 = db.exec("DELETE FROM users WHERE id = $1", id3)
			// puts id4 # => 2
			//
			// ```
			//
			// @return [Integer]
			//
			Name: "exec",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) < 1 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect at least 1 argument.")
					}

					conn, err := getDBConn(t, receiver)

					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
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
						return t.vm.initErrorObject(errors.InternalError, err.Error())
					}

					return t.vm.initIntegerObject(id)
				}
			},
		},
		{
			// The query method queries the result of the data set
			//
			// ```ruby
			// require "db"
			//
			// # Assume that there is a User table with name and age column
			//
			// db = DB.open("postgres", "user=postgres dbname=goby_doc sslmode=disable")
			// id = db.exec("INSERT INTO users (name, age) VALUES ('Stan', 23)")
			// puts id # => 1
			//
			// id2 = db.exec("INSERT INTO users (name, age) VALUES ('Maxwell', 21)")
			// puts id # => 2
			//
			// results = db.query("SELECT * FROM users WHERE id = $1", id)
			// results.size          # => 1
			// results.first[:name]  # => 'Stan'
			// results.first[:age]   # => 23
			//
			// age = 21
			// results2 = db.query("SELECT * FROM users WHERE age = $1", age)
			// results2.size         # => 1
			// results2.first[:name] # => 'Maxwell'
			// results2.first[:age]  # => 21
			//
			// ```
			//
			// @return [Array]
			//
			Name: "query",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) < 1 {
						return t.vm.initErrorObject(errors.ArgumentError, "Expect at least 1 argument.")
					}

					conn, err := getDBConn(t, receiver)

					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
					}

					queryString := args[0].(*StringObject).value
					execArgs := []interface{}{}

					for _, arg := range args[1:] {
						execArgs = append(execArgs, arg.(builtInType).Value())
					}

					rows, err := conn.Queryx(queryString, execArgs...)

					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
					}

					results := []Object{}

					for rows.Next() {
						row := make(map[string]interface{})

						err = rows.MapScan(row)

						if err != nil {
							return t.vm.initErrorObject(errors.InternalError, err.Error())
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
