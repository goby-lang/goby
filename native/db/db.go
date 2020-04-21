package db

import (
	"fmt"

	"github.com/goby-lang/goby/vm"
	"github.com/goby-lang/goby/vm/errors"
	"github.com/jmoiron/sqlx"

	// all packages imported by this need postgres
	_ "github.com/lib/pq"
)

type (
	// Object is an imported object from vm
	Object = vm.Object
	// GoObject is an imported object from vm
	GoObject = vm.GoObject
	// VM is an imported object from vm
	VM = vm.VM
	// Thread is an imported object from vm
	Thread = vm.Thread
	// Method is an imported object from vm
	Method = vm.Method
	// StringObject is an imported object from vm
	StringObject = vm.StringObject
)

func init() {
	vm.RegisterExternalClass("db", vm.NewExternalClassLoader("DB", "db.gb",
		// class methods
		map[string]vm.Method{
			"get_connection": getConnection,
		},
		// instance methods
		map[string]vm.Method{
			"query": query,
			"close": closeDB,
			"exec":  exec,
			"run":   run,
		},
	))
}

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
func getConnection(receiver vm.Object, sourceLine int, t *vm.Thread, args []vm.Object) vm.Object {
	if len(args) != 2 {
		return t.VM().InitErrorObject(errors.ArgumentError, sourceLine, errors.WrongNumberOfArgument, 2, len(args))
	}

	driverName, ok := args[0].(*vm.StringObject)

	if !ok {
		return t.VM().InitErrorObject(errors.ArgumentError, sourceLine, "Expect database's driver name to be a String object. got: %s", args[0].Class().Name)
	}

	dataSource, ok := args[1].(*vm.StringObject)

	if !ok {
		return t.VM().InitErrorObject(errors.ArgumentError, sourceLine, "Expect database's data source to be a String object. got: %s", args[1].Class().Name)
	}

	conn, err := sqlx.Open(driverName.Value().(string), dataSource.Value().(string))

	if err != nil {
		return t.VM().InitErrorObject(errors.InternalError, sourceLine, err.Error())
	}

	connObj := t.VM().InitObjectFromGoType(conn)
	return connObj

}

func closeDB(receiver vm.Object, sourceLine int, t *vm.Thread, args []vm.Object) vm.Object {
	conn, err := getDBConn(t, receiver)

	if err != nil {
		return t.VM().InitErrorObject(errors.InternalError, sourceLine, err.Error())
	}

	err = conn.Close()

	if err != nil {
		if err != nil {
			return t.VM().InitErrorObject(errors.InternalError, sourceLine, "Error happens when closing DB connection: %s", err.Error())
		}
	}

	return vm.TRUE

}

func run(receiver Object, sourceLine int, t *Thread, args []Object) Object {
	v := t.VM()
	if len(args) < 1 {
		return v.InitErrorObject(errors.ArgumentError, sourceLine, "Expect at least 1 argument.")
	}

	conn, err := getDBConn(t, receiver)

	if err != nil {
		return v.InitErrorObject(errors.InternalError, sourceLine, err.Error())
	}

	queryString := args[0].(*vm.StringObject).Value().(string)
	execArgs := []interface{}{}

	for _, arg := range args[1:] {
		execArgs = append(execArgs, arg.Value())
	}

	_, err = conn.Exec(queryString, execArgs...)

	if err != nil {
		return v.InitErrorObject(errors.InternalError, sourceLine, err.Error())
	}

	return vm.TRUE

}

// 		{
// 			Name: "get_connection",
// 			Fn: 		},
// 	}
// }

