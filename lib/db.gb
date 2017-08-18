#
# Generally the DB class handles the connection and the interaction to the database. It
# establish a connection object which enables developer to interact with the database with
# simple and basic method such as DB#exec and DB#query
#
# *Currently supported DB driver is Postgres
#
class DB
  attr_reader :connection

  def initialize(connection)
    @connection = connection
  end

  #
  # The DB.open method requires a driver_name and a data_source which specifies the type
  # of the DB driver and its' parameter such as the database name and the username ... etc
  #
  # Currently supported DB driver is 'postgres' type DB driver
  #
  # @return[Connection]
  #
  def self.open(driver_name, data_source)
    conn_obj = get_connection(driver_name, data_source)
    connection = Connection.new(conn_obj)
    new(connection)
  end

  #
  # The DB#ping method check whether the connection is established. If the connection is
  # established, it returns true.
  #
  # ```
  #	require "db"
  #
  # db = DB.open("postgres", "user=postgres sslmode=disable") # => Establish a connection
  #	db.ping  # => true
  #
  # db.close # => Closing DB
  # db.ping  # => false
  # ```
  #
  # @return [Boolean]
  #
  def ping
    connection.ping
  end

  #
  # The DB#conn_obj method returns the connection object
  #
  # ```
  #	require "db"
  #
  # db = DB.open("postgres", "user=postgres sslmode=disable") # => Establish a connection
  #	db.conn_obj # Returns connection object
  # ```
  #
  # @return [Object]
  #
  def conn_obj
    connection.conn_obj
  end

  #
  # The Connection class is handles the core connection part to the DB class. It requires a
  # connection object which specifies the information of the DB connection.
  #
  class Connection
    attr_reader :conn_obj

    def initialize(conn_obj)
      @conn_obj = conn_obj
    end

    #
    # The Connection#ping method checks the connection. It returns true if connection has
    # established.
    #
    # @return[Boolean]
    #
    def ping
      err = conn_obj.go_func("Ping")

      if err
        puts(err.go_func("Error"))
        false
      else
        true
      end
    end
  end
end