// // Instance methods -----------------------------------------------------
// func builtinDBInstanceMethods() []*BuiltinMethodObject {
// 	return []*BuiltinMethodObject{
// 		{
// 			// The close method closes the connection of the DB instance
// 			//
// 			// ```ruby
// 			// require "db"
// 			//
// 			// db = DB.open("postgres", "user=postgres sslmode=disable")
// 			// db.ping  # => true
// 			//
// 			// db.close # Close the DB connection
// 			// db.ping  # => false
// 			// ```
// 			//
// 			Name: "close",
// 		{
// 			Name: "run",
// 		},
// 		{
// 			// The exec method executes the Psql and automatically returns the data's primary key value
// 			//
// 			// ```ruby
// 			// require "db"
// 			//
// 			// # Assume that there is a User table with name and age column
// 			//
// 			// # Create Opcode
// 			// db = DB.open("postgres", "user=postgres dbname=goby_doc sslmode=disable")
// 			// id = db.exec("INSERT INTO users (name, age) VALUES ('Stan', 23)")
// 			// puts id # => 1
// 			//
// 			// # Update Opcode
// 			//
// 			// id2 = db.exec("INSERT INTO users (name, age) VALUES ('Maxwell', 21)")
// 			// puts id2 # => 2
// 			// id3 = db.exec("UPDATE users SET age=18 WHERE id = $1", id)
// 			// puts id3 # => 2
// 			//
// 			// # Delete Opcode
// 			// id4 = db.exec("DELETE FROM users WHERE id = $1", id3)
// 			// puts id4 # => 2
// 			//
// 			// ```
// 			//
// 			// @return [Integer]
// 			//
// 			Name: "exec",
func exec(receiver Object, sourceLine int, t *Thread, args []Object) Object {
	v := t.VM()
	if len(args) < 1 {
		return v.InitErrorObject(errors.ArgumentError, sourceLine, "Expect at least 1 argument.")
	}

	conn, err := getDBConn(t, receiver)

	if err != nil {
		return v.InitErrorObject(errors.InternalError, sourceLine, err.Error())
	}

	queryString := args[0].(*vm.StringObject).Value().(string)
	execArgs := []interface{}{}

	for _, arg := range args[1:] {
		execArgs = append(execArgs, arg.Value())
	}

	// The reason I implement this way: https://github.com/lib/pq/issues/24
	var id int

	err = conn.QueryRow(fmt.Sprintf("%s RETURNING id", queryString), execArgs...).Scan(&id)

	if err != nil {
		return v.InitErrorObject(errors.InternalError, sourceLine, err.Error())
	}

	return v.InitIntegerObject(id)

}

// 		},
// 		{
// 			// The query method queries the result of the data set
// 			//
// 			// ```ruby
// 			// require "db"
// 			//
// 			// # Assume that there is a User table with name and age column
// 			//
// 			// db = DB.open("postgres", "user=postgres dbname=goby_doc sslmode=disable")
// 			// id = db.exec("INSERT INTO users (name, age) VALUES ('Stan', 23)")
// 			// puts id # => 1
// 			//
// 			// id2 = db.exec("INSERT INTO users (name, age) VALUES ('Maxwell', 21)")
// 			// puts id # => 2
// 			//
// 			// results = db.query("SELECT * FROM users WHERE id = $1", id)
// 			// results.size          # => 1
// 			// results.first[:name]  # => 'Stan'
// 			// results.first[:age]   # => 23
// 			//
// 			// age = 21
// 			// results2 = db.query("SELECT * FROM users WHERE age = $1", age)
// 			// results2.size         # => 1
// 			// results2.first[:name] # => 'Maxwell'
// 			// results2.first[:age]  # => 21
// 			//
// 			// ```
// 			//
// 			// @return [Array]
// 			//
// 			Name: "query",
func query(receiver Object, sourceLine int, t *Thread, args []Object) Object {
	if len(args) < 1 {
		return t.VM().InitErrorObject(errors.ArgumentError, sourceLine, "Expect at least 1 argument.")
	}

	conn, err := getDBConn(t, receiver)

	if err != nil {
		return t.VM().InitErrorObject(errors.InternalError, sourceLine, err.Error())
	}

	queryString := args[0].(*StringObject).Value().(string)
	execArgs := []interface{}{}

	for _, arg := range args[1:] {
		execArgs = append(execArgs, arg.Value())
	}

	rows, err := conn.Queryx(queryString, execArgs...)

	if err != nil {
		return t.VM().InitErrorObject(errors.InternalError, sourceLine, err.Error())
	}

	results := []Object{}

	for rows.Next() {
		row := make(map[string]interface{})

		err = rows.MapScan(row)

		if err != nil {
			return t.VM().InitErrorObject(errors.InternalError, sourceLine, err.Error())
		}

		data := map[string]Object{}

		for k, v := range row {
			data[k] = t.VM().InitObjectFromGoType(v)
		}

		result := t.VM().InitHashObject(data)
		results = append(results, result)
	}

	return t.VM().InitArrayObject(results)

}

func getDBConn(t *vm.Thread, receiver Object) (*sqlx.DB, error) {
	connection, _ := receiver.InstanceVariableGet("@connection")
	connObj, _ := connection.InstanceVariableGet("@conn_obj")

	if connObj == vm.NULL {
		return nil, fmt.Errorf("DB connection is nil")
	}

	conn, ok := connObj.(*GoObject).Value().(*sqlx.DB)

	if !ok {
		return nil, fmt.Errorf("Connection is not *sql.DB")
	}

	return conn, nil
}
